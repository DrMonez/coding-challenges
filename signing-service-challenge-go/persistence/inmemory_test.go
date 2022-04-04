package persistence

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"testing"
)

var storage = LocalStorage{
	UserDevices: make(map[string]map[string]struct{}),
	Devices:     make(map[string]*domain.Device),
	Signatures:  make(map[string]map[int]*domain.Signature),
}

func TestLocalStorage_CreateSignatureDevice(t *testing.T) {
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

func TestLocalStorage_GetDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := storage.GetDevice(deviceId)
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, "label")
}

func TestLocalStorage_UpdateSignatureCounter(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	storage.UpdateSignatureCounter(deviceId)
	actualDevice := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, actualDevice, nil)
	helpers.ShouldBe(t, actualDevice.Id, deviceId)
	helpers.ShouldBe(t, actualDevice.SignatureCounter, 1)
}

func TestLocalStorage_GetDeviceSignaturesCount(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	defaultCount := storage.GetDeviceSignaturesCount(deviceId)
	helpers.ShouldBe(t, defaultCount, 0)
	storage.UpdateSignatureCounter(deviceId)
	actualCount := storage.GetDeviceSignaturesCount(deviceId)
	helpers.ShouldBe(t, actualCount, 1)
}

func TestLocalStorage_AddSignature(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	storage.AddSignature(deviceId, make([]byte, 10), make([]byte, 10), make([]byte, 10))
	signatures := storage.Signatures[deviceId]
	helpers.ShouldBe(t, len(signatures), 1)
}
