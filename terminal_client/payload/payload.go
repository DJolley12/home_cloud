package payload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/DJolley12/home_cloud/protos"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type PayloadClient struct {
	client    pb.PayloadClient
	chunkSize int
}

func NewPayloadClient(client pb.PayloadClient, chunkSize int) (PayloadClient, error) {
	if client == nil {
		err := fmt.Errorf("client cannot be nil")
		glog.Error(err)
		return PayloadClient{}, err
	}
	if chunkSize < 256000 {
		err := fmt.Errorf("cannot have chunk size smaller thank 250kB")
		glog.Error(err)
		return PayloadClient{}, err
	}
	return PayloadClient{
		client:    client,
		chunkSize: chunkSize,
	}, nil
}

func (c *PayloadClient) Authorize(passphrase string) (*pb.AuthResult, error) {
	// generate keys
	// send keys and passphrase to server
	// receive server keys, user id
	// decrypt token
	// save keys and token
}

func (c *PayloadClient) GetAccess(ctx context.Context, in *pb.RefreshRequest, opts ...grpc.CallOption) (*pb.Access, error) {
	// get keys and refresh token
	// encrypt and sign refresh token
	// send token, user id to server
	// get access token, decrypt, verify
}

func (c *PayloadClient) UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file %v - err: %v\n", filePath, err)
	}
	defer file.Close()
	p, err := c.client.ReceivePayload(context.Background())
	if err != nil {
		fmt.Printf("error getting payload upload: %v\n", err)
		return err
	}

	chunkCount := 0
	glog.Info("reading file")
	for {

		r := io.LimitReader(file, int64(c.chunkSize))
		b := make([]byte, 0, c.chunkSize)
		buf := bytes.NewBuffer(b)
		n, err := io.Copy(buf, r)
		if err != nil {
			fmt.Printf("error reading file %v\n", err)
		} else if n < 1 {
			fmt.Println("done reading file")
			break
		}
		chunkCount++
		glog.Infof("read chunk: %v, bytes: %v\n", chunkCount, n)

		req := &pb.DataChunk{
			Id:    filepath.Base(file.Name()),
			Chunk: buf.Bytes(),
		}
		err = p.Send(req)
		if err != nil {
			fmt.Printf("error sending data chunk: %v\n", err)
			break
		}
		fmt.Println("sent chunk")
	}
	res, err := p.CloseAndRecv()
	if err != nil {
		glog.Errorf("error received during close send %v", err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if !res.GetIsSuccess() || int64(res.GetRecvSize()) != fi.Size() {
		fmt.Printf("file upload unsuccessful - recv size %v, file size %v", res.GetRecvSize(), fi.Size())
	}

	return err
}
