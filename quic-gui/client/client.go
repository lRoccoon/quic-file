package  main
import (

	"bufio"
	"crypto/tls"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"strconv"
	"io"
	"log"
	"os"
	"context"
	quic "github.com/lucas-clemente/quic-go"
	"github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
)

// Client 启动客户端
func Client(address ,file string, test bool) {
	rs.Append(address)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos: []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err != nil {
		log.Fatalf("connect server error: %v\n", err)
	}
	if test {
		cxt :=context.TODO()
		stream, err := session.OpenStreamSync(cxt)
		if err != nil {
			log.Fatalf("open stream error: %v\n", err)
		}
		//发送文件名
		_, err = stream.Write([] byte(file))
		if err != nil {
			log.Fatalf("name send error")
		}
        //发送文件
		data,_:= ReadFile(file)
		sendBytes, err := io.Copy(stream, data)
		sendByte:= strconv.FormatInt(sendBytes,10)
		if err != nil {
			log.Fatalf("write stream error: %v\n", err)
		}
		result ="file send the "+sendByte+" bytes\n"
		rs.Append(result)
		stream.Close()
		//MD5加密算法
        var m []byte
		_,err =data.Read(m)
		h := md5.New()
		h.Write(m)
		s := hex.EncodeToString(h.Sum(nil))
		//接受密文
		w := make([]byte, 2048*10)
		streams,err :=session.AcceptStream(cxt)
		wc,err :=streams.Read(w)
		var x=w[0:wc]
		t :=strings.Compare(s,string(x))
		if t!=0 {
			log.Fatalf("file transmissison err:%v\n",err)
		} 
		streams.Close()
		session.Close()	
	} else if test == false {
		//发送文件名
		cxt :=context.TODO()
		streams,err :=session.OpenStreamSync(cxt)
		_, err = streams.Write([] byte(file))
		streams.Close()
		//接受密文
		stream, err := session.AcceptStream(cxt)
		if err != nil {
			log.Printf("accept stream error: %v\n", err)
		}
		data := make([]byte, 1024)
		wc, err := stream.Read(data)
		//接收文件
		buf, err := WriteFile(file)
		recvBytes, err := io.Copy(buf, stream)
		recvByte:= strconv.FormatInt(recvBytes,10)
		buf.Flush()
		if err != nil {
			log.Printf("receive stream data error: %v\n", err)
		}
		result ="recv the "+recvByte+" bytes\n"
		rs.Append(result)
		stream.Close()
		//MD5加密算法，接受密文
		var n []byte
		fp,err:=os.Open(file)
		_,err =fp.Read(n)
		h := md5.New()
		h.Write(n)
		s := hex.EncodeToString(h.Sum(nil))
		var x=data[0:wc]
		t :=strings.Compare(s,string(x))
		if t!=0 {
			log.Fatalf("file transmissison err:%v\n",err)
		} 
		session.Close()	
	}
	ui.MsgBox(mainwin,
		"congratulation ,operation completed",
		"please first close this window")
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
			log.Fatalf("the file is not exist \n")
		} else {
			log.Fatalf("get file info error: %v\n", err)
		}
	}
	return bufio.NewReader(fp), fileInfo.Size()
}

// WriteFile 写入文件
func WriteFile(file string)(*bufio.Writer, error){
	fp, err := os.Create(file)
	if err != nil {
		log.Fatalf("file create err: %v\n",err)
	}
	return bufio.NewWriter(fp), nil
}
