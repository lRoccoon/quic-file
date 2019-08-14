package server

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
//	"time"
	"strings"
	"github.com/lucas-clemente/quic-go"
)

// Server 启动服务端
func Server(address string, test bool)  {
	listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	} else {
		fmt.Println("the server is listening")
	}
	for {
			sess, err := listener.Accept()
			defer sess.Close()
			if err != nil {
			log.Printf("accept session error: %v\n", err)
			}
			if test {
				stream, err := sess.AcceptStream()
				if err != nil {
					log.Printf("accept stream error: %v\n", err)
					return
				}
				//接受文件名
				data := make([]byte, 1024)
				wc, err := stream.Read(data)
				fmt.Println("the name of file is " + string(data[0:wc]))
				//接受文件
				files :=string(data[0:wc])
				buf, err := WriteFile(files)
				recvByte, err := io.Copy(buf, stream)
				buf.Flush()
				if err != nil {
					log.Printf("write file error: %v\n", err)
				}
				fmt.Printf("recv %d bytes\n", recvByte)
				stream.Close()
				//MD5加密算法，发送密文
				streams,err :=sess.OpenStreamSync()
				var n []byte
				fp,err:=os.Open(files)
				_,err =fp.Read(n)
				h := md5.New()
				h.Write(n)
				s := hex.EncodeToString(h.Sum(nil))
				_,err =streams.Write([] byte(s))
				if err !=nil{
				log.Printf("secret send err: %v\n",err)
				}
				streams.Close()	
			} else if test == false {
				//接受文件名
				streams,err :=sess.AcceptStream()
				data := make([]byte, 1024)
				wc, err := streams.Read(data)
				fmt.Println("the name of file is " + string(data[0:wc]))
				streams.Close()
				//MD5加密算法，发送密文
				stream,err :=sess.OpenStreamSync()
				if err != nil {
					log.Fatalf("open stream error: %v\n", err)
				}
				var n []byte
				fp,err:=os.Open(string(data[0:wc]))
				_,err =fp.Read(n)
				h := md5.New()
				h.Write(n)
				s := hex.EncodeToString(h.Sum(nil))
				_,err =stream.Write([] byte(s))
				if err !=nil{
				log.Printf("secret send err: %v\n",err)
				}
				//发送文件
				datas, _ := ReadFile(string(data[0:wc]))
				sendBytes, err := io.Copy(stream, datas)
				if err != nil {
					log.Fatalf("write stream error: %v\n", err)
				}
				fmt.Printf("send %d bytes\n", sendBytes)
				stream.Close()
			}
		fmt.Println("file transmission end,and if you want to transmission again,please input ok,else input exit")
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		var ex string
		ex ="exit\r\n"
		ts :=strings.Compare(input,ex)
		if ts==0{
			return 
		} else{
			fmt.Println("the server is listening")
		}
	}
}

// WriteFile 写入文件
func WriteFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(".\\server\\"+file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
}

// ReadFile 读取文件
func ReadFile(file string) (*bufio.Reader, int64) {
	fp, err := os.Open(".\\server\\"+file)
	if err != nil {
		log.Fatalf("open file error: %v\n", err)
	}
	fileInfo, err := fp.Stat()
	if err != nil {
		log.Fatalf("get file info error: %v\n", err)
	}
	return bufio.NewReader(fp), fileInfo.Size()
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
