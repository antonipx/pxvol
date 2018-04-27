package main

import (
    "io/ioutil"
    "os"
    "log"
    "fmt"
    "bufio"
    "strings"

)

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
                    if f[0] == "/dev/pxd/pxd" + vol && ! strings.HasPrefix(f[1], "/var/lib/osd/mounts") {
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
    if len(os.Args) < 2 {
        panic("usage: pxvol volid")
    }

    cids := findvol(os.Args[1])

    for key, _ := range cids {
        fmt.Println(key)
    }
}