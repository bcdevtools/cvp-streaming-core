package codec

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bcdevtools/cvp-streaming-core/constants"
	"github.com/bcdevtools/cvp-streaming-core/types"
	"reflect"
	"strings"
	"testing"
	"time"
)

var cvpV3CodecImpl = getCvpCodecV3()

func bufferFromHex(hexString string) []byte {
	bz, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return bz
}

func Test_cvpCodecV3_EncodeDecodeStreamingLightValidators(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                               string
		validators                         types.StreamingLightValidators
		wantPanicEncode                    bool
		wantEncodedData                    []byte
		wantErrDecode                      bool
		wantErrDecodeContains              string
		wantDecodedOrUseInputAsWantDecoded types.StreamingLightValidators // if missing, use input as expect
	}{
		{
			name: "normal, 2 validators",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 10.11,
					Moniker:                   "Val1",
				},
				{
					Index:                     1,
					VotingPowerDisplayPercent: 01.02,
					Moniker:                   "Val2",
				},
			},
			wantPanicEncode: false,
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa6160e0e20ecb752bf60d764cf774c6c0b6350c8c8c4c600599d81500020000ffff80dcd11643000000"),
			wantErrDecode:   false,
		},
		{
			name: "normal, 1 validator",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 10.11,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: false,
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa6160e0e20ecb752bf60d764cf774c6c0b680000000ffffc73c489022000000"),
			wantErrDecode:   false,
		},
		{
			name: "truncate before encode then decode correct moniker UTF-8",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 10.10,
					Moniker:                   "✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅",
				},
				{
					Index:                     1,
					VotingPowerDisplayPercent: 01.02,
					Moniker:                   "❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌",
				},
			},
			wantPanicEncode: false,
			wantErrDecode:   false,
			wantDecodedOrUseInputAsWantDecoded: []types.StreamingLightValidator{
				// moniker of validators are truncated to max 20 bytes of runes
				{
					Index:                     0,
					VotingPowerDisplayPercent: 10.10,
					Moniker:                   "✅✅✅✅✅✅",
				},
				{
					Index:                     1,
					VotingPowerDisplayPercent: 01.02,
					Moniker:                   "❌❌❌❌❌❌",
				},
			},
		},
		{
			name: "normal, validator with 100% VP",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 100,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: false,
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa6160486108cb752bf60d764cf774c6c0b680000000ffff7548b1d322000000"),
			wantErrDecode:   false,
		},
		{
			name:                  "not accept empty validator list",
			validators:            []types.StreamingLightValidator{},
			wantPanicEncode:       false,
			wantErrDecode:         true,
			wantErrDecodeContains: "invalid empty validator raw data",
		},
		{
			name: "not accept validator negative index",
			validators: []types.StreamingLightValidator{
				{
					Index:                     -1,
					VotingPowerDisplayPercent: 99,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "not accept validator index greater than 998",
			validators: []types.StreamingLightValidator{
				{
					Index:                     999,
					VotingPowerDisplayPercent: 99,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "not accept validator negative voting power percent",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: -0.01,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "not accept validator voting power percent greater than 100%",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 100.01,
					Moniker:                   "Val1",
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "validator list size larger than cap",
			validators: func() types.StreamingLightValidators {
				var validators types.StreamingLightValidators
				for v := 1; v <= constants.MAX_VALIDATORS+1; v++ {
					validators = append(validators, types.StreamingLightValidator{
						Index:                     v - 1,
						VotingPowerDisplayPercent: 99,
						Moniker:                   fmt.Sprintf("Val%d", v),
					})
				}
				return validators
			}(),
			wantPanicEncode: true,
		},
		{
			name: "keep only first 20 bytes of moniker",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 99,
					Moniker:                   "123456789012345678901234567890",
				},
			},
			wantPanicEncode: false,
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa61604866f00df1acf2730935f2ab4a37f57571adf4ad0a34f4cb4a36f10f71b405040000ffff1a0dbe2022000000"),
			wantDecodedOrUseInputAsWantDecoded: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 99,
					Moniker:                   "12345678901234567890", // truncated
				},
			},
			wantErrDecode: false,
		},
		{
			name: "sanitize moniker",
			validators: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 99,
					Moniker:                   `<he'llo">`,
				},
			},
			wantPanicEncode: false,
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa616048660870cfc8f132aa284eaaf4d4f674764c47c2b680000000ffffc22b3a5b22000000"),
			wantDecodedOrUseInputAsWantDecoded: []types.StreamingLightValidator{
				{
					Index:                     0,
					VotingPowerDisplayPercent: 99,
					Moniker:                   "(he`llo`)",
				},
			},
			wantErrDecode: false,
		},
		{
			name: "collision of separator byte and bytes index",
			validators: func() types.StreamingLightValidators {
				var result types.StreamingLightValidators
				for i := 0; i < constants.MAX_VALIDATORS; i++ {
					result = append(result, types.StreamingLightValidator{
						Index:                     i,
						VotingPowerDisplayPercent: 99,
						Moniker:                   fmt.Sprintf("Val%d", i+1),
					})
				}
				return result
			}(),
			wantPanicEncode: false,
			wantErrDecode:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := func() (bz []byte) {
				defer func() {
					err := recover()
					if err != nil {
						if !tt.wantPanicEncode {
							t.Errorf("EncodeStreamingLightValidators() panic = %v but not wanted", err)
						}
					} else {
						if tt.wantPanicEncode {
							t.Errorf("EncodeStreamingLightValidators() panic = %v but wanted panic", err)
						}
					}
				}()
				bz = cvpV3CodecImpl.EncodeStreamingLightValidators(tt.validators)
				return
			}()

			if tt.wantPanicEncode {
				return
			}

			if len(tt.wantEncodedData) > 0 {
				if !bytes.Equal(gotEncoded, tt.wantEncodedData) {
					t.Errorf("EncodeStreamingLightValidators()\n%v (got)\n%v (want)", hex.EncodeToString(gotEncoded), hex.EncodeToString(tt.wantEncodedData))
					return
				}
			}

			gotDecoded, err := cvpV3CodecImpl.DecodeStreamingLightValidators(gotEncoded)
			if (err != nil) != tt.wantErrDecode {
				t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr %v", err, tt.wantErrDecode)
				return
			}
			if err == nil {
				if tt.wantDecodedOrUseInputAsWantDecoded == nil {
					tt.wantDecodedOrUseInputAsWantDecoded = tt.validators
				}
				if !reflect.DeepEqual(gotDecoded, tt.wantDecodedOrUseInputAsWantDecoded) {
					t.Errorf("DecodeStreamingLightValidators()\ngot = %v,\nwant %v", gotDecoded, tt.wantDecodedOrUseInputAsWantDecoded)
				}
			} else {
				if tt.wantErrDecodeContains == "" {
					t.Errorf("missing setup check error content, actual error: %v", err)
				} else {
					if !strings.Contains(err.Error(), tt.wantErrDecodeContains) {
						t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr contains %v", err, tt.wantErrDecodeContains)
					}
				}
			}
		})
	}
}

