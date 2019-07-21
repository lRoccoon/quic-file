package server

import (
	"bufio"

	"crypto/tls"
	"io"
	"log"
	"os"

	"github.com/873314461/quic-file/common"

	"github.com/lucas-clemente/quic-go"
)

// FileServer 文件服务端
type FileServer struct {
	Address    string
	TLSConfig  *tls.Config
	QuicConfig *quic.Config
	Sessions   []*quic.Session
	Listener   quic.Listener
	Streams    map[quic.StreamID]*quic.Stream
}

// NewFileServer 创建FileServer对象
func NewFileServer(address string, tlsConfig *tls.Config, quicConfig *quic.Config) *FileServer {
	return &FileServer{
		Address:    address,
		TLSConfig:  tlsConfig,
		QuicConfig: quicConfig,
		Sessions:   make([]*quic.Session, 0),
		Streams:    make(map[quic.StreamID]*quic.Stream, 0),
	}
}

// Run 启动服务端
func (s *FileServer) Run() error {
	var err error
	s.Listener, err = quic.ListenAddr(s.Address, s.TLSConfig, s.QuicConfig)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		sess, err := s.Listener.Accept()
		if err != nil {
			log.Printf("accept session error: %v\n", err)
			continue
		}
		s.Sessions = append(s.Sessions, &sess)

		go s.handler(&sess)
	}
}

func (s *FileServer) handler(session *quic.Session) {
	sess := *session
	cmdStream, err := sess.AcceptStream()
	if err != nil && err.Error() != common.NoError || cmdStream == nil {
		log.Printf("accept cmd stream error: %v, cmd: %v\n", err, cmdStream)
		return
	}
	cmdReader := bufio.NewReader(cmdStream)
	p := 0
	data := make([]byte, 1024)
	for {
		n, err := cmdReader.Read(data[p:])
		if err != nil {
			if err == io.EOF || err.Error() == common.NoError {
				break
			}
			log.Printf("read cmd stream error: %v\n", err)
			break
		}
		p += n
		log.Printf("recv cmd:%s\n", data)
	}
	handlerCMD(string(data), bufio.NewWriter(cmdStream))
}

func writeFile(file string) (*bufio.Writer, error) {
	fp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(fp), nil
}
