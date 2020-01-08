package v1

import (
  "context"
  //"errors"
  //"fmt"
  //"log"
  //"strconv"
  //"time"

  //"github.com/golang/protobuf/ptypes"
  //"encoding/json"
  //"github.com/ThreeDotsLabs/watermill"
  //"github.com/ThreeDotsLabs/watermill/message"
  // "github.com/go-redis/cache/v7"
  //"google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  // messageProto "github.com/ckbball/dev-message/pkg/api/v1"

  v1 "github.com/ckbball/dev-user/pkg/api/v1"
)

const (
  apiVersion = "v1"
  eventName  = "account_purchased"
)

// this handler must satisfy the UserServiceServer interface in user.pb.go
/*
CreateUser(context.Context, *UpsertRequest) (*UpsertResponse, error) X
UpdateUser(context.Context, *UpsertRequest) (*UpsertResponse, error) X
DeleteUser(context.Context, *DeleteRequest) (*DeleteResponse, error) X
GetById(context.Context, *FindRequest) (*FindResponse, error)
FilterUsers(context.Context, *FindRequest) (*FindResponse, error)
*/
type handler struct {
  repo repository
  //subscriber message.Subscriber
  //publisher  message.Publisher
}

func NewUserServiceServer(repo repository) *handler {
  //subscriber message.Subscriber, publisher message.Publisher) *handler {
  return &handler{
    repo: repo,
    //subscriber: subscriber,
    //publisher:  publisher,
  }
}

func (s *handler) checkAPI(api string) error {
  if len(api) > 0 {
    if apiVersion != api {
      return status.Errorf(codes.Unimplemented,
        "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
    }
  }
  return nil
}

func (s *handler) CreateUser(ctx context.Context, req *v1.UpsertRequest) (*v1.UpsertResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // add in hashing later
  /*
     user := &User{
       Email:      req.User.Email,
       Password:   req.User.Password,
       Username:   req.User.Username,
       LastActive: int(req.User.LastActive),
       Experience: req.User.Experience,
       Languages:  req.User.Languages,
     }*/

  id, err := s.repo.Create(req.User)
  if err != nil {
    return nil, err
  }

  // return
  return &v1.UpsertResponse{
    Api:    apiVersion,
    Status: "Created",
    Id:     id,
    // maybe in future add more data to response about the added user.
  }, nil
}

func (s *handler) UpdateUser(ctx context.Context, req *v1.UpsertRequest) (*v1.UpsertResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // add in hashing later
  /*
     user := &User{
       Email:      req.User.Email,
       Password:   req.User.Password,
       Username:   req.User.Username,
       LastActive: int(req.User.LastActive),
       Experience: req.User.Experience,
       Languages:  req.User.Languages,
     }*/

  /*match, modified, err := s.repo.Update(req.User, req.Id)
    if err != nil {
      return nil, err
    }*/

  // return
  return &v1.UpsertResponse{
    Api:    apiVersion,
    Status: "test",
    // Matched:  match,
    // Modified: modified,
    // maybe in future add more data to response about the added user.
  }, nil
}

func (s *handler) DeleteUser(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  count, err := s.repo.Delete(req.Id)
  if err != nil {
    return nil, err
  }

  return &v1.DeleteResponse{
    Api:    req.Api,
    Status: "Deleted",
    Count:  count,
  }, nil
}

func (s *handler) FilterUsers(ctx context.Context, req *v1.FindRequest) (*v1.FindResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  users, err := s.repo.FilterUsers(req)
  if err != nil {
    return nil, err
  }

  protoUsers := exportUserModel(users)

  return &v1.FindResponse{
    Api:    req.Api,
    Status: "Found Users",
    Users:  protoUsers,
  }, nil
}

func (s *handler) GetById(ctx context.Context, req *v1.FindRequest) (*v1.FindResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  return &v1.FindResponse{
    Api:    req.Api,
    Status: "Placeholder for test",
  }, nil
}

func exportUserModel(users []*User) []*v1.User {
  out := []*v1.User{}
  for _, element := range users {
    user := &v1.User{
      LastActive: int32(element.LastActive),
      Username:   element.Username,
      Experience: element.Experience,
      Languages:  element.Languages,
    }
    out = append(out, user)
  }
  return out
}
