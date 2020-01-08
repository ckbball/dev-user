package main

import (
  "context"
  "flag"
  "log"
  "time"

  //"github.com/golang/protobuf/ptypes"
  "google.golang.org/grpc"

  v1 "github.com/ckbball/dev-user/pkg/api/v1"
)

const (
  // apiVersion is version of API is provided by server
  apiVersion = "v1"
)

func main() {
  // get configuration
  address := flag.String("server", "", "gRPC server in format host:port")
  flag.Parse()

  // Set up a connection to the server.
  conn, err := grpc.Dial(*address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()

  c := v1.NewUserServiceClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  // Call Create
  req1 := v1.UpsertRequest{
    Api: apiVersion,
    User: &v1.User{
      Email:      "carl@yahoo.com",
      Password:   "blah",
      Username:   "carlyahoo",
      LastActive: 0,
      Experience: "middle",
      Languages:  []string{"java", "c++", "f#"},
    },
  }
  res1, err := c.CreateUser(ctx, &req1)
  if err != nil {
    log.Fatalf("CreateUser failed: %v", err)
  }
  log.Printf("CreateUser result: <%+v>\n\n", res1)

  // id := res1.Id

  /*
     // Read
     req2 := v1.ListRequest{
       Api:   apiVersion,
       Page:  1,
       Limit: 20,
     }
     res2, err := c.ListItems(ctx, &req2)
     if err != nil {
       log.Fatalf("Read failed: %v", err)
     }
     log.Printf("Read result: <%+v>\n\n", res2)

     // Update
     /*
        req3 := v1.UpdateRequest{
          Api: apiVersion,
          ToDo: &v1.ToDo{
            Id:          res2.ToDo.Id,
            Title:       res2.ToDo.Title,
            Description: res2.ToDo.Description + " + updated",
            Reminder:    res2.ToDo.Reminder,
          },
        }
        res3, err := c.Update(ctx, &req3)
        if err != nil {
          log.Fatalf("Update failed: %v", err)
        }
        log.Printf("Update result: <%+v>\n\n", res3)
  */

  /*
     // Call FindItems
     req4 := v1.Specification{
       Api:  apiVersion,
       Solo: 2,
     }
     res4, err := c.FindItems(ctx, &req4)
     if err != nil {
       log.Fatalf("FindItems failed: %v", err)
     }
     log.Printf("FindItems result: <%+v>\n\n", res4)

     // Delete
     req5 := v1.RemoveRequest{
       Api: apiVersion,
       Id:  id,
     }
     res5, err := c.RemoveItem(ctx, &req5)
     if err != nil {
       log.Fatalf("Delete failed: %v", err)
     }
     log.Printf("Delete result: <%+v>\n\n", res5)

  */
}
