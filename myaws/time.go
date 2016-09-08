package myaws

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
)

func FormatTime(t *time.Time) string {
	location, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		panic(err)
	}

	if viper.GetBool("humanize") {
		return humanize.Time(t.In(location))
	} else {
		return t.In(location).Format("2006-01-02 15:04:05")
	}
}
