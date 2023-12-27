package codec

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/bcdevtools/cvp-streaming-core/constants"
	"github.com/bcdevtools/cvp-streaming-core/types"
	"github.com/pkg/errors"
	"io"
)

//goland:noinspection SpellCheckingInspection

var _ CvpCodec = (*cvpCodecV3)(nil)

const cvpCodecV3Separator byte = '|'

var prefixDataEncodedByCvpCodecV3 = []byte{0x3, cvpCodecV3Separator}

type cvpCodecV3 struct {
	v2Codec CvpCodec
}

// GetCvpCodecV3 returns new instance of CvpCodec that actually encode data using v2 codec then gzip it.
// Procedures smaller data than v2 codec in most cases with large data size.
// But slower than v2 codec, ofc.
func GetCvpCodecV3() CvpCodec {
	return cvpCodecV3{
		v2Codec: GetCvpCodecV2(),
	}
}

func (c cvpCodecV3) EncodeStreamingLightValidators(validators types.StreamingLightValidators) []byte {
	if len(validators) > constants.MAX_VALIDATORS {
		panic(fmt.Errorf("too many validators: %d/%d", len(validators), constants.MAX_VALIDATORS))
	}

	bzByV2 := c.v2Codec.EncodeStreamingLightValidators(validators)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(bzByV2)
	if err != nil {
		panic(errors.Wrap(err, "failed to write gzipped content"))
	}
	err = w.Close()
	if err != nil {
		panic(errors.Wrap(err, "failed to close gzip writer"))
	}

	return append(prefixDataEncodedByCvpCodecV3, b.Bytes()...)
}

func (c cvpCodecV3) DecodeStreamingLightValidators(bz []byte) (types.StreamingLightValidators, error) {
	if !bytes.HasPrefix(bz, prefixDataEncodedByCvpCodecV3) {
		return nil, fmt.Errorf("bad encoding prefix")
	}

	gzipr, err := gzip.NewReader(bytes.NewReader(bz[2:]))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gzip reader")
	}

	bzByV2, err := io.ReadAll(gzipr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read gzipped content")
	}

	err = gzipr.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close gzip reader")
	}

	return c.v2Codec.DecodeStreamingLightValidators(bzByV2)
}

func (c cvpCodecV3) EncodeStreamingNextBlockVotingInformation(inf *types.StreamingNextBlockVotingInformation) []byte {
	if len(inf.ValidatorVoteStates) > constants.MAX_VALIDATORS {
		panic(fmt.Errorf("too many validators: %d/%d", len(inf.ValidatorVoteStates), constants.MAX_VALIDATORS))
	}

	bzByV2 := c.v2Codec.EncodeStreamingNextBlockVotingInformation(inf)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(bzByV2)
	if err != nil {
		panic(errors.Wrap(err, "failed to write gzipped content"))
	}
	err = w.Close()
	if err != nil {
		panic(errors.Wrap(err, "failed to close gzip writer"))
	}

	return append(prefixDataEncodedByCvpCodecV3, b.Bytes()...)
}

func (c cvpCodecV3) DecodeStreamingNextBlockVotingInformation(bz []byte) (*types.StreamingNextBlockVotingInformation, error) {
	if !bytes.HasPrefix(bz, prefixDataEncodedByCvpCodecV3) {
		return nil, fmt.Errorf("bad encoding prefix")
	}

	gzipr, err := gzip.NewReader(bytes.NewReader(bz[2:]))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gzip reader")
	}

	bzByV2, err := io.ReadAll(gzipr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read gzipped content")
	}

	err = gzipr.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close gzip reader")
	}

	return c.v2Codec.DecodeStreamingNextBlockVotingInformation(bzByV2)
}

func (c cvpCodecV3) GetVersion() CvpCodecVersion {
	return CvpCodecVersionV3
}
