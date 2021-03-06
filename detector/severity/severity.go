package severity

import (
	"fmt"
	"strings"
)

var severityMap = map[Severity]string{
	Low:    "low",
	Medium: "medium",
	High:   "high",
}

func String(severity Severity) string {
	return severityMap[severity]
}

func FromString(severity string) (Severity, error) {
	severityInLowerCase := strings.ToLower(severity)
	for k, v := range severityMap {
		if v == severityInLowerCase {
			return k, nil
		}
	}
	return 0, fmt.Errorf("unknown severity %v", severity)
}
type Severity int

func (s *Severity) String() string {
	return String(*s)
}

func (s *Severity) ExceedsThreshold(threshold Severity) bool {
	return *s >= threshold
}

func (s *Severity) UnmarshalYAML(get func(interface{}) error) error {
	in := ""
	err := get(&in)
	if err != nil {
		return fmt.Errorf("Severity.Umarshal error: %v\n", err)
	}
	*s, err = FromString(in)
	return err
}

func (s *Severity) MarshalYAML() (interface{}, error) {
	return String(*s), nil
}

func (s *Severity) UnmarshalJSON(input []byte) error {
	v, err := FromString(string(input))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

func (s *Severity) MarshalJSON() ([]byte, error) {
	return []byte(String(*s)), nil
}

const (
	Low = Severity(iota + 1)
	Medium
	High
)