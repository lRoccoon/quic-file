package server

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/lucas-clemente/quic-go"
)

// Server 启动服务端
func Server(address string) {
	listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		sess, err := listener.Accept()
		if err != nil {
			log.Printf("accept session error: %v\n", err)
			continue
		}
		stream, err := sess.AcceptStream()
		if err != nil {
			log.Printf("accept stream error: %v\n", err)
			return
		}
		buf, err := WriteFile("output.bin")
		if err != nil {
			log.Printf("open file error: %v\n", err)
			continue
		}
		recvByte, err := io.Copy(buf, stream)
		buf.Flush()
		if err != nil {
			log.Printf("write file error: %v\n", err)
		}
		log.Printf("recv %d bytes\n", recvByte)
	}
}

func WriteFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
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
		NextProtos:   []string{"quic-echo-example"},
	}
}
