package payload

import (
	"fmt"
	"io"
	"os"

	pb "github.com/DJolley12/home_cloud/protos"
)

type PayloadServer struct {
	basePath string
	pb.PayloadServer
}

func NewPaylodServer(basePath string) (*PayloadServer, error) {
	if _, err := os.Stat(basePath); err != nil {
		return &PayloadServer{}, fmt.Errorf("path error: %v", err)
	}
	return &PayloadServer{
		basePath: basePath,
	}, nil
}

func (s *PayloadServer) ReceivePayload(stream pb.Payload_ReceivePayloadServer) error {
	size := 0
	in, err := stream.Recv()
	if err != nil {
		return err
	}

	if err := createFile(in.GetId()); err != nil {
		return err
	}

	if err := writeFile(in.GetId(), in.GetChunk()); err != nil {
		return err
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := writeFile(in.GetId(), in.GetChunk()); err != nil {
			return err
		}

		size += len(in.GetChunk())
	}

	return err
}

func (s *PayloadServer) SendPayload(req *pb.DownloadRequest, sendServer pb.Payload_SendPayloadServer) error {
	return nil
}
