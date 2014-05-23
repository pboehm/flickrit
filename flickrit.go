package main

import (
    "encoding/json"
    "github.com/pboehm/flickrit/api"
    "fmt"
)

func main() {
    api := &api.API{
        ApiKey: "122cc483be92bd806b696e7d458596ac",
    }
    api.Setup()

    photos, _ := api.GetPhotosForUser("phboehm")

    b, _ := json.Marshal(photos)
    fmt.Println(string(b))
}
