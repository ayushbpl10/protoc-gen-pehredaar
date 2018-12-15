// Code generated by protoc-gen-defaults. DO NOT EDIT.

package rightsval

import (
	"context"

	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user "../../users"

	"../rights"
	"github.com/ayushbpl10/protoc-gen-rights/example/pb/example"

	right "github.com/ayushbpl10/protoc-gen-rights/example/rights"
)

const UsersResourcePaths = [...]string{

	"/users/{id}/cards.read/{blocked}",

	"/users/{id}/cards/user.write",

	"/{user_email.email}/users/{user_id}/cards/{tent_id.tent}/email/{user_email.email.checks.check.check_id.val_id}",

	"/users/{user_id}/cards/{tent_id.tent}/ex.write",
}

type RightsUsersServer struct {
	example.UsersServer
	rightsCli rights.RightValidatorsClient
	user      user.UserIDer
}

func init() {
	options = append(options, fx.Provide(NewRightsUsersClient))
}

type RightsUsersClientResult struct {
	fx.Out
	UsersClient example.UsersClient `name:"r"`
}

func NewRightsUsersClient(c rights.RightValidatorsClient, s example.UsersServer) RightsUsersClientResult {
	return RightsUsersClientResult{UsersClient: example.NewLocalUsersClient(NewRightsUsersServer(c, s))}
}
func NewRightsUsersServer(c rights.RightValidatorsClient, s example.UsersServer, u right.UserIDer) example.UsersServer {
	return &RightsUsersServer{
		s,
		c,
		u,
	}
}

func (s *RightsUsersServer) AddUser(ctx context.Context, rightsvar *example.User) (*empty.Empty, error) {

	ResourcePathOR := make([]string, 0)
	ResourcePathAND := make([]string, 0)

	for _, Blocked := range rightsvar.GetBlocked() {

		ResourcePathAND = append(ResourcePathAND,

			fmt.Sprintf("/users/%s/cards.read/%s",

				rightsvar.GetId(),

				Blocked,
			),
		)

	}

	ResourcePathOR = append(ResourcePath,

		fmt.Sprintf("/users/%s/cards/user.write",

			rightsvar.GetId(),
		),
	)

	res, err := s.rightsCli.IsValid(ctx, &rights.IsValidReq{
		ResourcePathOr:  ResourcePathOR,
		ResourcePathAnd: ResourcePathAND,
		UserId:          s.user.UserID(ctx),
		ModuleName:      "Users",
	})
	if err != nil {
		return nil, err
	}

	if !res.IsValid {
		return nil, status.Errorf(codes.PermissionDenied, res.Reason)
	}

	return s.UsersServer.AddUser(ctx, rightsvar)
}

func (s *RightsUsersServer) GetUser(ctx context.Context, rightsvar *example.GetUserReq) (*example.User, error) {

	ResourcePathOR := make([]string, 0)
	ResourcePathAND := make([]string, 0)

	for _, UserEmail := range rightsvar.GetUserEmail() {

		for _, Checks := range UserEmail.GetChecks() {

			for _, CheckId := range Checks.GetCheckId() {

				ResourcePathAND = append(ResourcePathAND,

					fmt.Sprintf("/%s/users/%s/cards/%s/email/%s",

						UserEmail.GetEmail(),

						rightsvar.GetUserId(),

						rightsvar.GetTentId().GetTent(),

						Checks.GetCheck(),

						CheckId.GetValId(),
					),
				)

			}

		}

	}

	ResourcePathOR = append(ResourcePath,

		fmt.Sprintf("/users/%s/cards/%s/ex.write",

			rightsvar.GetUserId(),

			rightsvar.GetTentId().GetTent(),
		),
	)

	res, err := s.rightsCli.IsValid(ctx, &rights.IsValidReq{
		ResourcePathOr:  ResourcePathOR,
		ResourcePathAnd: ResourcePathAND,
		UserId:          s.user.UserID(ctx),
		ModuleName:      "Users",
	})
	if err != nil {
		return nil, err
	}

	if !res.IsValid {
		return nil, status.Errorf(codes.PermissionDenied, res.Reason)
	}

	return s.UsersServer.GetUser(ctx, rightsvar)
}

func (s *RightsUsersServer) UpdateUser(ctx context.Context, rightsvar *example.UpdateUserReq) (*empty.Empty, error) {

	return s.UsersServer.UpdateUser(ctx, rightsvar)
}