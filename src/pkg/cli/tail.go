package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/defang-io/defang/src/pkg"
	"github.com/defang-io/defang/src/pkg/cli/client"
	"github.com/defang-io/defang/src/pkg/spinner"
	"github.com/defang-io/defang/src/pkg/term"
	"github.com/defang-io/defang/src/pkg/types"
	defangv1 "github.com/defang-io/defang/src/protos/io/defang/v1"
	"github.com/muesli/termenv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ansiCyan      = "\033[36m"
	ansiReset     = "\033[0m"
	replaceString = ansiCyan + "$0" + ansiReset
	RFC3339Micro  = "2006-01-02T15:04:05.000000Z07:00" // like RFC3339Nano but with 6 digits of precision
)

var (
	colorKeyRegex = regexp.MustCompile(`"(?:\\["\\/bfnrt]|[^\x00-\x1f"\\]|\\u[0-9a-fA-F]{4})*"\s*:|[^\x00-\x20"=&?]+=`) // handles JSON, logfmt, and query params
	DoVerbose     = false
)

type TailDetectStopEventFunc func(service string, host string, eventlog string) bool
type ServiceUpNotifier <-chan types.ServiceStatus

type TailOptions struct {
	Service            string
	Etag               types.ETag
	Since              time.Time
	Raw                bool
	EndEventDetectFunc TailDetectStopEventFunc
	ServiceUpNotifier  ServiceUpNotifier
}

type P = client.Property // shorthand for tracking properties

// ParseTimeOrDuration parses a time string or duration string (e.g. 1h30m) and returns a time.Time.
// At a minimum, this function supports RFC3339Nano, Go durations, and our own TimestampFormat (local).
func ParseTimeOrDuration(str string) (time.Time, error) {
	if strings.ContainsAny(str, "TZ") {
		return time.Parse(time.RFC3339Nano, str)
	}
	if strings.Contains(str, ":") {
		local, err := time.ParseInLocation("15:04:05.999999", str, time.Local)
		if err != nil {
			return time.Time{}, err
		}
		// Replace the year, month, and day of t with today's date
		now := time.Now().Local()
		sincet := time.Date(now.Year(), now.Month(), now.Day(), local.Hour(), local.Minute(), local.Second(), local.Nanosecond(), local.Location())
		if sincet.After(now) {
			sincet = sincet.AddDate(0, 0, -1) // yesterday; subtract 1 day
		}
		return sincet, nil
	}
	dur, err := time.ParseDuration(str)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().Add(-dur), nil // - because we want to go back in time
}

type CancelError struct {
	Service string
	Etag    string
	Last    time.Time
	error
}

func (cerr *CancelError) Error() string {
	cmd := "tail --since " + cerr.Last.UTC().Format(time.RFC3339Nano)
	if cerr.Service != "" {
		cmd += " --name " + cerr.Service
	}
	if cerr.Etag != "" {
		cmd += " --etag " + cerr.Etag
	}
	if DoVerbose {
		cmd += " --verbose"
	}
	return cmd
}

func (cerr *CancelError) Unwrap() error {
	return cerr.error
}

func onServiceUpHandler(notifier ServiceUpNotifier, cancel context.CancelFunc) {
	for event := range notifier {
		switch event {
		case types.ServicesStarting:
			fmt.Println("Service starting")
		case types.ServicesAllStarted:
			fmt.Println("Service all started")
			cancel()
			return
		case types.ServicesAllStopped:
			fmt.Println("Service all stopped")
		case types.ServicesStopping:
			fmt.Println("Service stoping")
		}
	}
}

