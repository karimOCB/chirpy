package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{name: "well formed auth", headers: http.Header{"Authorization": {"Bearer abc123"}}, want: "abc123", wantErr: false},
		{name: "missing auth header", headers: http.Header{}, want: "", wantErr: true},
		{name: "malformed header", headers: http.Header{"Authorization": {"Basic abc123"}}, want: "", wantErr: true},
		{name: "well formed auth whitespace case", headers: http.Header{"Authorization": {" Bearer  abc123  "}}, want: "abc123", wantErr: false},
		{name: "missing token", headers: http.Header{"Authorization": {" Bearer  "}}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken error = %q, wantErr %v", err, tt.wantErr)
				return
			}

			if tokenString != tt.want {
				t.Errorf("GetBearerToken() = %q, want %q", tokenString, tt.want)
			}
		})
	}

}
