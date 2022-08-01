package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"bufio"
	"io"
	"os"

	pb "github.com/prasad-joshi-ntx/net_file_copy/file-copy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr       = flag.String("addr", "localhost:50051", "the address to connect to")
	name       = flag.String("name", defaultName, "Name to greet")
	file_name  = flag.String("file_name", "", "Name of the file to copy")
	block_size = flag.Int64("block_size", 512*1024, "block size")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFileCopyClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fd, err := os.Open(*file_name)
	if err != nil {
		log.Fatalf("Failed to open file ", *file_name, err)
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	buf := make([]byte, *block_size)

	fmt.Printf("Sending file ", *file_name)
	count := 0
	for {
		read, readErr := reader.Read(buf)
		if readErr != nil {
			if readErr != io.EOF {
				log.Fatalf("error while reading file ", readErr)
			}
			break
		}

		var err error
		var bytesCopied int64
		var errorCode int32

		fmt.Printf("count ", count, read, "\n")
		count += 1
		failed := true
		for ii := 0; ii < 10; ii++ {
			resp, writeErr := c.Write(ctx,
				&pb.WriteArgs{FileName: "1", Data: hex.EncodeToString(buf[0:read]), Offset: 0})
			if err != nil || resp.GetError() != 0 || resp.GetByesCopied() != int64(read) {
				err = writeErr
				errorCode = resp.GetError()
				bytesCopied = resp.GetByesCopied()
				time.Sleep(250 * time.Millisecond)
				continue
			}
			failed = false
			break
		}
		if failed == true {
			log.Fatalf("Failed.", count, err, errorCode, bytesCopied)
		}
	}
}
