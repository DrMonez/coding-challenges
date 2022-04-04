package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/persistence"
	"math/big"
	"testing"
)

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

var storage persistence.Storage = &persistence.LocalStorage{
	UserDevices: make(map[string]map[string]struct{}),
	Devices:     make(map[string]*domain.Device),
	Signatures:  make(map[string]map[int]*domain.Signature),
}

var rsaSigner = RSASigner{
	Storage:      &storage,
	RsaMarshaler: NewRSAMarshaler(),
	RsaGenerator: RSAGenerator{},
}

func TestRSASigner_Sign(t *testing.T) {
	deviceId, _ := (*rsaSigner.Storage).CreateSignatureDevice("test", "RSA", "")
	rsaSigner.Device = (*rsaSigner.Storage).GetDevice(deviceId)
	dataToBeSigned := []byte("some data")
	signedData, _ := rsaSigner.Sign(dataToBeSigned)
	lastSignature, _ := (*rsaSigner.Storage).GetLastDeviceSignature(rsaSigner.Device.Id)
	keyPair, _ := rsaSigner.RsaMarshaler.Unmarshal(lastSignature.PrivateKey)
	err := rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, GetHash(dataToBeSigned), signedData)
	helpers.ShouldBe(t, err, nil)
}

var eccSigner = ECCSigner{
	Storage:      &storage,
	EccMarshaler: NewECCMarshaler(),
	EccGenerator: ECCGenerator{},
}

func TestECCSigner_Sign(t *testing.T) {
	deviceId, _ := (*eccSigner.Storage).CreateSignatureDevice("test", "ECC", "")
	eccSigner.Device = (*eccSigner.Storage).GetDevice(deviceId)
	dataToBeSigned := []byte("some data")
	signedData, _ := eccSigner.Sign(dataToBeSigned)
	lastSignature, _ := (*eccSigner.Storage).GetLastDeviceSignature(eccSigner.Device.Id)
	keyPair, _ := eccSigner.EccMarshaler.Decode(lastSignature.PrivateKey)

	var esig struct {
		R, S *big.Int
	}
	asn1.Unmarshal(signedData, &esig)

	err := ecdsa.Verify(keyPair.Public, GetHash(dataToBeSigned), esig.R, esig.S)
	helpers.ShouldBe(t, err, nil)
}
