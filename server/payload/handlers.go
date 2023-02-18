package payload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/DJolley12/home_cloud/protos"
	services "github.com/DJolley12/home_cloud/server/payload/services"
	"github.com/DJolley12/home_cloud/server/persist/ports"
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
	keys       services.KeyService
	persist    ports.UserPersist
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

func (s *PayloadServer) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResult, error) {
	// verify
	userId, ok, err := s.keys.VerifyPassphrase(req.GetPassphrase())
	if err != nil || userId < 1 {
		glog.Error(err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	if !ok {
		glog.Errorf("unable to match passphrase %v", req.GetPassphrase())
		return nil, grpc.Errorf(codes.Unauthenticated, "unable to verify passphrase")
	}
	// generate keys
	kSet, err := s.keys.GenKeyPairForUser(userId, []byte(req.GetKeys().GetEncryptionKey()), []byte(req.GetKeys().GetSignKey()))
	if err != nil {
		glog.Errorf("error generating key pair")
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	// refresh token
	tokenSig, err := s.keys.MakeRefreshToken(userId, req.GetKeys().GetEncryptionKey(), kSet.PrivSign)
	if err != nil {
		glog.Errorf("unable to make refresh token: %v", err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}

	return &pb.AuthResult{
		UserId: userId,
		Keys: &pb.KeySet{
			EncryptionKey: kSet.PubEncrKey,
			SignKey:       kSet.PubSignKey,
		},
		TokenSet: &pb.TokenSet{
			Token:     tokenSig.Token,
			Signature: tokenSig.Signature,
		},
	}, nil
}

func (s *PayloadServer) GetAccess(ctx context.Context, req *pb.RefreshRequest) (*pb.Access, error) {
	userId := ctx.Value("user-id").(int64)
	ts := services.TokenSig{
		Token:     req.GetTokenSet().GetToken(),
		Signature: req.TokenSet.GetSignature(),
	}
	if err := s.keys.VerifyRefreshToken(userId, ts); err != nil {
		return nil, err
	}
	tokenSig, plainTxtTkn, err := s.keys.MakeAccessToken(userId)
	if err != nil {
		return nil, err
	}

	keys, err := s.persist.GetKeys(userId)
	s.tokenCache.add(userId, keys.UserSignKey, keys.PrivSignKey, plainTxtTkn)
	if err != nil {
		return nil, err
	}
	return &pb.Access{
		TokenSet: &pb.TokenSet{
			Token:     tokenSig.Token,
			Signature: tokenSig.Signature,
		},
	}, nil
}

func (s *PayloadServer) ReceivePayload(stream pb.Payload_ReceivePayloadServer) error {
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

func (s *PayloadServer) SendPayload(ctx context.Context, req *pb.DownloadRequest, sendServer pb.Payload_SendPayloadServer) error {
	if !s.tokenCache.tokenIsValid(ctx) {
		return grpc.Errorf(codes.Unauthenticated, "invalid access token")
	}
	return nil
}
