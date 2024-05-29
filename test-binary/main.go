package main

import (
	"syscall"
	"unsafe"
)

func MessageBox(hwnd uintptr, caption, title string, flags uint) (int, error) {
	captionPtr, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		return -1, err
	}

	titlePrt, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return -1, err
	}

	ret, _, err := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(captionPtr)),
		uintptr(unsafe.Pointer(titlePrt)),
		uintptr(flags),
	)

	return int(ret), err
}

func main() {
	_, _ = MessageBox(0, "Weeee", "Works!", 0)
}
