//go:build windows

package windows

import "golang.org/x/sys/windows"

var (
	virtualDesktopAccessor        = windows.MustLoadDLL("lib/windows/VirtualDesktopAccessor.dll")
	getCurrentDesktopNumber       = virtualDesktopAccessor.MustFindProc("GetCurrentDesktopNumber")
	getDesktopCount               = virtualDesktopAccessor.MustFindProc("GetDesktopCount")
	moveWindowToDesktopNumber     = virtualDesktopAccessor.MustFindProc("MoveWindowToDesktopNumber")
	goToDesktopNumber             = virtualDesktopAccessor.MustFindProc("GoToDesktopNumber")
	restartVirtualDesktopAccessor = virtualDesktopAccessor.MustFindProc("RestartVirtualDesktopAccessor")
)

func GetCurrentDesktopNumber() int {
	ret, _, _ := getCurrentDesktopNumber.Call()
	return int(ret)
}

func GetDesktopCount() int {
	ret, _, _ := getDesktopCount.Call()
	return int(ret)
}

func MoveWindowToDesktopNumber(handle uintptr, number int) error {
	_, _, err := moveWindowToDesktopNumber.Call(handle, uintptr(number))
	if err == windows.ERROR_SUCCESS {
		return nil
	}
	return err
}

func GoToDesktopNumber(number int) error {
	_, _, err := goToDesktopNumber.Call(uintptr(number))
	if err == windows.ERROR_SUCCESS {
		return nil
	}
	return err
}

func RestartVirtualDesktopAccessor() {
	restartVirtualDesktopAccessor.Call()
}
