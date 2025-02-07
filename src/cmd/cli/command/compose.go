package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DefangLabs/defang/src/pkg/cli"
	"github.com/DefangLabs/defang/src/pkg/cli/compose"
	"github.com/DefangLabs/defang/src/pkg/term"
	"github.com/DefangLabs/defang/src/pkg/track"
	"github.com/DefangLabs/defang/src/pkg/types"
	defangv1 "github.com/DefangLabs/defang/src/protos/io/defang/v1"
	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
)

func isManagedService(service compose.ServiceConfig) bool {
	if service.Extensions == nil {
		return false
	}

	return service.Extensions["x-defang-static-files"] != nil || service.Extensions["x-defang-redis"] != nil || service.Extensions["x-defang-postgres"] != nil
}

func splitManagedAndUnmanagedServices(serviceInfos compose.Services) ([]string, []string) {
	var managedServices []string
	var unmanagedServices []string
	for _, service := range serviceInfos {
		if isManagedService(service) {
			managedServices = append(managedServices, service.Name)
		} else {
			unmanagedServices = append(unmanagedServices, service.Name)
		}
	}

	return managedServices, unmanagedServices
}

func makeComposeUpCmd() *cobra.Command {
	mode := Mode(defangv1.DeploymentMode_DEVELOPMENT)
	composeUpCmd := &cobra.Command{
		Use:         "up",
		Annotations: authNeededAnnotation,
		Args:        cobra.NoArgs, // TODO: takes optional list of service names
		Short:       "Reads a Compose file and deploy a new project or update an existing project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var force, _ = cmd.Flags().GetBool("force")
			var detach, _ = cmd.Flags().GetBool("detach")
			var waitTimeout, _ = cmd.Flags().GetInt("wait-timeout")

			upload := compose.UploadModeDigest
			if force {
				upload = compose.UploadModeForce
			}

			since := time.Now()
			loader := configureLoader(cmd)
			provider, err := getProvider(cmd.Context())
			if err != nil {
				return err
			}
			deploy, project, err := cli.ComposeUp(cmd.Context(), loader, client, provider, upload, mode.Value())

			if err != nil {
				if !nonInteractive && strings.Contains(err.Error(), "maximum number of projects") {
					if resp, err2 := provider.GetServices(cmd.Context(), &defangv1.GetServicesRequest{Project: project.Name}); err2 == nil {
						term.Error("Error:", prettyError(err))
						if _, err := cli.InteractiveComposeDown(cmd.Context(), provider, resp.Project); err != nil {
							term.Debug("ComposeDown failed:", err)
							printDefangHint("To deactivate a project, do:", "compose down --project-name "+resp.Project)
						} else {
							printDefangHint("To try deployment again, do:", "compose up")
						}
						return nil
					}
				}
				if errors.Is(err, types.ErrComposeFileNotFound) {
					printDefangHint("To start a new project, do:", "new")
				}
				return err
			}

			if len(deploy.Services) == 0 {
				return errors.New("no services being deployed")
			}

			printPlaygroundPortalServiceURLs(deploy.Services)

			managedServices, unmanagedServices := splitManagedAndUnmanagedServices(project.Services)

			if len(managedServices) > 0 {
				term.Warnf("Defang cannot monitor status of the following managed service(s): %v.\n   To check if the managed service is up, check the status of the service which depends on it.", managedServices)
			}

			if detach {
				term.Info("Detached.")
				return nil
			}

			tailCtx, cancelTail := context.WithCancelCause(cmd.Context())
			defer cancelTail(nil) // to cancel WaitServiceState and clean-up context

			if waitTimeout >= 0 {
				var cancelTimeout context.CancelFunc
				tailCtx, cancelTimeout = context.WithTimeout(tailCtx, time.Duration(waitTimeout)*time.Second)
				defer cancelTimeout()
			}

			errCompleted := errors.New("deployment succeeded") // tail canceled because of deployment completion
			const targetState = defangv1.ServiceState_DEPLOYMENT_COMPLETED

			go func() {
				if err := cli.WaitServiceState(tailCtx, provider, targetState, deploy.Etag, unmanagedServices); err != nil {
					var errDeploymentFailed cli.ErrDeploymentFailed
					if errors.As(err, &errDeploymentFailed) {
						cancelTail(err)
					} else if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
						term.Warnf("failed to wait for service status: %v", err) // TODO: don't print in Go-routine
					}
				} else {
					cancelTail(errCompleted)
				}
			}()

			// show users the current streaming logs
			tailSource := "all services"
			if deploy.Etag != "" {
				tailSource = "deployment ID " + deploy.Etag
			}

			term.Info("Tailing logs for", tailSource, "; press Ctrl+C to detach:")
			tailParams := cli.TailOptions{
				Etag:    deploy.Etag,
				Since:   since,
				Raw:     false,
				Verbose: verbose,
			}

			// blocking call to tail
			if err := cli.Tail(tailCtx, loader, provider, tailParams); err != nil {
				term.Debug("Tail stopped with", err)

				if connect.CodeOf(err) == connect.CodePermissionDenied {
					// If tail fails because of missing permission, we wait for the deployment to finish
					term.Warn("Unable to tail logs. Waiting for the deployment to finish.")
					<-tailCtx.Done()
					// Get the actual error from the context so we won't print "Error: missing tail permission"
					err = context.Cause(tailCtx)
				} else if !(errors.Is(tailCtx.Err(), context.Canceled) || errors.Is(tailCtx.Err(), context.DeadlineExceeded)) {
					return err // any error other than cancelation
				}

				// The tail was canceled; check if it was because of deployment failure or explicit cancelation or wait-timeout reached
				if errors.Is(context.Cause(tailCtx), context.Canceled) {
					// Tail was canceled by the user before deployment completion/failure; show a warning and exit with an error
					term.Warn("Deployment is not finished. Service(s) might not be running.")
					return err
				} else if errors.Is(context.Cause(tailCtx), context.DeadlineExceeded) {
					// Tail was canceled when wait-timeout is reached; show a warning and exit with an error
					term.Warn("Wait-timeout exceeded, detaching from logs. Deployment still in progress.")
					return err
				}

				var errDeploymentFailed cli.ErrDeploymentFailed
				if errors.As(context.Cause(tailCtx), &errDeploymentFailed) {
					// Tail got canceled because of deployment failure: prompt to show the debugger
					term.Warn(errDeploymentFailed)
					if !nonInteractive {
						failedServices := []string{errDeploymentFailed.Service}
						track.Evt("Debug Prompted", P("failedServices", failedServices), P("etag", deploy.Etag), P("reason", errDeploymentFailed))
						// Call the AI debug endpoint using the original command context (not the tailCtx which is canceled)
						_ = cli.InteractiveDebug(cmd.Context(), loader, client, provider, deploy.Etag, project, failedServices)
					}
					return err
				}
			}

			// Print the current service states of the deployment
			if errors.Is(context.Cause(tailCtx), errCompleted) {
				for _, service := range deploy.Services {
					service.State = targetState
				}

				printEndpoints(deploy.Services)
			}

			term.Info("Done.")
			return nil
		},
	}
	composeUpCmd.Flags().BoolP("detach", "d", false, "run in detached mode")
	composeUpCmd.Flags().Bool("force", false, "force a build of the image even if nothing has changed")
	composeUpCmd.Flags().Bool("tail", false, "tail the service logs after updating") // obsolete, but keep for backwards compatibility
	_ = composeUpCmd.Flags().MarkHidden("tail")
	composeUpCmd.Flags().VarP(&mode, "mode", "m", "deployment mode, possible values: "+strings.Join(allModes(), ", "))
	composeUpCmd.Flags().Bool("build", true, "build the image before starting the service") // docker-compose compatibility
	_ = composeUpCmd.Flags().MarkHidden("build")
	composeUpCmd.Flags().Bool("wait", true, "wait for services to be running|healthy") // docker-compose compatibility
	_ = composeUpCmd.Flags().MarkHidden("wait")
	composeUpCmd.Flags().Int("wait-timeout", -1, "maximum duration to wait for the project to be running|healthy") // docker-compose compatibility
	return composeUpCmd
}

