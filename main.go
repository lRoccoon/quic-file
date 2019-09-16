package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"sync"

	"github.com/873314461/quic-file/client"
	"github.com/873314461/quic-file/server"
)

func main() {
	serverAddr := flag.String("s", "", "server listen address")
	clientAddr := flag.String("c", "", "client connect server address")
	downloadFlag := flag.Bool("d", false, "download file from server")
	flag.Parse()
	files := flag.Args()
	if (len(*serverAddr) == 0 && len(*clientAddr) == 0) || (len(*serverAddr) != 0 && len(*clientAddr) != 0) {
		log.Fatalln("server or client?")
	}
	if len(*serverAddr) > 0 {
		s := server.NewFileServer(*serverAddr, generateTLSConfig(), nil)
		s.Run()
	}
	if len(*clientAddr) > 0 {
		c := client.NewFileClient(*clientAddr)
		var wg sync.WaitGroup
		for _, file := range files {
			wg.Add(1)
			go func(file string) {
				var err error
				if *downloadFlag {
					err = c.Download(file)
				} else {
					err = c.Upload(file)
				}
				if err != nil {
					log.Printf("upload/download file error: %v\n", err)
				} else {
					log.Printf("upload/download file success: %s\n", file)
				}
				wg.Done()
			}(file)
		}
		wg.Wait()
		c.Close()
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-file"},
	}
}
