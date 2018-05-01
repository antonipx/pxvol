
// go get github.com/docker/docker/client
// go get golang.org/x/net/context
// https://godoc.org/github.com/docker/docker/api/types#ContainerJSON
// https://godoc.org/github.com/docker/docker/client
// https://docs.docker.com/develop/sdk/examples/

package main

import (
//    "io/ioutil"
    "os"
//    "log"
    "fmt"
//    "bufio"
//    "strings"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"

)

var docker *client.Client

func dockerinspect(cid string) {

    i, err := docker.ContainerInspect(context.Background(), cid)
    if err != nil {
        panic(err)
    }

    fmt.Println("Name:\t", i.Name, "\nImg:\t", i.Config.Image, "\nArgs:\t", i.Args, "\nCmd:\t", "\nPath:\t", i.Path)

    for _, m := range i.Mounts {
        fmt.Println("Mount:\t", m.Driver, ":", m.Source, ":", m.Destination)
    }

    
}

func main() {
    var err error

    docker, err = client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    if len(os.Args) < 2 {
        cids, err := docker.ContainerList(context.Background(), types.ContainerListOptions{})
        if err != nil {
            panic(err)
        }
    
        for _, c := range cids {
            dockerinspect(c.ID)
        }
    } else {
        dockerinspect(os.Args[1])
    }
    

}