func makeComposeStartCmd() *cobra.Command {
	composeStartCmd := &cobra.Command{
		Use:         "start",
		Aliases:     []string{"deploy"},
		Annotations: authNeededAnnotation,
		Args:        cobra.NoArgs, // TODO: takes optional list of service names
		Short:       "Reads a Compose file and deploys services to the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Command 'start' is deprecated, use 'up' instead")
		},
	}
	composeStartCmd.Flags().Bool("force", false, "force a build of the image even if nothing has changed")
	return composeStartCmd
}

func makeComposeRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "restart",
		Annotations: authNeededAnnotation,
		Args:        cobra.NoArgs, // TODO: takes optional list of service names
		Short:       "Reads a Compose file and restarts its services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Command 'restart' is deprecated, use 'up' instead")
		},
	}
}

func makeComposeStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "stop",
		Annotations: authNeededAnnotation,
		Args:        cobra.NoArgs, // TODO: takes optional list of service names
		Short:       "Reads a Compose file and stops its services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Command 'stop' is deprecated, use 'down' instead")
		},
	}
}

func makeComposeDownCmd() *cobra.Command {
	composeDownCmd := &cobra.Command{
		Use:         "down [SERVICE...]",
		Aliases:     []string{"rm", "remove"}, // like docker stack
		Annotations: authNeededAnnotation,
		Short:       "Reads a Compose file and deprovisions its services",
		RunE: func(cmd *cobra.Command, args []string) error {
			var detach, _ = cmd.Flags().GetBool("detach")

			loader := configureLoader(cmd)
			provider, err := getProvider(cmd.Context())
			if err != nil {
				return err
			}
			since := time.Now()
			etag, err := cli.ComposeDown(cmd.Context(), loader, client, provider, args...)
			if err != nil {
				if connect.CodeOf(err) == connect.CodeNotFound {
					// Show a warning (not an error) if the service was not found
					term.Warn(prettyError(err))
					return nil
				}
				return err
			}

			term.Info("Deleted services, deployment ID", etag)

			if detach {
				printDefangHint("To track the update, do:", "tail --etag "+etag)
				return nil
			}

			endLogConditions := []cli.EndLogConditional{
				{Service: "cd", Host: "pulumi", EventLog: "Destroy succeeded in "},
				{Service: "cd", Host: "pulumi", EventLog: "Update succeeded in "},
			}

			endLogDetectFunc := cli.CreateEndLogEventDetectFunc(endLogConditions)
			tailParams := cli.TailOptions{
				Etag:               etag,
				Since:              since,
				Raw:                false,
				EndEventDetectFunc: endLogDetectFunc,
				Verbose:            verbose,
			}

			err = cli.Tail(cmd.Context(), loader, provider, tailParams)
			if err != nil {
				return err
			}
			term.Info("Done.")
			return nil
		},
	}
	composeDownCmd.Flags().BoolP("detach", "d", false, "run in detached mode")
	composeDownCmd.Flags().Bool("tail", false, "tail the service logs after deleting") // obsolete, but keep for backwards compatibility
	_ = composeDownCmd.Flags().MarkHidden("tail")
	return composeDownCmd
}

func makeComposeConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Args:  cobra.NoArgs, // TODO: takes optional list of service names
		Short: "Reads a Compose file and shows the generated config",
		RunE: func(cmd *cobra.Command, args []string) error {
			loader := configureLoader(cmd)
			provider, err := getProvider(cmd.Context())
			if err != nil {
				return err
			}
			if _, _, err := cli.ComposeUp(cmd.Context(), loader, client, provider, compose.UploadModeIgnore, defangv1.DeploymentMode_UNSPECIFIED_MODE); !errors.Is(err, cli.ErrDryRun) {
				return err
			}
			return nil
		},
	}
}

func makeComposeLsCmd() *cobra.Command {
	getServicesCmd := &cobra.Command{
		Use:         "ps",
		Annotations: authNeededAnnotation,
		Args:        cobra.NoArgs,
		Aliases:     []string{"getServices", "services"},
		Short:       "Get list of services in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			long, _ := cmd.Flags().GetBool("long")

			loader := configureLoader(cmd)
			provider, err := getProvider(cmd.Context())
			if err != nil {
				return err
			}

			if err := cli.GetServices(cmd.Context(), loader, provider, long); err != nil {
				if errNoServices := new(cli.ErrNoServices); !errors.As(err, errNoServices) {
					return err
				}

				term.Warn(err)

				printDefangHint("To start a new project, do:", "new")
				return nil
			}

			if !long {
				printDefangHint("To see more information about your services, do:", cmd.CalledAs()+" -l")
			}
			return nil
		},
	}
	getServicesCmd.Flags().BoolP("long", "l", false, "show more details")
	return getServicesCmd
}

