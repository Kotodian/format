package zqformat

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Kotodian/format"
	"github.com/tidwall/gjson"
)

var (
	imeiMatch  = regexp.MustCompile(`(.*?):\s*"%i"\s*`)
	iccidMatch = regexp.MustCompile(`(.*?):\s*"%c"\s*`)
	timeMatch  = regexp.MustCompile(`(.*?):\s*%t\s*`)
	dataMatch  = regexp.MustCompile(`\s*%d\s*\}`)
)

type formatter struct {
	format   string
	imeiKey  string
	iccidKey string
	timeKey  string
}

func ParseFormat(format string) (format.Formatter, error) {
	format = strings.TrimSuffix(strings.TrimPrefix(format, "{"), "}")
	frt := &formatter{format: format}

	frt.match(imeiMatch, func(s string) {
		frt.imeiKey = strings.TrimSpace(s)
	})

	frt.match(iccidMatch, func(s string) {
		frt.iccidKey = strings.TrimSpace(s)
	})

	frt.match(timeMatch, func(s string) {
		frt.timeKey = strings.TrimSpace(s)
	})

	if dataMatch.FindStringSubmatch(format) == nil {
		return nil, errors.New("format must contain {data: %d}")
	}

	return frt, nil
}

func (f *formatter) Match(data []byte) error {
	if f.imeiKey != "" {
		if gjson.GetBytes(data, f.imeiKey).String() == "" {
			return errors.New("imei is not string")
		}
	}

	if f.iccidKey != "" {
		if gjson.GetBytes(data, f.iccidKey).String() == "" {
			return errors.New("iccid is not string")
		}
	}

	if f.timeKey != "" {
		if t := gjson.GetBytes(data, f.timeKey).String(); t == "" {
			return errors.New("time is not string")
		} else if _, err := time.Parse("2006-01-02 15:04:05", t); err != nil {
			if err != nil {
				return errors.New("invalid time format")
			}
		}
	}
	return nil
}

// It should be like this: imei, iccid, time, data
func (f *formatter) Format(args ...interface{}) ([]byte, error) {

	imei, ok := args[0].(string)
	m := make(map[string]interface{})
	if ok && f.imeiKey != "" {
		m[f.imeiKey] = imei
	}

	iccid, ok := args[1].(string)
	if ok && f.iccidKey != "" {
		m[f.iccidKey] = iccid
	}

	time, ok := args[2].(string)
	if ok && f.timeKey != "" {
		m[f.timeKey] = time
	}

	data, ok := args[3].(map[string]interface{})
	if ok {
		for k, v := range data {
			m[k] = v
		}
	}

	return json.Marshal(m)
}

func (f *formatter) match(reg *regexp.Regexp, fn func(string)) {
	if m := reg.FindStringSubmatch(f.format); len(m) > 0 {
		fn(m[1])
	}
}
