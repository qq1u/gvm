//go:build windows

package util

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Setenv(name, value string) error {
	text := fmt.Sprintf(`if (-not [Environment]::GetEnvironmentVariable("%[1]s", "User")) { setx %[1]s %[2]s }`, name, value)
	return execute(exec.Command("powershell", "-Command", text))
}

func ForceSetenv(name, value string) error {
	// setx NAME VALUE: 默认是设置 用户级 的环境变量（会直接覆盖旧的）
	return execute(exec.Command("setx", name, value))
}

func SetLink(source, target string) error {
	text := fmt.Sprintf(`New-Item -ItemType Junction -Path %s -Value %s`, source, target)
	return execute(exec.Command("powershell", "-Command", text))
}

func AppendPath(value string) error {
	text := fmt.Sprintf(`
		$OLD = [Environment]::GetEnvironmentVariable('PATH', 'User');
		$PATHS = $OLD -split ';';
		$expandedPath = [Environment]::ExpandEnvironmentVariables('%s');
		if (-not ($PATHS -contains $expandedPath)) { 
			setx PATH "$OLD;%s"
		}`, value, value,
	)
	return execute(exec.Command("powershell", "-Command", text))
}

func Source() {
	fmt.Println("Please restart terminal")
}

func GetTarget(src string) (string, error) {
	text := fmt.Sprintf("(Get-Item %s).Target", src)
	cmd := exec.Command("powershell", "-Command", text)
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("no link")
	}
	return strings.TrimSpace(string(data)), nil
}
