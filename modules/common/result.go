package common

import (
    "os"
    "log"
    "fmt"
    "github.com/golang-collections/collections/set"
)

type Page struct {
    Link string
    Source string
    Status int
}

func SaveResult(links *set.Set) {
    // TODO: move output filename to config
    file, err := os.Create("result.txt")
    if err != nil {
        log.Fatal("Cannot create file", err)
    }
    defer file.Close()

    links.Do(func(item interface{}) {
        fmt.Fprintln(file, item.(string))
    })
}