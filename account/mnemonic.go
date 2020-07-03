package account

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pborman/uuid"
)

func GetPrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		return nil, err
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func PrivateKeyToKeystore(privateKey *ecdsa.PrivateKey, auth string) ([]byte, error) {
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
	keyJson, err := keystore.EncryptKey(key, auth, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	return keyJson, nil
}
