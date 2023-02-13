package main

import (
	"flag"
	"fmt"
	"log"
	"net"

  pb "github.com/DJolley12/home_cloud/protos"
	"github.com/DJolley12/home_cloud/server/payload"
	"google.golang.org/grpc"
)

func main() {
  storeDir := flag.String("base-dir", "", "specify the directory for the file store")
  addr := flag.String("ip-address", "127.0.0.1", "ip address to serve on: defaults to localhost")
  port := flag.Int("port", 50051, "the server port: defaults to 50051")

  flag.Parse()
  lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", *addr, *port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  log.Printf("payload server listening on %v", lis.Addr())

  s := grpc.NewServer()
  p, err := payload.NewPaylodServer(*storeDir)
  if err != nil {
    log.Fatal(err)
  }
  pb.RegisterPayloadServer(s, p)
  if err := s.Serve(lis); err != nil {
  	log.Fatalf("failed to serve: %v", err)
  }

  if *storeDir == "" {
    log.Fatal("cannot start file service without a file store directory")
  }

  if err != nil {
    log.Fatal(err)
  }

  for {
  	in := ""
  	println("enter --exit-- to quit")
  	if _, err := fmt.Scanln(&in); err != nil {
  		log.Fatal("oops something went wrong :((")
  	}
  	if in == "exit" {
  		break
  	}
  }
}
