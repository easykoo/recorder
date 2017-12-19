package recorder

import "unsafe"

func strptr(s string) uintptr {
	return uintptr(unsafe.Pointer(StringToINT8Ptr(s)))
}

func INT8FromString(s string) ([]byte, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return nil, nil
		}
	}
	return []byte(s), nil
}

func StringToINT8(s string) []byte {
	a, err := INT8FromString(s)
	if err != nil {
		panic("syscall: string with NUL passed to StringToINT8")
	}
	return a
}

func StringToINT8Ptr(s string) *byte {
	return &StringToINT8(s)[0]
}
