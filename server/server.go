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
	"time"
	"strings"
	"math"
	"github.com/lucas-clemente/quic-go"
)

// Server 启动服务端
func Server(address, file string, test bool) {
	listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	sess, err := listener.Accept()
	if err != nil {
		log.Printf("accept session error: %v\n", err)
	}
	for ts:=0;math.Abs(float64(ts))<=0;{
		if test {
			stream, err := sess.AcceptStream()
			if err != nil {
				log.Printf("accept stream error: %v\n", err)
				return
			}
			buf, err := WriteFile(file)
			if err != nil {
				log.Printf("open file error: %v\n", err)
				continue
			}
			//接受文件
			recvByte, err := io.Copy(buf, stream)
			buf.Flush()
			if err != nil {
				log.Printf("write file error: %v\n", err)
			}
			fmt.Printf("recv %d bytes\n", recvByte)
			//MD5加密算法，发送密文
			var n []byte
			_,err =buf.Write(n)
			h := md5.New()
			h.Write(n)
			s := hex.EncodeToString(h.Sum(nil))
            streams,err :=sess.OpenStreamSync()
		    p :=strings.NewReader(s)
		    _, err = io.Copy(streams, p)
		    if err !=nil{
              log.Printf("secret send err: %v\n",err)
		    }
			streams.Close()	
			fmt.Println("file upload end,and if you want to upload again,please input ok,else input exit")	
		} else if test == false {
			stream, err := sess.OpenStreamSync()
			if err != nil {
				log.Fatalf("open stream error: %v\n", err)
			}
			data, _ := ReadFile(file)
			//MD5加密算法，验证文件完整性
			var m []byte
			_,err =data.Read(m)
			h := md5.New()
			h.Write(m)
			s := hex.EncodeToString(h.Sum(nil))
	        //发送文件
			sendBytes, err := io.Copy(stream, data)
			if err != nil {
				log.Fatalf("write stream error: %v\n", err)
			}
			fmt.Printf("send %d \n", sendBytes)
			stream.Close()
			//接受密文
			streams,err :=sess.AcceptStream()
			str :=bufio.NewReader(streams)
			var w []byte
			_, err =str.Read(w)
			t :=strings.Compare(s,string(w))
			if t!=0 {
				log.Fatalf("file transmissison err:%v\n",err)
			} else{
				time.Sleep(time.Millisecond * 1)
				streams.Close()
			}	
			fmt.Println("file download end,and if you want to download again,please input ok,else input exit")
		}
		var end string
		fmt.Scanf(end)
		ts =strings.Compare(end,"ok")
	}
	sess.Close()
	return 
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
