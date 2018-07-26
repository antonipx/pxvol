package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

func findvol(vol string) map[string]string {
    cids := make(map[string]string)

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
                    if (f[0] == "/dev/pxd/pxd"+vol || f[0] == "pxfs"+vol) && !strings.HasPrefix(f[1], "/var/lib/osd/mounts") && !strings.Contains(f[1], "kubernetes.io~portworx-volume/") && !strings.HasPrefix(f[1], "/pxmounts/") {
                        //fmt.Println(">>> ", pid.Name(), f[0], f[1])
                        id := getcdockercid(pid.Name())
                        //cids[id] = f[1]
                        if !strings.Contains(cids[id], f[1]) { //TODO: rewrite as slice contains instead of string
                            cids[id] = cids[id] + "\n  " + f[1]
                        }
                    }
                }
            }
        }
    }

    return cids
}

func getcdockercid(pid string) string {
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

    return "unknown"
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: pxvol volume_id")
        os.Exit(1)
    }

    cids := findvol(os.Args[1])

    for key, val := range cids {
        if key == "host" {
            fmt.Println("host mounted:", val)
        } else {
            fmt.Println("Docker ID:", key, val)
        }
    }
}
