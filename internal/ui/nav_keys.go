// internal/ui/nav_keys.go
//
// The keymap preset (default or Colemak) and the arrow-key translation
// that lets dialogs and chords follow it.
package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/bubbles/confirmprompt"
)

// The key names dialogs and chords use for the four directions.
const (
	arrowUp    = "up"
	arrowDown  = "down"
	arrowLeft  = "left"
	arrowRight = "right"
)

// Keys that also move between a confirm prompt's buttons, whatever the
// keymap says left and right are.
const (
	nextButtonKey     = "tab"
	previousButtonKey = "shift+tab"
)

// SetKeyMap installs the key bindings normal mode, the help overlay and
// the dialogs use. cmd/slk calls it with the preset the config names.
func (a *App) SetKeyMap(keyMap KeyMap) {
	a.keys = keyMap
	a.confirmPrompt.KeyMap = confirmPromptKeyMap(keyMap)
}

// confirmPromptKeyMap is the confirm prompt's default bindings, with its
// button movement taken from the keymap's left and right.
func confirmPromptKeyMap(keyMap KeyMap) confirmprompt.KeyMap {
	promptKeys := confirmprompt.DefaultKeyMap()
	promptKeys.Next = key.NewBinding(key.WithKeys(append(keyMap.Right.Keys(), nextButtonKey)...))
	promptKeys.Prev = key.NewBinding(key.WithKeys(append(keyMap.Left.Keys(), previousButtonKey)...))
	return promptKeys
}

// directionKey names the arrow a key press stands for under the active
// keymap: "up", "down", "left" or "right" when it matches a movement
// binding, and the key as typed otherwise. Dialogs and chords that only
// understand arrows go through it, so they move with h/n/e/i on Colemak
// as they do with h/j/k/l by default.
func (a *App) directionKey(msg tea.KeyMsg) string {
	switch {
	case key.Matches(msg, a.keys.Up):
		return arrowUp
	case key.Matches(msg, a.keys.Down):
		return arrowDown
	case key.Matches(msg, a.keys.Left):
		return arrowLeft
	case key.Matches(msg, a.keys.Right):
		return arrowRight
	}
	return msg.String()
}
