package matcher

import (
	"sort"
	"strings"

	"github.com/carmel/microservices/midware"
)

// Matcher is a midware matcher.
type Matcher interface {
	Use(ms ...midware.Midware)
	Add(selector string, ms ...midware.Midware)
	Match(operation string) []midware.Midware
}

// New new a midware matcher.
func New() Matcher {
	return &matcher{
		matchs: make(map[string][]midware.Midware),
	}
}

type matcher struct {
	prefix   []string
	defaults []midware.Midware
	matchs   map[string][]midware.Midware
}

func (m *matcher) Use(ms ...midware.Midware) {
	m.defaults = ms
}

func (m *matcher) Add(selector string, ms ...midware.Midware) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.prefix = append(m.prefix, selector)
		// sort the prefix:
		//  - /foo/bar
		//  - /foo
		sort.Slice(m.prefix, func(i, j int) bool {
			return m.prefix[i] > m.prefix[j]
		})
	}
	m.matchs[selector] = ms
}

func (m *matcher) Match(operation string) []midware.Midware {
	ms := make([]midware.Midware, 0, len(m.defaults))
	if len(m.defaults) > 0 {
		ms = append(ms, m.defaults...)
	}
	if next, ok := m.matchs[operation]; ok {
		return append(ms, next...)
	}
	for _, prefix := range m.prefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.matchs[prefix]...)
		}
	}
	return ms
}
