package client

import (
	"bufio"
	"crypto/tls"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
)

// Client 启动客户端
func Client(address ,file string, test bool) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos: []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	if test {
		stream, err := session.OpenStreamSync()
		if err != nil {
			log.Fatalf("open stream error: %v\n", err)
		}
		data,_:= ReadFile(file)
		//MD5加密算法
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
		fmt.Printf("file send %d bytes\n",sendBytes)
		stream.Close()
		//接受密文
		streams,err :=session.AcceptStream()
		str :=bufio.NewReader(streams)
		var w []byte
		_, err =str.Read(w)
		t :=strings.Compare(s,string(w))
		if t!=0 {
			log.Fatalf("file transmissison err:%v\n",err)
		} else{
			time.Sleep(time.Millisecond * 1)
			streams.Close()
			session.Close()
		}	
	} else if test == false {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Printf("accept stream error: %v\n", err)
			return
		}
		buf, err := WriteFile(file)
		if err != nil {
			log.Printf("write file error: %v\n", err)
		}
		//接收文件
		recvByte, err := io.Copy(buf, stream)
		buf.Flush()
		if err != nil {
			log.Printf("receive stream data error: %v\n", err)
		}
		fmt.Printf("recv %d bytes\n", recvByte)
		//MD5加密算法，发送密文
		var n []byte
		_,err =buf.Write(n)
		h := md5.New()
		h.Write(n)
		s := hex.EncodeToString(h.Sum(nil))
		streams,err :=session.OpenStreamSync()
		p :=strings.NewReader(s)
		_, err = io.Copy(streams, p)
		if err !=nil{
           log.Printf("secret send err: %v\n",err)
		}
		time.Sleep(time.Millisecond * 1)
		streams.Close()
		session.Close()
	}
}

// ReadFile 读取文件
func ReadFile(file string) (*bufio.Reader, int64){
	fp, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file error: %v\n", err)
	}
	fileInfo, err := fp.Stat()
	if err != nil  {
		if os.IsNotExist(err){
			fmt.Printf("the file is not exist \n")
		} else {
			log.Fatalf("get file info error: %v\n", err)
		}
	}
	return bufio.NewReader(fp), fileInfo.Size()
}

// WriteFile 写入文件
func WriteFile(file string)(*bufio.Writer, error){
	fp, err := os.Create(".\\server\\"+file)
	if err != nil {
		log.Fatalf("file create err: %v\n",err)
	}
	return bufio.NewWriter(fp), nil
}
