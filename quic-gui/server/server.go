package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"os"
	"context"
	"path/filepath"
	"strconv"
	quic "github.com/lucas-clemente/quic-go"
)

// Server 接入连接
func Server(address string,test bool) {
	listener, errs := quic.ListenAddr(address, generateTLSConfig(), nil)
	cxt :=context.TODO()
	if errs != nil {
		log.Fatalf("listen error: %v\n", errs)
	} else {
		result="the server is listening\n"+"the address is "+address+"\n"
		rs.SetText(result)
	}
	for {
			sess, err := listener.Accept(cxt)
			if err != nil {
				log.Fatal("accept session error: %v\n", err)
			}
			if test {
				stream, err := sess.AcceptStream(cxt)
				if err != nil {
					log.Fatalf("accept stream error: %v\n", err)
				}
				//接受文件名
				data := make([]byte, 1024)
				wc, err := stream.Read(data)
				//接受文件
				file :=string(data[0:wc])
				filename :=filepath.Base(file)
				rs.Append(result)
				buf, err := WriteFile(filename)
				recvBytes, err := io.Copy(buf, stream)
				recvByte:= strconv.FormatInt(recvBytes,10)
				buf.Flush()
				if err != nil {
					log.Fatalf("write file error: %v\n", err)
				}
				result="the name of file is " +filename+" and recv "+ recvByte +" bytes\n"
				rs.Append(result)
				stream.Close()
				//MD5加密算法，发送密文
				streams,err :=sess.OpenStreamSync(cxt)
				var n []byte
				fp,err:=os.Open(file)
				_,err =fp.Read(n)
				h := md5.New()
				h.Write(n)
				s := hex.EncodeToString(h.Sum(nil))
				_,err =streams.Write([] byte(s))
				if err !=nil{
					log.Fatalf("secret send err: %v\n",err)
				}
				streams.Close()	
			} else if test == false {
				//接受文件名
				streams,err :=sess.AcceptStream(cxt)
				data := make([]byte, 1024)
				wc, err := streams.Read(data)
				streams.Close()
				//MD5加密算法，发送密文
				stream,err :=sess.OpenStreamSync(cxt)
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
					log.Fatalf("secret send err: %v\n",err)
				}
				//发送文件
				datas, _ := ReadFile(string(data[0:wc]))
				sendBytes, err := io.Copy(stream, datas)
				sendByte:= strconv.FormatInt(sendBytes,10)
				if err != nil {
					log.Fatalf("write stream error: %v\n", err)
				}
				result="the name of file is "+ string(data[0:wc])+" and send "+ sendByte +" bytes\n"
				rs.Append(result)
				stream.Close()
			}	
	}
}
// WriteFile 写入文件
func WriteFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
}

// ReadFile 读取文件
func ReadFile(file string) (*bufio.Reader, int64) {
	fp, err := os.Open(file)
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
