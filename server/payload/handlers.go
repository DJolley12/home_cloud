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
	"github.com/DJolley12/home_cloud/server/persist/adapters"
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
	keyService services.KeyService
	persist    ports.UserPersist
}

func NewPaylodServer(basePath string, dbConfig ports.DBConfig) (*PayloadServer, error) {
	if _, err := os.Stat(basePath); err != nil {
		return nil, fmt.Errorf("path error: %v", err)
	}

	c := newTokenCache(100, 5)
	p := newPassCache(100, 5)
	ks, err := services.NewKeyService(adapters.NewUserPersist(dbConfig))
	if err != nil {
		return nil, err
	}
	server := PayloadServer{
		basePath:   basePath,
		tokenCache: c,
		passCache:  p,
		keyService: *ks,
	}
	go func() {
		server.tokenCache.collectTokens()
		server.passCache.collectPass()
		time.Sleep(15 * time.Second)
	}()

	return &server, nil
}

func (s *PayloadServer) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResult, error) {
	// verify
	userId, ok, err := s.keyService.VerifyPassphrase(req.GetPassphrase())
	if err != nil || userId < 1 {
		glog.Error(err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	if !ok {
		glog.Errorf("unable to match passphrase %v", req.GetPassphrase())
		return nil, grpc.Errorf(codes.Unauthenticated, "unable to verify passphrase")
	}
	// generate keys
	kSet, err := s.keyService.GenKeyPairForUser(userId, []byte(req.GetKeys().GetEncryptionKey()), []byte(req.GetKeys().GetSignKey()))
	if err != nil {
		glog.Errorf("error generating key pair")
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}
	// refresh token
	tokenSig, err := s.keyService.MakeRefreshToken(userId, req.GetKeys().GetEncryptionKey(), kSet.PrivSign)
	if err != nil {
		glog.Errorf("unable to make refresh token: %v", err)
		return nil, grpc.Errorf(codes.Internal, "internal error, unable to complete authorization")
	}

	return &pb.AuthResult{
		UserId: userId,
		Keys: &pb.KeySet{
			EncryptionKey: kSet.Recipient,
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
	keys, err := s.persist.GetKeys(userId)
	if err != nil {
		return nil, err
	}
	if err := s.keyService.VerifyRefreshToken(userId, ts, keys.UserSignKey); err != nil {
		return nil, err
	}
	tokenSig, plainTxtTkn, err := s.keyService.MakeAccessToken(userId)
	if err != nil {
		return nil, err
	}

	s.tokenCache.add(userId, keys.UserSignKey, keys.PrivSignKey, plainTxtTkn, keys.UserSignKey)

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
	size := 0
	n, err := writeFile(f, in.GetChunk())
	if err != nil {
		glog.Error(err)
		return err
	}
	size += n

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			glog.Error(err)
			return err
		}
		n, err := writeFile(f, in.GetChunk())
		if err != nil {
			glog.Error(err)
			return err
		}

		size += n
	}

	return stream.SendAndClose(&pb.UploadResult{
		RecvSize:  int32(size),
		IsSuccess: true,
	})
}

func (s *PayloadServer) SendPayload(req *pb.DownloadRequest, sendServer pb.Payload_SendPayloadServer) error {
	ctx := sendServer.Context()
	if !s.tokenCache.tokenIsValid(ctx) {
		return grpc.Errorf(codes.Unauthenticated, "invalid access token")
	}
	return grpc.Errorf(codes.Unimplemented, "unimplemented")
}
