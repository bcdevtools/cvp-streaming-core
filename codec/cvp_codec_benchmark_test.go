package codec

import (
	"fmt"
	"github.com/bcdevtools/cvp-streaming-core/types"
	"testing"
	"time"
)

var benchmarkDataSizes = []int{1, 8, 60, 100, 150, 180}

func BenchmarkEncodeLightValidators(b *testing.B) {
	for _, benchmarkDataSize := range benchmarkDataSizes {
		validators := types.StreamingLightValidators{}
		for v := 1; v <= benchmarkDataSize; v++ {
			validators = append(validators, types.StreamingLightValidator{
				Index:                     v - 1,
				VotingPowerDisplayPercent: 99.98,
				Moniker:                   fmt.Sprintf("Val%d✅✅✅✅✅✅✅", v),
			})
		}

		b.Run(fmt.Sprintf("codec v2 encode %d validators", benchmarkDataSize), func(b *testing.B) {
			_ = cvpV2CodecImpl.EncodeStreamingLightValidators(validators)
		})

		b.Run(fmt.Sprintf("codec v3 encode %d validators", benchmarkDataSize), func(b *testing.B) {
			_ = cvpV3CodecImpl.EncodeStreamingLightValidators(validators)
		})
	}
}

func BenchmarkDecodeLightValidators(b *testing.B) {
	for _, benchmarkDataSize := range benchmarkDataSizes {
		validators := types.StreamingLightValidators{}
		for v := 1; v <= benchmarkDataSize; v++ {
			validators = append(validators, types.StreamingLightValidator{
				Index:                     v - 1,
				VotingPowerDisplayPercent: 99.98,
				Moniker:                   fmt.Sprintf("Val%d✅✅✅✅✅✅✅", v),
			})
		}

		encodedV2 := cvpV2CodecImpl.EncodeStreamingLightValidators(validators)
		encodedV3 := cvpV3CodecImpl.EncodeStreamingLightValidators(validators)

		b.Run(fmt.Sprintf("codec v2 decode %d validators", benchmarkDataSize), func(b *testing.B) {
			_, _ = cvpV2CodecImpl.DecodeStreamingLightValidators(encodedV2)
		})

		b.Run(fmt.Sprintf("codec v3 decode %d validators", benchmarkDataSize), func(b *testing.B) {
			_, _ = cvpV3CodecImpl.DecodeStreamingLightValidators(encodedV3)
		})
	}
}

func BenchmarkEncodeNextBlockPreVoteInfo(b *testing.B) {
	for _, benchmarkDataSize := range benchmarkDataSizes {
		inf := types.StreamingNextBlockVotingInformation{
			HeightRoundStep:       "999999999/9999/9999",
			Duration:              365 * 2 * 24 * time.Hour,
			PreVotedPercent:       99.98,
			PreCommitVotedPercent: 99.98,
			ValidatorVoteStates:   nil,
		}
		for v := 1; v <= benchmarkDataSize; v++ {
			inf.ValidatorVoteStates = append(inf.ValidatorVoteStates, types.StreamingValidatorVoteState{
				ValidatorIndex:    v - 1,
				PreVotedBlockHash: "C0FF",
				PreVoted:          true,
				VotedZeroes:       false,
				PreCommitVoted:    true,
			})
		}

		b.Run(fmt.Sprintf("codec v2 encode %d votes", benchmarkDataSize), func(b *testing.B) {
			_ = cvpV2CodecImpl.EncodeStreamingNextBlockVotingInformation(&inf)
		})

		b.Run(fmt.Sprintf("codec v3 encode %d votes", benchmarkDataSize), func(b *testing.B) {
			_ = cvpV3CodecImpl.EncodeStreamingNextBlockVotingInformation(&inf)
		})
	}
}

func BenchmarkDecodeNextBlockPreVoteInfo(b *testing.B) {
	for _, benchmarkDataSize := range benchmarkDataSizes {
		inf := types.StreamingNextBlockVotingInformation{
			HeightRoundStep:       "999999999/9999/9999",
			Duration:              365 * 2 * 24 * time.Hour,
			PreVotedPercent:       99.98,
			PreCommitVotedPercent: 99.98,
			ValidatorVoteStates:   nil,
		}
		for v := 1; v <= benchmarkDataSize; v++ {
			inf.ValidatorVoteStates = append(inf.ValidatorVoteStates, types.StreamingValidatorVoteState{
				ValidatorIndex:    v - 1,
				PreVotedBlockHash: "C0FF",
				PreVoted:          true,
				VotedZeroes:       false,
				PreCommitVoted:    true,
			})
		}

		encodedV2 := cvpV2CodecImpl.EncodeStreamingNextBlockVotingInformation(&inf)
		encodedV3 := cvpV3CodecImpl.EncodeStreamingNextBlockVotingInformation(&inf)

		b.Run(fmt.Sprintf("codec v2 decode %d votes", benchmarkDataSize), func(b *testing.B) {
			_, _ = cvpV2CodecImpl.DecodeStreamingNextBlockVotingInformation(encodedV2)
		})

		b.Run(fmt.Sprintf("codec v3 decode %d votes", benchmarkDataSize), func(b *testing.B) {
			_, _ = cvpV3CodecImpl.DecodeStreamingNextBlockVotingInformation(encodedV3)
		})
	}
}