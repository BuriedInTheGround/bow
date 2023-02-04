//go:build windows

package bow

import "interrato.dev/bow/internal/windows"

func FocusDesktop(num int) error {
	if num <= 0 {
		panic("bow: desktop number must be greater than 0")
	}
	windows.RestartVirtualDesktopAccessor()
	return windows.GoToDesktopNumber(num - 1)
}

func SendToDesktop(num int) error {
	if num <= 0 {
		panic("bow: desktop number must be greater than 0")
	}
	h, err := windows.GetForegroundWindow()
	if err != nil {
		return err
	}
	windows.RestartVirtualDesktopAccessor()
	return windows.MoveWindowToDesktopNumber(h, num-1)
}

func DesktopCount() int {
	return windows.GetDesktopCount()
}
