package myaws

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
)

// FormatTime returns a localized time string.
// If humanize flag is true, it is converted to human frendly representation.
func FormatTime(t *time.Time) string {
	location, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		panic(err)
	}

	if viper.GetBool("humanize") {
		// humanized format
		return humanize.Time(t.In(location))
	}

	// default format
	return t.In(location).Format("2006-01-02 15:04:05")
}
