package types

type StreamingLightValidators []StreamingLightValidator

type StreamingLightValidator struct {
	Index                     int     `json:"i"`
	VotingPowerDisplayPercent float64 `json:"vdp"`
	Moniker                   string  `json:"m"`
}
