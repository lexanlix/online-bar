package repeatable

import (
	"strings"
	"time"
)

func DoWithTries(fn func() error, attemts int, delay time.Duration) (err error) {
	for attemts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemts--

			continue
		}
		return nil
	}

	return
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}
