package ui

import (
	"reflect"
	"slices"
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/ui/wintree"
)

// bindingsByKey maps every key a keymap binds to the names of the
// bindings that hold it.
func bindingsByKey(keyMap KeyMap) map[string][]string {
	byKey := map[string][]string{}
	value := reflect.ValueOf(keyMap)
	for i := range value.NumField() {
		binding := value.Field(i).Interface().(key.Binding)
		for _, k := range binding.Keys() {
			byKey[k] = append(byKey[k], value.Type().Field(i).Name)
		}
	}
	return byKey
}

// TestColemakKeyMap_AddsNoSharedKeys: moving keys around must not leave
// one key on two bindings, or one of them would never fire. A key the
// default map already shares is left to the default's own design.
func TestColemakKeyMap_AddsNoSharedKeys(t *testing.T) {
	defaultShared := bindingsByKey(DefaultKeyMap())
	for k, names := range bindingsByKey(ColemakKeyMap()) {
		if len(names) > 1 && len(names) != len(defaultShared[k]) {
			t.Errorf("key %q is bound to %v in the Colemak preset", k, names)
		}
	}
}

func TestColemakKeyMap_MovesTheVimKeys(t *testing.T) {
	keyMap := ColemakKeyMap()
	cases := []struct {
		binding key.Binding
		name    string
		want    string
	}{
		{keyMap.Up, "Up", "e"},
		{keyMap.Down, "Down", "n"},
		{keyMap.Left, "Left", "h"},
		{keyMap.Right, "Right", "i"},
		{keyMap.SearchNext, "SearchNext", "k"},
		{keyMap.SearchPrev, "SearchPrev", "K"},
		{keyMap.InsertMode, "InsertMode", "l"},
		{keyMap.UserProfile, "UserProfile", "P"},
	}
	for _, c := range cases {
		if !slices.Contains(c.binding.Keys(), c.want) {
			t.Errorf("%s keys = %v, want them to include %q", c.name, c.binding.Keys(), c.want)
		}
	}
}

func TestDirectionKey_FollowsThePreset(t *testing.T) {
	cases := []struct {
		name   string
		keyMap KeyMap
		press  rune
		want   string
	}{
		{"default j", DefaultKeyMap(), 'j', arrowDown},
		{"default l", DefaultKeyMap(), 'l', arrowRight},
		{"default n is not movement", DefaultKeyMap(), 'n', "n"},
		{"colemak n", ColemakKeyMap(), 'n', arrowDown},
		{"colemak e", ColemakKeyMap(), 'e', arrowUp},
		{"colemak h", ColemakKeyMap(), 'h', arrowLeft},
		{"colemak i", ColemakKeyMap(), 'i', arrowRight},
		{"colemak j is not movement", ColemakKeyMap(), 'j', "j"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := newTestApp(t)
			a.SetKeyMap(c.keyMap)
			if got := a.directionKey(keyPress(c.press)); got != c.want {
				t.Errorf("directionKey(%q) = %q, want %q", c.press, got, c.want)
			}
		})
	}
}

// TestColemak_WindowChordMovesWithI: ctrl+w i focuses the window to the
// right, as ctrl+w l does with the default keys.
func TestColemak_WindowChordMovesWithI(t *testing.T) {
	a := newTestApp(t, withSize(200, 50))
	a.SetKeyMap(ColemakKeyMap())
	leftWindow := a.focusedWin
	_ = a.splitWindow(wintree.SplitSideBySide) // focus moves to the new window
	rightWindow := a.focusedWin
	if rightWindow == leftWindow {
		t.Fatal("precondition: the split did not focus a new window")
	}

	_ = a.handleWindowChord(keyPress('h'))
	if a.focusedWin != leftWindow {
		t.Fatalf("after ctrl+w h focused %v, want the left window %v", a.focusedWin, leftWindow)
	}
	_ = a.handleWindowChord(keyPress('i'))
	if a.focusedWin != rightWindow {
		t.Errorf("after ctrl+w i focused %v, want the right window %v", a.focusedWin, rightWindow)
	}
}

// TestColemak_HelpScrollsWithNButTypesItInSearch: n moves down the help
// overlay, but once / has opened its search, n is part of the query.
func TestColemak_HelpScrollsWithNButTypesItInSearch(t *testing.T) {
	a := newTestApp(t)
	a.SetKeyMap(ColemakKeyMap())
	a.openHelpOverlay()

	_ = dispatchModeKey(a, keyPress('n'))
	if got := a.help.Selected(); got != 1 {
		t.Fatalf("after n the help row is %d, want 1", got)
	}

	_ = dispatchModeKey(a, keyPress('/'))
	_ = dispatchModeKey(a, keyPress('n'))
	if got := a.help.Selected(); got != 0 {
		t.Errorf("n while searching moved the row to %d; want it typed into the query", got)
	}
}

func TestColemak_ProfileClosesWithPNotK(t *testing.T) {
	a := openUserProfileWith(t)
	a.SetKeyMap(ColemakKeyMap())

	_ = dispatchModeKey(a, keyPress('K'))
	if !a.userProfile.IsVisible() {
		t.Fatal("K closed the profile; on Colemak K is previous match, P is the profile key")
	}
	_ = dispatchModeKey(a, keyPress('P'))
	if a.userProfile.IsVisible() || a.mode != ModeNormal {
		t.Errorf("P left the profile open (visible %v, mode %v)", a.userProfile.IsVisible(), a.mode)
	}
}

func TestConfirmPromptKeyMap_FollowsLeftAndRight(t *testing.T) {
	cases := []struct {
		name           string
		keyMap         KeyMap
		next, previous tea.KeyMsg
	}{
		{"default", DefaultKeyMap(), keyPress('l'), keyPress('h')},
		{"colemak", ColemakKeyMap(), keyPress('i'), keyPress('h')},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			promptKeys := confirmPromptKeyMap(c.keyMap)
			if !key.Matches(c.next, promptKeys.Next) || !key.Matches(keyCode(tea.KeyTab), promptKeys.Next) {
				t.Errorf("Next keys = %v", promptKeys.Next.Keys())
			}
			if !key.Matches(c.previous, promptKeys.Prev) {
				t.Errorf("Prev keys = %v", promptKeys.Prev.Keys())
			}
		})
	}
}
