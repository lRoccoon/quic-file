package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/873314461/quic-file/common"
	"github.com/lucas-clemente/quic-go"
)

type FileClient struct {
	Session quic.Session
	Ctx     context.Context
}

func NewFileClient(address string) *FileClient {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-file"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	return &FileClient{
		Session: session,
		Ctx:     context.Background(),
	}
}

func (c *FileClient) Close() {
	time.Sleep(time.Second)
	c.Session.Close()
	time.Sleep(time.Second)
}

func (c *FileClient) Upload(file string) error {
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		return fmt.Errorf("open stream error: %v", err)
	}
	defer stream.Close()

	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(1))
	if err != nil {
		return fmt.Errorf("write op error: %v", err)
	}
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writen, err := writer.Write(pathLenBytes)
	if err != nil {
		return fmt.Errorf("write path len error: %v", err)
	}
	if writen != 2 {
		return errors.New("path len != 2")
	}
	writen, err = writer.WriteString(file)
	if err != nil {
		return fmt.Errorf("write path error: %v", err)
	}
	if writen != len(file) {
		return fmt.Errorf("writen != path len, %d, %d", writen, len(file))
	}
	fileReader, size := ReadFile(file)
	defer fileReader.Close()
	dataLenBytes := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(dataLenBytes, size)
	writen, err = writer.Write(dataLenBytes)
	if err != nil {
		return fmt.Errorf("write path len error: %v", err)
	}
	if writen != 8 {
		return errors.New("data len != 8")
	}
	writeFileN, err := writer.ReadFrom(fileReader)
	if err != nil {
		return fmt.Errorf("write data error: %v", err)
	}
	if uint64(writeFileN) != size {
		return errors.New("write file n != file size")
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("writer flush error: %v", err)
	}
	return nil
}

func ReadFile(file string) (*os.File, uint64) {
	fp, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file error: %v\n", err)
	}
	fileInfo, err := fp.Stat()
	if err != nil {
		log.Fatalf("get file info error: %v\n", err)
	}
	return fp, uint64(fileInfo.Size())
}

func (c *FileClient) Download(file string) error {
	stream, err := c.Session.OpenStreamSync(c.Ctx)
	if err != nil {
		return fmt.Errorf("open stream error: %v", err)
	}
	defer stream.Close()

	writer := bufio.NewWriter(stream)
	err = writer.WriteByte(byte(2))
	if err != nil {
		return fmt.Errorf("write op error: %v", err)
	}
	pathLenBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pathLenBytes, uint16(len(file)))
	writeN, err := writer.Write(pathLenBytes)
	if err != nil {
		return fmt.Errorf("write path len error: %v", err)
	}
	if writeN != 2 {
		return errors.New("path len != 2")
	}
	writeN, err = writer.WriteString(file)
	if err != nil {
		return fmt.Errorf("write path error: %v", err)
	}
	if writeN != len(file) {
		return fmt.Errorf("writeN != path len, %d, %d", writeN, len(file))
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("writer flush error: %v", err)
	}
	tmpAbsPath, err := filepath.Abs(file + common.TempFileSuffix)
	if err != nil {
		return fmt.Errorf("get tmp abs path error: %v", err)
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("get abs path error: %v", err)
	}
	tmpFile, err := os.Create(tmpAbsPath)
	if err != nil {
		return fmt.Errorf("creat file error: %v", err)
	}
	recvN, err := io.Copy(tmpFile, stream)
	if err != nil {
		return fmt.Errorf("write file error: %v", err)
	}
	log.Printf("recv file: %s[%d bytes] from server", file, recvN)
	tmpFile.Close()
	err = os.Rename(tmpAbsPath, absPath)
	if err != nil {
		return fmt.Errorf("rename file error: %v", err)
	}
	return nil
}
