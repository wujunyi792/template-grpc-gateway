//go:build windows
// +build windows

package fsx

import "os"

func CloseOnExec(*os.File) {
}
