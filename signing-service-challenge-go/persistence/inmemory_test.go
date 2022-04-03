package persistence

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"testing"
)

var storage = LocalStorage{
	UserDevices: map[string]*map[string]bool{},
	Devices:     map[string]*domain.Device{},
}

func TestCreateSignatureDevice(t *testing.T) {
	deviceId, label := storage.CreateSignatureDevice("test", "RSA", "")
	helpers.ShouldNotBe(t, label, "")
	helpers.ShouldNotBe(t, deviceId, "")
	userDevices := *(storage.UserDevices["test"])
	helpers.ShouldNotBe(t, userDevices[deviceId], nil)
	device := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, label)
	helpers.ShouldBe(t, device.SignatureCounter, uint64(0))
}

func TestGetDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := storage.GetDevice(deviceId)
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, "label")
	helpers.ShouldBe(t, device.SignatureCounter, uint64(0))
}

func TestUpdateDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := domain.Device{
		Id:               deviceId,
		Algorithm:        "ECC",
		Label:            "new label",
		SignatureCounter: 1,
	}
	storage.UpdateDevice(&device)
	actualDevice := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, actualDevice, nil)
	helpers.ShouldBe(t, actualDevice.Id, deviceId)
	helpers.ShouldBe(t, actualDevice.Algorithm.String(), "ECC")
	helpers.ShouldBe(t, actualDevice.Label, "new label")
	helpers.ShouldBe(t, actualDevice.SignatureCounter, uint64(1))
}
