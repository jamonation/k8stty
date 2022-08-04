package service_test

import (
	"context"
	"testing"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/service"

	"k8stty/internal/pkg/objectmanager"
)

type mockServiceManager interface {
	Create(context.Context, map[string]string) error
	Delete(context.Context, string) error
}

type mockServiceManagerImpl objectmanager.ManagerImpl

func newMockServiceManager() mockServiceManager {
	return &mockServiceManagerImpl{}
}

func TestCreateService(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.CreateServiceReq
		resp      *pb.CreateServiceResp
		expectErr bool
		errMsg    string
	}{
		{
			name:      "request ok",
			req:       &pb.CreateServiceReq{ServiceId: "test"},
			resp:      &pb.CreateServiceResp{Success: true},
			expectErr: false, // function returns nil resp on error
		},
		{
			name:      "empty request returns error",
			req:       &pb.CreateServiceReq{},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			errMsg:    "missing request id",
		},
	}

	mockServiceManager := newMockServiceManager()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			serviceService := service.NewServiceServer(mockServiceManager)
			resp, err := serviceService.CreateService(ctx, tc.req)
			if tc.expectErr {
				if resp.GetSuccess() != tc.resp.GetSuccess() {
					t.Errorf("incorrect status message:\ngot: '%#v'\nwant: '%#v'", resp.GetSuccess(), tc.resp.GetSuccess())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v\n", err)
				}
			}
		})
	}

}

func TestDeleteService(t *testing.T) {}
