package payload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	pb "github.com/DJolley12/home_cloud/protos"
	"github.com/golang/glog"
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
		glog.Error(err)
		return err
	}

	f, err := createFile(filepath.Join(s.basePath, in.GetId()))
	if err != nil {
		glog.Error(err)
		return err
	}

	glog.Errorf("chunk: %v", string(in.GetChunk()))
	if err := writeFile(f, in.GetChunk()); err != nil {
		glog.Error(err)
		return err
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			glog.Error(err)
			return err
		}

		if err := writeFile(f, in.GetChunk()); err != nil {
			glog.Error(err)
			return err
		}

		size += len(in.GetChunk())
	}

	return stream.SendAndClose(&pb.UploadResult{
		RecvSize:  int32(size),
		IsSuccess: true,
	})
}

func (s *PayloadServer) SendPayload(req *pb.DownloadRequest, sendServer pb.Payload_SendPayloadServer) error {
	return nil
}
