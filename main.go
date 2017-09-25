package main

import (
	"bufio"
	"os"

	"github.com/antonikonovalov/how-match-go/config"
	"github.com/antonikonovalov/how-match-go/executor"
	"github.com/antonikonovalov/how-match-go/log"
	"github.com/antonikonovalov/how-match-go/match"
	"github.com/antonikonovalov/how-match-go/source"
)

func main() {
	var (
		cfg    = config.New()
		logger = log.New()

		sourcer     source.Sourcer
		loggerMatch log.Logger
	)

	if cfg.SourceType == config.SourceType_File {
		sourcer = source.NewFileSourcer()
	} else {
		sourcer = source.NewUrlSourcer()
	}

	if cfg.Verbal {
		loggerMatch = logger
	}

	matcher, err := match.New(sourcer, loggerMatch)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	exec := executor.New(cfg.PoolSize)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		exec.Do(matcher.Calc, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Error("reading standard input:", err)
	}
	exec.Wait()

	logger.Printf("Total: %d\n", matcher.Total())
}
