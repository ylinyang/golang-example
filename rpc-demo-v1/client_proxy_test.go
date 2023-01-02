package rpc_demo_v1

import (
	"context"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name    string
		service *UserService
		wantErr error
	}{
		{
			name:    "user service",
			service: &UserService{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service)
			assert.Equal(t, tc.wantErr, err)
			_, _ = tc.service.GetById(context.Background(), &GetByIdReq{})
		})
	}
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
}
type GetByIdResp struct {
}