func Test_cvpCodecV3_DecodeStreamingLightValidators(t *testing.T) {
	// The codec v3 is use v2 underlying and gzip so not much to test here.

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                  string
		inputEncodedData      []byte
		wantDecoded           types.StreamingLightValidators
		wantErrDecode         bool
		wantErrDecodeContains string
	}{
		{
			name: "icorrect codec version",
			inputEncodedData: mergeBuffers(
				[]byte{'1', cvpCodecV3Separator},
				[]byte{0x0, 0x0}, []byte{0x0a, 0x0a}, b64bz(fssut("Val1", 20)),
			),
			wantErrDecode:         true,
			wantErrDecodeContains: "bad encoding prefix",
		},
		{
			name: "icorrect codec version",
			inputEncodedData: mergeBuffers(
				[]byte{'2', cvpCodecV3Separator},
				[]byte{0x0, 0x0}, []byte{0x0a, 0x0a}, b64bz(fssut("Val1", 20)),
			),
			wantErrDecode:         true,
			wantErrDecodeContains: "bad encoding prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDecoded, err := cvpV3CodecImpl.DecodeStreamingLightValidators(tt.inputEncodedData)

			if (err != nil) != tt.wantErrDecode {
				if err == nil {
					fmt.Println("Un-expected result:", gotDecoded)
				}
				t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr %v", err, tt.wantErrDecode)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(gotDecoded, tt.wantDecoded) {
					t.Errorf("DecodeStreamingLightValidators()\ngot = %v,\nwant %v", gotDecoded, tt.wantDecoded)
				}
			} else {
				if tt.wantErrDecodeContains == "" {
					t.Errorf("missing setup check error content, actual error: %v", err)
				} else {
					if !strings.Contains(err.Error(), tt.wantErrDecodeContains) {
						t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr contains %v", err, tt.wantErrDecodeContains)
					}
				}
			}
		})
	}
}

