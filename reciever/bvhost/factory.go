package bvhost

import "certbot/reciever"

func Init() {
	reciever.RegisterFactory(&factory{})
}

type factory struct{}

func (f *factory) Name() string {
	return "bvhost"
}

func (f *factory) NewReciever(cfg reciever.Unmarshaler) (reciever.Reciever, error) {
	r := new(Reciever)
	err := cfg.Unmarshal(&r.cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}
