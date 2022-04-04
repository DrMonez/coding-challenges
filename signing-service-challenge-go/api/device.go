package api

import (
	"encoding/json"
	"fmt"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/crypto"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"net/http"
)

type CreateSignatureDeviceResponse struct {
	DeviceId string `json:"device_id"`
	Label    string `json:"label"`
}

type CreateSignatureDeviceRequest struct {
	Id        string                     `json:"id"`
	Algorithm domain.CryptoAlgorithmType `json:"algorithm"`
	Label     string                     `json:"label"`
}

type SignTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type SignTransactionRequest struct {
	DeviceId string `json:"device_id"`
	Data     string `json:"data"`
}

func PostMethodTemplate[T any](request *http.Request, body *T) (isValid bool, errors []string) {
	if request.Method != http.MethodPost {
		return false, []string{http.StatusText(http.StatusMethodNotAllowed)}
	}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return false, []string{http.StatusText(http.StatusBadRequest)}
	}
	return true, nil
}

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var body CreateSignatureDeviceRequest
	isValidRequest, errors := PostMethodTemplate(request, &body)
	if !isValidRequest {
		WriteErrorResponse(response, http.StatusBadRequest, errors)
		return
	}
	deviceId, label := s.storage.CreateSignatureDevice(body.Id, body.Algorithm, body.Label)
	createSignatureDeviceResponse := CreateSignatureDeviceResponse{
		DeviceId: deviceId,
		Label:    label,
	}

	WriteAPIResponse(response, http.StatusOK, createSignatureDeviceResponse)
}

func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	var body SignTransactionRequest
	isValidRequest, errors := PostMethodTemplate(request, &body)
	if !isValidRequest {
		WriteErrorResponse(response, http.StatusBadRequest, errors)
		return
	}

	device := s.storage.GetDevice(body.DeviceId)
	if device == nil {
		WriteInternalError(response)
		return
	}

	lastSignature, err := s.storage.GetLastDeviceSignature(body.DeviceId)
	if err != nil && device.SignatureCounter != 0 {
		WriteInternalError(response)
		return
	}

	var lastSignedData = []byte(body.DeviceId)
	if lastSignature != nil && lastSignature.SignedData != nil {
		lastSignedData = lastSignature.SignedData
	}

	signer, err := s.GetSigner(device)
	if err != nil {
		WriteInternalError(response)
		return
	}

	signedData, err := signer.Sign([]byte(body.Data))
	if err != nil {
		WriteInternalError(response)
		return
	}
	signatureCounter := s.storage.GetDeviceSignaturesCount(body.DeviceId)

	signTransactionResponse := SignTransactionResponse{
		Signature:  string(signedData),
		SignedData: fmt.Sprintf("%d_%s_%s", signatureCounter, body.Data, string(lastSignedData)),
	}

	WriteAPIResponse(response, http.StatusOK, signTransactionResponse)
}

func (s *Server) GetSigner(device *domain.Device) (crypto.Signer, error) {
	switch device.Algorithm {
	case domain.RSA:
		return &crypto.RSASigner{
			Storage:      s.storage,
			RsaMarshaler: crypto.NewRSAMarshaler(),
			RsaGenerator: crypto.RSAGenerator{},
			Device:       device,
		}, nil
	case domain.ECC:
		return crypto.ECCSigner{
			Storage:      s.storage,
			EccMarshaler: crypto.NewECCMarshaler(),
			EccGenerator: crypto.ECCGenerator{},
			Device:       device,
		}, nil
	default:
		return nil, fmt.Errorf("algorithm is not implemented")
	}
}