func Test_cvpCodecV3_EncodeAndDecodeStreamingNextBlockVotingInformation(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                               string
		inf                                types.StreamingNextBlockVotingInformation
		wantPanicEncode                    bool
		wantEncodedData                    []byte
		wantDecodedOrUseInputAsWantDecoded *types.StreamingNextBlockVotingInformation // if missing, use input as expect
		wantErrDecode                      bool
		wantErrDecodeContains              string
	}{
		{
			name: "normal, 4 validators",
			inf: types.StreamingNextBlockVotingInformation{
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
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31ac61646032ab6160707472767166603400010626102f8c81595757573702100000ffff5854bb842b000000"),
			wantErrDecode:   false,
		},
		{
			name: "normal, 1 validators",
			inf: types.StreamingNextBlockVotingInformation{
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
				},
			},
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31ac61646032ab6160707472767106040000ffff92aa111416000000"),
			wantErrDecode:   false,
		},
		{
			name: "can not decode zero validators vote state",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              1 * time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2.54,
				ValidatorVoteStates:   []types.StreamingValidatorVoteState{},
			},
			wantEncodedData:       bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31ac61646032ab01040000ffff950029b10f000000"),
			wantErrDecode:         true,
			wantErrDecodeContains: "missing validator vote states",
		},
		{
			name: "duration will be corrected to zero if negative",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              -1 * time.Second,
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
				},
			},
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31a861646032ab6160707472767106040000ffffe44b1e8916000000"),
			wantDecodedOrUseInputAsWantDecoded: &types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              0,
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
				},
			},
		},
		{
			name: "percent computed correctly",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              2 * time.Second,
				PreVotedPercent:       99.98,
				PreCommitVotedPercent: 97.96,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "ABCD",
						PreVoted:          true,
						VotedZeroes:       false,
						PreCommitVoted:    true,
					},
				},
			},
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31aa494e4a4ca86160707472767106040000ffffa8d5399f16000000"),
		},
		{
			name: "panic encode if negative validator index",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex: -1,
					},
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "panic encode if validator index greater than 998",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex: 999,
					},
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "panic encode if validator list size larger than cap",
			inf: func() types.StreamingNextBlockVotingInformation {
				inf := types.StreamingNextBlockVotingInformation{
					HeightRoundStep:       "1/2/3",
					Duration:              time.Second,
					PreVotedPercent:       1,
					PreCommitVotedPercent: 2,
				}

				for v := 1; v <= constants.MAX_VALIDATORS+1; v++ {
					inf.ValidatorVoteStates = append(inf.ValidatorVoteStates, types.StreamingValidatorVoteState{
						ValidatorIndex: v - 1,
					})
				}

				return inf
			}(),
			wantPanicEncode: true,
		},
		{
			name: "panic encode if block hash length is not 0 or 4",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "123",
						PreVoted:          true,
					},
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "panic encode if block hash length is not 0 or 4",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "12345",
						PreVoted:          true,
					},
				},
			},
			wantPanicEncode: true,
		},
		{
			name: "automatically fill prevoted block hash if empty",
			inf: types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              3 * time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "",
					},
				},
			},
			wantEncodedData: bufferFromHex("037c1f8b08000000000000ff62aa31d437d237ae31ae61646062a86160d0d5d5d58d00040000ffff5655c0eb16000000"),
			wantDecodedOrUseInputAsWantDecoded: &types.StreamingNextBlockVotingInformation{
				HeightRoundStep:       "1/2/3",
				Duration:              3 * time.Second,
				PreVotedPercent:       1,
				PreCommitVotedPercent: 2,
				ValidatorVoteStates: []types.StreamingValidatorVoteState{
					{
						ValidatorIndex:    0,
						PreVotedBlockHash: "----",
					},
				},
			},
		},
		{
			name: "collision of separator byte with bytes index",
			inf: func() types.StreamingNextBlockVotingInformation {
				nextBlockVotingInfo := types.StreamingNextBlockVotingInformation{
					HeightRoundStep:       "1/2/3",
					Duration:              time.Second,
					PreVotedPercent:       1,
					PreCommitVotedPercent: 2,
				}

				for i := 0; i < constants.MAX_VALIDATORS; i++ {
					nextBlockVotingInfo.ValidatorVoteStates = append(nextBlockVotingInfo.ValidatorVoteStates, types.StreamingValidatorVoteState{
						ValidatorIndex:    i,
						PreVotedBlockHash: "C0FF",
						PreVoted:          true,
					})
				}

				return nextBlockVotingInfo
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded := func() (bz []byte) {
				defer func() {
					err := recover()
					if err != nil {
						if !tt.wantPanicEncode {
							t.Errorf("EncodeStreamingNextBlockVotingInformation() panic = %v but not wanted", err)
						}
					} else {
						if tt.wantPanicEncode {
							t.Errorf("EncodeStreamingNextBlockVotingInformation() panic = %v but wanted panic", err)
						}
					}
				}()
				bz = cvpV3CodecImpl.EncodeStreamingNextBlockVotingInformation(&tt.inf)
				return
			}()

			if tt.wantPanicEncode {
				return
			}

			if len(tt.wantEncodedData) > 0 {
				if !bytes.Equal(gotEncoded, tt.wantEncodedData) {
					t.Errorf("EncodeStreamingNextBlockVotingInformation()\n%v (got)\n%v (want)", hex.EncodeToString(gotEncoded), hex.EncodeToString(tt.wantEncodedData))
					return
				}
			}

			gotDecoded, err := cvpV3CodecImpl.DecodeStreamingNextBlockVotingInformation(gotEncoded)
			if (err != nil) != tt.wantErrDecode {
				t.Errorf("DecodeStreamingNextBlockVotingInformation() error = %v, wantErr %v", err, tt.wantErrDecode)
				return
			}
			if err == nil {
				if tt.wantDecodedOrUseInputAsWantDecoded == nil {
					tt.wantDecodedOrUseInputAsWantDecoded = &tt.inf
				}
				if !reflect.DeepEqual(gotDecoded, tt.wantDecodedOrUseInputAsWantDecoded) {
					t.Errorf("DecodeStreamingNextBlockVotingInformation()\ngot = %v,\nwant %v", gotDecoded, tt.wantDecodedOrUseInputAsWantDecoded)
				}
			} else {
				if tt.wantErrDecodeContains == "" {
					t.Errorf("missing setup check error content, actual error: %v", err)
				} else {
					if !strings.Contains(err.Error(), tt.wantErrDecodeContains) {
						t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr contains %v", err, tt.wantErrDecodeContains)
					}
				}
			}
		})
	}
}

