package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/persistence"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	Device       *domain.Device
	Storage      persistence.Storage
	RsaGenerator RSAGenerator
	RsaMarshaler RSAMarshaler
}

func (s RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.RsaGenerator.Generate()
	if err != nil {
		return nil, err
	}
	digits := GetSha256Hash(dataToBeSigned)

	signedData, err := keyPair.Private.Sign(rand.Reader, digits, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	marshaledPublicKey, marshaledPrivateKey, err := s.RsaMarshaler.Marshal(*keyPair)
	if err != nil {
		return nil, err
	}
	err = s.Storage.AddSignature(s.Device.Id, marshaledPublicKey, marshaledPrivateKey, signedData)
	if err != nil {
		return nil, err
	}
	return signedData, nil
}

type ECCSigner struct {
	Device       *domain.Device
	Storage      persistence.Storage
	EccGenerator ECCGenerator
	EccMarshaler ECCMarshaler
}

func (s ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.EccGenerator.Generate()
	if err != nil {
		return nil, err
	}
	digits := GetSha256Hash(dataToBeSigned)

	signedData, err := keyPair.Private.Sign(rand.Reader, digits, crypto.SHA256)
	if err != nil {
		return nil, err
	}

	marshaledPublicKey, marshaledPrivateKey, err := s.EccMarshaler.Encode(*keyPair)
	if err != nil {
		return nil, err
	}
	err = s.Storage.AddSignature(s.Device.Id, marshaledPublicKey, marshaledPrivateKey, signedData)
	if err != nil {
		return nil, err
	}
	return signedData, nil
}

func GetSha256Hash(dataToBeSigned []byte) []byte {
	hash := sha256.New()
	hash.Write(dataToBeSigned)
	return hash.Sum(nil)
}
