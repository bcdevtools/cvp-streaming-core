package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPreVoteStreamingSession(t *testing.T) {
	tests := []struct {
		name    string
		chainId string
		wantErr bool
	}{
		{
			name:    "normal",
			chainId: "cosmoshub-4",
			wantErr: false,
		},
		{
			name:    "bad chain id",
			chainId: " 8poles",
			wantErr: true,
		},
		{
			name:    "empty chain id",
			chainId: "",
			wantErr: true,
		},
		{
			name:    "blank chain id",
			chainId: "      ",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotKey, err := NewPreVoteStreamingSession(tt.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPreVoteStreamingSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			require.Nil(t, gotId.ValidateBasic())
			require.Nil(t, gotKey.ValidateBasic())
		})
	}
}
