package main

import (
    "crypto/tls"
    "io/ioutil"
    "log"
    "os"
    "github.com/zhangyoufu/lumina"
)

func main() {
    f, err := os.Open("ida.key")
    if err != nil {
        log.Fatal("unable to open ida.key: ", err)
    }
    data, err := ioutil.ReadAll(f)
    if err != nil {
        log.Fatal("unable to read license: ", err)
    }
    licKey := lumina.LicenseKey(data)
    idaInfo := licKey.GetIDAInfo()
    licOS = idaInfo.OS
    licID := idaInfo.Id
    proxy := NewProxy(licKey, licID)

    cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        log.Fatal("unable to load X509 key pair: ", err)
    }
    config := &tls.Config{
        Certificates: []tls.Certificate{ cert },
    }
    ln, err := tls.Listen("tcp", ":65432", config)
    if err != nil {
        log.Fatal("unable to listen: ", err)
    }
    log.Print("proxy is listening on ", ln.Addr())
    proxy.Serve(ln)
    // TODO: check return value of Serve and support graceful shutdown
}
