package persistence

import (
	"fmt"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"sync"
)

const DEFAULT_LABEL = "Signing transaction..."

type Storage interface {
	CreateSignatureDevice(
		userId string, algorithm domain.CryptoAlgorithmType, label string,
	) (DeviceId string, Label string)

	GetDevice(deviceId string) *domain.Device
	UpdateSignatureCounter(deviceId string) error
	AddSignature(deviceId string, publicKey []byte, privateKey []byte, signedData []byte) error
	GetDeviceSignaturesCount(deviceId string) int
	GetLastDeviceSignature(deviceId string) (*domain.Signature, error)
}

type LocalStorage struct {
	UserDevicesMutex sync.Mutex
	UserDevices      map[string]map[string]struct{}
	DevicesMutex     sync.Mutex
	Devices          map[string]*domain.Device
	SignaturesMutex  sync.Mutex
	Signatures       map[string]map[int]*domain.Signature
}

func (s *LocalStorage) CreateSignatureDevice(
	userId string,
	algorithm domain.CryptoAlgorithmType,
	label string,
) (DeviceId string, Label string) {
	s.UserDevicesMutex.Lock()
	userDevices := s.UserDevices[userId]
	if userDevices == nil {
		userDevices = make(map[string]struct{})
	}
	s.UserDevices[userId] = userDevices
	s.UserDevicesMutex.Unlock()

	deviceId := uuid.New().String()

	actualLabel := DEFAULT_LABEL
	if label != "" {
		actualLabel = label
	}
	device := domain.Device{
		Id:               deviceId,
		Algorithm:        algorithm,
		Label:            actualLabel,
		SignatureCounter: 0,
	}
	s.DevicesMutex.Lock()
	s.Devices[deviceId] = &device
	s.DevicesMutex.Unlock()

	s.SignaturesMutex.Lock()
	s.Signatures[deviceId] = make(map[int]*domain.Signature)
	s.SignaturesMutex.Unlock()

	return deviceId, actualLabel
}

func (s *LocalStorage) GetDevice(deviceId string) *domain.Device {
	s.DevicesMutex.Lock()
	device := s.Devices[deviceId]
	s.DevicesMutex.Unlock()
	return device
}

func (s *LocalStorage) UpdateSignatureCounter(deviceId string) error {
	s.DevicesMutex.Lock()
	if s.Devices[deviceId] == nil {
		return fmt.Errorf("Device with Id=\"%s\" does not exist", deviceId)
	}
	s.Devices[deviceId].SignatureCounter++
	s.DevicesMutex.Unlock()
	return nil
}

func (s *LocalStorage) GetDeviceSignaturesCount(deviceId string) int {
	s.DevicesMutex.Lock()
	deviceSignaturesCount := s.Devices[deviceId].SignatureCounter
	s.DevicesMutex.Unlock()
	return deviceSignaturesCount
}

func (s *LocalStorage) AddSignature(deviceId string, publicKey []byte, privateKey []byte, signedData []byte) error {
	s.SignaturesMutex.Lock()
	deviceSignatures := s.Signatures[deviceId]
	if deviceSignatures == nil {
		deviceSignatures = make(map[int]*domain.Signature)
	}
	signatureCount := s.GetDeviceSignaturesCount(deviceId)
	signature := domain.Signature{
		Id:         signatureCount,
		SignedData: signedData,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	deviceSignatures[signature.Id] = &signature
	s.Signatures[deviceId] = deviceSignatures
	s.SignaturesMutex.Unlock()
	err := s.UpdateSignatureCounter(deviceId)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) GetLastDeviceSignature(deviceId string) (*domain.Signature, error) {
	s.SignaturesMutex.Lock()
	deviceSignatures := s.Signatures[deviceId]
	if deviceSignatures == nil {
		return nil, fmt.Errorf("Signatures of device with Id=\"%s\" do not exist", deviceId)
	}
	signaturesCount := s.GetDeviceSignaturesCount(deviceId)
	lastSignature := deviceSignatures[signaturesCount-1]
	if lastSignature == nil {
		return nil, fmt.Errorf("Signatures of device with Id=\"%s\" do not exist", deviceId)
	}
	s.SignaturesMutex.Unlock()
	return lastSignature, nil
}
