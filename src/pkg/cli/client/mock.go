package client

import (
	"context"

	defangv1 "github.com/DefangLabs/defang/src/protos/io/defang/v1"
	compose "github.com/compose-spec/compose-go/v2/types"
)

type MockClient struct {
	Client
	UploadUrl string
	Project   *compose.Project
}

var _ Client = (*MockClient)(nil)

func (m MockClient) CreateUploadURL(ctx context.Context, req *defangv1.UploadURLRequest) (*defangv1.UploadURLResponse, error) {
	return &defangv1.UploadURLResponse{Url: m.UploadUrl + req.Digest}, nil
}

func (m MockClient) ServiceDNS(service string) string {
	return service
}

func (m MockClient) LoadProject() (*compose.Project, error) {
	return m.Project, nil
}
