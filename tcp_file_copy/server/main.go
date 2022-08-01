package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	addr      = flag.String("listen", "localhost:8080", "port to listen to")
	writeFile = flag.Bool("write_file", false, "controlls whether data should be written to file.")
)

func Server(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleIncoming(conn)
	}
}

type NullWriter struct {
}

func (self *NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (self *NullWriter) Flush() (err error) {
	return nil
}

func HandleIncoming(conn io.ReadWriteCloser) {
	defer conn.Close()

	var filenameSize int64
	err := binary.Read(conn, binary.LittleEndian, &filenameSize)
	if err != nil {
		log.Println(err)
		return
	}

	filename := make([]byte, int(filenameSize))
	if _, err = io.ReadFull(conn, filename); err != nil {
		log.Println(err)
		return
	}

	var fileSize int64

	if err = binary.Read(conn, binary.LittleEndian, &fileSize); err != nil {
		log.Println(err)
		return
	}

	file, err := os.Create(string(filename) + ".server")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	if err := file.Truncate(fileSize); err != nil {
		log.Println(err)
		return
	}

	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(file)
	defer bw.Flush()

	if _, err = io.CopyN(bw, br, fileSize); err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	Server(*addr)
}
