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
	"github.com/DJolley12/home_cloud/shared/encryption"
	"github.com/golang/glog"
)

type PayloadClient struct {
	ctx    context.Context
	client pb.PayloadClient

	chunkSize    int
	keyService   KeyService
	passphrase   string
	userKeySet   UserKeySet
	serverKeySet ServerKeySet
	token        Token
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
		ctx:       context.Background(),
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
	res, err := c.client.Authorize(c.ctx, &pb.AuthRequest{
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
	c.token = Token{}
	c.token.Expiry = ts.Expiry.AsTime()
	c.token.RefreshToken, err = encryption.DecryptAndVerify(ts.GetToken(),
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
	return c.keyService.SaveToken(c.token)
}

func (c *PayloadClient) GetAccess() error {
	// TODO check expiry time 
	// get keys and refresh token
	ts, err := encryption.EncryptAndSign(c.token.RefreshToken, []byte(c.serverKeySet.Recipient), []byte(c.serverKeySet.PubSignKey))
	if err != nil {
	  return err
	}
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
