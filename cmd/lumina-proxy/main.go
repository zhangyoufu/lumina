package main

import (
    "io/ioutil"
    "log"
    "net"
    "os"
    "github.com/palantir/stacktrace"
    "github.com/zhangyoufu/lumina"
)

func main() {
    f, err := os.Open("ida.key")
    if err != nil {
        log.Fatal(stacktrace.Propagate(err, "unable to open ida.key under working directory"))
    }
    data, err := ioutil.ReadAll(f)
    if err != nil {
        log.Fatal(stacktrace.Propagate(err, "unable to read license"))
    }
    licKey := lumina.LicenseKey(data)
    idaInfo := licKey.GetIDAInfo()
    licOS = idaInfo.OS
    licID := idaInfo.Id
    proxy := NewProxy(licKey, licID)

    // TODO: TLS listener
    ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 65432})
    if err != nil {
        log.Fatal(stacktrace.Propagate(err, "unable to listen"))
    }
    log.Print("proxy is listening on ", ln.Addr())
    proxy.Serve(ln)
    // TODO: check return value of Serve and support graceful shutdown
}
