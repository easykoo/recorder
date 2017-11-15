package recorder

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

// const BufSize = 1024 * 100 //100k

// const BufSize = 1024 * 10 //10k
const BufSize = 3200
const BufNum = 4
const SampleRate = 16000

type Record struct {
	hwnd        uintptr
	stopped     bool
	closed      bool
	waveIn      uintptr
	buffers     [BufNum][BufSize]byte
	waveHdrs    [BufNum]WaveHdr
	handlerFunc HandlerFunc
	stopChan    chan int
}

type HandlerFunc func(data []byte, length int)

func NewRecord(callback HandlerFunc) *Record {
	r := Record{}
	r.handlerFunc = callback
	r.window()
	return &r
}

func (r *Record) OpenDevice() error {
	fmx := WaveFormatX{}
	fmx.WFormatTag = WAVE_FORMAT_PCM
	fmx.NChannels = 1
	fmx.NSamplesPerSec = SampleRate
	fmx.WBitsPerSample = 16
	fmx.NBlockAlign = fmx.WBitsPerSample * fmx.NChannels / 8
	fmx.NAvgBytesPerSec = uint32(fmx.WBitsPerSample * fmx.NBlockAlign)
	fmx.CbSize = 0

	ret := WaveInOpen(&r.waveIn, WAVE_MAPPER, &fmx, r.hwnd, uintptr(0), CALLBACK_WINDOW)
	// ret := WaveInOpen(&r.waveIn, WAVE_MAPPER, &fmx, uintptr(r.event.h), CALLBACK_EVENT)
	// ret := WaveInOpenFunction(&r.waveIn, WAVE_MAPPER, &fmx, r.CallBack, CALLBACK_FUNCTION)
	if ret != 0 {
		fmt.Println("ret: ", ret)
		r.release()
		return Error_OpenDevice
	}

	// prepare wave header

	r.waveHdrs = [4]WaveHdr{}
	r.buffers = [4][BufSize]byte{}
	for i := 0; i < BufNum; i++ {
		r.buffers[i] = [BufSize]byte{}
		r.waveHdrs[i] = WaveHdr{}
		r.waveHdrs[i].LpData = uintptr(unsafe.Pointer(&r.buffers[i][0]))
		r.waveHdrs[i].DwBufferLength = BufSize
		r.waveHdrs[i].DwLoops = 1

		ret = WaveInPrepareHeader(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
		if ret != 0 {
			fmt.Println("WaveInAddBuffer ret: ", ret)
			r.release()
			return Error_AddBuffer
		}
		ret = WaveInAddBuffer(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
		if ret != 0 {
			fmt.Println("WaveInAddBuffer ret: ", ret)
			r.release()
			return Error_AddBuffer
		}

		// fmt.Printf("%+v\n", r.waveHdrs[i])
	}

	return nil
}

func (r *Record) CloseDevice() error {

	r.closed = true
	fmt.Println("close device")
	if r.waveIn == 0 {
		return Error_InvalidHandle
	}

	if WaveInReset(r.waveIn) != 0 {
		return Error_Reset
	}

	r.release()

	if ret := WaveInClose(r.waveIn); ret != 0 {

		return Error_CloseDevice
	}
	return nil
}

func (r *Record) ProcessMsg() {
	var msg Msg
	for GetMessage(&msg, r.hwnd, 0, 0) > 0 && !r.stopped {
		fmt.Printf("%x\n", msg.Message)
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

// must be executed on main thread
func (r *Record) StartRecord() error {
	if r.waveIn == 0 {
		return Error_InvalidHandle
	}
	ret := WaveInStart(r.waveIn)
	if ret != 0 {
		return Error_StartRecord
	}
	// fmt.Println("StartRecord")
	r.ProcessMsg()
	return nil
}

func (r *Record) StopRecord() error {
	if r.waveIn == 0 {
		return Error_InvalidHandle
	}
	r.stopped = true
	ret := WaveInStop(r.waveIn)
	if ret != 0 {
		return Error_StopRecord
	}
	fmt.Println("StopRecord")
	return nil
}

func (r *Record) release() {
	for _, waveHdr := range r.waveHdrs {
		fmt.Println("release")
		WaveInUnprepareHeader(r.waveIn, &waveHdr, unsafe.Sizeof(WaveHdr{}))
	}
}

func (r *Record) window() {
	ins := GetModuleHandle(nil)
	pc, _ := syscall.UTF16PtrFromString("TESTWIN")
	pt, _ := syscall.UTF16PtrFromString("TEST")

	wc := WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})), //必須
		LpfnWndProc:   syscall.NewCallback(r.windowProc),   //必須
		LpszClassName: pc,
	}
	if int(RegisterClassEx(&wc)) == 0 {
		return
	}
	r.hwnd = uintptr(CreateWindowEx(
		WS_EX_TRANSPARENT,
		pc,
		pt,
		WS_POPUP,
		0, 0,
		0, 0,
		0, 0, ins, nil,
	))
}

func (r *Record) windowProc(hWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case WIM_OPEN:
		// fmt.Println("windowProc open")
	case WIM_CLOSE:
		// fmt.Println("windowProc close")
	case WIM_DATA:
		// fmt.Println("windowProc data")
		if r.stopped || r.closed {
			return 0
		}

		for i := 0; i < BufNum; i++ {
			if lParam == uintptr(unsafe.Pointer(&r.waveHdrs[i])) {
				y := reflect.SliceHeader{
					Len:  int(r.waveHdrs[i].DwBytesRecorded),
					Cap:  int(r.waveHdrs[i].DwBytesRecorded),
					Data: r.waveHdrs[i].LpData,
				}
				data := *(*[]byte)(unsafe.Pointer(&y))
				r.handlerFunc(data, int(r.waveHdrs[i].DwBytesRecorded))
				ret := WaveInPrepareHeader(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
				if ret != 0 {
					fmt.Println("WaveInPrepareHeader ret: ", ret)
					r.release()
					return 1
				}
				ret = WaveInAddBuffer(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
				if ret != 0 {
					fmt.Println("WaveInAddBuffer failed: ", ret)
					r.release()
					return 1
				}
			}
		}

	default:
		return DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
