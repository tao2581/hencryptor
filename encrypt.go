package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"

	paillier "github.com/tao2581/go-go-gadget-paillier"

	"fyne.io/fyne"
)

type PurePrivKey struct {
	p         *big.Int
	pp        *big.Int
	pminusone *big.Int
	q         *big.Int
	qq        *big.Int
	qminusone *big.Int
	pinvq     *big.Int
	hp        *big.Int
	hq        *big.Int
	n         *big.Int
}

func Key2str(privKey *paillier.PrivateKey) string {
	asByteSlice := paillier.Marshal(privKey)
	return base64.RawStdEncoding.EncodeToString(asByteSlice)
}

func Pubkey2str(PubKey *paillier.PublicKey) string {
	bytes, _ := json.Marshal(PubKey)
	return base64.RawStdEncoding.EncodeToString(bytes)
}

func Str2key(keyStr string) (privKey *paillier.PrivateKey, err error) {
	decodeKey, err := base64.RawStdEncoding.DecodeString(keyStr)
	return paillier.UnMarshal(decodeKey), nil
}

func NewKey(app fyne.App) *paillier.PrivateKey {
	// Generate a 128-bit private key.
	privKey, _ := paillier.GenerateKey(rand.Reader, 32)
	keyStr := Key2str(privKey)
	// Store in preference
	app.Preferences().SetString("privKey", keyStr)
	return privKey
}

func LoadKey(app fyne.App) *paillier.PrivateKey {
	savedKey := app.Preferences().String("privKey")
	if savedKey == "nil" {
		privKey := NewKey(app)
		return privKey
	}
	privKey, err := Str2key(savedKey)
	if err != nil {
		return nil
	}
	return privKey
}

func RestoreKey(keyStr string, app fyne.App) (*paillier.PrivateKey, error) {
	privKey, err := Str2key(keyStr)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
