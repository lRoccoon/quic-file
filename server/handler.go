package server

import (
	"bufio"
	"log"
	"strings"
)

func handlerCMD(cmd string, w *bufio.Writer) error {
	args := strings.Split(cmd, " ")
	switch args[0] {
	case "PUT":
		log.Println("path: ", args[1])
	default:
		log.Println("unknow cmd: ", cmd)
	}
	return nil
}

func putHandler(path string, w *bufio.Writer) error {
	// absPath, err := filepath.Abs(path)
	// if err != nil {
	// 	return err
	// }
	// file, err := os.Create(absPath)
	// if err != nil {
	// 	return err
	// }
	// sess.AcceptStream()
	return nil
}
