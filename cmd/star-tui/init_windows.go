//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleOutputCP         = kernel32.NewProc("SetConsoleOutputCP")
	procSetConsoleCP               = kernel32.NewProc("SetConsoleCP")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

func initPlatform() {
	// Força UTF-8 (Code Page 65001) para acentuação e bordas perfeitas
	_, _, _ = procSetConsoleOutputCP.Call(uintptr(65001))
	_, _, _ = procSetConsoleCP.Call(uintptr(65001))

	// Habilita processamento VT (Virtual Terminal) para sequências ANSI
	stdout := syscall.Handle(syscall.Stdout)
	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(uintptr(stdout), uintptr(unsafe.Pointer(&mode)))
	if ret != 0 {
		mode |= enableVirtualTerminalProcessing
		_, _, _ = procSetConsoleMode.Call(uintptr(stdout), uintptr(mode))
	}
}
