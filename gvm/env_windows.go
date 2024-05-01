//go:build windows

package gvm

var (
	GOPATH    = `%USERPROFILE%\go`
	GOROOTBin = `%GOROOT%\bin`
	GOPATHBin = `%GOPATH%\bin`
)
