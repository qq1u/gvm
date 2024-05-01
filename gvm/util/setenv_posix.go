//go:build linux || darwin

package util

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func export(name, value string) string {
	return fmt.Sprintf(
		`echo export %s=%s >> %s`,
		name, value, profile,
	)
}

func Setenv(name, value string) error {
	text := fmt.Sprintf(
		`if [ -z "$%s" ]; then %s; fi`,
		name, export(name, value),
	)
	if err := execute(exec.Command("bash", "-c", text)); err != nil {
		return err
	}
	return os.Setenv(name, value)
}

func ForceSetenv(name, value string) error {
	text := fmt.Sprintf(
		`if [ "$%s" != "%s" ]; then %s; fi`,
		name, value, export(name, value),
	)
	if err := execute(exec.Command("bash", "-c", text)); err != nil {
		return err
	}
	return os.Setenv(name, value)
}

func SetLink(source, target string) error {
	return os.Symlink(target, source)
}

func recursionExpandEnv(value string) string {
	for strings.Contains(value, "$") {
		value = os.ExpandEnv(value)
	}
	return value
}

func AppendPath(value string) error {
	envPaths := strings.Split(os.Getenv("PATH"), ":")
	expandedValue := recursionExpandEnv(value)
	if !slices.Contains(envPaths, expandedValue) {
		if strings.HasPrefix(value, "$") {
			value = `\` + value
		}
		text := fmt.Sprintf(`echo "export PATH=\$PATH:%s" >> %s`, value, profile)
		return execute(exec.Command("bash", "-c", text))
	}

	return nil
}

func Source() {
	fmt.Printf("Please execute\nsource %s\n", profile)
}

func GetTarget(src string) (string, error) {
	return os.Readlink(src)
}
