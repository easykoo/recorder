package recorder

import (
	"errors"
)

var (
	Error_OpenDevice      = errors.New("open device failed")
	Error_StartRecord     = errors.New("start record failed")
	Error_StopRecord      = errors.New("stop record failed")
	Error_Reset           = errors.New("reset failed")
	Error_CloseDevice     = errors.New("close device failed")
	Error_InvalidHandle   = errors.New("invalid handle")
	Error_PrepareHeader   = errors.New("prepare header failed")
	Error_UnprepareHeader = errors.New("unprepare header failed")
	Error_AddBuffer       = errors.New("add buffer failed")
)
