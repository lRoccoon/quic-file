package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
)

func Client(address, file string) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-file"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	defer session.Close()
	go func() {
		s, err := session.AcceptStream()
		if err != nil {
			log.Fatalf("accept data stream error: %v", err)
		}
		defer s.Close()
		s.Write([]byte("This is a test file."))
		time.Sleep(time.Second)
	}()
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
	fmt.Printf("send %d bytes\n", sendBytes)
	time.Sleep(time.Second)
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
