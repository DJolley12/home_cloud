package payload

import (
	"bytes"
	"context"
	"io"
	"os"

	pb "github.com/DJolley12/home_cloud/protos"
	"github.com/golang/glog"
)

type PayloadClient struct {
	client    pb.PayloadClient
	chunkSize int
}

func (c *PayloadClient) UploadFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		glog.Errorf("error opening file %v - err: %v", filePath, err)
	}
	defer file.Close()
	p, err := c.client.ReceivePayload(context.Background())
	if err != nil {
		glog.Errorf("error getting payload upload: %v", err)
	}

	for {
		r := io.LimitReader(file, int64(c.chunkSize))
		b := make([]byte, c.chunkSize)
		buf := bytes.NewBuffer(b)
		n, sendErr := io.Copy(buf, r)
		if sendErr != nil {
		  glog.Errorf("error reading file %v", sendErr)
		} else if n < 1 {
		  glog.Info("done reading file")
		  break
		}
		req := &pb.DataChunk{
		  Id: file.Name(),
		  Chunk: b,
		}
		sendErr = p.Send(req)
		if sendErr != nil {
		  glog.Errorf("error sending data chunk: %v", sendErr)
		  break
		}
	}
}
