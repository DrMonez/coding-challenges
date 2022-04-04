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
	device       *domain.Device
	storage      *persistence.Storage
	rsaGenerator RSAGenerator
	rsaMarshaler RSAMarshaler
}

func (s RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.rsaGenerator.Generate()
	if err != nil {
		return nil, err
	}
	digits := GetHash(dataToBeSigned)

	signedData, err := keyPair.Private.Sign(rand.Reader, digits, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	marshaledPublicKey, marshaledPrivateKey, err := s.rsaMarshaler.Marshal(*keyPair)
	if err != nil {
		return nil, err
	}
	err = (*s.storage).AddSignature(s.device.Id, marshaledPublicKey, marshaledPrivateKey, signedData)
	if err != nil {
		return nil, err
	}
	return signedData, nil
}

type ECCSigner struct {
	device       *domain.Device
	storage      *persistence.Storage
	eccGenerator ECCGenerator
	eccMarshaler ECCMarshaler
}

func (s ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.eccGenerator.Generate()
	if err != nil {
		return nil, err
	}
	digits := GetHash(dataToBeSigned)

	signedData, err := keyPair.Private.Sign(rand.Reader, digits, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	marshaledPublicKey, marshaledPrivateKey, err := s.eccMarshaler.Encode(*keyPair)
	if err != nil {
		return nil, err
	}
	err = (*s.storage).AddSignature(s.device.Id, marshaledPublicKey, marshaledPrivateKey, signedData)
	if err != nil {
		return nil, err
	}
	return signedData, nil
}

func GetHash(dataToBeSigned []byte) []byte {
	hash := sha256.New()
	hash.Write(dataToBeSigned)
	return hash.Sum(nil)
}
