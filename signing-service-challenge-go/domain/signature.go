package domain

import "crypto"

type Signature struct {
	Id         string
	DeviceId   string
	Number     uint64
	PrivateKey crypto.PrivateKey
}
