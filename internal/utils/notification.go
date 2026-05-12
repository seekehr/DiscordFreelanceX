package utils

import "fyne.io/fyne/v2"

// Notify sends an OS-level toast notification (Windows).
func Notify(a fyne.App, title, body string) {
	a.SendNotification(&fyne.Notification{
		Title:   title,
		Content: body,
	})
}
