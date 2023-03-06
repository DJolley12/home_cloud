package encryption

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

type TokenSig struct {
	Token     []byte
	Signature []byte
}

func decrypt(privK, data []byte) ([]byte, error) {
	identity, err := age.ParseX25519Identity(string(privK))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)

	r, err := age.Decrypt(buf, identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func DecryptAndVerify(data, cryptoKey, sig, sigKey []byte) ([]byte, error) {
	t, err := decrypt(cryptoKey, data)
	if err != nil {
		return nil, err
	}
	token := t
	if ok := ed25519.Verify(sigKey, token, sig); !ok {
		return nil, fmt.Errorf("token does not match")
	}
	return token, nil
}

func Encrypt(pubK, data []byte) ([]byte, error) {
	rec, err := age.ParseX25519Recipient(string(pubK))
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	dBuff := bytes.NewBuffer(data)

	fmt.Fprintf(os.Stderr, "in enc %v\n", string(data))
	w, err := age.Encrypt(out, rec)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	if _, err := io.Copy(w, dBuff); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func EncryptAndSign(data []byte, encryptionKey, signatureKey []byte) (*TokenSig, error) {
	out := &bytes.Buffer{}
	sig := ed25519.Sign(signatureKey, data)
	signer := bytes.NewBuffer(data)

	rec, err := age.ParseX25519Recipient(string(encryptionKey))
	w, err := age.Encrypt(out, rec)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, signer); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return &TokenSig{
		Token:     out.Bytes(),
		Signature: sig,
	}, nil
}
