package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/873314461/quic-file/common"
	"github.com/lucas-clemente/quic-go"
)

type StreamHandler struct {
	Ctx    context.Context
	Stream quic.Stream
	Reader io.Reader
	Writer io.Writer
}

func NewStreamHandler(stream *quic.Stream) *StreamHandler {
	return &StreamHandler{
		Stream: *stream,
		Reader: io.Reader(*stream),
		Writer: io.Writer(*stream),
		Ctx:    context.Background(),
	}
}

func (h *StreamHandler) Run() {
	defer h.Stream.Close()
	tmp := make([]byte, 1, 1)
	len, err := h.Reader.Read(tmp)
	if err != nil {
		log.Printf("read byte error: %v", err)
		return
	}
	if len != 1 {
		log.Printf("read byte len != 1")
		return
	}
	op := uint8(tmp[0])
	switch op {
	case 1:
		log.Printf("upload op: %d", op)
		err := h.handlerUpload()
		if err != nil {
			log.Printf("handler upload error: %v", err)
		}
	case 2:
		log.Printf("download op: %d", op)
	default:
		log.Printf("unknow op: %d", op)
	}

}

func (h *StreamHandler) handlerUpload() error {
	lenBytes := make([]byte, 2, 2)
	readn, err := h.Reader.Read(lenBytes)
	if err != nil {
		return fmt.Errorf("read path len error: %v", err)
	}
	if readn != 2 {
		return errors.New("readn != 2")
	}
	pathLen := binary.BigEndian.Uint16(lenBytes)

	path := make([]byte, pathLen, pathLen)
	readn, err = h.Reader.Read(path)
	if err != nil {
		return fmt.Errorf("read path error: %v", err)
	}
	if readn != int(pathLen) {
		return errors.New("readn != path len")
	}

	tmpAbsPath, err := filepath.Abs(string(path) + common.TempFileSuffix)
	if err != nil {
		return fmt.Errorf("get tmp abs path error: %v", err)
	}
	absPath, err := filepath.Abs(string(path))
	if err != nil {
		return fmt.Errorf("get abs path error: %v", err)
	}
	dataLenBytes := make([]byte, 8, 8)
	readn, err = h.Reader.Read(dataLenBytes)
	if err != nil {
		return fmt.Errorf("read data len error: %v", err)
	}
	if readn != 8 {
		return errors.New("readn != 8")
	}
	dataLen := binary.BigEndian.Uint64(dataLenBytes)
	file, err := os.Create(tmpAbsPath)
	if err != nil {
		return fmt.Errorf("creat file error: %v", err)
	}
	writen, err := io.Copy(file, h.Reader)
	if err != nil {
		return fmt.Errorf("write file error: %v", err)
	}
	if dataLen != uint64(writen) {
		return errors.New("data len != writen")
	}
	file.Close()

	err = os.Rename(tmpAbsPath, absPath)
	if err != nil {
		return fmt.Errorf("rename file error: %v", err)
	}
	return nil
}
