//go:build linux || darwin
// +build linux darwin

package fsx

import (
	"os"
	"syscall"
)

// CloseOnExec makes sure closing the file on process forking.
func CloseOnExec(file *os.File) {
	if file != nil {
		syscall.CloseOnExec(int(file.Fd()))
	}
}
