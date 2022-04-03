package api

import (
	"encoding/json"
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

func PostMethodTemplate[T any](request *http.Request, body *T) (isValid bool, errors *[]string) {
	if request.Method != http.MethodPost {
		return false, &[]string{http.StatusText(http.StatusMethodNotAllowed)}
	}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return false, &[]string{http.StatusText(http.StatusBadRequest)}
	}
	return true, nil
}

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var body CreateSignatureDeviceRequest
	isValidRequest, errors := PostMethodTemplate(request, &body)
	if isValidRequest {
		WriteErrorResponse(response, http.StatusBadRequest, *errors)
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
	if isValidRequest {
		WriteErrorResponse(response, http.StatusBadRequest, *errors)
		return
	}

	signTransactionResponse := SignTransactionResponse{}

	WriteAPIResponse(response, http.StatusOK, signTransactionResponse)
}
