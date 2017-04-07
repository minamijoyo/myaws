package myaws

import (
	"time"

	"github.com/dustin/go-humanize"
)

// FormatTime returns a localized time string.
// If humanize flag is true, it is converted to human frendly representation.
func (client *Client) FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	location, err := time.LoadLocation(client.timezone)
	if err != nil {
		panic(err)
	}

	if client.format == "json" {
		return t.Format("2006-01-02T15:04:05Z")
	}

	if client.humanize {
		// humanized format
		return humanize.Time(t.In(location))
	}

	// default format
	return t.In(location).Format("2006-01-02 15:04:05")
}
