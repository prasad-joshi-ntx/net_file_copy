package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/snappy"
	"github.com/prasad-joshi-ntx/net_file_copy/compression-benchmark/common"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	addr           = flag.String("listen", "localhost:8080", "port to listen to")
	filename       = flag.String("filename", "", "file to send")
	blockSize      = flag.Int("block_size", 32*1024, "block size")
	compressBuffer = flag.Bool("compress", true, "compress buffers over wire")
	compressRatio  = flag.Int("compress_percent", 50, "expected compression percetange")
	numBytesToSend = flag.Int("byes_to_send_in_gb", 100, "number of bytes to send")
	verbose        = flag.Bool("verbose", false, "verbose output")
)

func main() {
	flag.Parse()

	var buffer []byte
	if *compressBuffer {
		buffer = generateCompressibleBuffer(*blockSize, *compressRatio)
	} else {
		buffer = generateRandomBuffer(*blockSize)
	}
	if buffer == nil || len(buffer) <= 0 {
		log.Println("No buffer to send")
		return
	}

	numBlocks := (*numBytesToSend * 1024 * 1024 * 1024) / *blockSize
	benchmarkSendBuffer(*addr, buffer, *blockSize, *compressBuffer, numBlocks)
}

func generateCompressibleBuffer(blockSize int, compressRatio int) []byte {
	if compressRatio >= 100 {
		return bytes.Repeat([]byte("A"), blockSize)
	} else if compressRatio <= 0 {
		return generateRandomBuffer(blockSize)
	}

	unCompressibleRatio := 100 - compressRatio
	unCompressibleSz := (blockSize * unCompressibleRatio) / 100
	buffer := generateRandomBuffer(unCompressibleSz)
	sz := blockSize - unCompressibleSz
	return append(buffer, bytes.Repeat([]byte("A"), sz)...)
}

func generateRandomBuffer(blockSize int) []byte {
	rand.Seed(time.Now().UnixNano())
	buffer := make([]byte, blockSize)
	rand.Read(buffer)
	return buffer
}

func getCompressedBuffer(in []byte) []byte {
	return snappy.Encode(nil, in)
}

func benchmarkSendBuffer(
	address string,
	buffer []byte,
	blockSize int,
	shouldCompressBuffer bool,
	numBlocks int) {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	dst := bufio.NewWriter(conn)
	defer dst.Flush()

	var compressedSize int64 = int64(blockSize)
	if shouldCompressBuffer == true {
		compressedSize = int64(len(getCompressedBuffer(buffer)))
	}
	header := &common.TransferHeader{
		BlockSize:      int64(blockSize),
		CompressedSize: compressedSize,
		NumBlocks:      int64(numBlocks),
		Compression:    shouldCompressBuffer,
	}

	if err := WriteHeader(conn, header); err != nil {
		log.Println("Failed to write header ", err)
		return
	}

	var start time.Time = time.Now()
	for ii := 0; ii < numBlocks; ii += 1 {
		var buf []byte
		if shouldCompressBuffer == true {
			buf = getCompressedBuffer(buffer)
		} else {
			buf = buffer
		}
		if err := WriteBuffer(dst, buf); err != nil {
			log.Println(err)
			return
		}
	}
	elapsed := time.Since(start)
	var seconds float64 = elapsed.Seconds()
	bytesWrote := float64(header.BlockSize * header.NumBlocks)
	bw := float64(bytesWrote) / seconds
	// fmt.Println(fmt.Sprintf("%.2f,%.2f,%.2f", GB(bytesWrote), seconds, MB(bw)))

	// block size, compressed buffer size, num blocks, compress data,
	// compression ratio, total bytes sent, time, bw
	if *verbose == false {
		msg := fmt.Sprintf("%v,%v,%v,%v,%v,%.2f,%.2f,%.2f",
			header.BlockSize,
			header.CompressedSize,
			header.NumBlocks,
			header.Compression,
			*compressRatio,
			GB(bytesWrote),
			seconds,
			MB(bw))
		fmt.Println(msg)
	} else {
		msg := fmt.Sprintf("block_size=%v,compressed_size=%v,num_blocks=%v,"+
			"compressed_enabled=%v,compression_ratio=%v,bytes_transferred_gb=%.2f,"+
			"time_in_seconds=%.2f,bw_in_mb=%.2f",
			header.BlockSize,
			header.CompressedSize,
			header.NumBlocks,
			header.Compression,
			*compressRatio,
			GB(bytesWrote),
			seconds,
			MB(bw))
		fmt.Println(msg)
	}
}

func GB(bytes float64) float64 {
	return MB(bytes) / 1024
}

func MB(bytes float64) float64 {
	return bytes / 1024 / 1024
}

func WriteHeader(conn net.Conn, header *common.TransferHeader) error {
	err := binary.Write(conn, binary.LittleEndian, header.BlockSize)
	if err != nil {
		log.Println("failed to write block size ", err)
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, header.CompressedSize)
	if err != nil {
		log.Println("failed to write compressed size ", err)
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, header.NumBlocks)
	if err != nil {
		log.Println("failed to write number of blocks ", err)
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, header.Compression)
	if err != nil {
		log.Println("failed to write compression ", err)
		return err
	}
	return nil
}

func WriteBuffer(dst *bufio.Writer, buf []byte) error {
	wrote, err := dst.Write(buf)
	if err != nil {
		return err
	}
	if wrote != len(buf) {
		return errors.New("short write")
	}
	return nil
}
