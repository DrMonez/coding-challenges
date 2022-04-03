package domain

type Device struct {
	Id               string
	Algorithm        CryptoAlgorithmType
	Label            string
	SignatureCounter uint64
}
