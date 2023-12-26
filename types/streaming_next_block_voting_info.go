package types

import "time"

type StreamingNextBlockVotingInformation struct {
	HeightRoundStep       string                        `json:"hrs"`
	Duration              time.Duration                 `json:"d,omitempty"`
	PreVotedPercent       float64                       `json:"pv,omitempty"`
	PreCommitVotedPercent float64                       `json:"pc,omitempty"`
	ValidatorVoteStates   []StreamingValidatorVoteState `json:"v,omitempty"`
}

type StreamingValidatorVoteState struct {
	ValidatorIndex    int    `json:"i"`
	PreVotedBlockHash string `json:"hash,omitempty"`
	PreVoted          bool   `json:"pv,omitempty"`
	VotedZeroes       bool   `json:"vz,omitempty"`
	PreCommitVoted    bool   `json:"pc,omitempty"`
}
