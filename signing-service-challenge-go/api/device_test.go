package api

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostMethodTemplateWithOptionalParameter(t *testing.T) {
	var body CreateSignatureDeviceRequest
	json := `{ "id": "123456", "algorithm":"RSA", "label":"some label" }`
	request := http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(json)),
	}
	isValid, errors := PostMethodTemplate(&request, &body)
	var expectedErrors *[]string = nil
	assert.ShouldBe(t, isValid, true)
	assert.ShouldBe(t, errors, expectedErrors)
	assert.ShouldBe(t, body, CreateSignatureDeviceRequest{Id: "123456", Algorithm: "RSA", Label: "some label"})
}

func TestPostMethodTemplateWithoutOptionalParameter(t *testing.T) {
	var body CreateSignatureDeviceRequest
	json := `{ "id": "123456", "algorithm":"RSA" }`
	request := http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(json)),
	}
	isValid, errors := PostMethodTemplate(&request, &body)
	var expectedErrors *[]string = nil
	assert.ShouldBe(t, isValid, true)
	assert.ShouldBe(t, errors, expectedErrors)
	assert.ShouldBe(t, body, CreateSignatureDeviceRequest{Id: "123456", Algorithm: "RSA", Label: ""})
}
