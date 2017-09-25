package config

import (
	"errors"
	"flag"
	"strings"
)

var (
	ErrUnknownSourceType = errors.New(`unknown source type`)
)

type SourceType int

const (
	SourceType_URL SourceType = iota
	SourceType_File
)

var namesToSourceTypes = map[string]SourceType{
	`url`:  SourceType_URL,
	`file`: SourceType_File,
}

var sourceTypesToNames = map[SourceType]string{
	SourceType_URL:  `url`,
	SourceType_File: `file`,
}

func (t *SourceType) Set(raw string) (err error) {
	var (
		value = strings.ToLower(strings.Trim(raw, ` `))
		ok    bool
	)

	*t, ok = namesToSourceTypes[value]
	if !ok {
		err = ErrUnknownSourceType
	}
	return err
}

func (t SourceType) String() string {
	str, ok := sourceTypesToNames[t]
	if !ok {
		return ""
	}
	return str
}

type Config struct {
	SourceType SourceType
	PoolSize   int
	Verbal     bool
}

func New() *Config {
	cfg := &Config{}

	flag.Var(&cfg.SourceType, `type`, `set source type 'file' or 'url'`)
	flag.IntVar(&cfg.PoolSize, `k`, 5, `set max concurrency size of pool workers`)
	flag.BoolVar(&cfg.Verbal, `v`, true, `write log each of match by resource`)

	flag.Parse()

	return cfg
}
