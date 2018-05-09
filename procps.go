package main

import (
    "io/ioutil"
    "os"
    "log"
    "fmt"
    "bufio"
    "strings"

    "github.com/docker/docker/client"
    "golang.org/x/net/context"

)

var docker *client.Client

func findvol(vol string) map[string]int {
    cids := make(map[string]int)

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
                    if (f[0] == "/dev/pxd/pxd" + vol && ! strings.HasPrefix(f[1], "/var/lib/osd/mounts")) || strings.Contains(f[1], "kubernetes.io~portworx-volume/" + vol) {
						fmt.Println(">>> ", pid.Name())
                        cids[getcdockercid(pid.Name())] = 1
                    }
                }
            }
        }
    }

    return cids
}

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
                if d[1] == "docker" && len(d[2]) > 0 {
                    return d[2]
                } else if d[1] == "kubepods" && len(d[4]) > 0 {
                    return d[4]
                } else {
                    return "host"
                }

            }
        }
    }


    return "host"
}

func dockerinspect(cid string) {

    i, err := docker.ContainerInspect(context.Background(), cid)
    if err != nil {
        panic(err)
    }

    for _, m := range i.Mounts {
        if m.Driver == "pxd" || strings.Contains(m.Source, "kubernetes.io~portworx-volume")  {
            fmt.Println("Name:\t", i.Name, "\nImg:\t", i.Config.Image, "\nArgs:\t", i.Args, "\nCmd:\t", "\nPath:\t", i.Path)
            fmt.Println("Mount:\t", m.Name, ":", m.Driver, ":", m.Source, ":", m.Destination)
        }
    }
    
}

func main() {
    if len(os.Args) < 2 {
        panic("usage: pxvol volid")
    }

    var err error

    docker, err = client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    cids := findvol(os.Args[1])

    for key, _ := range cids {
        if key == "host" {
            fmt.Println("host mounted")
        } else {
            fmt.Println("ID: \t", key)
                dockerinspect(key)
        }
    }
}
