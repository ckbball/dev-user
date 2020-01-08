package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
  "time"
)

func main() {
  // get configuration
  address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
  flag.Parse()

  t := time.Now().In(time.UTC)
  pfx := t.Format(time.RFC3339Nano)

  var body string

  // Call CreateUser
  resp, err := http.Post(*address+"/v1/users", "application/json", strings.NewReader(fmt.Sprintf(`
    {
      "api":"v1",
      "user": {
        "email": "bobby@gmail.com",
        "password": "haha",
        "username": "btotheg",
        "last_active": 0,
        "experience": "beginner",
        "languages": ["golang", "ruby", "javascript"]
      }
    }
  `, pfx, pfx, pfx)))
  if err != nil {
    log.Fatalf("failed to call CreateUser method: %v", err)
  }
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read CreateUser response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("CreateUser response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  // parse status of CreateUser
  var created struct {
    API    string `json:"api"`
    Status string `json:"status"`
    Id     string `json:"id"`
  }
  err = json.Unmarshal(bodyBytes, &created)
  if err != nil {
    log.Fatalf("failed to unmarshal JSON response of CreateUser method: %v", err)
    fmt.Println("error:", err)
  }
  log.Printf("created struct: %s\n", created)
}
