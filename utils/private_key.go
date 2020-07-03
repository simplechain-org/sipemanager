package utils

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pborman/uuid"
)

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

//根据私钥获取以太坊地址
func GetAddress(privateKey *ecdsa.PrivateKey) (string, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}
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
