// internal/ui/mode_user_profile.go
//
// User-profile-dialog key handler. The profile key (K by default, P
// on Colemak), esc and q all close the modal and return to Normal
// mode; e copies the loaded email and c copies
// the loaded huddle link (the same actions a click on the 📋 beside
// each synthesises). Every other key is swallowed, since the dialog is
// read-only and has no scroll state of its own.
package ui

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func handleUserProfileMode(a *App, msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, a.keys.UserProfile) {
		closeUserProfile(a)
		return nil
	}
	switch msg.String() {
	case "esc", "q":
		closeUserProfile(a)
	case "e":
		return copyProfileEmail(a)
	case "c":
		return copyProfileHuddle(a)
	}
	return nil
}

// closeUserProfile closes the dialog and returns to normal mode.
func closeUserProfile(a *App) {
	a.userProfile.Close()
	a.SetMode(ModeNormal)
}

// copyProfileEmail writes the dialog's loaded email to the clipboard
// and toasts the outcome. The dialog stays open.
func copyProfileEmail(a *App) tea.Cmd {
	email := a.userProfile.Email()
	if email == "" {
		return toastWithClear(a, "No email to copy", 2*time.Second)
	}
	return tea.Batch(
		a.clipboardWrite(email),
		toastWithClear(a, "Copied email", 2*time.Second),
	)
}

// copyProfileHuddle writes the dialog's huddle deep link to the
// clipboard and toasts the outcome. The dialog stays open. Nothing is
// offered when the profile did not name a huddle channel, or when the
// live status says the huddle has ended — the same conditions that
// hide the row's 📋, so a cached profile never copies a dead link.
func copyProfileHuddle(a *App) tea.Cmd {
	if !a.userProfile.HuddleLinkOffered() {
		return toastWithClear(a, "No huddle link", 2*time.Second)
	}
	return tea.Batch(
		a.clipboardWrite(a.userProfile.HuddleURL()),
		toastWithClear(a, "Copied huddle link", 2*time.Second),
	)
}
