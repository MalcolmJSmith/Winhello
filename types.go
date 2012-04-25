package main

// Windows Types
type WNDCLASSEX struct {
	cbSize        uint
	style         uint
	lpfnWndProc   uintptr
	cbClsExtra    int
	cbWndExtra    int
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type POINT struct {
	X,
	Y int
}

type RECT struct {
	Left,
	Top,
	Right,
	Bottom int
}

type MSG struct {
	Hwnd    uintptr
	Message uint
	WParam  uintptr
	LParam  uintptr
	Time    uint
	Pt      POINT
}

type PAINTSTRUCT struct {
	hdc         uintptr
	fErase      int
	rcPaint     RECT
	fRestore    int
	fIncUpdate  int
	rgbReserved [32]byte
}
