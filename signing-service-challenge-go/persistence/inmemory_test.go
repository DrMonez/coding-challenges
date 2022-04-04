package persistence

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"testing"
)

var storage = LocalStorage{
	UserDevices: make(map[string]map[string]struct{}),
	Devices:     make(map[string]*domain.Device),
	Signatures:  make(map[string]map[string]*domain.Signature),
}

func TestCreateSignatureDevice(t *testing.T) {
	deviceId, label := storage.CreateSignatureDevice("test", "RSA", "")
	helpers.ShouldNotBe(t, label, "")
	helpers.ShouldNotBe(t, deviceId, "")
	userDevices := storage.UserDevices["test"]
	helpers.ShouldNotBe(t, userDevices[deviceId], nil)
	device := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, label)
}

func TestGetDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := storage.GetDevice(deviceId)
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, "label")
}

func TestUpdateDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := domain.Device{
		Id:        deviceId,
		Algorithm: "ECC",
		Label:     "new label",
	}
	storage.UpdateDevice(&device)
	actualDevice := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, actualDevice, nil)
	helpers.ShouldBe(t, actualDevice.Id, deviceId)
	helpers.ShouldBe(t, actualDevice.Algorithm.String(), "ECC")
	helpers.ShouldBe(t, actualDevice.Label, "new label")
}

func TestAddSignature(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	storage.AddSignature(deviceId, make([]byte, 10), make([]byte, 10))
	signatures := storage.Signatures[deviceId]
	helpers.ShouldBe(t, len(signatures), 1)
}

func TestGetDeviceSignaturesCount(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	defaultLength := storage.GetDeviceSignaturesCount(deviceId)
	helpers.ShouldBe(t, defaultLength, 0)
	storage.AddSignature(deviceId, make([]byte, 10), make([]byte, 10))
	actualLength := storage.GetDeviceSignaturesCount(deviceId)
	helpers.ShouldBe(t, actualLength, 1)
}