func Test_cvpCodecV3_DecodeStreamingNextBlockVotingInformation(t *testing.T) {
	// The codec v3 is use v2 underlying and gzip so not much to test here.

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                  string
		inputEncodedData      []byte
		wantDecoded           *types.StreamingNextBlockVotingInformation
		wantErrDecode         bool
		wantErrDecodeContains string
	}{
		{
			name:                  "icorrect codec version",
			inputEncodedData:      []byte("1|1/2/3|1000|100|254|000ABCDC"),
			wantErrDecode:         true,
			wantErrDecodeContains: "bad encoding prefix",
		},
		{
			name: "icorrect codec version",
			inputEncodedData: mergeBuffers(
				[]byte{0x2, cvpCodecV2Separator},
				[]byte("1/2/3"), []byte{cvpCodecV2Separator},
				[]byte("1"), []byte{cvpCodecV2Separator},
				[]byte{0x01, 0x00}, []byte{0x02, 0x36}, []byte{cvpCodecV2Separator},
				[]byte{0x00, 0x00}, []byte("ABCD"), []byte("C"),
			),
			wantErrDecode:         true,
			wantErrDecodeContains: "bad encoding prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDecoded, err := cvpV3CodecImpl.DecodeStreamingNextBlockVotingInformation(tt.inputEncodedData)

			if (err != nil) != tt.wantErrDecode {
				if err == nil {
					fmt.Println("Un-expected result:", gotDecoded)
				}
				t.Errorf("DecodeStreamingNextBlockVotingInformation() error = %v, wantErr %v", err, tt.wantErrDecode)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(gotDecoded, tt.wantDecoded) {
					t.Errorf("DecodeStreamingNextBlockVotingInformation()\ngot = %v,\nwant %v", gotDecoded, tt.wantDecoded)
				}
			} else {
				if tt.wantErrDecodeContains == "" {
					t.Errorf("missing setup check error content, actual error: %v", err)
				} else {
					if !strings.Contains(err.Error(), tt.wantErrDecodeContains) {
						t.Errorf("DecodeStreamingLightValidators() error = %v, wantErr contains %v", err, tt.wantErrDecodeContains)
					}
				}
			}
		})
	}
}

func Test_cvpCodecV3_Base64EncodedMonikerBufferSize(t *testing.T) {
	bz := make([]byte, cvpCodecV2MonikerBufferSize)
	for r := 1; r <= 100; r++ {
		size, err := rand.Read(bz)
		if err != nil {
			t.Fatal(err)
			return
		}
		if size != cvpCodecV2MonikerBufferSize {
			t.Fatalf("bad rand read size: %v", size)
			return
		}
		base64EncodedMonikerBuffer := []byte(base64.StdEncoding.EncodeToString(bz))
		if len(base64EncodedMonikerBuffer) != cvpCodecV2Base64EncodedMonikerBufferSize {
			t.Fatalf("bad base64 encoded moniker buffer size: %v", len(base64EncodedMonikerBuffer))
			return
		}
	}
}
