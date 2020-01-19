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

  // Call CreateUser------------------------------------------------------------------------------------------------------------------
  resp, err := http.Post(*address+"/v1/users", "application/json", strings.NewReader(fmt.Sprintf(`
    {
      "api":"v1",
      "user": {
        "email": "loola@gmail.com",
        "password": "haha",
        "username": "loolala",
        "last_active": 0,
        "experience": "beginner",
        "languages": ["haskell", "python", "csharp"]
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

  // ------------------------------------------------------------------------------------------------------------------LOGIN
  req, err = http.NewRequest("POST", fmt.Sprintf("%s%s", *address, "/v1/login"), strings.NewReader(fmt.Sprintf(`
       {
         "api":"v1",
          "email": "loola@gmail.com",
          "password": "haha",
       }
     `, pfx, pfx, pfx)))
  if err != nil {
    log.Fatalf("failed to call Login method: %v", err)
  }
  /*
    req.Header["Authorization"][0] = loggedToken
  */
  resp, err = http.DefaultClient.Do(req)
  if err != nil {
    log.Fatalf("failed to call Login method: %v", err)
  }
  bodyBytes, err = ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read Login response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("Login response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  var logged struct {
    API    string `json:"api"`
    Status string `json:"status"`
    Token  string `json:"token"`
  }
  err = json.Unmarshal(bodyBytes, &updated)
  if err != nil {
    log.Fatalf("failed to unmarshal JSON response of UpdateUser method: %v", err)
    fmt.Println("error:", err)
  }
  log.Printf("updated struct: %s\n", updated)

  loggedToken := updated.Token

  // Call UpdateUser with correct info------------------------------------------------------------------------------------------------------------------
  resp, err = http.Post(*address+"/v1/users/"+created.Id, "application/json", strings.NewReader(fmt.Sprintf(`
       {
         "api":"v1",
         "user": {
           "email": "loola@gmail.com",
           "password": "haha",
           "username": "brewhaha",
           "last_active": 100,
           "experience": "senior",
           "languages": ["haskell", "python", "csharp"]
         }
       }
     `, pfx, pfx, pfx)))
  if err != nil {
    log.Fatalf("failed to call UpdateUser method: %v", err)
  }
  /*
    req.Header["Authorization"][0] = loggedToken
  */
  bodyBytes, err = ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read UpdateUser response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("UpdateUser response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  // parse status of UpdateUser
  var updated struct {
    API      string `json:"api"`
    Status   string `json:"status"`
    Matched  string `json:"matched"`
    Modified string `json:"modified"`
  }
  err = json.Unmarshal(bodyBytes, &updated)
  if err != nil {
    log.Fatalf("failed to unmarshal JSON response of UpdateUser method: %v", err)
    fmt.Println("error:", err)
  }
  log.Printf("updated struct: %s\n", updated)

  // Call FilterUsers------------------------------------------------------------------------------------------------------------------
  resp, err = http.Post(*address+"/v1/search", "application/json", strings.NewReader(fmt.Sprintf(`
          {
            "api":"v1",
            "experience": "beginner",
            "page": 1,
            "limit": 20
          }
        `, pfx, pfx, pfx)))
  if err != nil {
    log.Fatalf("failed to call FilterUsers method: %v", err)
  }
  bodyBytes, err = ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read FilterUsers response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("FilterUsers searching for users who know java\n")
  log.Printf("FilterUsers response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  var users struct {
    LastActive string `json:"last_active"`
    Experience string `json:"experience"`
    Languages  string `json:"languages"`
    Username   string `json:"username"`
  }
  bobby := users
  err = json.Unmarshal(bodyBytes, &bobby)
  if err != nil {
    log.Fatalf("failed to unmarshal JSON response of FilterUsers method: %v", err)
    fmt.Println("error:", err)
  }

  // Call GetById------------------------------------------------------------------------------------------------------------------
  req, err = http.NewRequest("GET", fmt.Sprintf("%s%s/%s", *address, "/v1/users", created.Id), nil)
  resp, err = http.DefaultClient.Do(req)
  if err != nil {
    log.Fatalf("failed to call GetById method: %v", err)
  }
  bodyBytes, err = ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read GetById response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("GetById response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  // Call DeleteUser------------------------------------------------------------------------------------------------------------------
  req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/users", created.Id), nil)
  resp, err = http.DefaultClient.Do(req)
  if err != nil {
    log.Fatalf("failed to call Delete method: %v", err)
  }
  bodyBytes, err = ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read Delete response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
