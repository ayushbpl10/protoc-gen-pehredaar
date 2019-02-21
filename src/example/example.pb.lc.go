package example

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type localUsersClient struct {
	UsersServer
}

func (lc *localUsersClient) AddUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*empty.Empty, error) {
	return lc.UsersServer.AddUser(ctx, in)
}
func (lc *localUsersClient) GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*User, error) {
	return lc.UsersServer.GetUser(ctx, in)
}
func (lc *localUsersClient) UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...grpc.CallOption) (*empty.Empty, error) {
	return lc.UsersServer.UpdateUser(ctx, in)
}

func NewLocalUsersClient(s UsersServer) UsersClient {
	return &localUsersClient{s}
}
