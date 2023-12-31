package codec

//goland:noinspection SpellCheckingInspection
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// DetectEncodingVersion will try to detect the encoding version of the given byte array based on the very first bytes.
// The returned version is 'possible' because it is not guaranteed to be the correct version without actual decode it.
func DetectEncodingVersion(bz []byte) (possible CvpCodecVersion, detected bool) {
	if bytes.HasPrefix(bz, prefixDataEncodedByCvpCodecV3) {
		return CvpCodecVersionV3, true
	}
	if bytes.HasPrefix(bz, prefixDataEncodedByCvpCodecV2) {
		return CvpCodecVersionV2, true
	}
	if bytes.HasPrefix(bz, []byte(prefixDataEncodedByCvpCodecV1)) {
		return CvpCodecVersionV1, true
	}
	return CvpCodecVersionUnknown, false
}

func toUint16Buffer(num int) []byte {
	if num < 0 || num > math.MaxUint16 {
		panic(fmt.Errorf("overflow uint16: %d", num))
	}
	bz := make([]byte, 2)
	binary.BigEndian.PutUint16(bz, uint16(num))
	return bz
}

func fromUint16Buffer(bz []byte) int {
	if len(bz) != 2 {
		panic(fmt.Errorf("invalid uint16 buffer length: %d, require 2", len(bz)))
	}
	return int(binary.BigEndian.Uint16(bz))
}

func toPercentBuffer(percent float64) []byte {
	if percent < 0 || percent > 100 {
		panic(fmt.Errorf("overflow percent: %f", percent))
	}
	var pi, pf byte
	str := fmt.Sprintf("%.2f", percent)
	parts := strings.SplitN(str, ".", 2)
	ipi, _ := strconv.Atoi(parts[0])
	ipf, _ := strconv.Atoi(parts[1])
	pi = byte(ipi)
	pf = byte(ipf)
	return []byte{pi, pf}
}

func fromPercentBuffer(bz []byte) float64 {
	if len(bz) != 2 {
		panic(fmt.Errorf("invalid percent buffer length: %d, require 2", len(bz)))
	}
	return float64(bz[0]) + float64(bz[1])/100
}

func tryTakeNBytesFrom(bz []byte, fromIndex, size int) ([]byte, bool) {
	if size < 1 {
		panic("invalid size")
	}
	if fromIndex < 0 {
		panic("invalid beginning index")
	}
	if fromIndex+size > len(bz) {
		return nil, false
	}
	return bz[fromIndex : fromIndex+size], true
}

func mustTakeNBytesFrom(bz []byte, fromIndex, size int) []byte {
	bz, ok := tryTakeNBytesFrom(bz, fromIndex, size)
	if !ok {
		panic(fmt.Errorf("failed to take %d bytes from %d", size, fromIndex))
	}
	return bz
}

func takeUntilSeparatorOrEnd(bz []byte, fromIndex int, separator byte) (taken []byte) {
	for i := fromIndex; i < len(bz); i++ {
		if bz[i] == separator {
			break
		}
		taken = append(taken, bz[i])
	}
	return
}

func sanitizeMoniker(moniker string) string {
	moniker = strings.ReplaceAll(moniker, "<", "(")
	moniker = strings.ReplaceAll(moniker, ">", ")")
	moniker = strings.ReplaceAll(moniker, "'", "`")
	moniker = strings.ReplaceAll(moniker, "\"", "`")
	return moniker
}
