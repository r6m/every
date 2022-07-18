package every

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// every minute
// every hour
// 		hour at 30 minutes
// every day at 3 am
// 		day at 3:10 pm
// every Fri at 3 pm
// 		Mon,Fri
// every month on Fri at 3 am
// every 2 minutes
// every 2 hours on Fri
// every 2 days on Des,Nov
// every Mon,Fri on Des,Nov
// every 5 days

var (
	reMinute   = regexp.MustCompile(`(?i)^(?P<min>[0-5]?[0-9] )?(?:minutes|minute|min)`)
	reHour     = regexp.MustCompile(`(?i)^(?P<hour>[1-2]?[0-9] )?(?:hours|hour)`)
	reDay      = regexp.MustCompile(`(?i)^(?P<day>[1-3]?[0-9] )?(?:days|day)`)
	reMonth    = regexp.MustCompile(`(?i)(?:in )?(?P<month>(?:(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)(?:,)?)+)`)
	reWeekdays = regexp.MustCompile(`(?i)(?P<weekday>(?:(?:Sun|Mon|Tus|Wed|Thu|Fri|Sat)(?:,)?)+)`)
	reAtTime   = regexp.MustCompile(`(?i)(?:at (?P<time>(?:(?:[1-9]|[1][0-2])|[0-2][0-3]:[0-5][0-9]) (?:am|pm)))`)
)

// Config represents Everyfile config
type Config struct {
	Path    string
	Everies []*Every `hcl:"every,block"`
}

// Every block data
type Every struct {
	Every string `hcl:",label"`
	User  string `hcl:"user"`
	Run   string `hcl:"run"`
}

// Parse reads config file
func Parse(name string) (*Config, error) {
	bytes, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read config '%s': %v", name, err)
	}

	var diags hcl.Diagnostics
	file, diags := hclsyntax.ParseConfig(bytes, name, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse config: %w", diags)
	}

	config := &Config{}
	diags = gohcl.DecodeBody(file.Body, nil, config)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse config: %w", diags)
	}

	config.Path = name
	return config, nil
}

// Cronjob returns a crontab expression
func (e *Every) Cronjob() (string, error) {
	if e.Run == "" {
		return "", errors.New("run is empty")
	}

	expr, err := e.CronExpr()
	if err != nil {
		return "", err
	}

	cronjob := expr + " " + e.Run

	return cronjob, nil
}

// CronExpr parses every to crontab time expression
func (e *Every) CronExpr() (string, error) {
	var min, hour, day, month, weekday = "*", "*", "*", "*", "*"

	groups, matched := regexMatchMap(reMinute, e.Every)
	if minValue, ok := groups["min"]; ok && matched {
		min = "0"
		if minValue != "" && minValue != "1" {
			min = "*/" + strings.TrimSpace(minValue)
		}
	}

	groups, matched = regexMatchMap(reHour, e.Every)
	if hourValue, ok := groups["hour"]; ok && matched {
		min = "0"
		if hourValue != "" && hourValue != "1" {
			hour = "*/" + strings.TrimSpace(hourValue)
		}
	}

	groups, matched = regexMatchMap(reDay, e.Every)
	if dayValue, ok := groups["day"]; ok && matched {
		hour = "0"
		min = "0"
		if dayValue != "" && dayValue != "1" {
			day = "*/" + strings.TrimSpace(dayValue)
		}
	}

	groups, matched = regexMatchMap(reMonth, e.Every)
	if monthValue, ok := groups["month"]; ok && matched {
		month = strings.TrimSpace(monthValue)
	}

	groups, matched = regexMatchMap(reWeekdays, e.Every)
	if weekdayValue, ok := groups["weekday"]; ok && matched {
		weekday = strings.TrimSpace(weekdayValue)
	}

	groups, matched = regexMatchMap(reAtTime, e.Every)
	if atValue, ok := groups["time"]; ok && matched {
		if atValue != "" {
			atValue = strings.ToUpper(atValue)
			t, err := time.Parse("3 PM", atValue)
			if err != nil {
				t, err = time.Parse("3:4 PM", atValue)
				if err != nil {
					return "", fmt.Errorf("can't parse '%s': %v", atValue, err)
				}
			}
			hour = t.Format("15")
			min = t.Format("4")
		}
	}

	return fmt.Sprintf("%s %s %s %s %s", min, hour, day, month, weekday), nil
}

// regexMatchMap returns matched groups in map
func regexMatchMap(r *regexp.Regexp, str string) (map[string]string, bool) {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string, 0)
	if len(match) == 0 {
		return subMatchMap, false
	}

	for i, name := range r.SubexpNames() {
		subMatchMap[name] = match[i]
	}

	return subMatchMap, true
}
