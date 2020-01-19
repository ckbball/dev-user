package v1

import (
  "context"
  //"errors"
  //"strconv"
  //"time"
  "bcrypt"

  //"github.com/golang/protobuf/ptypes"
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
  repo          repository
  loggerAddress string
  tokenService  Authable
  //subscriber message.Subscriber
  //publisher  message.Publisher
}

func NewUserServiceServer(repo repository, loggerAddress string, tokenService Authable) *handler {
  //subscriber message.Subscriber, publisher message.Publisher) *handler {
  return &handler{
    repo:          repo,
    loggerAddress: loggerAddress,
    tokenService:  tokenService,
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

  // generate hash of password
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("error hashing password: %v", err))
  }
  req.User.Password = string(hashedPass)

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

  // check user is updating their own profile.
  // grab http headers from metadata
  md, ok := metadata.FromIncomingContext(ctx)
  // grab user token from metadata
  reqToken := md["Authorization"]
  // validate the token user and request user
  claims, err := s.tokenService.Decode(reqToken)
  if err != nil {
    return nil, err
  }

  // if token User != req User or there is no user id in claims return error
  if claims.User.Id != req.Id || claims.User.Id == "" {
    return nil, errors.New("Invalid Token")
  }

  // add in password hashing later
  if req.User.Password != "" {
    hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("error hashing password: %v", err))
    }
    req.User.Password = string(hashedPass)
  }

  match, modified, err := s.repo.Update(req.User, req.Id)
  if err != nil {
    return nil, err
  }

  // return
  return &v1.UpsertResponse{
    Api:      apiVersion,
    Status:   "test",
    Matched:  match,
    Modified: modified,
    // maybe in future add more data to response about the added user.
  }, nil
}

func (s *handler) DeleteUser(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // check user is updating their own profile.
  // grab http headers from metadata
  md, ok := metadata.FromIncomingContext(ctx)
  // grab user token from metadata
  reqToken := md["Authorization"]
  // validate the token user and request user
  claims, err := s.tokenService.Decode(reqToken)
  if err != nil {
    return nil, err
  }

  // if token User != req User or there is no user id in claims return error
  if claims.User.Id != req.Id || claims.User.Id == "" {
    return nil, errors.New("Invalid Token")
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

  protoUsers := exportUserModels(users)

  return &v1.FindResponse{
    Api:    req.Api,
    Status: "Does it load local code",
    Users:  protoUsers,
  }, nil
}

func (s *handler) GetById(ctx context.Context, req *v1.FindRequest) (*v1.FindResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  user, err := s.repo.GetById(req.Id)
  if err != nil {
    return nil, err
  }

  exportedUser := &v1.User{
    LastActive: int32(user.LastActive),
    Username:   user.Username,
    Experience: user.Experience,
    Languages:  user.Languages,
  }

  return &v1.FindResponse{
    Api:    req.Api,
    Status: "Found User",
    User:   exportedUser,
  }, nil
}

func exportUserModels(users []*User) []*v1.User {
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

type logMessage struct {
  Message []byte `json:"message"`
  Process string `json:"process"`
}
