package hotkeytest

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/hotkeytest/internal/dll"
	"syscall"
	"unicode/utf16"
)

func CreateWindow(title string, opt *w32.WindowOptions) (hwnd w32.HWND, unRegisterClassFunc func(), err error) {
	hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))

	if opt.ClassName == "" {
		opt.ClassName = "example"
	}

	var hIcon w32.HANDLE
	if opt.IconPath != "" {
		hIcon, _ = dll.User.LoadImage(0, // hInstance must be NULL when loading from a file
			opt.IconPath,
			w32.IMAGE_ICON, 0, 0, w32.LR_LOADFROMFILE|w32.LR_DEFAULTSIZE|w32.LR_SHARED)
	}

	if atom, errno := dll.User.RegisterClass(&w32.WNDCLASS{
		Style:         opt.ClassStyle,
		HbrBackground: w32.COLOR_WINDOW,
		WndProc:       syscall.NewCallback(opt.WndProc),
		HInstance:     hInstance,
		HIcon:         w32.HICON(hIcon),
		ClassName:     &utf16.Encode([]rune(opt.ClassName + "\x00"))[0],
	}); atom == 0 {
		return 0, nil, fmt.Errorf("[RegisterClass Error] %w", errno)
	}

	width := opt.Width
	if width == 0 {
		width = w32.CW_USEDEFAULT
	}
	height := opt.Height
	if height == 0 {
		height = w32.CW_USEDEFAULT
	}
	posX := opt.X
	if posX == 0 {
		posX = w32.CW_USEDEFAULT
	}
	posY := opt.Y
	if posY == 0 {
		posY = w32.CW_USEDEFAULT
	}

	if opt.Style == 0 {
		opt.Style = w32.WS_OVERLAPPEDWINDOW
	}

	// Create window
	hwnd, errno := dll.User.CreateWindowEx(
		w32.DWORD(opt.ExStyle),
		opt.ClassName,
		title,
		w32.DWORD(opt.Style),

		// Size and position
		posX, posY, width, height,

		0, // Parent window
		0, // Menu
		hInstance,
		0, // Additional application data
	)

	unRegisterClassFunc = func() {
		if errno2 := dll.User.UnregisterClass(opt.ClassName, hInstance); errno2 != 0 {
			fmt.Printf("Error UnregisterClass: %s", errno2)
		}
	}

	if errno != 0 {
		unRegisterClassFunc()
		return 0, nil, errno
	}
	return hwnd, unRegisterClassFunc, nil
}
