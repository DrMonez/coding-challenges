package persistence

import (
	"errors"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	uuid2 "github.com/google/uuid"
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
	UserDevices map[string]*map[string]bool
	Devices     map[string]*domain.Device
}

func (s LocalStorage) CreateSignatureDevice(
	userId string, algorithm domain.CryptoAlgorithmType, label string,
) (DeviceId string, Label string) {

	userDevices := s.UserDevices[userId]
	if userDevices == nil {
		userDevices = &map[string]bool{}
	}
	deviceId := uuid2.New().String()
	(*userDevices)[deviceId] = true
	s.UserDevices[userId] = userDevices

	var actualLabel string
	if label != "" {
		actualLabel = label
	} else {
		actualLabel = DEFAULT_LABEL
	}
	device := domain.Device{
		Id:               deviceId,
		Algorithm:        algorithm,
		Label:            actualLabel,
		SignatureCounter: 0,
	}
	s.Devices[deviceId] = &device

	return deviceId, actualLabel
}

func (s LocalStorage) GetDevice(deviceId string) *domain.Device {
	// In general, we should check that required device is active. In our case it's not necessary
	return s.Devices[deviceId]
}

func (s LocalStorage) UpdateDevice(device *domain.Device) error {
	if s.Devices[device.Id] == nil {
		return errors.New("Device with such Id is not exist")
	}
	s.Devices[device.Id] = device
	return nil
}
