// go build -ldflags "-s -w" -o test.exe --pkgdir=../..

package main

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/hotkeytest"
	"github.com/CarsonSlovoka/hotkeytest/internal/dll"
	"log"
)

const (
	HotkeyLocalF3 = w32.WM_APP + iota
	HotKeyGlobalF4
	HotKeyGlobalF6

	WMCommandM
)

func main() {
	opt := &w32.WindowOptions{Width: 100, Height: 100,
		ClassName: "hotkeyTest",
	}

	opt.WndProc = func(hwnd w32.HWND, uMsg uint32, wParam w32.WPARAM, lParam w32.LPARAM) uintptr {
		switch uMsg {
		case w32.WM_CREATE:
			dll.User.ShowWindow(hwnd, w32.SW_SHOW)
			if err := dll.User.RegisterHotKey(hwnd, HotkeyLocalF3, w32.MOD_CONTROL, w32.VK_F3); err != 0 {
				log.Println(err)
			}
			if err := dll.User.RegisterHotKey(0, HotKeyGlobalF4, w32.MOD_CONTROL, w32.VK_F4); err != 0 {
				log.Println(err)
			}

			// 定義全域熱鍵，此熱鍵沒辦法被WndProc定義WM_HOTKEY來收到消息，一定要在MsgLoop來捕獲
			if err := dll.User.RegisterHotKey(0, HotKeyGlobalF6, w32.MOD_WIN, w32.VK_F6); err != 0 {
				log.Println(err)
			}
		case w32.WM_DESTROY:
			for _, hotkeyID := range []int32{HotkeyLocalF3} {
				if en := dll.User.UnregisterHotKey(hwnd, hotkeyID); en != 0 {
					log.Printf("Error [UnregisterHotKey] %s", en)
				}
			}
			for _, hotkeyID := range []int32{HotKeyGlobalF4, HotKeyGlobalF6} {
				if en := dll.User.UnregisterHotKey(0, hotkeyID); en != 0 {
					log.Printf("Error [(global) UnregisterHotKey] %s", en)
				}
			}
			dll.User.PostQuitMessage(0)
			return 0
		case w32.WM_HOTKEY: // local 熱鍵
			log.Println("local 熱鍵")
			switch wParam {
			case HotKeyGlobalF4:
				log.Println("永遠無法顯示到此訊息")
			case HotkeyLocalF3:
				log.Println("Ctrl+F3已經被按下")
			}
			return 1
		case WMCommandM:
			log.Println("Command+F6已經被按下")
			return 1
		}
		return uintptr(dll.User.DefWindowProc(hwnd, w32.UINT(uMsg), wParam, lParam))
	}

	hwnd, unRegisterClassFunc, err := hotkeytest.CreateWindow("test command key", opt)
	if err != nil {
		log.Fatal(err)
	}
	defer unRegisterClassFunc()

	fmt.Println("hwnd:", hwnd)

	var msg w32.MSG
	for {
		if bRet, eno := dll.User.GetMessage(&msg, 0, 0, 0); bRet == 0 {
			// WM_QUIT
			break
		} else if bRet == -1 {
			log.Printf("GetMessage error:%s\n", eno)
		} else {
			if doTranslate := MsgLoop(hwnd, &msg); !doTranslate {
				continue
			}
			dll.User.TranslateMessage(&msg)
			dll.User.DispatchMessage(&msg)
		}
	}
}

func MsgLoop(hwnd w32.HWND, msg *w32.MSG) bool {
	if msg.Message == w32.WM_HOTKEY {
		log.Println("WM_HOTKEY received")
		switch msg.WParam {
		case HotKeyGlobalF6:
			_ = dll.User.PostMessage(hwnd, WMCommandM, 0, 0)
			return false
		}
	}
	return true
}
