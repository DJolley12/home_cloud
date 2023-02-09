package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main()  {
  serverAddr := flag.String("server-addr", "localhost:50051", "ip address and port for the file store server")
  filePath := flag.String("file-path", "", "file-path of file to upload")
  if *filePath == "" {
  	log.Fatal("cannot have empty file path")
  }
  conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
  	log.Fatalf("could not connect to server: %v", err)
  }
  defer conn.Close()

}
