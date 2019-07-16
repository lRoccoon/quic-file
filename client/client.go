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

func Client(address, file string) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Fatalf("open stream error: %v\n", err)
	}
	data, size := ReadFile(file)
	if size > 4096*1024 {
		size = 4096 * 1024
	}
	sendBytes, err := io.Copy(stream, data)
	if err != nil {
		log.Fatalf("write stream error: %v\n", err)
	}
	fmt.Printf("send %d bytes", sendBytes)
	time.Sleep(time.Millisecond * 1)
	stream.Close()
	session.Close()
}

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
