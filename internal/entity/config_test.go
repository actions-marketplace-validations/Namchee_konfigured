package entity

import (
	"os"
	"testing"

	"github.com/Namchee/konfigured/internal/constant"
	"github.com/stretchr/testify/assert"
)

func TestCreateConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    *Configuration
		wantErr error
	}{
		{
			name:    "missing access token",
			env:     map[string]string{},
			want:    nil,
			wantErr: constant.ErrMissingToken,
		},
		{
			name: "invalid pattern",
			env: map[string]string{
				"INPUT_TOKEN":   "access-token",
				"INPUT_NEWLINE": "true",
				"INPUT_INCLUDE": "\\",
			},
			want:    nil,
			wantErr: constant.ErrInvalidGlob,
		},
		{
			name: "default include pattern",
			env: map[string]string{
				"INPUT_TOKEN":   "access-token",
				"INPUT_NEWLINE": "true",
			},
			want: &Configuration{
				Token:   "access-token",
				Newline: true,
				Include: defaultPattern,
			},
			wantErr: nil,
		},
		{
			name: "success",
			env: map[string]string{
				"INPUT_TOKEN":   "access-token",
				"INPUT_NEWLINE": "true",
				"INPUT_INCLUDE": "**/*.{json,ini}",
			},
			want: &Configuration{
				Token:   "access-token",
				Newline: true,
				Include: "**/*.{json,ini}",
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			defer os.Clearenv()

			got, err := CreateConfiguration()

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
