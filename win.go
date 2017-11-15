package recorder

import (
	"fmt"
	"syscall"
	"unsafe"
)

// "user32.dll"

type (
	HANDLE uintptr
	HWND   HANDLE
)
type POINT struct {
	X, Y int32
}
type Msg struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

const (
	WS_EX_TRANSPARENT = 0X00000020
	WS_POPUP          = 0X80000000
)

var (
	libuser32        *syscall.LazyDLL
	getMessage       *syscall.LazyProc
	registerClassEx  *syscall.LazyProc
	createWindowEx   *syscall.LazyProc
	translateMessage *syscall.LazyProc
	dispatchMessage  *syscall.LazyProc
	defWindowProc    *syscall.LazyProc

	libkernel32     *syscall.LazyDLL
	getModuleHandle *syscall.LazyProc
	virtualAlloc    *syscall.LazyProc
	createThread    *syscall.LazyProc
	waitSingle      *syscall.LazyProc
)

func init() {
	libuser32 = syscall.NewLazyDLL("user32.dll")
	getMessage = libuser32.NewProc("GetMessageW")
	registerClassEx = libuser32.NewProc("RegisterClassExW")
	createWindowEx = libuser32.NewProc("CreateWindowExW")
	translateMessage = libuser32.NewProc("TranslateMessage")
	dispatchMessage = libuser32.NewProc("DispatchMessageW")
	defWindowProc = libuser32.NewProc("DefWindowProcW")

	libkernel32 = syscall.NewLazyDLL("kernel32.dll")
	getModuleHandle = libkernel32.NewProc("GetModuleHandleW")
	virtualAlloc = libkernel32.NewProc("VirtualAlloc")
	createThread = libkernel32.NewProc("CreateThread")
	waitSingle = libkernel32.NewProc("WaitForSingleObject")
}

func GetMessage(msg *Msg, hWnd uintptr, msgFilterMin, msgFilterMax uint32) uintptr {
	ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(msg)),
		uintptr(hWnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return ret
}

func RegisterClassEx(windowClass *WNDCLASSEX) uintptr {
	ret, _, _ := registerClassEx.Call(uintptr(unsafe.Pointer(windowClass)))
	return ret
}

func CreateWindowEx(dwExStyle uint32, lpClassName, lpWindowName *uint16, dwStyle uint32, x, y, nWidth, nHeight int32, hWndParent uintptr, hMenu uintptr, hInstance uintptr, lpParam unsafe.Pointer) uintptr {
	ret, _, _ := createWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(lpClassName)),
		uintptr(unsafe.Pointer(lpWindowName)),
		uintptr(dwStyle),
		uintptr(x),
		uintptr(y),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	return ret
}

func TranslateMessage(msg *Msg) uintptr {
	ret, _, _ := translateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}
func DispatchMessage(msg *Msg) uintptr {
	ret, _, _ := dispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}

func DefWindowProc(hWnd uintptr, Msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := defWindowProc.Call(
		uintptr(hWnd),
		uintptr(Msg),
		wParam,
		lParam)
	return ret
}

func GetModuleHandle(lpModuleName *uint16) uintptr {
	ret, _, _ := getModuleHandle.Call(uintptr(unsafe.Pointer(lpModuleName)))
	return ret
}

func VirtualAlloc(lpModuleName *uint16) uintptr {
	ret, _, _ := virtualAlloc.Call(uintptr(unsafe.Pointer(lpModuleName)))
	return ret
}

// HANDLE CreateThread(
// 	LPSECURITY_ATTRIBUTES lpThreadAttributes, 　　　　　// pointer to security attributes
// 	DWORD dwStackSize,　　　　　　　　　　　　　　　　　　// initial thread stack size
// 	LPTHREAD_START_ROUTINE lpStartAddress, 　　　　　　// pointer to thread function
// 	LPVOID lpParameter,　　　　　　　　　　　　　　　　　　　// argument for new thread
// 	DWORD dwCreationFlags,　　　　　　　　　　　　　　　　// creation flags
// 	LPDWORD lpThreadId　　　　　　　　　　　　　　　　　　// pointer to receive thread ID
//   );
func CreateThread(threadFunction interface{}, soundIn uintptr, threadId *uintptr) uintptr {
	ret, _, _ := createThread.Call(0, 0, syscall.NewCallback(threadFunction), soundIn, 0, uintptr(unsafe.Pointer(threadId)))
	fmt.Println(threadId)
	fmt.Println(ret)
	return ret
}

// func WaitSingle() uintptr {
// 	ret, _, _ := waitSingle.Call(uintptr(unsafe.Pointer(lpModuleName)))
// 	return ret
// }
