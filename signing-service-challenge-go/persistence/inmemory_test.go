package persistence

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/helpers"
	"testing"
)

var storage = LocalStorage{
	UserDevices: make(map[string]map[string]struct{}),
	Devices:     make(map[string]*domain.Device),
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
	helpers.ShouldBe(t, len(device.Signatures), 0)
}

func TestGetDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := storage.GetDevice(deviceId)
	helpers.ShouldNotBe(t, device, nil)
	helpers.ShouldBe(t, device.Id, deviceId)
	helpers.ShouldBe(t, device.Algorithm.String(), "RSA")
	helpers.ShouldBe(t, device.Label, "label")
	helpers.ShouldBe(t, len(device.Signatures), 0)
}

func TestUpdateDevice(t *testing.T) {
	deviceId, _ := storage.CreateSignatureDevice("test", "RSA", "label")
	device := domain.Device{
		Id:         deviceId,
		Algorithm:  "ECC",
		Label:      "new label",
		Signatures: []domain.Signature{{DeviceId: deviceId, Id: "1", Number: 1, PrivateKey: nil}},
	}
	storage.UpdateDevice(&device)
	actualDevice := storage.Devices[deviceId]
	helpers.ShouldNotBe(t, actualDevice, nil)
	helpers.ShouldBe(t, actualDevice.Id, deviceId)
	helpers.ShouldBe(t, actualDevice.Algorithm.String(), "ECC")
	helpers.ShouldBe(t, actualDevice.Label, "new label")
	helpers.ShouldBe(t, len(actualDevice.Signatures), 1)
}
