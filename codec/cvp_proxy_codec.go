package codec

import (
	"fmt"
	"github.com/bcdevtools/cvp-streaming-core/types"
	"regexp"
)

var _ CvpCodec = (*proxyCvpCodec)(nil)

// proxyCvpCodec is an implementation of CvpCodec.
//
// The proxy automatically detect version of encoded data and forward to the corresponding implementation for decoding.
//
// When invoking encode functions, it forward to default CvpCodec.
type proxyCvpCodec struct {
	cvpCodecImpl CvpCodec
}

// NewProxyCvpCodec returns new instance of proxy CvpCodec.
//
// The proxy automatically detect version of encoded data and forward to the corresponding implementation for decoding.
//
// When invoking encode functions, it forward to default CvpCodec version.
func NewProxyCvpCodec() CvpCodec {
	return proxyCvpCodec{
		cvpCodecImpl: GetCvpCodecV3(),
	}
}

// WrapCvpCodecInProxy wraps a CvpCodec into a proxy CvpCodec.
//
// The proxy automatically detect version of encoded data and forward to the corresponding implementation for decoding.
//
// When invoking encode functions, it forward to the provided version.
func WrapCvpCodecInProxy(inner CvpCodec) CvpCodec {
	if _, ok := inner.(proxyCvpCodec); ok {
		panic(fmt.Errorf("can not wrap proxy CvpCodec into another proxy CvpCodec"))
	}
	return proxyCvpCodec{
		cvpCodecImpl: inner,
	}
}

func (p proxyCvpCodec) EncodeStreamingLightValidators(validators types.StreamingLightValidators) []byte {
	return p.cvpCodecImpl.EncodeStreamingLightValidators(validators)
}

func (p proxyCvpCodec) DecodeStreamingLightValidators(bz []byte) (types.StreamingLightValidators, error) {
	possibleVersion, detected := DetectEncodingVersion(bz)
	if detected {
		switch possibleVersion {
		case CvpCodecVersionV3:
			return GetCvpCodecV3().DecodeStreamingLightValidators(bz)
		case CvpCodecVersionV2:
			return GetCvpCodecV2().DecodeStreamingLightValidators(bz)
		case CvpCodecVersionV1:
			//goland:noinspection GoDeprecation
			return GetCvpCodecV1().DecodeStreamingLightValidators(bz)
		}
	}

	return nil, fmt.Errorf("unable to detect encoder version")
}

func (p proxyCvpCodec) EncodeStreamingNextBlockVotingInformation(information *types.StreamingNextBlockVotingInformation) []byte {
	return p.cvpCodecImpl.EncodeStreamingNextBlockVotingInformation(information)
}

var regexpHeightRoundStep = regexp.MustCompile(`^\d+/\d+/\d+$`)
var regexpPreVotedFingerprintBlockHash = regexp.MustCompile(`^[a-fA-F\d]{4}$`)

func (p proxyCvpCodec) DecodeStreamingNextBlockVotingInformation(bz []byte) (*types.StreamingNextBlockVotingInformation, error) {
	possibleVersion, detected := DetectEncodingVersion(bz)
	if detected {
		switch possibleVersion {
		case CvpCodecVersionV3:
			return GetCvpCodecV3().DecodeStreamingNextBlockVotingInformation(bz)
		case CvpCodecVersionV2:
			return GetCvpCodecV2().DecodeStreamingNextBlockVotingInformation(bz)
		case CvpCodecVersionV1:
			//goland:noinspection GoDeprecation
			return GetCvpCodecV1().DecodeStreamingNextBlockVotingInformation(bz)
		}
	}

	return nil, fmt.Errorf("unable to detect encoder version")
}

func (p proxyCvpCodec) GetVersion() CvpCodecVersion {
	return p.cvpCodecImpl.GetVersion()
}
