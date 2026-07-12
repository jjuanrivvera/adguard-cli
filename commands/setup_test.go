package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeSecret(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips bracketed-paste markers",
			in:   "\x1b[200~KEY\x1b[201~\n",
			want: "KEY",
		},
		{
			name: "clean secret is unchanged",
			in:   "KEY",
			want: "KEY",
		},
		{
			name: "trims surrounding whitespace",
			in:   "  KEY  ",
			want: "KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sanitizeSecret(tt.in))
		})
	}
}
