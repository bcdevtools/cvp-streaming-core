package utils

import (
	"github.com/bcdevtools/cvp-streaming-core/constants"
	"strings"
)

func GetRemoteUrlRegisterPreVoteStreamingSession(baseUrl, chainId string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/" + strings.ReplaceAll(constants.STREAMING_PATH_REGISTER_PRE_VOTE, ":chainId", chainId)
}

func GetRemoteUrlResumePreVoteStreamingSession(baseUrl, sessionId string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/" + strings.ReplaceAll(constants.STREAMING_PATH_RESUME_PRE_VOTE, ":sessionId", sessionId)
}

func GetRemoteUrlBroadcastPreVoteDuringStreamingSession(baseUrl, sessionId string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/" + strings.ReplaceAll(constants.STREAMING_PATH_BROADCAST_PRE_VOTE, ":sessionId", sessionId)
}

func GetPublicUrlViewPreVoteStreamingSession(baseUrl, sessionId string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/" + strings.ReplaceAll(constants.STREAMING_PATH_VIEW_PRE_VOTE, ":sessionId", sessionId)
}

func GetUrlFetchPreVoteStreamingSessionUpdate(baseUrl, sessionId string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/" + strings.ReplaceAll(constants.STREAMING_PATH_VIEW_PRE_VOTE_FETCH_UPDATE, ":sessionId", sessionId)
}
