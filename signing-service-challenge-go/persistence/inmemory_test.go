package persistence

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/assert"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"testing"
)

var storage = LocalStorage{
	UserDevices: make(map[string]map[string]struct{}),
	Devices:     make(map[string]*domain.Device),
	Signatures:  make(map[string]map[int]*domain.Signature),
}

func TestLocalStorage_CreateSignatureDevice(t *testing.T) {
	deviceId, label := storage.CreateSignatureDevice("test", "RSA", "")
	assert.ShouldNotBe(t, label, "")
	assert.ShouldNotBe(t, deviceId, "")
	userDevices := storage.UserDevices["test"]
	assert.ShouldNotBe(t, userDevices[deviceId], nil)
	device := storage.Devices[deviceId]
	assert.ShouldNotBe(t, device, nil)
	assert.ShouldBe(t, device.Id, deviceId)
	assert.ShouldBe(t, device.Algorithm, domain.RSA)
	assert.ShouldBe(t, device.Label, label)
}

func TestLocalStorage_GetDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := storage.GetDevice(deviceId)
	assert.ShouldNotBe(t, device, nil)
	assert.ShouldBe(t, device.Id, deviceId)
	assert.ShouldBe(t, device.Algorithm, domain.RSA)
	assert.ShouldBe(t, device.Algorithm, domain.RSA)
	assert.ShouldBe(t, device.Label, "label")
}

func TestLocalStorage_UpdateSignatureCounter(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	storage.UpdateSignatureCounter(deviceId)
	actualDevice := storage.Devices[deviceId]
	assert.ShouldNotBe(t, actualDevice, nil)
	assert.ShouldBe(t, actualDevice.Id, deviceId)
	assert.ShouldBe(t, actualDevice.SignatureCounter, 1)
}

func TestLocalStorage_GetDeviceSignaturesCount(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	defaultCount := storage.GetDeviceSignaturesCount(deviceId)
	assert.ShouldBe(t, defaultCount, 0)
	storage.UpdateSignatureCounter(deviceId)
	actualCount := storage.GetDeviceSignaturesCount(deviceId)
	assert.ShouldBe(t, actualCount, 1)
}

func TestLocalStorage_AddSignature(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	storage.AddSignature(deviceId, make([]byte, 10), make([]byte, 10), make([]byte, 10))
	signatures := storage.Signatures[deviceId]
	assert.ShouldBe(t, len(signatures), 1)
}

func TestLocalStorage_GetLastDeviceSignature(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	signedData := []byte("some data")
	storage.AddSignature(deviceId, make([]byte, 10), make([]byte, 10), signedData)
	lastSignature, _ := storage.GetLastDeviceSignature(deviceId)
	assert.ShouldBe(t, string(lastSignature.SignedData), string(signedData))
}
