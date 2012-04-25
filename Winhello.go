package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"syscall"
	"unsafe"
)

// Windows functions
var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	GetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	GetConsoleWindow = kernel32.NewProc("GetConsoleWindow")

	user32           = syscall.NewLazyDLL("user32.dll")
	LoadIcon         = user32.NewProc("LoadIconW")
	LoadCursor       = user32.NewProc("LoadCursorW")
	RegisterClassEx  = user32.NewProc("RegisterClassExW")
	CreateWindowEx   = user32.NewProc("CreateWindowExW")
	ShowWindow       = user32.NewProc("ShowWindow")
	UpdateWindow     = user32.NewProc("UpdateWindow")
	PostQuitMessage  = user32.NewProc("PostQuitMessage")
	DefWindowProc    = user32.NewProc("DefWindowProcW")
	GetMessage       = user32.NewProc("GetMessageW")
	PostMessage      = user32.NewProc("PostMessageW")
	TranslateMessage = user32.NewProc("TranslateMessage")
	DispatchMessage  = user32.NewProc("DispatchMessageW")
	BeginPaint       = user32.NewProc("BeginPaint")
	EndPaint         = user32.NewProc("EndPaint")
	CreateMenu       = user32.NewProc("CreateMenu")
	CreatePopupMenu  = user32.NewProc("CreatePopupMenu")
	AppendMenu       = user32.NewProc("AppendMenuW")
	SetMenu          = user32.NewProc("SetMenu")

	gdi32          = syscall.NewLazyDLL("gdi32.dll")
	GetStockObject = gdi32.NewProc("GetStockObject")
	TextOut        = gdi32.NewProc("TextOutW")
)

// Resource definition
const (
	ID_FILE_EXIT = 100
)

// If program was linked with -Hwindowsgui there won't be a console window to display a panic
func errorHandler() {
	console, _, _ := GetConsoleWindow.Call()
	if console == 0 {
		p := recover()
		if p != nil {
			pf, _ := os.Create("panic.txt")
			fmt.Fprintln(pf, "panic: ", p, "\n")
			fmt.Fprintln(pf, string(debug.Stack()))
			pf.Close()
		}
	}
}

func main() {
	defer errorHandler()

	hInstance, _, _ := GetModuleHandle.Call(0)

	// Combine all args except the program name into a string
	szCmdLine := ""
	for i := 1; i < len(os.Args); i++ {
		szCmdLine = szCmdLine + os.Args[i] + " "
	}
	szCmdLine = strings.TrimSpace(szCmdLine)

	startupInfo := new(syscall.StartupInfo)
	_ = syscall.GetStartupInfo(startupInfo)

	exitCode := WinMain(hInstance, 0, szCmdLine, int(startupInfo.ShowWindow))
	os.Exit(exitCode)
}

func WinMain(hInstance, hPrevInstance uintptr, szCmdLine string, iCmdShow int) int {

	var (
		lastError   error
		returnValue uintptr
	)

	// Initialise window class
	wndClass := new(WNDCLASSEX)
	wndClass.cbSize = uint(unsafe.Sizeof(*wndClass))
	wndClass.style = CS_HREDRAW | CS_VREDRAW
	wndClass.lpfnWndProc = syscall.NewCallback(WndProc)
	wndClass.cbClsExtra = 0
	wndClass.cbWndExtra = 0
	wndClass.hInstance = hInstance
	wndClass.hIcon, _, lastError = LoadIcon.Call(0, uintptr(IDI_APPLICATION))
	if wndClass.hIcon == 0 {
		panic(lastError.Error())
	}
	wndClass.hCursor, _, lastError = LoadCursor.Call(0, uintptr(IDC_ARROW))
	if wndClass.hCursor == 0 {
		panic(lastError.Error())
	}
	wndClass.hbrBackground, _, _ = GetStockObject.Call(WHITE_BRUSH)
	if wndClass.hbrBackground == 0 {
		panic("Failed")
	}
	wndClass.lpszClassName = syscall.StringToUTF16Ptr("WinhelloClass")
	wndClass.lpszMenuName = nil
	wndClass.hIconSm = wndClass.hIcon

	// Register class
	class, _, lastError := RegisterClassEx.Call(uintptr(unsafe.Pointer(wndClass)))
	if class == 0 {
		panic(lastError.Error())
	}

	// Create window
	hWnd, _, lastError := CreateWindowEx.Call(uintptr(WS_EX_CLIENTEDGE),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("WinhelloClass"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Winhello"))),
		uintptr(WS_OVERLAPPEDWINDOW),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		0,
		0,
		hInstance,
		0)
	if hWnd == 0 {
		panic(lastError.Error())
	}

	// Show window
	_, _, _ = ShowWindow.Call(hWnd, uintptr(iCmdShow))

	// Update window
	returnValue, _, lastError = UpdateWindow.Call(hWnd)
	if returnValue == 0 {
		panic(lastError.Error())
	}

	message := new(MSG)

	// message loop
messageLoop:
	for {
		returnValue, _, lastError = GetMessage.Call(uintptr(unsafe.Pointer(message)), 0, 0, 0)
		switch int(returnValue) {
		case -1:
			panic(lastError.Error())
		case 0:
			break messageLoop
		default:
			_, _, _ = TranslateMessage.Call(uintptr(unsafe.Pointer(message)))
			_, _, _ = DispatchMessage.Call(uintptr(unsafe.Pointer(message)))
		}
	}

	return int(message.WParam)
}

func WndProc(hWnd uintptr, uMsg uint, wparam, lparam uintptr) uintptr {

	var (
		lastError   error
		returnValue uintptr
	)

	switch uMsg {

	case WM_CREATE:
		hMenu, _, lastError := CreateMenu.Call()
		if hMenu == 0 {
			panic(lastError.Error())
		}
		hSubMenu, _, lastError := CreatePopupMenu.Call()
		if hSubMenu == 0 {
			panic(lastError.Error())
		}
		returnValue, _, lastError = AppendMenu.Call(hSubMenu, MF_STRING, ID_FILE_EXIT, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("E&xit"))))
		if returnValue == 0 {
			panic(lastError.Error())
		}
		returnValue, _, lastError = AppendMenu.Call(hMenu, MF_STRING|MF_POPUP, hSubMenu, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("&File"))))
		if returnValue == 0 {
			panic(lastError.Error())
		}
		returnValue, _, lastError = SetMenu.Call(hWnd, hMenu)
		if returnValue == 0 {
			panic(lastError.Error())
		}

	case WM_DESTROY:
		_, _, _ = PostQuitMessage.Call(0)

	case WM_PAINT:
		paintStruct := new(PAINTSTRUCT)
		hdc, _, _ := BeginPaint.Call(hWnd, uintptr(unsafe.Pointer(paintStruct)))
		if hdc == 0 {
			panic("Failed")
		}
		text := "Hello, 世界"
		returnValue, _, _ = TextOut.Call(hdc, 100, 100, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))), uintptr(len([]rune(text))))
		if returnValue == 0 {
			panic("Failed")
		}
		_, _, _ = EndPaint.Call(hWnd, uintptr(unsafe.Pointer(paintStruct)))

	case WM_COMMAND:
		switch int(wparam) {
		case ID_FILE_EXIT:
			returnValue, _, lastError = PostMessage.Call(hWnd, WM_CLOSE, 0, 0)
			if returnValue == 0 {
				panic(lastError.Error())
			}

		}

	default:
		returnValue, _, _ = DefWindowProc.Call(hWnd, uintptr(uMsg), wparam, lparam)
		return returnValue
	}
	return 0
}
