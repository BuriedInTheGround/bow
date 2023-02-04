//go:build windows

package main

import (
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
	"interrato.dev/bow"
)

func main() {
	mainthread.Init(run)
}

func run() {
	if n := bow.DesktopCount(); n > 10 {
		errorWithHint("too many virtual desktops",
			"a maximum of ten virtual desktops is allowed",
			"you may want to press Win+Tab to remove some virtual desktops",
		)
	}

	// We register all hotkeys for ten virtual desktops even if less than ten
	// virtual desktops are available. This is to avoid confusion that may
	// arise by having hotkeys that are unregistered only for some numbers.
	hotkeys := make(map[string]*hotkey.Hotkey)
	for i := 1; i < 10; i++ {
		hotkeys["focus-desktop-"+strconv.Itoa(i)] = hotkey.New([]hotkey.Modifier{hotkey.ModWin}, hotkey.Key(0x30+i))
		hotkeys["send-to-desktop-"+strconv.Itoa(i)] = hotkey.New([]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift}, hotkey.Key(0x30+i))
	}
	hotkeys["focus-desktop-10"] = hotkey.New([]hotkey.Modifier{hotkey.ModWin}, hotkey.Key0)
	hotkeys["send-to-desktop-10"] = hotkey.New([]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift}, hotkey.Key0)

	withExplorerRestart(func() {
		for action, hotkey := range hotkeys {
			err := hotkey.Register()
			if err != nil {
				warningf("failed to register hotkey for %q: %v", action, err)
			}
		}
	})

	var (
		actions []string
		cases   []reflect.SelectCase
	)
	for purpose, hotkey := range hotkeys {
		actions = append(actions, purpose)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(hotkey.Keydown()),
		})
	}

	for {
		chosen, _, _ := reflect.Select(cases)

		if strings.HasPrefix(actions[chosen], "focus-desktop") {
			n, err := strconv.Atoi(strings.TrimPrefix(actions[chosen], "focus-desktop-"))
			if err != nil {
				errorf("action %q is invalid", actions[chosen])
			}
			if err := bow.FocusDesktop(n); err != nil {
				warningf("%v", err)
			}
		}

		if strings.HasPrefix(actions[chosen], "send-to-desktop") {
			n, err := strconv.Atoi(strings.TrimPrefix(actions[chosen], "send-to-desktop-"))
			if err != nil {
				errorf("action %q is invalid", actions[chosen])
			}
			if err := bow.SendToDesktop(n); err != nil {
				warningf("%v", err)
			}
		}
	}
}

func withExplorerRestart(f func()) {
	cmd := exec.Command("TASKKILL", "/IM", "explorer.exe", "/F")
	if err := cmd.Run(); err != nil {
		errorf("failed to run `TASKKILL /IM explorer.exe /F`: %v", err)
	}

	f()

	cmd = exec.Command("explorer.exe")
	if err := cmd.Start(); err != nil {
		errorf("failed to start `explorer.exe`: %v", err)
	}
	defer func() {
		// This call to Wait() will hang until explorer.exe is killed, so we
		// run it on another goroutine to be able to continue the program
		// execution.
		go cmd.Wait()
	}()
}

// l is a logger with no prefixes.
var l = log.New(os.Stderr, "", 0)

func printf(format string, v ...any) {
	l.Printf("bow: "+format, v...)
}

func errorf(format string, v ...any) {
	l.Fatalf("bow: error: "+format, v...)
}

func warningf(format string, v ...any) {
	l.Printf("bow: warning: "+format, v...)
}

func errorWithHint(error string, hints ...string) {
	l.Printf("bow: error: %s", error)
	for _, hint := range hints {
		l.Printf("bow: hint: %s", hint)
	}
	os.Exit(1)
}
