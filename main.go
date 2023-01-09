// package cobratype
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {

	dat, err := os.ReadFile("/home/gitpod/.checkpoints/test")

	s := strings.TrimSuffix(string(dat), "\n")
	fmt.Println(err)
	fmt.Println(string(dat))
	fmt.Println(s)

	t, err := time.Parse(time.RFC3339Nano, s)

	fmt.Println(err)
	fmt.Println(t)
}

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
		log.Error().Err(err).Send()
		os.Exit(1)
	}
	*f = (TimeFlag)(t)
	return nil
}

func (f *TimeFlag) Type() string {
	return "time"
}

type IntervalFlag struct {
	Start *time.Time
	End   *time.Time
	name  string
}

func NewIntervalValue(start *time.Time, end *time.Time) *IntervalFlag {
	if start == nil {
		start = &time.Time{}
	}
	if end == nil {
		end = &time.Time{}
	}

	return &IntervalFlag{Start: start, End: end}
}

func (f *IntervalFlag) String() string {
	return f.name + "|" + f.Start.Format(time.RFC3339Nano) + "|" + f.End.Format(time.RFC3339Nano)
}

func (f *IntervalFlag) Set(v string) error {

	f.name = v

	home, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Send()
	}

	path := home + "/.checkpoints"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	dat, err := os.ReadFile(home + "/.checkpoints/" + v)

	if err != nil {
		*f.Start = time.Now().UTC()
	} else {
		t, err := time.Parse(time.RFC3339Nano, strings.TrimSuffix(string(dat), "\n"))
		if err != nil {
			log.Error().Err(err).Send()
			*f.Start = time.Now().UTC()
		}
		*f.Start = t
	}
	*f.End = time.Now().UTC()

	log.Debug().Msg("Start Time: " + f.Start.String())
	log.Debug().Msg("End Time: " + f.End.String())
	return nil
}

func SaveInterval(cmd *cobra.Command, args []string) error {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Send()
	}

	cmd.Flags().Visit(func(flag *pflag.Flag) {

		if flag.Value.Type() == "interval" {

			parts := strings.Split(flag.Value.String(), "|")

			name := parts[0]
			end := parts[2]

			file, err := os.Create(home + "/.checkpoints/" + name)

			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			_, err = file.WriteString(end)

			if err != nil {
				log.Error().Err(err).Send()
				return
			}
		}

	})

	return nil

}

func (f *IntervalFlag) Type() string {
	return "interval"
}

func ExclusiveRequireGroups(combinations [][]string) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {
		setOfNames := make(map[string]bool)
		var bucket *int

		for _, v := range combinations {
			for _, v2 := range v {
				setOfNames[v2] = true
			}
		}

		error := false

		cmd.Flags().Visit(func(flag *pflag.Flag) {

			for i, v := range combinations {
				for _, v2 := range v {
					if v2 == flag.Name {
						if bucket == nil {
							bucket = &i
						} else if bucket != &i {
							error = true
						}
					}
				}
			}
		})

		if error || bucket == nil {
			b, err := json.Marshal(combinations)
			if err != nil {
				fmt.Println(err)
			}
			return errors.New("Only include one of the following groups: " + string(b))
		} else {
			return nil
		}
	}
}
