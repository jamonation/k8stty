package networkpolicy_test

import (
	"context"
	"testing"

	"k8stty/internal/networkpolicy"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type mockNetworkpolicyManager interface {
	Create(context.Context, map[string]string) error
	Delete(context.Context, string) error
}

type mockNetworkpolicyManagerImpl objectmanager.ManagerImpl

func NewNetworkpolicyServer() mockNetworkpolicyManager {
	return &mockNetworkpolicyManagerImpl{}
}

func (k *mockNetworkpolicyManagerImpl) Create(ctx context.Context, reqInfo map[string]string) error {
	return nil
}

func (k *mockNetworkpolicyManagerImpl) Delete(ctx context.Context, reqInfo string) error {
	return nil
}

func TestCreateNetworkpolicy(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.CreateNetworkpolicyReq
		resp      *pb.CreateNetworkpolicyResp
		expectErr bool
		err       string
	}{
		{
			name:      "request ok",
			req:       &pb.CreateNetworkpolicyReq{NetworkpolicyId: "test"},
			resp:      &pb.CreateNetworkpolicyResp{Success: true},
			expectErr: false,
		},
		{
			name:      "empty request returns error",
			req:       &pb.CreateNetworkpolicyReq{},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "missing request id",
		},
	}

	mockNetworkpolicyManager := NewNetworkpolicyServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			NetworkpolicyService := networkpolicy.NewNetworkpolicyServer(mockNetworkpolicyManager)
			resp, err := NetworkpolicyService.CreateNetworkpolicy(ctx, tc.req)
			if tc.expectErr {
				if resp != tc.resp {
					t.Error(err)
				}
				if err.Error() != tc.err {
					t.Errorf("incorrect error message:\ngot: '%v'\nwant: '%v'", err.Error(), tc.err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func TestDeleteNetworkpolicy(t *testing.T) {}
