package recorder

import (
	"C"
	"fmt"
	"reflect"
	"unsafe"
)

const BufSize = 2048
const BufNum = 10

type Record struct {
	hwnd          uintptr
	stopped       bool
	closed        bool
	waveIn        uintptr
	buffers       [BufNum][BufSize]byte
	waveHdrs      [BufNum]WaveHdr
	handlerFunc   HandlerFunc
	stopChan      chan int
	sampleRate    int
	channel       int
	bitsPerSample int
}

type HandlerFunc func(data []byte, length int)

func NewRecord(sampleRate, channel, bitsPerSample int, callback HandlerFunc) *Record {
	r := Record{}
	r.sampleRate = sampleRate
	r.channel = channel
	r.bitsPerSample = bitsPerSample
	r.handlerFunc = callback
	return &r
}

func (r *Record) OpenDevice() error {
	fmx := WaveFormatX{}
	fmx.WFormatTag = WAVE_FORMAT_PCM
	fmx.NChannels = uint16(r.channel)
	fmx.NSamplesPerSec = uint32(r.sampleRate)
	fmx.WBitsPerSample = uint16(r.bitsPerSample)
	fmx.NBlockAlign = fmx.WBitsPerSample * fmx.NChannels / 8
	fmx.NAvgBytesPerSec = uint32(fmx.WBitsPerSample * fmx.NBlockAlign)
	fmx.CbSize = 0

	r.waveHdrs = [BufNum]WaveHdr{}
	r.buffers = [BufNum][BufSize]byte{}
	for i := 0; i < BufNum; i++ {
		r.buffers[i] = [BufSize]byte{}
		r.waveHdrs[i] = WaveHdr{}
		r.waveHdrs[i].LpData = uintptr(unsafe.Pointer(&r.buffers[i][0]))
		r.waveHdrs[i].DwBufferLength = BufSize
		r.waveHdrs[i].DwLoops = 1
	}

	ret := WaveInOpenFunction(&r.waveIn, WAVE_MAPPER, &fmx, r.waveInProc, CALLBACK_FUNCTION)
	if ret != 0 {
		fmt.Println("WaveInOpenFunction failed: ", ret)
		r.release()
		return Error_OpenDevice
	}

	// prepare wave header
	for i := 0; i < BufNum; i++ {
		ret = WaveInPrepareHeader(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
		if ret != 0 {
			fmt.Println("WaveInAddBuffer failed: ", ret)
			r.release()
			return Error_AddBuffer
		}
		ret = WaveInAddBuffer(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
		if ret != 0 {
			fmt.Println("WaveInAddBuffer failed: ", ret)
			r.release()
			return Error_AddBuffer
		}
	}

	return nil
}

func (r *Record) CloseDevice() error {

	r.closed = true
	// fmt.Println("close device")
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

func (r *Record) StartRecord() error {
	if r.waveIn == 0 {
		return Error_InvalidHandle
	}
	r.stopped = false
	ret := WaveInStart(r.waveIn)
	if ret != 0 {
		return Error_StartRecord
	}
	// fmt.Println("StartRecord")
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
	// fmt.Println("StopRecord")
	return nil
}

func (r *Record) release() {
	for _, waveHdr := range r.waveHdrs {
		// fmt.Println("release")
		WaveInUnprepareHeader(r.waveIn, &waveHdr, unsafe.Sizeof(WaveHdr{}))
	}
}

func (r *Record) waveInProc(hWnd uintptr, msg uint32, instance uint32, wParam uintptr, lParam uintptr) int {
	// fmt.Printf("waveInProc 0x%08x \n", msg)
	switch msg {
	case WIM_OPEN:
		// fmt.Println("waveInProc open")
	case WIM_CLOSE:
		// fmt.Println("waveInProc close")
	case WIM_DATA:
		// fmt.Println("waveInProc data")
		if r.stopped || r.closed {
			return 0
		}

		for i := 0; i < BufNum; i++ {
			if wParam == uintptr(unsafe.Pointer(&r.waveHdrs[i])) {
				y := reflect.SliceHeader{
					Len:  int(r.waveHdrs[i].DwBytesRecorded),
					Cap:  int(r.waveHdrs[i].DwBytesRecorded),
					Data: r.waveHdrs[i].LpData,
				}
				data := *(*[]byte)(unsafe.Pointer(&y))
				r.handlerFunc(data, int(r.waveHdrs[i].DwBytesRecorded))
				ret := WaveInUnprepareHeader(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
				if ret != 0 {
					fmt.Println("WaveInUnprepareHeader failed: ", ret)
					r.release()
					return 1
				}
				ret = WaveInPrepareHeader(r.waveIn, &r.waveHdrs[i], unsafe.Sizeof(r.waveHdrs[i]))
				if ret != 0 {
					fmt.Println("WaveInPrepareHeader failed: ", ret)
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
	}
	return 0
}
