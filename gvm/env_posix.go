//go:build linux || darwin

package gvm

var (
	GOPATH    = `$HOME/go`
	GOROOTBin = `$GOROOT/bin`
	GOPATHBin = `$GOPATH/bin`
)
