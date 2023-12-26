package utils

import "testing"

func TestGetRemoteUrlRegisterPreVoteStreamingSession(t *testing.T) {
	tests := []struct {
		name    string
		baseUrl string
		chainId string
		want    string
	}{
		{
			name:    "normal",
			baseUrl: "https://cvp.bcdev.tools",
			chainId: "cosmoshub-4",
			want:    "https://cvp.bcdev.tools/register-session/pre-vote/cosmoshub-4",
		},
		{
			name:    "normal with suffix slash",
			baseUrl: "http://localhost:8080/",
			chainId: "cosmoshub-4",
			want:    "http://localhost:8080/register-session/pre-vote/cosmoshub-4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRemoteUrlRegisterPreVoteStreamingSession(tt.baseUrl, tt.chainId); got != tt.want {
				t.Errorf("GetRemoteUrlRegisterPreVoteStreamingSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteUrlResumePreVoteStreamingSession(t *testing.T) {
	tests := []struct {
		name      string
		baseUrl   string
		sessionId string
		want      string
	}{
		{
			name:      "normal",
			baseUrl:   "https://cvp.bcdev.tools",
			sessionId: "sample-session-id-1",
			want:      "https://cvp.bcdev.tools/resume-session/pre-vote/sample-session-id-1",
		},
		{
			name:      "normal with suffix slash",
			baseUrl:   "http://localhost:8080/",
			sessionId: "sample-session-id-2",
			want:      "http://localhost:8080/resume-session/pre-vote/sample-session-id-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRemoteUrlResumePreVoteStreamingSession(tt.baseUrl, tt.sessionId); got != tt.want {
				t.Errorf("GetRemoteUrlResumePreVoteStreamingSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteUrlBroadcastPreVoteDuringStreamingSession(t *testing.T) {
	tests := []struct {
		name      string
		baseUrl   string
		sessionId string
		want      string
	}{
		{
			name:      "normal",
			baseUrl:   "https://cvp.bcdev.tools",
			sessionId: "sample-session-id-1",
			want:      "https://cvp.bcdev.tools/broadcast/pre-vote/sample-session-id-1",
		},
		{
			name:      "normal with suffix slash",
			baseUrl:   "http://localhost:8080/",
			sessionId: "sample-session-id-2",
			want:      "http://localhost:8080/broadcast/pre-vote/sample-session-id-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRemoteUrlBroadcastPreVoteDuringStreamingSession(tt.baseUrl, tt.sessionId); got != tt.want {
				t.Errorf("GetRemoteUrlBroadcastPreVoteDuringStreamingSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicUrlViewPreVoteStreamingSession(t *testing.T) {
	tests := []struct {
		name      string
		baseUrl   string
		sessionId string
		want      string
	}{
		{
			name:      "normal",
			baseUrl:   "https://cvp.bcdev.tools",
			sessionId: "sample-session-id-1",
			want:      "https://cvp.bcdev.tools/pvtop/sample-session-id-1",
		},
		{
			name:      "normal with suffix slash",
			baseUrl:   "http://localhost:8080/",
			sessionId: "sample-session-id-2",
			want:      "http://localhost:8080/pvtop/sample-session-id-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPublicUrlViewPreVoteStreamingSession(tt.baseUrl, tt.sessionId); got != tt.want {
				t.Errorf("GetPublicUrlViewPreVoteStreamingSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUrlFetchPreVoteStreamingSessionUpdate(t *testing.T) {
	tests := []struct {
		name      string
		baseUrl   string
		sessionId string
		want      string
	}{
		{
			name:      "normal",
			baseUrl:   "https://cvp.bcdev.tools",
			sessionId: "sample-session-id-1",
			want:      "https://cvp.bcdev.tools/pvtop/sample-session-id-1/update",
		},
		{
			name:      "normal with suffix slash",
			baseUrl:   "http://localhost:8080/",
			sessionId: "sample-session-id-2",
			want:      "http://localhost:8080/pvtop/sample-session-id-2/update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUrlFetchPreVoteStreamingSessionUpdate(tt.baseUrl, tt.sessionId); got != tt.want {
				t.Errorf("GetUrlFetchPreVoteStreamingSessionUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
