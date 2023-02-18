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
	"google.golang.org/grpc/status"
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

	return &pb.KeyResult{
		Key: k,
	}, nil
}

func (s *PayloadServer) RefreshToken(ctx context.Context, req *pb.AuthRequest) (*pb.Refresh, error) {
	userId := req.GetUserId()
	
	if encr, err := s.passCache.passIsValid(userId); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
		// should be encrypted with server pub key
	} else if encr != req.GetPassphrase() {
		return nil, status.Errorf(codes.Unauthenticated, "password not valid")
	}

	token, err := MakeRefreshToken(userId)
	if err != nil {
		return nil, err
	}
	return &pb.Refresh{
		Token: token,
	}, nil
}

func (s *PayloadServer) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResult, error) {
	// verify
	userId, ok, err := VerifyPassphrase(req.GetPassphrase())
	if err != nil || userId < 1 {
		glog.Error(err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	if !ok {
		glog.Errorf("unable to match passphrase %v", req.GetPassphrase())
		return nil, grpc.Errorf(codes.Unauthenticated, "unable to verify passphrase")
	}
	// generate keys
	kSet, err := GenKeyPairForUser(userId, []byte(req.GetKeys().GetEncryptionKey()), []byte(req.GetKeys().GetSignKey()))
	if err != nil {
		glog.Errorf("error generating key pair")
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}

	refToken, err := MakeRefreshToken(userId, kSet.pubEncrKey, kSet.privSign)
	if err != nil {
		glog.Errorf("unable to make refresh token: %v", err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	return &pb.AuthResult{
		Keys: &pb.KeySet{
			EncryptionKey: kSet.pubEncrKey,
			SignKey: kSet.pubSignKey,
		},
		TokenSet: &pb.TokenSet{
			UserId: userId,
			RefreshToken: refToken,
		},
	}, nil
}

func (s *PayloadServer) GetAccess(ctx context.Context, req *pb.RefreshRequest) (*pb.Access, error) {
	if err := VerifyRefreshToken(req.GetTokenSet().GetUserId(), ); err != nil {
		return nil, err
	}
	token := GenerateToken(50)
	encr, err := Encrypt(req.GetTokenSet().GetUserId(), []byte(token))
	// TODO how to deal with cache?
	// s.tokenCache.add(, req.GetUserId())
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
