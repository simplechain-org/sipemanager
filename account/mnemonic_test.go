package account

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Error(errors.New("not ecdsa.PublicKey"))
		return
	}
	//生成地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	t.Log("address from privateKey:", address)
}

func TestGetPrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("mnemonic=", mnemonic)

	privateKey, err := GetPrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		t.Error(err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Error(errors.New("not ecdsa.PublicKey"))
		return
	}
	//生成地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	t.Log("address from privateKey:", address)
}

func TestPrivateKeyToKeystore(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
		return
	}
	auth := "123456"
	kjson, err := PrivateKeyToKeystore(privateKey, auth)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("keystore file content=", string(kjson))
}
