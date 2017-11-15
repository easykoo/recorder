package recorder

import (
	"syscall"
	"unsafe"
)

const (
	WAVE_FORMAT_PCM = 0x0001
	WAVE_MAPPER     = 0xFFFFFFFF

	CALLBACK_TYPEMASK uint32 = 0x00070000    /* callback type mask */
	CALLBACK_NULL     uint32 = 0x00000000    /* no callback */
	CALLBACK_WINDOW   uint32 = 0x00010000    /* dwCallback is a HWND */
	CALLBACK_TASK     uint32 = 0x00020000    /* dwCallback is a HTASK */
	CALLBACK_FUNCTION uint32 = 0x00030000    /* dwCallback is a FARPROC */
	CALLBACK_THREAD          = CALLBACK_TASK /* thread ID replaces 16 bit task */
	CALLBACK_EVENT    uint32 = 0x00050000    /* dwCallback is an EVENT Handle */

	WIM_OPEN  = 0x3BE
	WIM_CLOSE = 0x3BF
	WIM_DATA  = 0x3C0
)

type WaveHdr struct {
	LpData          uintptr
	DwBufferLength  uint32
	DwBytesRecorded uint32
	DwUser          uintptr
	DwFlags         uint32
	DwLoops         uint32
	LpNext          uintptr
	Reserved        uintptr
}

type WaveFormatX struct {
	WFormatTag      uint16
	NChannels       uint16
	NSamplesPerSec  uint32
	NAvgBytesPerSec uint32
	NBlockAlign     uint16
	WBitsPerSample  uint16
	CbSize          uint16
}

var (
	// Library
	libwinmm *syscall.LazyDLL

	// Functions
	waveInAddBuffer       *syscall.LazyProc
	waveInClose           *syscall.LazyProc
	waveInOpen            *syscall.LazyProc
	waveInPrepareHeader   *syscall.LazyProc
	waveInUnprepareHeader *syscall.LazyProc
	waveInReset           *syscall.LazyProc
	waveInStart           *syscall.LazyProc
	waveInStop            *syscall.LazyProc
)

func init() {
	// Library
	libwinmm = syscall.NewLazyDLL("winmm.dll")

	waveInAddBuffer = libwinmm.NewProc("waveInAddBuffer")
	waveInClose = libwinmm.NewProc("waveInClose")
	//	waveInGetNumDevs = libwinmm.NewProc("waveInGetNumDevs")
	waveInOpen = libwinmm.NewProc("waveInOpen")
	waveInPrepareHeader = libwinmm.NewProc("waveInPrepareHeader")
	waveInUnprepareHeader = libwinmm.NewProc("waveInUnprepareHeader")
	waveInReset = libwinmm.NewProc("waveInReset")
	waveInStart = libwinmm.NewProc("waveInStart")
	waveInStop = libwinmm.NewProc("waveInStop")
}

//MMRESULT waveInOpen( LPHWAVEIN phwi,  //phwi是返回的句柄存放地址
//UINT uDeviceID,   // uDeviceID是要打开的音频设备ID号，一般都指定为WAVE_MAPPER
//LPWAVEFORMATEX pwfx,
//DWORD dwCallback,  //dwCallback则为指定的回调函数或线程,窗口等的地址
//DWORD dwCallbackInstance,   // dwCallbackInstance为需要向回调函数或线程送入的用户参数
//DWORD fdwOpen  // fdwOpen指定回调方式：CALLBACK_FUNCTION, CALLBACK_THREAD和CALLBACK_WINDOW
//);
func WaveInOpen(hwaveIn *uintptr, uDeviceID uint32, pwfx *WaveFormatX, dwCallback uintptr, dwInstance uintptr, fdwOpen uint32) uint32 {
	ret, _, _ := waveInOpen.Call(uintptr(unsafe.Pointer(hwaveIn)), uintptr(uDeviceID),
		uintptr(unsafe.Pointer(pwfx)), dwCallback,
		dwInstance, uintptr(fdwOpen))
	return uint32(ret)
}

func WaveInOpenFunction(hwaveIn *uintptr, uDeviceID uint32, pwfx *WaveFormatX, dwCallback interface{}, fdwOpen uint32) uint32 {
	ret, _, _ := waveInOpen.Call(uintptr(unsafe.Pointer(hwaveIn)), uintptr(uDeviceID),
		uintptr(unsafe.Pointer(pwfx)), syscall.NewCallback(dwCallback),
		0, uintptr(fdwOpen))
	return uint32(ret)
}

func WaveInPrepareHeader(hwaveIn uintptr, pwh *WaveHdr, cbwh uintptr) uint32 {
	ret, _, _ := waveInPrepareHeader.Call(hwaveIn,
		uintptr(unsafe.Pointer(pwh)), cbwh)
	return uint32(ret)
}

func WaveInUnprepareHeader(hwaveIn uintptr, pwh *WaveHdr, cbwh uintptr) uint32 {
	ret, _, _ := waveInUnprepareHeader.Call(hwaveIn,
		uintptr(unsafe.Pointer(pwh)), cbwh)
	return uint32(ret)
}

func WaveInAddBuffer(hwaveIn uintptr, pwh *WaveHdr, cbwh uintptr) uint32 {
	ret, _, _ := waveInAddBuffer.Call(hwaveIn,
		uintptr(unsafe.Pointer(pwh)), cbwh)
	return uint32(ret)
}

func WaveInReset(hwaveIn uintptr) uint32 {
	ret, _, _ := waveInReset.Call(hwaveIn)
	return uint32(ret)
}

func WaveInStart(hwaveIn uintptr) uint32 {
	ret, _, _ := waveInStart.Call(hwaveIn)
	return uint32(ret)
}

func WaveInStop(hwaveIn uintptr) uint32 {
	ret, _, _ := waveInStop.Call(hwaveIn)
	return uint32(ret)
}

func WaveInClose(hwaveIn uintptr) uint32 {
	ret, _, _ := waveInClose.Call(hwaveIn)
	return uint32(ret)
}
