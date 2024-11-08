package cli

import (
	"context"

	"github.com/DefangLabs/defang/src/pkg/cli/client"
	"github.com/DefangLabs/defang/src/pkg/term"
	defangv1 "github.com/DefangLabs/defang/src/protos/io/defang/v1"
)

func DeploymentsList(ctx context.Context, loader client.Loader, provider client.Provider) error {
	projectName, err := LoadProjectName(ctx, loader, provider)
	if err != nil {
		return err
	}
	term.Debugf("Listing deployments in project %q", projectName)

	deployments, err := provider.ListDeployments(ctx, &defangv1.ListDeploymentsRequest{Project: projectName})
	if err != nil {
		return err
	}

	return PrintObject("", deployments)
}
