package payload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/DJolley12/home_cloud/protos"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	// "google.golang.org/grpc/metadata"
)

type PayloadServer struct {
	basePath string
	pb.PayloadServer
	tokenCache tokenCache
	passCache  passCache
	tokenSize  int
}

func NewPaylodServer(basePath string) (*PayloadServer, error) {
	if _, err := os.Stat(basePath); err != nil {
		return &PayloadServer{}, fmt.Errorf("path error: %v", err)
	}

	c := newTokenCache(100, 5)
	p := newPassCache(100, 5)
	server := &PayloadServer{
		basePath:   basePath,
		tokenCache: c,
		passCache:  p,
	}
	go func() {
		server.tokenCache.collectTokens()
		server.passCache.collectPass()
		time.Sleep(15 * time.Second)
	}()

	return server, nil
}

func (s *PayloadServer) RequestKey(ctx context.Context, req *pb.KeyRequest) (*pb.KeyResult, error) {
	k, err := GenKeyPairForUser(req.GetUserId(), req.GetKey())
	if err != nil {
		return nil, err
	}

	return &pb.KeyResult{
		Key: k,
	}, nil
}

func (s *PayloadServer) RefreshToken(ctx context.Context, req *pb.AuthRequest) (*pb.Refresh, error) {
	userId := req.GetUserId()
	encr, err := s.passCache.passIsValid(userId)
	if err != nil {
		return nil, err
	}
	match, err := VerifyPassphrase(req.GetPassphrase(), encr, userId)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, grpc.Errorf(codes.Unauthenticated, "unable to verify passphrase")
	}
	token, err := MakeRefreshToken(userId)
	if err != nil {
		return nil, err
	}
	return &pb.Refresh{
		Token: token,
	}, nil
}

func (s *PayloadServer) Authorize(ctx context.Context, req *pb.Refresh) (*pb.Access, error) {
	if err := VerifyRefreshToken(req.GetUserId(), req.GetToken()); err != nil {
		return nil, err
	}
	token := GenerateToken(50)
	encr, err := Encrypt(req.GetUserId(), []byte(token))
	// TODO how to deal with cache?
	s.tokenCache.add(, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.Access{
		Token: string(encr),
	}, nil
}

func (s *PayloadServer) ReceivePayload(stream pb.Payload_ReceivePayloadServer) error {
	size := 0
	if !s.tokenCache.tokenIsValid(stream.Context()) {
		return grpc.Errorf(codes.Unauthenticated, "invalid access token")
	}
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
