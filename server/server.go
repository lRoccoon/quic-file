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

	"github.com/873314461/quic-file/common"

	"github.com/lucas-clemente/quic-go"
)

type FileServer struct {
	sess quic.Session
}

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
		cmd, err := sess.AcceptStream()
		if err != nil && err.Error() != common.NoError || cmd == nil {
			log.Printf("accept cmd stream error: %v, cmd: %v\n", err, cmd)
			continue
		}
		cmdReader := bufio.NewReader(cmd)
		p := 0
		data := make([]byte, 1024)
		for {
			n, err := cmdReader.Read(data[p:])
			if err != nil {
				if err == io.EOF || err.Error() == common.NoError {
					break
				}
				log.Printf("read cmd stream error: %v\n", err)
				break
			}
			p += n
			log.Printf("recv cmd:%s\n", data)
		}
		handlerCMD(string(data), bufio.NewWriter(cmd))
	}
}

func acceptClient(listener *quic.Listener) {

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
