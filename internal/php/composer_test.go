package php

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectVersionFromComposer(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		want     string
	}{
		{
			name:     "platform php pinned",
			contents: `{"config":{"platform":{"php":"8.3"}}}`,
			want:     "8.3",
		},
		{
			name:     "platform php with patch",
			contents: `{"config":{"platform":{"php":"8.2.15"}}}`,
			want:     "8.2",
		},
		{
			name:     "platform takes precedence over require",
			contents: `{"config":{"platform":{"php":"8.3"}},"require":{"php":"~8.2"}}`,
			want:     "8.3",
		},
		{
			name:     "require tilde",
			contents: `{"require":{"php":"~8.3"}}`,
			want:     "8.3",
		},
		{
			name:     "require caret",
			contents: `{"require":{"php":"^8.2"}}`,
			want:     "8.2",
		},
		{
			name:     "require gte",
			contents: `{"require":{"php":">=8.2"}}`,
			want:     "8.2",
		},
		{
			name:     "require wildcard",
			contents: `{"require":{"php":"8.4.*"}}`,
			want:     "8.4",
		},
		{
			name:     "require range picks first bound",
			contents: `{"require":{"php":">=8.2 <8.4"}}`,
			want:     "8.2",
		},
		{
			name:     "unsupported version ignored",
			contents: `{"require":{"php":"~7.4"}}`,
			want:     "",
		},
		{
			name:     "no php entry",
			contents: `{"require":{"ext-intl":"*"}}`,
			want:     "",
		},
		{
			name:     "empty object",
			contents: `{}`,
			want:     "",
		},
		{
			name:     "invalid json",
			contents: `{not json`,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "composer.json")
			if err := os.WriteFile(path, []byte(tt.contents), 0o644); err != nil {
				t.Fatalf("write composer.json: %v", err)
			}
			if got := DetectVersionFromComposer(path); got != tt.want {
				t.Errorf("DetectVersionFromComposer() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectVersionFromComposer_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "composer.json")
	if got := DetectVersionFromComposer(path); got != "" {
		t.Errorf("DetectVersionFromComposer() = %q, want empty for missing file", got)
	}
}
