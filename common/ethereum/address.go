package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Address struct {
	Address    string
	PrivateKey string
	PublicKey  string
}

func IsValidateAddressFail(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return !re.MatchString(addr)
}

func GenerateAddress() (Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return Address{}, fmt.Errorf("crypto.GenerateKey error: %s", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyBytesStr := hexutil.Encode(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Address{}, fmt.Errorf("publicKey aes encrypt error: %s", err)
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyBytesStr := hexutil.Encode(publicKeyBytes)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return Address{
		Address:    address,
		PrivateKey: privateKeyBytesStr,
		PublicKey:  publicKeyBytesStr,
	}, nil
}

func HexToECDSA(privateKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(privateKey[2:])
}

func HexECDSAPub(publicKey string) (*ecdsa.PublicKey, error) {
	publicKeyByte, err := hexutil.Decode(publicKey)
	if err != nil {
		return nil, fmt.Errorf("hexutil.Decode do publicKey error: %s", err)
	}

	publicKeyECDSA, err := crypto.UnmarshalPubkey(publicKeyByte)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalPubkey error: %s", err)
	}

	return publicKeyECDSA, nil
}
