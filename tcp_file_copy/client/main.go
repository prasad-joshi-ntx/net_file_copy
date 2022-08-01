package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
	// "time"
)

var (
	addr     = flag.String("listen", "localhost:8080", "port to listen to")
	filename = flag.String("filename", "", "file to send")
)

func SendFile(address string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	filenameSize := int64(len(filename))

	if err = binary.Write(conn, binary.LittleEndian, filenameSize); err != nil {
		log.Println(err)
		return
	}

	if _, err = io.WriteString(conn, filename); err != nil {
		log.Println(err)
		return
	}

	stat, _ := file.Stat()
	if err = binary.Write(conn, binary.LittleEndian, stat.Size()); err != nil {
		log.Println(err)
		return
	}

	br := bufio.NewReader(file)
	bw := bufio.NewWriter(conn)
	defer bw.Flush()

	log.Println("Starting file transfer size=", stat.Size())
	if _, err = io.CopyN(bw, br, stat.Size()); err != nil {
		log.Println(err)
		return
	}
	log.Println("File sent.")
}

func main() {
	flag.Parse()
	SendFile(*addr, *filename)
}