func makeComposeLogsCmd() *cobra.Command {
	var logsCmd = &cobra.Command{
		Use:         "logs",
		Annotations: authNeededAnnotation,
		Aliases:     []string{"tail"},
		Args:        cobra.NoArgs,
		Short:       "Tail logs from one or more services",
		RunE: func(cmd *cobra.Command, args []string) error {
			var name, _ = cmd.Flags().GetString("name")
			var etag, _ = cmd.Flags().GetString("etag")
			var raw, _ = cmd.Flags().GetBool("raw")
			var since, _ = cmd.Flags().GetString("since")
			var utc, _ = cmd.Flags().GetBool("utc")

			if utc {
				os.Setenv("TZ", "") // used by Go's "time" package, see https://pkg.go.dev/time#Location
			}

			ts, err := cli.ParseTimeOrDuration(since, time.Now())
			if err != nil {
				return fmt.Errorf("invalid duration or time: %w", err)
			}

			ts = ts.UTC()
			sinceStr := ""
			if ts.Year() > 1970 {
				sinceStr = " since " + ts.Format(time.RFC3339Nano) + " "
			}
			term.Infof("Showing logs%s; press Ctrl+C to stop:", sinceStr)
			services := []string{}
			if len(name) > 0 {
				services = strings.Split(name, ",")
			}
			tailOptions := cli.TailOptions{
				Services: services,
				Etag:     etag,
				Since:    ts,
				Raw:      raw,
				Verbose:  true, // always verbose for explicit tail command
			}

			loader := configureLoader(cmd)
			provider, err := getProvider(cmd.Context())
			if err != nil {
				return err
			}
			return cli.Tail(cmd.Context(), loader, provider, tailOptions)
		},
	}
	logsCmd.Flags().StringP("name", "n", "", "name of the service")
	logsCmd.Flags().String("etag", "", "deployment ID (ETag) of the service")
	logsCmd.Flags().BoolP("raw", "r", false, "show raw (unparsed) logs")
	logsCmd.Flags().StringP("since", "S", "", "show logs since duration/time")
	logsCmd.Flags().Bool("utc", false, "show logs in UTC timezone (ie. TZ=UTC)")
	return logsCmd
}

func setupComposeCommand() *cobra.Command {
	var composeCmd = &cobra.Command{
		Use:     "compose",
		Aliases: []string{"stack"},
		Args:    cobra.NoArgs,
		Short:   "Work with local Compose files",
		Long: `Define and deploy multi-container applications with Defang. Most compose commands require
a "compose.yaml" file. The simplest "compose.yaml" file with a single service is:

services:
  app:              # the name of the service
    build: .        # the folder with the Dockerfile and app sources (. means current folder)
    ports:
      - 80          # the port the service listens on for HTTP requests
`,
	}
	// Compose Command
	// composeCmd.Flags().Bool("compatibility", false, "Run compose in backward compatibility mode"); TODO: Implement compose option
	// composeCmd.Flags().String("env-file", "", "Specify an alternate environment file."); TODO: Implement compose option
	// composeCmd.Flags().Int("parallel", -1, "Control max parallelism, -1 for unlimited (default -1)"); TODO: Implement compose option
	// composeCmd.Flags().String("profile", "", "Specify a profile to enable"); TODO: Implement compose option
	// composeCmd.Flags().String("project-directory", "", "Specify an alternate working directory"); TODO: Implement compose option
	composeCmd.AddCommand(makeComposeUpCmd())
	composeCmd.AddCommand(makeComposeConfigCmd())
	composeCmd.AddCommand(makeComposeDownCmd())
	composeCmd.AddCommand(makeComposeLsCmd())
	composeCmd.AddCommand(makeComposeLogsCmd())

	// deprecated, will be removed in future releases
	composeCmd.AddCommand(makeComposeStartCmd())
	composeCmd.AddCommand(makeComposeRestartCmd())
	composeCmd.AddCommand(makeComposeStopCmd())
	return composeCmd
}
