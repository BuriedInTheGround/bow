//go:build windows

package windows

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32              = windows.NewLazySystemDLL("User32.dll")
	findWindow          = user32.NewProc("FindWindowW")
	getActiveWindow     = user32.NewProc("GetActiveWindow")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	setActiveWindow     = user32.NewProc("SetActiveWindow")
)

func FindWindow(class string) (uintptr, error) {
	class16, err := windows.UTF16PtrFromString(class)
	if err != nil {
		return 0, fmt.Errorf("invalid class value: %q: %w", class, err)
	}
	ret, _, err := findWindow.Call(
		uintptr(unsafe.Pointer(class16)),
		0, // Do not set a value for lpWindowName.
	)
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func GetActiveWindow() (uintptr, error) {
	ret, _, err := getActiveWindow.Call()
	if err == windows.ERROR_SUCCESS {
		return ret, nil
	}
	return 0, err
}

func GetForegroundWindow() (uintptr, error) {
	ret, _, err := getForegroundWindow.Call()
	if err == windows.ERROR_SUCCESS {
		return ret, nil
	}
	return 0, err
}

func SetActiveWindow(handle uintptr) (uintptr, error) {
	ret, _, err := setActiveWindow.Call(handle)
	if err == windows.ERROR_SUCCESS {
		return ret, nil
	}
	return 0, err
}
