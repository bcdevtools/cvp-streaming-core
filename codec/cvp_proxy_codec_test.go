package codec

import (
	"github.com/bcdevtools/cvp-streaming-core/types"
	"reflect"
	"strings"
	"testing"
	"time"
)

//goland:noinspection GoVarAndConstTypeMayBeOmitted
var cvpProxyCodecImpl CvpCodec = NewProxyCvpCodec()

func Test_proxyCvpCodec_WrapCvpCodecInProxy(t *testing.T) {
	t.Run("can wrap non-proxy codec", func(t *testing.T) {
		_ = WrapCvpCodecInProxy(GetCvpCodecV1())
		_ = WrapCvpCodecInProxy(GetCvpCodecV2())
		_ = WrapCvpCodecInProxy(GetCvpCodecV3())
	})
	t.Run("can not wrap proxy codec", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("WrapCvpCodecInProxy() did not panic when wrap a proxy codec")
			}
		}()
		_ = WrapCvpCodecInProxy(NewProxyCvpCodec())
	})
}

func Test_proxyCvpCodec_EncodeDecodeStreamingLightValidators(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		input types.StreamingLightValidators
		want  types.StreamingLightValidators
	}{
		{
			name: "can codec",
			input: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 0.1,
					Moniker:                   "Val1",
				},
				{
					Index:                     1,
					VotingPowerDisplayPercent: 2.5,
					Moniker:                   "Val2",
				},
			},
			want: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 0.1,
					Moniker:                   "Val1",
				},
				{
					Index:                     1,
					VotingPowerDisplayPercent: 2.5,
					Moniker:                   "Val2",
				},
			},
		},
		{
			name: "sanitize moniker",
			input: []types.StreamingLightValidator{
				{
					VotingPowerDisplayPercent: 0.1,
					Moniker:                   `<he'llo">`,
				},
			},
			want: []types.StreamingLightValidator{
				{
					VotingPowerDisplayPercent: 0.1,
					Moniker:                   "(he`llo`)",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := cvpProxyCodecImpl.EncodeStreamingLightValidators(tt.input)
			gotDecoded, err := cvpProxyCodecImpl.DecodeStreamingLightValidators(gotEncoded)
			if err != nil {
				t.Errorf("DecodeStreamingLightValidators() error = %v", err)
				return
			}
			if reflect.DeepEqual(gotDecoded, tt.want) {
				testDetect := func(codec CvpCodec) {
					gotEncoded := codec.EncodeStreamingLightValidators(tt.input)
					_, err := cvpProxyCodecImpl.DecodeStreamingLightValidators(gotEncoded)
					if err != nil {
						t.Errorf("proxy not forward %T encoded data correctly, error = %v", codec, err)
					}
				}
				testDetect(cvpV1CodecImpl)
				testDetect(cvpV2CodecImpl)
				testDetect(cvpV3CodecImpl)
			} else {
				t.Errorf("DecodeStreamingLightValidators()\ngotDecoded = %v\nwant %v", gotDecoded, tt.want)
			}
		})
	}

	t.Run("want error if decoder not able to detect encoder version", func(t *testing.T) {
		_, err := cvpProxyCodecImpl.DecodeStreamingLightValidators([]byte("invalid data"))
		if err == nil {
			t.Errorf("DecodeStreamingLightValidators() did not return error")
			return
		}
		if !strings.Contains(err.Error(), "unable to detect encoder version") {
			t.Errorf("DecodeStreamingLightValidators() error = %v, want contains %v", err, "unable to detect encoder version")
		}
	})
}

func Test_proxyCvpCodec_EncodeDecodeStreamingNextBlockVotingInformation(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name  string
		input *types.StreamingNextBlockVotingInformation
	}{
		{
			name: "can codec",
			input: &types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              1 * time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2.54,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "ABCD",
						PreVoted:          true,
						VotedZeroes:       false,
						PreCommitVoted:    true,
					},
					{
						ValidatorIndex:    1,
						PreVotedBlockHash: "0000",
						PreVoted:          true,
						VotedZeroes:       true,
						PreCommitVoted:    false,
					},
					{
						ValidatorIndex:    2,
						PreVotedBlockHash: "ABCD",
						PreVoted:          true,
						VotedZeroes:       false,
						PreCommitVoted:    false,
					},
					{
						ValidatorIndex:    3,
						PreVotedBlockHash: "----",
						PreVoted:          false,
						VotedZeroes:       false,
						PreCommitVoted:    false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := cvpProxyCodecImpl.EncodeStreamingNextBlockVotingInformation(tt.input)
			gotDecoded, err := cvpProxyCodecImpl.DecodeStreamingNextBlockVotingInformation(gotEncoded)
			if err != nil {
				t.Errorf("DecodeStreamingNextBlockVotingInformation() error = %v", err)
				return
			}
			if reflect.DeepEqual(gotDecoded, tt.input) {
				testDetect := func(codec CvpCodec) {
					gotEncoded := codec.EncodeStreamingNextBlockVotingInformation(tt.input)
					_, err := cvpProxyCodecImpl.DecodeStreamingNextBlockVotingInformation(gotEncoded)
					if err != nil {
						t.Errorf("proxy not forward %T encoded data correctly, error = %v", codec, err)
					}
				}
				testDetect(cvpV1CodecImpl)
				testDetect(cvpV2CodecImpl)
				testDetect(cvpV3CodecImpl)
			} else {
				t.Errorf("DecodeStreamingNextBlockVotingInformation()\ngotDecoded = %v\nwant %v", gotDecoded, tt.input)
			}
		})
	}

	t.Run("want error if decoder not able to detect encoder version", func(t *testing.T) {
		_, err := cvpProxyCodecImpl.DecodeStreamingNextBlockVotingInformation([]byte("invalid data"))
		if err == nil {
			t.Errorf("DecodeStreamingNextBlockVotingInformation() did not return error")
			return
		}
		if !strings.Contains(err.Error(), "unable to detect encoder version") {
			t.Errorf("DecodeStreamingNextBlockVotingInformation() error = %v, want contains %v", err, "unable to detect encoder version")
		}
	})
}
