package gvm

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"gvm/util"
)

var VERSION = "" // gvm version

func generateUrl(version string) (string, error) {
	var suffix string
	switch runtime.GOOS {
	case "windows":
		suffix = "zip"
	case "darwin", "linux":
		suffix = "tar.gz"
	default:
		return "", fmt.Errorf("not supported os: %s", runtime.GOOS)
	}

	filename := fmt.Sprintf("go%s.%s-%s.%s", version, runtime.GOOS, runtime.GOARCH, suffix)
	return url.JoinPath(c.BaseURL, filename)
}

func find(version string) (string, bool) {
	targetPath := filepath.Join(GOVersions, version)
	return targetPath, util.Exists(targetPath)
}

func Install(version string) (err error) {
	mkdirVersions()

	// 查找是否存在
	foundPath, ok := find(version)
	if ok {
		fmt.Printf("Found version: %q\n", version)
		return use(foundPath)
	}

	// 下载
	var tmpPath string
	tmpPath, err = Download(version)
	if err != nil {
		return err
	}

	// 解压
	util.PrintHeader("extract")
	var extractFunc func(string) (string, error)
	if strings.HasSuffix(tmpPath, ".tar.gz") {
		extractFunc = util.ExtractTarGz
	} else if filepath.Ext(tmpPath) == ".zip" {
		extractFunc = util.ExtractZip
	} else {
		return errors.New("not support compress format")
	}
	ch, done := util.Progress()
	var dest string
	dest, err = extractFunc(tmpPath)
	close(ch)
	<-done
	if err != nil {
		util.Removes(tmpPath, dest)
		return fmt.Errorf("extract failed: %w", err)
	} else {
		util.Removes(tmpPath)
	}

	var targetPath = filepath.Join(GOVersions, version)
	oldGoPath := filepath.Join(dest, "go")
	if err = os.Rename(oldGoPath, targetPath); err != nil {
		return fmt.Errorf("rename %q to %q failed: %w", oldGoPath, targetPath, err)
	}
	if err = os.Remove(dest); err != nil {
		return fmt.Errorf("remove %q failed: %w", dest, err)
	}

	if err = use(targetPath); err != nil {
		return err
	}

	if _, err = exec.LookPath("go"); err != nil {
		source()
	}

	return nil
}

func setEnv() error {
	var err error
	if err = util.ForceSetenv("GOROOT", GO); err != nil {
		return fmt.Errorf("forece setenv GOROOT %q failed: %w", GO, err)
	} else {
		fmt.Printf("set env GOROOT: %s\n", GO)
	}

	if err = util.Setenv("GOPATH", GOPATH); err != nil {
		return fmt.Errorf("setenv GOPATH %q failed: %w", GOPATH, err)
	} else {
		fmt.Printf("set env GOPATH: %s\n", GOPATH)
	}

	if err = util.AppendPath(GOROOTBin); err != nil {
		return fmt.Errorf("append path GOPATH %q failed: %w", GOROOTBin, err)
	} else {
		fmt.Printf("append PATH: %s\n", GOROOTBin)
	}

	if err = util.AppendPath(GOPATHBin); err != nil {
		return fmt.Errorf("append path GOPATHBin %q failed: %w", GOPATHBin, err)
	} else {
		fmt.Printf("append PATH: %s\n", GOPATHBin)
	}

	return nil
}

func VerifyVersion(version string) bool {
	tmp := strings.Split(version, ".")
	if len(tmp) != 3 {
		return false
	}

	for _, v := range tmp {
		if _, err := strconv.Atoi(v); err != nil {
			return false
		}
	}

	return true
}

func use(target string) error {
	if util.Exists(GO) {
		if err := os.RemoveAll(GO); err != nil {
			return fmt.Errorf("remove %q failed: %w", GO, err)
		}
	}
	return util.SetLink(GO, target)
}

func Use(version string) error {
	target, ok := find(version)
	if ok {
		return use(target)
	}

	return errors.New("not found version, please install first")
}

func List() error {
	entries, err := os.ReadDir(GOVersions)
	if err != nil {
		return fmt.Errorf("list %s failed: %w", GOVersions, err)
	}

	var target, targetVersion string
	if target, err = util.GetTarget(GO); err == nil {
		targetVersion = filepath.Base(target)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if name == targetVersion {
				fmt.Print("* ")
			}
			fmt.Println(entry.Name())
		}
	}

	return nil
}

func Download(version string) (string, error) {
	util.PrintHeader("download")
	var ch, done = util.Progress()
	defer func() {
		close(ch)
		<-done
	}()
	downloadUrl, err := generateUrl(version)
	if err != nil {
		return "", err
	}

	// 如果存在的话，就直接用
	tmpPath := filepath.Join(GOVersions, path.Base(downloadUrl))
	if util.Exists(tmpPath) {
		if stat, e := os.Stat(tmpPath); e == nil {
			if stat.Size() != 0 {
				return tmpPath, nil
			}
		}
	}

	err = util.Download(tmpPath, downloadUrl)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	return tmpPath, nil
}

func Mirror() {
	util.PrintHeader("mirror")
	found := false
	for _, mirror := range mirrors {
		if c.BaseURL == mirror {
			found = true
			fmt.Print("* ")
		}
		fmt.Println(mirror)
	}

	if !found {
		fmt.Printf("* %s\n", c.BaseURL)
	}
	fmt.Println()
}

func SetMirror(url string) error {
	url = strings.TrimSpace(url)
	for strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	c.BaseURL = url
	if err := c.Save(); err != nil {
		return fmt.Errorf("set mirror %s failed: %w", url, err)
	}
	return nil
}

func Setup() error {
	_, err := exec.LookPath("gvm")
	if err != nil {
		err = util.AppendPath(dir)
		if err != nil {
			fmt.Printf("append %s to PATH failed: %v\n", dir, err)
			fmt.Println("please manual addition")
		} else {
			fmt.Printf("append %s to PATH\n", dir)
		}
	}

	mkdirVersions()

	util.PrintHeader("set env")
	if err = setEnv(); err != nil {
		return err
	}

	source()

	return nil
}

func source() {
	util.Source()
	fmt.Println()
}
