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

type Crontab struct {
	ConfigPath string
	User       string
}

func (c *Crontab) WriteCrontab(items ...*Every) (string, error) {
	crontab, err := c.read()
	if err != nil {
		return "", err
	}

	return c.writeCrontab(crontab, items...)
}

func (c *Crontab) writeCrontab(crontab string, items ...*Every) (string, error) {
	configPath, err := filepath.Abs(c.ConfigPath)
	if err != nil {
		return "", err
	}

	header := fmt.Sprintf("# Begin every generated jobs for %s", configPath)
	footer := fmt.Sprintf("# End every generated jobs for %s", configPath)

	reBlock, err := regexp.Compile(fmt.Sprintf(`(?m)^%s$(?:.*\n)+^%s$`, regexp.QuoteMeta(header), regexp.QuoteMeta(footer)))
	if err != nil {
		return "", err
	}

	buf := &strings.Builder{}
	buf.WriteString(header)
	buf.WriteString("\n")

	for _, item := range items {
		cronjob, err := item.Cronjob()
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

	return crontab, nil
}

func (c *Crontab) read() (string, error) {
	args := []string{"-l"}

	if c.User != "" {
		args = append(args, "-u", c.User)
	}

	output, err := exec.Command(crontabCommand, args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
