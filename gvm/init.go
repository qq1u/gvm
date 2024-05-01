package gvm

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"gvm/util"
)

var (
	dir        string
	GO         string
	GOVersions string

	mirrors = []string{
		"https://go.dev/dl",
		"https://mirrors.aliyun.com/golang",
		"https://golang.google.cn/dl",
	}
)

type config struct {
	path    string
	BaseURL string `json:"baseUrl"`
}

func newConfig() *config {
	u, err := user.Current()
	if err != nil {
		util.PrintlnExit("Get current user failed: %v", err)
	}

	p := filepath.Join(u.HomeDir, "gvm.json")
	conf := &config{
		path:    p,
		BaseURL: mirrors[0],
	}

	var data []byte
	if !util.Exists(p) {
		_ = conf.Save()
	} else {
		data, err = os.ReadFile(p)
		if err == nil {
			_ = json.Unmarshal(data, conf)
		}
	}

	return conf
}

func (c *config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}

var c = newConfig()

func init() {
	var err error
	dir, err = util.CurrentDir()
	if err != nil {
		util.PrintlnExit("Get current dir path failed: %v", err)
	}

	GOVersions = filepath.Join(dir, "go_versions")

	GO = filepath.Join(dir, "go")
}

func mkdirVersions() {
	if !util.Exists(GOVersions) {
		if err := os.Mkdir(GOVersions, 0755); err != nil {
			util.PrintlnExit("Mkdir %s failed: %v", GOVersions, err)
		}
	}
}
