package match

import (
	"bytes"
	"errors"
	"sync/atomic"

	"fmt"
	"github.com/antonikonovalov/how-match-go/log"
	"github.com/antonikonovalov/how-match-go/source"
)

var (
	Pattern                 = []byte(`Go`)
	ErrNotSetSourceResolver = errors.New(`matcher: not set source resolver`)
)

func New(sourcer source.Sourcer, logger log.Logger) (Matcher, error) {
	m := &matcher{
		sourceResolver: sourcer,
		log:            logger,
	}

	if m.log == nil {
		m.log = &log.NoopLogger{}
	}

	if m.sourceResolver == nil {
		return nil, ErrNotSetSourceResolver
	}

	return m, nil
}

// Matcher - interface for collect match by pattern on input source and collect total match
type Matcher interface {
	// Calc - calculate of count matches in source
	Calc(sourcePath string)
	// Total - return total count of find matches
	Total() int
}

type matcher struct {
	sourceResolver source.Sourcer

	log log.Logger

	total int32
}

func (m *matcher) calc(c int) {
	atomic.AddInt32(&m.total, int32(c))
}

func (m *matcher) Calc(sourcePath string) {
	data, err := m.sourceResolver.Get(sourcePath)
	if err != nil {
		m.log.Error(fmt.Errorf("can't get source %s: %s\n", sourcePath, err))
		return
	}

	count := bytes.Count(data, Pattern)
	m.log.Printf("Count for %s: %d\n", sourcePath, count)
	if count != 0 {
		m.calc(count)
	}
}

func (m *matcher) Total() int {
	return int(atomic.LoadInt32(&m.total))
}
