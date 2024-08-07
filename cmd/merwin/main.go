package main

import (
	merwin "github.com/jgrove2/Merwin"
)

func main() {
	logDir := "dir/logs/"
	merwinInstance := merwin.NewMerwin(logDir)
	merwinInstance.Run()
}
