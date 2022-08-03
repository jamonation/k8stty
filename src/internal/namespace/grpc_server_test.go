package namespace_test

import (
	"context"
	"testing"

	"k8stty/internal/namespace"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type mockNamespaceManger interface {
	Create(context.Context, map[string]string) error
	Delete(context.Context, string) error
}

type mockNamespaceManagerImpl objectmanager.ManagerImpl

func newMockNamespaceManager() mockNamespaceManger {
	return &mockNamespaceManagerImpl{}
}

func (k *mockNamespaceManagerImpl) Create(ctx context.Context, reqInfo map[string]string) error {
	return nil
}

func (k *mockNamespaceManagerImpl) Delete(ctx context.Context, reqInfo string) error {
	return nil
}

func TestCreateNamespace(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.CreateNamespaceReq
		resp      *pb.CreateNamespaceResp
		expectErr bool
		err       string
	}{
		{
			name:      "request ok",
			req:       &pb.CreateNamespaceReq{NamespaceId: "test"},
			resp:      &pb.CreateNamespaceResp{Success: true},
			expectErr: false,
		},
		{
			name:      "empty request returns error",
			req:       &pb.CreateNamespaceReq{},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "missing create namespace request id",
		},
	}

	mockNamespaceManager := newMockNamespaceManager()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			namespaceService := namespace.NewNamespaceServer(mockNamespaceManager)
			resp, err := namespaceService.CreateNamespace(ctx, tc.req)
			if tc.expectErr {
				if resp != tc.resp {
					t.Error(err)
				}
				if err.Error() != tc.err {
					t.Errorf("incorrect error message:\ngot: '%#v'\nwant: '%#v'", err, tc.err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func TestDeleteNamespace(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.DeleteNamespaceReq
		resp      *pb.DeleteNamespaceResp
		expectErr bool
		err       string
	}{
		{
			name:      "request ok",
			req:       &pb.DeleteNamespaceReq{NamespaceId: "test"},
			resp:      &pb.DeleteNamespaceResp{Success: true},
			expectErr: false,
		},
		{
			name:      "empty request returns error",
			req:       &pb.DeleteNamespaceReq{},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "missing delete namespace request id",
		},
	}

	mockNamespaceManager := newMockNamespaceManager()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			namespaceService := namespace.NewNamespaceServer(mockNamespaceManager)
			resp, err := namespaceService.DeleteNamespace(ctx, tc.req)
			if tc.expectErr {
				if resp != tc.resp {
					t.Error(err)
				}
				if err.Error() != tc.err {
					t.Errorf("incorrect error message:\ngot: '%#v'\nwant: '%#v'", err, tc.err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}
