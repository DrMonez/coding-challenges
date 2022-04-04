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
	UpdateDevice(device *domain.Device) error
}

type LocalStorage struct {
	UserDevicesMutex sync.Mutex
	UserDevices      map[string]map[string]struct{}
	DevicesMutex     sync.Mutex
	Devices          map[string]*domain.Device
	SignaturesMutex  sync.Mutex
	Signatures       map[string]map[string]*domain.Signature
}

func (s *LocalStorage) CreateSignatureDevice(
	userId string, algorithm domain.CryptoAlgorithmType, label string,
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
		Id:        deviceId,
		Algorithm: algorithm,
		Label:     actualLabel,
	}
	s.DevicesMutex.Lock()
	s.Devices[deviceId] = &device
	s.DevicesMutex.Unlock()

	s.SignaturesMutex.Lock()
	s.Signatures[deviceId] = make(map[string]*domain.Signature)
	s.SignaturesMutex.Unlock()

	return deviceId, actualLabel
}

func (s *LocalStorage) GetDevice(deviceId string) *domain.Device {
	s.DevicesMutex.Lock()
	device := s.Devices[deviceId]
	s.DevicesMutex.Unlock()
	return device
}

func (s *LocalStorage) UpdateDevice(device *domain.Device) error {
	s.DevicesMutex.Lock()
	if s.Devices[device.Id] == nil {
		return fmt.Errorf("Device with Id=\"%s\" does not exist", device.Id)
	}
	s.Devices[device.Id] = device
	s.DevicesMutex.Unlock()
	return nil
}

func (s *LocalStorage) AddSignature(deviceId string, publicKey []byte, privateKey []byte) error {
	newId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("internal error during Id creation")
	}

	s.SignaturesMutex.Lock()
	deviceSignatures := s.Signatures[deviceId]
	if deviceSignatures == nil {
		deviceSignatures = make(map[string]*domain.Signature)
	}
	signatureCount := len(deviceSignatures) + 1
	signature := domain.Signature{
		Id:         newId.String(),
		DeviceId:   deviceId,
		Number:     signatureCount,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	deviceSignatures[signature.Id] = &signature
	s.Signatures[deviceId] = deviceSignatures
	s.SignaturesMutex.Unlock()
	return nil
}

func (s *LocalStorage) GetDeviceSignaturesCount(deviceId string) int {
	s.SignaturesMutex.Lock()
	deviceSignaturesCount := len(s.Signatures[deviceId])
	s.SignaturesMutex.Unlock()
	return deviceSignaturesCount
}
