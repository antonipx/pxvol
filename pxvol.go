// go get github.com/docker/docker/client
// go get golang.org/x/net/context
// https://godoc.org/github.com/docker/docker/api/types#ContainerJSON
// https://docs.docker.com/develop/sdk/examples/

package main

import (
    "io/ioutil"
    "os"
    "log"
    "fmt"
    "bufio"
    "strings"
    //"encoding/json"

    //"github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"
)

func getcdockercid(pid string) string  {
    cgroup, err := os.Open("/proc/" + pid + "/cgroup")
    if err == nil {
        defer cgroup.Close()
        scanner := bufio.NewScanner(cgroup)
        scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
            f := strings.Split(scanner.Text(), ":")
            // this can be any other cgroup
            if f[1] == "pids" {
                d := strings.Split(f[2], "/")
                // todo k8s
                if d[1] == "docker" && len(d[2]) > 0 {
                    //fmt.Println(">>> DOCKER:", pid, f[1], d[0], d[1])
                    return d[2]
                } else {
                    return "host"
                }

            }
        }
    }



    return "host"
}

func main() {
    var vol string = "888651093940211871"
//    var vol string = "848783366117993172"


    docker, err := client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    pids, err := ioutil.ReadDir("/proc")
    if err != nil {
        log.Fatal(err)
    }

    for _, pid := range pids {
        if pid.IsDir() && pid.Name()[0] >= '1' && pid.Name()[0] <= '9' {

            //fmt.Println("/proc/" + pid.Name() + "/mounts")
            mounts, err := os.Open("/proc/" + pid.Name() + "/mounts")
            if err == nil {
                defer mounts.Close()
                scanner := bufio.NewScanner(mounts)
                scanner.Split(bufio.ScanLines)
                for scanner.Scan() {
                    f := strings.Fields(scanner.Text())
                    //fmt.Println("  ", f[0], f[1])
                    if f[0] == "/dev/pxd/pxd" + vol && ! strings.HasPrefix(f[1], "/var/lib/osd/mounts") {
                        cid := getcdockercid(pid.Name())
                        if cid != "host" {
                            fmt.Println("mnt:", pid.Name(), f[1], cid)
                            inspect, err := docker.ContainerInspect(context.Background(), cid)
                            if err == nil {

                                fmt.Println("Inspect:", inspect.Mounts[0].Driver)

                            }
                        }
                    }
                }
            }
        }
    }
}
