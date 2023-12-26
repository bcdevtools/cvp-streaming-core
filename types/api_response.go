package types

type PreVoteStreamingSessionRegistrationResponse struct {
	SessionId  PreVoteStreamingSessionId  `json:"session-id"`
	SessionKey PreVoteStreamingSessionKey `json:"session-key"`
}
