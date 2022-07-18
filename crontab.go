package every

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	crontabCommand = "crontab"
)

func WriteCrontab(config *Config) error {
	crontab, err := readCrontab()
	if err != nil {
		return err
	}

	crontab, err = UpdateCrontab(crontab, config)

	return writeCrontab(crontab)
}

func UpdateCrontab(crontab string, config *Config) (string, error) {
	configPath, err := filepath.Abs(config.Path)
	if err != nil {
		return "", fmt.Errorf("can't get config absolute path: %v", err)
	}

	header := fmt.Sprintf("# Begin every generated jobs for %s", configPath)
	footer := fmt.Sprintf("# End every generated jobs for %s", configPath)

	reBlock, err := regexp.Compile(fmt.Sprintf(`(?m)^%s$(?:.*\n)+^%s$`, regexp.QuoteMeta(header), regexp.QuoteMeta(footer)))
	if err != nil {
		return "", fmt.Errorf("crontab regex error: %v", err)
	}

	buf := &strings.Builder{}
	buf.WriteString(header)
	buf.WriteString("\n")

	for _, e := range config.Everies {
		cronjob, err := e.Cronjob()
		if err != nil {
			return "", err
		}

		buf.WriteString(cronjob)
		buf.WriteString("\n")
	}
	buf.WriteString(footer)

	matched := reBlock.MatchString(crontab)
	if matched {
		crontab = reBlock.ReplaceAllString(crontab, buf.String())
		return crontab, nil
	}

	crontab += "\n\n"
	crontab += buf.String()
	crontab += "\n\n"

	return crontab, nil
}

func readCrontab() (string, error) {
	args := []string{"-l"}

	output, err := exec.Command(crontabCommand, args...).Output()
	if err != nil {
		return "", fmt.Errorf("can't read crontab: %v", err)
	}

	return string(output), nil
}

func writeCrontab(content string) error {
	cmd := exec.Command(crontabCommand)
	cmd.Stdin = strings.NewReader(content)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can't write crontab: %v", err)
	}

	return nil
}

func CleanCrontab(config *Config) error {
	crontab, err := readCrontab()
	if err != nil {
		return err
	}

	configPath, err := filepath.Abs(config.Path)
	if err != nil {
		return fmt.Errorf("can't get config absolute path: %v", err)
	}

	header := fmt.Sprintf("# Begin every generated jobs for %s", configPath)
	footer := fmt.Sprintf("# End every generated jobs for %s", configPath)

	reBlock, err := regexp.Compile(fmt.Sprintf(`(?m)^%s$(?:.*\n)+^%s$`, regexp.QuoteMeta(header), regexp.QuoteMeta(footer)))
	if err != nil {
		return fmt.Errorf("crontab regex error: %v", err)
	}
	matched := reBlock.MatchString(crontab)
	if matched {
		crontab = reBlock.ReplaceAllString(crontab, "")
	}

	return nil
}
