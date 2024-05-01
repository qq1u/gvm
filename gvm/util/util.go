package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func FormatTime(t time.Time) string { return t.Format("2006/01/02 15:04:05") }

func PrintlnExit(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(0)
}

func PrintHeader(text string) {
	fmt.Print(strings.Repeat("=", 20))
	fmt.Print("      ", strings.ToUpper(text), "      ")
	fmt.Println(strings.Repeat("=", 20))
}
