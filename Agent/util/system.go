package util

import "runtime"

func DetermineShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "powershell", ""
	}

	return "bash", "-c"
}
