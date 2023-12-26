package constants

//goland:noinspection GoSnakeCaseUsage
const (
	STREAMING_BASE_URL       = "https://cvp.bcdev.tools"
	STREAMING_BASE_URL_LOCAL = "http://localhost:8080" // for development purpose only

	STREAMING_PATH_REGISTER_PRE_VOTE          = "register-session/pre-vote/:chainId"
	STREAMING_PATH_RESUME_PRE_VOTE            = "resume-session/pre-vote/:sessionId"
	STREAMING_PATH_BROADCAST_PRE_VOTE         = "broadcast/pre-vote/:sessionId"
	STREAMING_PATH_VIEW_PRE_VOTE              = "pvtop/:sessionId"
	STREAMING_PATH_VIEW_PRE_VOTE_FETCH_UPDATE = "pvtop/:sessionId/update"

	STREAMING_CONTENT_TYPE       = "application/octet-stream"
	STREAMING_HEADER_SESSION_KEY = "X-Session-Key"
)

//goland:noinspection GoSnakeCaseUsage
const (
	MAX_VALIDATORS                             = 250
	MAX_ENCODED_LIGHT_VALIDATORS_BYTES         = 12251 // 12251 v1, 8251 v2
	MAX_ENCODED_NEXT_BLOCK_PRE_VOTE_INFO_BYTES = 2044  // 2044 v1, 1786 v2
)
