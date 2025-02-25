package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/DefangLabs/defang/src/pkg/cli/client"
	"github.com/DefangLabs/defang/src/pkg/cli/compose"
	defangv1 "github.com/DefangLabs/defang/src/protos/io/defang/v1"
)

type mockGetServicesProvider struct {
	client.Provider
}

func (mockGetServicesProvider) LoadProjectName(ctx context.Context) (string, error) {
	return "TestGetServices", nil
}

func (mockGetServicesProvider) GetServices(ctx context.Context, req *defangv1.GetServicesRequest) (*defangv1.ListServicesResponse, error) {
	return &defangv1.ListServicesResponse{}, nil
}

func TestGetServices(t *testing.T) {
	ctx := context.Background()
	loader := client.MockLoader{Project: &compose.Project{Name: "TestGetServices"}}
	provider := mockGetServicesProvider{}

	t.Run("ErrNoServices", func(t *testing.T) {
		err := GetServices(ctx, loader, provider, false)
		if err == nil {
			t.Error("expected error")
		}
		var e ErrNoServices
		if !errors.As(err, &e) {
			t.Errorf("expected %T, got %T", e, err)
		}
	})
}
