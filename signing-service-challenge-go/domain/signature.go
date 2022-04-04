package domain

type Signature struct {
	Id         string
	DeviceId   string
	Number     int
	PrivateKey []byte
	PublicKey  []byte
}
