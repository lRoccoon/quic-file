package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/873314461/quic-file/common"
	"github.com/lucas-clemente/quic-go"
)

func handlerCMD(session *quic.Session, w *bufio.Writer, cmd string) error {
	args := strings.Split(cmd, " ")
	switch args[0] {
	case "PUT":
		log.Println("path: ", args[1])
		return putHandler(session, w, args[1])
	default:
		log.Println("unknow cmd: ", cmd)
	}
	return nil
}

func putHandler(session *quic.Session, w *bufio.Writer, path string) error {
	sess := *session

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get abs path error: %v", err)
	}
	tmpPath := absPath + common.TempFileSuffix

	dataStream, err := sess.OpenStreamSync()
	if err != nil {
		return fmt.Errorf("open stream error: %v", err)
	}

	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create tmp file error: %v", err)
	}

	response := fmt.Sprintf("200 %d", dataStream.StreamID())
	sendBytes, err := w.WriteString(response)
	if err != nil {
		return fmt.Errorf("write to cmd error: %v", err)
	}
	log.Printf("send response to [%s] %d bytes, msg: %s\n", sess.RemoteAddr(), sendBytes, response)

	recvBytes, err := io.Copy(file, dataStream)
	if err != nil {
		return fmt.Errorf("write to tmp file error: %v", err)
	}
	log.Printf("recv [%s] %d bytes.\n", sess.RemoteAddr(), recvBytes)

	err = os.Rename(tmpPath, absPath)
	if err != nil {
		return fmt.Errorf("rename file error: %v", err)
	}
	return nil
}
