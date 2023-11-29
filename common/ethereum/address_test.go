package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHexToECDSA(t *testing.T) {
	addr, err := GenerateAddress()
	if err != nil {
		t.Logf("get new address error: %s", err)
		return
	}

	type args struct {
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				privateKey: addr.PrivateKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToECDSA(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "HexToECDSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := hexutil.Encode(crypto.FromECDSA(got))

			assert.Equal(t, want, tt.args.privateKey)
		})
	}
}

func TestHexECDSAPub(t *testing.T) {
	addr, err := GenerateAddress()
	if err != nil {
		t.Logf("get new address error: %s", err)
		return
	}

	type args struct {
		publicKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				publicKey: addr.PublicKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexECDSAPub(tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err, "HexECDSAPub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := hexutil.Encode(crypto.FromECDSAPub(got))

			assert.Equal(t, want, tt.args.publicKey)
		})
	}
}
