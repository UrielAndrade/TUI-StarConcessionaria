//go:build !windows
// +build !windows

package main

func initPlatform() {
	// Em Linux e Unix-like, o terminal já opera com UTF-8 e sequências ANSI por padrão.
}
