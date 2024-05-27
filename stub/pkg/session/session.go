package session

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32, _             = syscall.LoadLibrary("kernel32.dll")
	processIDToSessionID, _ = syscall.GetProcAddress(kernel32, "ProcessIdToSessionId")
)

// プロセスセッションIDを取得
func GetProcessSessionID() uint32 {
	var sessionID uint32
	sessionID = 0xffffffff

	if processIDToSessionID != 0 {
		_, _, err := syscall.SyscallN(
			uintptr(processIDToSessionID),
			2,
			uintptr(os.Getpid()),
			uintptr(unsafe.Pointer(&sessionID)),
			0)
		if err != 0 {
			sessionID = 0xffffffff
		}
	} else {
		sessionID = 0xffffffff
	}
	return sessionID
}
