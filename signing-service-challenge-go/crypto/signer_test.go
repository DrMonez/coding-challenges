package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"testing"
)

type Opts struct{}

func (o Opts) HashFunc() crypto.Hash {
	return 0
}

func TestSignInfo(t *testing.T) {
	rsaGenerator := RSAGenerator{}
	keyPair, _ := rsaGenerator.Generate()
	dataToBeSigned := []byte("some data")
	h := sha256.New()
	h.Write(dataToBeSigned)
	digest := h.Sum(nil)
	signedData, _ := keyPair.Private.Sign(rand.Reader, digest, crypto.SHA256)
	err := rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, digest, signedData)
	helpers.ShouldBe(t, err, nil)
}
