package cobratype

import (
	"fmt"
	"log"
	"os"
	"time"
)

type TimeFlag time.Time

func NewTimeValue(p *time.Time) *TimeFlag {
	return (*TimeFlag)(p)
}

func (f *TimeFlag) String() string {
	layout := "2006-01-02T15:04:05-07:00"
	return fmt.Sprintf("%v", time.Time(*f).Format(layout))
}

func (f *TimeFlag) Set(v string) error {
	layout := "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layout, v)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	*f = (TimeFlag)(t)
	return nil
}

// Type is only used in help text
func (f *TimeFlag) Type() string {
	return "time"
}