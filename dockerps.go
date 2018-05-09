
// go get github.com/docker/docker/client
// go get golang.org/x/net/context
// https://godoc.org/github.com/docker/docker/api/types#ContainerJSON
// https://godoc.org/github.com/docker/docker/client
// https://docs.docker.com/develop/sdk/examples/

package main

import (
    "os"
    "fmt"
    "strings"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"

)

var docker *client.Client

func findcontainer(cid string, volname string) {

    i, err := docker.ContainerInspect(context.Background(), cid)
    if err != nil {
        panic(err)
    }


    for _, m := range i.Mounts {
        if (m.Driver == "pxd" && m.Name == volname) || strings.Contains(m.Source, "kubernetes.io~portworx-volume/" + volname)  {
            fmt.Println("ID:\t", i.ID, "\nName:\t", i.Name, "\nImg:\t", i.Config.Image, "\nArgs:\t", i.Args, "\nCmd:\t", "\nPath:\t", i.Path)
            fmt.Println("Mount:\t", m.Name, ":", m.Driver, ":", m.Source, ":", m.Destination)
        }
    }

    
}

func main() {
    if len(os.Args) < 2 {
        panic("Usage: dockerps <volname>")
    }

    var err error

    docker, err = client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    cids, err := docker.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    for _, c := range cids {
        findcontainer(c.ID, os.Args[1])
    }

}
