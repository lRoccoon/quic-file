package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
)

// Client 启动客户端
func Client(address, file string, test bool) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
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
		data, size := ReadFile(file)
		/*
			if size > 4096*1024 {
				size = 4096 * 1024
			}*/
		fmt.Printf("the length of file is : %d\n", size)
		sendBytes, err := io.Copy(stream, data)
		if err != nil {
			log.Fatalf("write stream error: %v\n", err)
		}
		fmt.Printf("send %d bytes\n", sendBytes)
		time.Sleep(time.Millisecond * 1)
		stream.Close()
		session.Close()
	} else if test == false {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Printf("accept stream error: %v\n", err)
			return
		}
		buf, err := WriteFile(file)
		if err != nil {
			log.Printf("create file error: %v\n", err)
		}
		recvByte, err := io.Copy(buf, stream)
		buf.Flush()
		if err != nil {
			log.Printf("write file error: %v\n", err)
		}
		fmt.Printf("recv %d bytes\n", recvByte)
		time.Sleep(time.Millisecond * 1)
		stream.Close()
		session.Close()
	}
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

// WriteFile 写入文件
func WriteFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
}
