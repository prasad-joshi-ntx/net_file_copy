package main

import (
	"encoding/binary"
	"flag"
	"github.com/golang/snappy"
	"github.com/prasad-joshi-ntx/net_file_copy/compression-benchmark/common"
	"io"
	"log"
	"net"
)

var (
	addr = flag.String("listen", "localhost:8080", "port to listen to")
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

func HandleIncoming(conn io.ReadWriteCloser) {
	defer conn.Close()

	header, err := ReadHeader(conn)
	if err != nil {
		log.Println("Failed to read header ", err)
		return
	}
	log.Println("Header ", *header)

	buffer := make([]byte, header.CompressedSize)
	for ii := int64(0); ii < header.NumBlocks; ii += 1 {
		_, err := ReadBuffer(conn, header, buffer)
		if err != nil {
			log.Println("Failed reading buffer ", err)
			return
		}
	}
}

func ReadHeader(conn io.ReadWriteCloser) (*common.TransferHeader, error) {
	var blockSize int64
	err := binary.Read(conn, binary.LittleEndian, &blockSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var compressedSize int64
	err = binary.Read(conn, binary.LittleEndian, &compressedSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var numBlocks int64
	err = binary.Read(conn, binary.LittleEndian, &numBlocks)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var compression bool
	err = binary.Read(conn, binary.LittleEndian, &compression)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &common.TransferHeader{
		BlockSize:      blockSize,
		CompressedSize: compressedSize,
		NumBlocks:      numBlocks,
		Compression:    compression,
	}, nil
}

func ReadBuffer(
	conn io.ReadWriteCloser,
	header *common.TransferHeader,
	buffer []byte) ([]byte, error) {

	err := ReadBlock(conn, buffer)
	if err != nil {
		return nil, err
	}
	if header.Compression == false {
		return buffer, nil
	}
	len, _ := snappy.DecodedLen(buffer)
	uncompressedBuf := make([]byte, len)
	uncompressedBuf, err = snappy.Decode(uncompressedBuf, buffer)
	if err != nil {
		return nil, err
	}
	return uncompressedBuf, nil
}

func ReadBlock(
	conn io.ReadWriteCloser,
	buffer []byte) error {

	toRead := len(buffer)
	read := 0
	for read < toRead {
		bytesRead, err := conn.Read(buffer[read:])
		if err != nil {
			return err
		}
		read += bytesRead
	}
	return nil
}

func main() {
	flag.Parse()
	Server(*addr)
}
