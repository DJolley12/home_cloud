package services

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"log"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestKeyEncrypt(t *testing.T) {
	// t.Run("test encryptAndSign, decryptAndVerify",
	// 	func(t *testing.T) {
	// 		ident, err := age.GenerateX25519Identity()
	// 		assert.Nil(t, err)
	//
	// 		pubSign, privSign, err := ed25519.GenerateKey(nil)
	//
	// 		tokenSig, err := encryptAndSign([]byte("testicles, 1, 2"), []byte(ident.Recipient().String()), privSign)
	// 		assert.Nil(t, err)
	//
	// 		b, err := DecryptAndVerify([]byte(ident.String()), tokenSig.Token, tokenSig.Signature, pubSign)
	// 		assert.Nil(t, err)
	//
	// 		assert.Equal(t, "testing, 1, 2", string(b))
	// 	})

	t.Run("works", func(t *testing.T) {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			log.Fatalf("Failed to generate key pair: %v", err)
		}
		recipient := identity.Recipient()

		out := &bytes.Buffer{}

		pub, priv, _ := ed25519.GenerateKey(nil)
		str := []byte("test test testing")
		sig := ed25519.Sign(priv, str)
		signer := bytes.NewBuffer(str)

		rec, err := age.ParseX25519Recipient(recipient.String())
		w, err := age.Encrypt(out, rec)
		// w, err := age.Encrypt(out, recipient)
		if err != nil {
			log.Fatalf("Failed to create encrypted file: %v", err)
		}
		if _, err := io.Copy(w, signer); err != nil {
			log.Fatalf("Failed to write to encrypted file: %v", err)
		}
		if err := w.Close(); err != nil {
			log.Fatalf("Failed to close encrypted file: %v", err)
		}

		fmt.Printf("Encrypted file size: %d\n", out.Len())

		f := bytes.NewBuffer(out.Bytes())
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}

		r, err := age.Decrypt(f, identity)
		if err != nil {
			log.Fatalf("Failed to open encrypted file: %v", err)
		}
		out2 := &bytes.Buffer{}
		if _, err := io.Copy(out2, r); err != nil {
			log.Fatalf("Failed to read encrypted file: %v", err)
		}

		c := out2.Bytes()
		fmt.Printf("File contents: %q\n", c)
		if ok := ed25519.Verify(pub, c, sig); ok {
			println("works")
		} else {
			println("it no work")
		}
		assert.Equal(t, "test test testing", string(c))
		fmt.Printf("c %#v", string(c))

	})
	t.Run("partial test", func(t *testing.T) {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			log.Fatalf("Failed to generate key pair: %v", err)
		}
		recipient := identity.Recipient()

		pub, priv, _ := ed25519.GenerateKey(nil)
		str := []byte("test test testing")

		sig := ed25519.Sign(priv, str)
		rec, err := age.ParseX25519Recipient(string(recipient.String()))
		assert.Nil(t, err)

		out := &bytes.Buffer{}
		dBuff := bytes.NewBuffer(str)

		w, err := age.Encrypt(out, rec)
		assert.Nil(t, err)
		defer w.Close()
		if _, err := io.Copy(w, dBuff); err != nil {
			assert.Nil(t, err)
		}
		token := out.Bytes()
		assert.Nil(t, err)

		tk := &TokenSig{
			Token:     token,
			Signature: sig,
		}

		f := bytes.NewBuffer(tk.Token)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}

		r, err := age.Decrypt(f, identity)
		if err != nil {
			log.Fatalf("Failed to open encrypted file: %v", err)
		}
		out2 := &bytes.Buffer{}
		if _, err := io.Copy(out2, r); err != nil {
			log.Fatalf("Failed to read encrypted file: %v", err)
		}

		c := out2.Bytes()
		fmt.Printf("File contents: %q\n", c)
		if ok := ed25519.Verify(pub, c, tk.Signature); ok {
			println("works")
		} else {
			assert.Fail(t, "it no work")
		}
		assert.Equal(t, "test test testing", string(c))

	})
}
