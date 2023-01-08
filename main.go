package cobratype

import (
	"errors"
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

func (f *TimeFlag) Type() string {
	return "time"
}

type IntervalFlag struct {
	Start time.Time
	End   time.Time
	name  string
}

func NewIntervalValue() *IntervalFlag {
	return &IntervalFlag{}
}

func (f *IntervalFlag) String() string {
	return fmt.Sprintf("%v - %v", time.Time(f.Start).Format(time.RFC3339Nano), time.Time(f.End).Format(time.RFC3339Nano))
}

func (f *IntervalFlag) Set(v string) error {

	f.name = v

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}

	path := home + "/.checkpoints"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	dat, err := os.ReadFile(home + "/.checkpoints/" + v)

	if err != nil {
		f.Start = time.Now().UTC()
	} else {
		t, err := time.Parse(time.RFC3339Nano, string(dat))
		if err != nil {
			log.Println(err)
			f.Start = time.Now().UTC()
		}
		f.Start = t
	}
	f.End = time.Now().UTC()
	return nil
}

func (f *IntervalFlag) Save() error {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(home + "/.checkpoints/" + f.name)

	if err != nil {
		return err
	}

	_, err = file.WriteString(f.End.Format(time.RFC3339Nano))

	if err != nil {
		return err
	}

	return nil
}

func (f *IntervalFlag) Type() string {
	return "interval"
}
