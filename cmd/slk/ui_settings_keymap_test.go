package main

import (
	"slices"
	"testing"

	"github.com/gammons/slk/internal/config"
)

func TestKeyMapForPreset(t *testing.T) {
	cases := []struct {
		preset   string
		wantDown string
	}{
		{config.KeyMapDefault, "j"},
		{"", "j"},
		{config.KeyMapColemak, "n"},
	}
	for _, c := range cases {
		if keys := keyMapForPreset(c.preset).Down.Keys(); !slices.Contains(keys, c.wantDown) {
			t.Errorf("preset %q: Down keys = %v, want %q among them", c.preset, keys, c.wantDown)
		}
	}
}
