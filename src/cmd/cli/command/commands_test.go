package command

import (
	"context"
	"testing"
)

func TestVersion(t *testing.T) {
	err := testCommand([]string{"version"})
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}
}

func testCommand(args []string) error {
	ctx := context.TODO()
	SetupCommands("test")
	RootCmd.SetArgs(args)
	return RootCmd.ExecuteContext(ctx)
}
