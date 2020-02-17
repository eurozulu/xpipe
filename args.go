package main

import (
	"strconv"
	"strings"
	"time"
)

type Arguments struct {
	Args []string
}

func (args Arguments) Parameters() []string {
	var pms []string
	for i := 0; i < len(args.Args); i++ {
		pm := args.Args[i]

		// Only non 'flags'
		if strings.HasPrefix(pm, "-") && pm != "-" {
			i++
			continue
		}
		pms = append(pms, pm)
	}
	return pms
}

func (args Arguments) Flag(keys ...string) (string, bool) {
	index := -1
	for i, arg := range args.Args {
		// Only 'flags'
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		arg = strings.Trim(arg, "-")
		var key string
		for _, k := range keys {
			if strings.EqualFold(k, arg) {
				key = k
				break
			}
		}
		if key == "" { // not this argument
			continue
		}
		index = i
		break
	}
	if index < 0 {
		return "", false
	}
	v := ""
	if  index + 1 < len(args.Args) {
		v = args.Args[index + 1]
	}
	return v, true
}

func (args Arguments) FlagInt(key string) (int64, bool) {
	var v int64
	s, ok := args.Flag(key)
	if !ok {
		return v, false
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return v, false
	}
	return i, true
}

func (args Arguments) FlagDuration(key string) (time.Duration, bool) {
	s, ok := args.Flag(key)
	if !ok {
		return 0, false
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false
	}
	return d, true
}
