package reciever

type Reciever interface {
	PushCert(name string, cert []byte) error
}
