package payload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "github.com/DJolley12/home_cloud/protos"
	"github.com/DJolley12/home_cloud/shared/encryption"
	"github.com/DJolley12/home_cloud/shared/utils"
	"github.com/golang/glog"
)

type PayloadClient struct {
	client pb.PayloadClient

	chunkSize    int
	keyService   KeyService
	passphrase   string
	userKeySet   UserKeySet
	serverKeySet ServerKeySet
	refreshToken Token
	accessToken  Token
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

func (c *PayloadClient) Authorize(passphrase string) error {
	// generate keys
	ks, err := c.keyService.GenKeyPair()
	if err != nil {
		return err
	}
	c.userKeySet = *ks
	// send keys and passphrase to server
	res, err := c.client.Authorize(context.Background(), &pb.AuthRequest{
		Keys: &pb.KeySet{
			EncryptionKey: c.userKeySet.Recipient,
			SignKey:       c.userKeySet.PubSignKey,
		},
		Passphrase: c.passphrase,
	})
	if err != nil {
		return err
	}

	ts, k := res.GetTokenSet(), res.GetKeys()
	// unencrypt token, verify
	c.refreshToken = Token{}
	c.refreshToken.Expiry = ts.Expiry.AsTime()
	c.refreshToken.Token, err = encryption.DecryptAndVerify(ts.GetToken(),
		[]byte(c.userKeySet.Identity), ts.GetSignature(), k.GetSignKey())

	if err != nil {
		return err
	}
	// save server keys
	c.serverKeySet = ServerKeySet{
		UserId:     res.GetUserId(),
		PubSignKey: res.GetKeys().GetSignKey(),
		Recipient:  res.GetKeys().GetEncryptionKey(),
	}
	err = c.keyService.SaveServerKeys(
		c.serverKeySet,
	)
	if err != nil {
		return err
	}
	// save refresh token
	return c.keyService.SaveToken(c.refreshToken)
}

func (c *PayloadClient) GetAccess() error {
	// TODO check expiry time
	// get keys and refresh token
	// encrypt and sign refresh token
	ts, err := encryption.EncryptAndSign(c.refreshToken.Token, []byte(c.serverKeySet.Recipient), []byte(c.serverKeySet.PubSignKey))
	if err != nil {
		return err
	}
	// send token, user id to server

	req := &pb.RefreshRequest{
		TokenSet: &pb.TokenSet{
			Token:     ts.Token,
			Signature: ts.Signature,
			Expiry:    utils.ToTimeStamppb(c.refreshToken.Expiry),
		},
	}
	ctx := context.WithValue(context.Background(), "user-id", c.serverKeySet.UserId)
	res, err := c.client.GetAccess(ctx, req)
	if err != nil {
		return err
	}

	rts := res.GetTokenSet()
	at, err := encryption.DecryptAndVerify(rts.GetToken(),
		[]byte(c.userKeySet.Identity), rts.GetSignature(), c.serverKeySet.PubSignKey)

	c.accessToken = Token{
		Token:  at,
		Expiry: rts.GetExpiry().AsTime(),
	}

	return nil
}

func (c *PayloadClient) UploadFile(filePath string) error {
	if c.accessToken.Expiry.Before(time.Now()) {
		return fmt.Errorf("access token expired at %v - please reauth and continue", c.accessToken.Expiry)
	}

	at, err := encryption.EncryptAndSign(c.accessToken.Token, []byte(c.serverKeySet.Recipient), c.userKeySet.PrivSignKey)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), "access-token", at)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file %v - err: %v\n", filePath, err)
	}
	defer file.Close()
	p, err := c.client.ReceivePayload(ctx)
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