func Tail(ctx context.Context, client client.Client, params TailOptions) error {
	if params.Service != "" {
		params.Service = NormalizeServiceName(params.Service)
		// Show a warning if the service doesn't exist (yet);; TODO: could do fuzzy matching and suggest alternatives
		if _, err := client.Get(ctx, &defangv1.ServiceID{Name: params.Service}); err != nil {
			switch connect.CodeOf(err) {
			case connect.CodeNotFound:
				term.Warn(" ! Service does not exist (yet):", params.Service)
			case connect.CodeUnknown:
				// Ignore unknown (nil) errors
			default:
				term.Warn(" !", err)
			}
		}
	}

	projectName, err := client.LoadProjectName()
	if err != nil {
		return err
	}
	term.Debug(" - Tailing logs in project", projectName)

	if DoDryRun {
		return ErrDryRun
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if params.ServiceUpNotifier != nil {
		onServiceUpHandler(params.ServiceUpNotifier, cancel)
	}

	serverStream, err := client.Tail(ctx, &defangv1.TailRequest{Service: params.Service, Etag: params.Etag, Since: timestamppb.New(params.Since)})
	if err != nil {
		return err
	}
	defer serverStream.Close() // this works because it takes a pointer receiver

	spin := spinner.New()
	doSpinner := !params.Raw && term.CanColor && term.IsTerminal

	if term.IsTerminal && !params.Raw {
		if doSpinner {
			term.Stdout.HideCursor()
			defer term.Stdout.ShowCursor()
		}

		if !DoVerbose {
			// Allow the user to toggle verbose mode with the V key
			if oldState, err := term.MakeUnbuf(int(os.Stdin.Fd())); err == nil {
				defer term.Restore(int(os.Stdin.Fd()), oldState)

				term.Info(" * Press V to toggle verbose mode")
				input := term.NewNonBlockingStdin()
				defer input.Close() // abort the read loop
				go func() {
					var b [1]byte
					for {
						if _, err := input.Read(b[:]); err != nil {
							return // exit goroutine
						}
						switch b[0] {
						case 3: // Ctrl-C
							cancel() // cancel the tail context
						case 10, 13: // Enter or Return
							fmt.Println(" ") // empty line, but overwrite the spinner
						case 'v', 'V':
							verbose := !DoVerbose
							DoVerbose = verbose
							modeStr := "OFF"
							if verbose {
								modeStr = "ON"
							}
							term.Info(" * Verbose mode", modeStr)
							go client.Track("Verbose Toggled", P{"verbose", verbose})
						}
					}
				}()
			}
		}
	}

	skipDuplicate := false
	for {
		if !serverStream.Receive() {
			if errors.Is(serverStream.Err(), context.Canceled) {
				return &CancelError{Service: params.Service, Etag: params.Etag, Last: params.Since, error: serverStream.Err()}
			}

			// TODO: detect ALB timeout (504) or Fabric restart and reconnect automatically
			code := connect.CodeOf(serverStream.Err())
			// Reconnect on Error: internal: stream error: stream ID 5; INTERNAL_ERROR; received from peer
			if code == connect.CodeUnavailable || (code == connect.CodeInternal && !connect.IsWireError(serverStream.Err())) {
				term.Debug(" - Disconnected:", serverStream.Err())
				if !params.Raw {
					term.Fprint(term.Stderr, term.WarnColor, " ! Reconnecting...\r") // overwritten below
				}
				time.Sleep(time.Second)
				serverStream, err = client.Tail(ctx, &defangv1.TailRequest{Service: params.Service, Etag: params.Etag, Since: timestamppb.New(params.Since)})
				if err != nil {
					term.Debug(" - Reconnect failed:", err)
					return err
				}
				if !params.Raw {
					term.Fprint(term.Stderr, term.WarnColor, " ! Reconnected!   \r") // overwritten with logs
				}
				skipDuplicate = true
				continue
			}

			return serverStream.Err() // returns nil on EOF
		}
		msg := serverStream.Msg()

		// Show a spinner if we're not in raw mode and have a TTY
		if doSpinner {
			fmt.Print(spin.Next())
		}

		// HACK: skip noisy CI/CD logs (except errors)
		isInternal := msg.Service == "cd" || msg.Service == "ci" || msg.Service == "kaniko" || msg.Service == "fabric" || msg.Host == "kaniko" || msg.Host == "fabric"
		onlyErrors := !DoVerbose && isInternal
		for _, e := range msg.Entries {
			if onlyErrors && !e.Stderr {
				if params.EndEventDetectFunc != nil && params.EndEventDetectFunc(msg.Service, msg.Host, e.Message) {
					return nil
				}
				continue
			}

			ts := e.Timestamp.AsTime()
			if skipDuplicate && ts.Equal(params.Since) {
				skipDuplicate = false
				continue
			}
			if ts.After(params.Since) {
				params.Since = ts
			}

			if params.Raw {
				out := term.Stdout
				if e.Stderr {
					out = term.Stderr
				}
				fmt.Fprintln(out, e.Message) // TODO: trim trailing newline because we're already printing one?
				continue
			}

			// Replace service progress messages with our own spinner
			if doSpinner && isProgressDot(e.Message) {
				continue
			}

			tsString := ts.Local().Format(RFC3339Micro)
			tsColor := termenv.ANSIWhite
			if e.Stderr {
				tsColor = termenv.ANSIBrightRed
			}
			var prefixLen int
			trimmed := strings.TrimRight(e.Message, "\t\r\n ")
			for i, line := range strings.Split(trimmed, "\n") {
				if i == 0 {
					prefixLen, _ = term.Print(tsColor, tsString, " ")
					if params.Etag == "" {
						l, _ := term.Print(termenv.ANSIYellow, msg.Etag, " ")
						prefixLen += l
					}
					if params.Service == "" {
						l, _ := term.Print(termenv.ANSIGreen, msg.Service, " ")
						prefixLen += l
					}
					if DoVerbose {
						l, _ := term.Print(termenv.ANSIMagenta, msg.Host, " ")
						prefixLen += l
					}
				} else {
					fmt.Print(strings.Repeat(" ", prefixLen))
				}
				if term.CanColor {
					if !strings.Contains(line, "\033[") {
						line = colorKeyRegex.ReplaceAllString(line, replaceString) // add some color
					}
					term.Stdout.Reset()
				} else {
					line = pkg.StripAnsi(line)
				}
				fmt.Println(line)

				// Detect end logging event
				if params.EndEventDetectFunc != nil && params.EndEventDetectFunc(msg.Service, msg.Host, line) {
					cancel()
					return nil
				}
			}
		}
	}
}

func isProgressDot(line string) bool {
	if len(line) <= 1 {
		return true
	}
	stripped := pkg.StripAnsi(line)
	for _, r := range stripped {
		if r != '.' {
			return false
		}
	}
	return true
}
