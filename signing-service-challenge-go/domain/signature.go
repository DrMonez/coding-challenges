package domain

type Signature struct {
	Id         int
	SignedData []byte
	PrivateKey []byte
	PublicKey  []byte
}
