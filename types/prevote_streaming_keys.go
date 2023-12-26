package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type PreVoteStreamingSessionId string
type PreVoteStreamingSessionKey string

var regexpChainId = regexp.MustCompile(`^[a-zA-Z\d][a-zA-Z\d_-]{2,41}$`)

// NewPreVoteStreamingSession generates a new session id and key pair with seed is given chain-id
func NewPreVoteStreamingSession(chainId string) (PreVoteStreamingSessionId, PreVoteStreamingSessionKey, error) {
	if !regexpChainId.MatchString(chainId) {
		return "", "", fmt.Errorf("invalid chain id")
	}

	bufferId := make([]byte, 32)
	_, err := rand.Read(bufferId)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to generate random bytes")
	}

	sid := PreVoteStreamingSessionId(fmt.Sprintf("%s_%X", chainId, bufferId))

	bufferKey := make([]byte, 32)
	_, err = rand.Read(bufferKey)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to generate random bytes")
	}

	sk := PreVoteStreamingSessionKey(hex.EncodeToString(bufferKey))

	return sid, sk, nil
}

var regexpPreVoteStreamingSessionId = regexp.MustCompile(`^[a-zA-Z\d_\-]+_[A-F\d]{64}$`)

// ValidateBasic returns an error if the session id is invalid format.
func (sid PreVoteStreamingSessionId) ValidateBasic() error {
	if len(sid) == 0 {
		return fmt.Errorf("empty")
	}

	if !regexpPreVoteStreamingSessionId.MatchString(string(sid)) {
		return fmt.Errorf("invalid format %s", sid)
	}

	return nil
}

// ForChainId returns true if the session id value is for the given chain id.
func (sid PreVoteStreamingSessionId) ForChainId(chainId string) bool {
	return strings.HasPrefix(string(sid), chainId+"_") && sid.ValidateBasic() == nil
}

var regexpPreVoteStreamingSessionKey = regexp.MustCompile(`^[a-f\d]{64}$`)

// ValidateBasic returns an error if the session key is invalid format.
func (sk PreVoteStreamingSessionKey) ValidateBasic() error {
	if len(sk) == 0 {
		return fmt.Errorf("empty")
	}

	if !regexpPreVoteStreamingSessionKey.MatchString(string(sk)) {
		return fmt.Errorf("invalid format %s", sk)
	}

	return nil
}
