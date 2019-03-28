package sdk

import "github.com/proximax-storage/nem2-crypto-go"

func NewAccountFactory() AccountFactory {
	return &accountFactoryImpl{}
}

type AccountFactory interface {
	NewAccount(networkType NetworkType) (*Account, error)
	NewAccountFromPrivateKey(pKey string, networkType NetworkType) (*Account, error)
	NewAccountFromPublicKey(pKey string, networkType NetworkType) (*PublicAccount, error)
}

type accountFactoryImpl struct{}

func (f *accountFactoryImpl) NewAccount(networkType NetworkType) (*Account, error) {
	kp, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.DefaultEngine)
	if err != nil {
		return nil, err
	}

	pa, err := f.NewAccountFromPublicKey(kp.PublicKey.String(), networkType)
	if err != nil {
		return nil, err
	}

	return &Account{pa, kp}, nil
}

func (f *accountFactoryImpl) NewAccountFromPrivateKey(pKey string, networkType NetworkType) (*Account, error) {
	k, err := crypto.NewPrivateKeyfromHexString(pKey)
	if err != nil {
		return nil, err
	}

	kp, err := crypto.NewKeyPair(k, nil, nil)
	if err != nil {
		return nil, err
	}

	pa, err := f.NewAccountFromPublicKey(kp.PublicKey.String(), networkType)
	if err != nil {
		return nil, err
	}

	return &Account{pa, kp}, nil
}

func (f *accountFactoryImpl) NewAccountFromPublicKey(pKey string, networkType NetworkType) (*PublicAccount, error) {
	ad, err := NewAddressFromPublicKey(pKey, networkType)
	if err != nil {
		return nil, err
	}
	return &PublicAccount{ad, pKey}, nil
}
