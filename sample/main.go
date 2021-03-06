package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/easykoo/recorder"
)

var block = make(chan []byte, 100)
var pcm []byte

func DataProc(data []byte, size int) {
	block <- data
}

func main() {
	r := recorder.NewRecord(16000, 2, 16, DataProc)
	r.OpenDevice()
	defer r.CloseDevice()
	r.StartRecord()
	handle(r)
}

func handle(r *recorder.Record) {
	for {
		select {
		case d := <-block:
			pcm = append(pcm, d...)
			if len(pcm) > 1024*200 { //200k
				writePCM(pcm)
				r.StopRecord()
				goto OUT
			}
		}
	}
OUT:
	fmt.Println("handle end")
}

func writePCM(pcm []byte) {
	resultWav, err := ConvertBytes(pcm, 1, 16000, 16)
	if err != nil {
		fmt.Printf("\n%s\n", err)
	}
	if err = ioutil.WriteFile(fmt.Sprintf("./wave_%s.wav",
		time.Now().Format("20060102_150405.000")), resultWav, 0666); err != nil {
		fmt.Printf("\n%s\n", err)
	}
}
