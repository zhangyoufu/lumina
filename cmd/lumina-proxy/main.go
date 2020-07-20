package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/c0va23/go-proxyprotocol"
	"github.com/zhangyoufu/lumina"
)

func main() {
	var (
		enableTls   bool
		enableProxy bool
		listenAddr  string
		tlsCertPath string
		tlsKeyPath  string
		idaKeyPath  string
		ln          net.Listener
	)
	flag.BoolVar(&enableTls, "tls", false, "enable TLS")
	flag.BoolVar(&enableProxy, "proxy", false, "enable PROXY protocol support")
	flag.StringVar(&listenAddr, "listen", ":8000", "listen address")
	flag.StringVar(&tlsCertPath, "tlsCert", "cert.pem", "path to TLS certificate (PEM format)")
	flag.StringVar(&tlsKeyPath, "tlsKey", "key.pem", "path to TLS certificate key (PEM format)")
	flag.StringVar(&idaKeyPath, "idaKey", "ida.key", "path to ida.key")
	flag.Parse()

	tcpListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("unable to listen: ", err)
	}
	ln = tcpListener

	if enableTls {
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
		if err != nil {
			log.Fatal("unable to load X509 key pair: ", err)
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		tlsListener := tls.NewListener(ln, config)
		ln = tlsListener
	}

	if enableProxy {
		// DefaultFallbackHeaderParserBuilder contains StubHeaderParserBuilder,
		// which accepts non-PROXY protocol traffic
		proxyListener := proxyprotocol.NewDefaultListener(ln)
		ln = proxyListener
	}

	f, err := os.Open(idaKeyPath)
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

	log.Print("proxy is listening on ", ln.Addr())
	log.Fatal("proxy stopped serving with error: ", proxy.Serve(ln))
	// TODO: graceful shutdown
}
