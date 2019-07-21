package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

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
	defer session.Close()
	cmdStream, err := session.OpenStreamSync()
	if err != nil {
		log.Fatalf("open stream error: %v\n", err)
	}
	defer cmdStream.Close()
	writer := bufio.NewWriter(cmdStream)
	sendBytes, err := writer.WriteString("PUT test.bin")
	if err != nil {
		log.Printf("write stream error: %v\n", err)
	}
	writer.Flush()
	fmt.Printf("send %d bytes", sendBytes)
	// time.Sleep(time.Microsecond * 10)
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
