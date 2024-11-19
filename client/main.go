package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/bhushan-aruto/file_streaming_grpc/proto/filestream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConn, err := grpc.NewClient("127.0.0.1:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("error occured -> % v ", err.Error())
		return
	}

	defer grpcConn.Close()

	client := filestream.NewFileStreamServiceClient(grpcConn)
	fileRequest := &filestream.FileRequest{
		Filename: os.Args[1],
	}

	stream, err := client.DownloadFile(context.Background(), fileRequest)
	if err != nil {
		fmt.Printf("error occured -> %v ", err.Error())
		return
	}
	file, err := os.Create(fileRequest.Filename)
	if err != nil {
		fmt.Printf("errorr occure -> %v ", err.Error())
		return
	}
	defer file.Close()
	for {
		chunk, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("file downloaded succesfullly")
			break
		}

		if err != nil {
			fmt.Printf("error occure -> %v ", err.Error())
			return
		}

		var buffer []byte
		for _, value := range chunk.Chunk {
			if value == 0 {
				break
			}
			buffer = append(buffer, value)
		}
		size, err := file.Write(buffer)

		if err != nil {
			fmt.Printf("error occured -> %v ", err.Error())
			return
		}
		fmt.Printf("succcesfully written %v bytes \n", size)

	}
}
