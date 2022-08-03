package pod_test

import (
	"context"
	"testing"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
	"k8stty/internal/pod"
)

type mockPodManger interface {
	Create(context.Context, map[string]string) error
	Delete(context.Context, string) error
}

type mockPodManagerImpl objectmanager.ManagerImpl

func newMockPodManager() mockPodManger {
	return &mockPodManagerImpl{}
}

func (k *mockPodManagerImpl) Create(ctx context.Context, reqInfo map[string]string) error {
	return nil
}

func (k *mockPodManagerImpl) Delete(ctx context.Context, reqInfo string) error {
	return nil
}

func TestCreatePod(t *testing.T) {
	allowedImages := make(map[string]struct{})
	allowedImages["alpine:3"] = struct{}{}
	registryUrl := "index.docker.io"
	testCases := []struct {
		name      string
		req       *pb.CreatePodReq
		resp      *pb.CreatePodResp
		expectErr bool
		err       string
	}{
		{
			name:      "request ok",
			req:       &pb.CreatePodReq{PodId: "test", ImageName: "alpine:3"},
			resp:      &pb.CreatePodResp{Success: true},
			expectErr: false,
		},
		{
			name:      "missing image",
			req:       &pb.CreatePodReq{PodId: "test"},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "missing create pod image",
		},
		{
			name:      "invalid image",
			req:       &pb.CreatePodReq{PodId: "test", ImageName: "invalid:image"},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "invalid image",
		},
		{
			name:      "empty request returns error",
			req:       &pb.CreatePodReq{},
			resp:      nil, // function returns nil resp on error
			expectErr: true,
			err:       "missing create pod request id",
		},
		//TODO: need a test for an error from CreatePod
	}

	mockPodManager := newMockPodManager()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			podService := pod.NewPodServer(mockPodManager, allowedImages, registryUrl)
			resp, err := podService.CreatePod(ctx, tc.req)
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

func TestDeletePod(t *testing.T) {}
