package datastructure

import "strings"

type StringSet map[string]struct{}

func (s StringSet) Add(keys ...string) {
	for _, k := range keys {
		s[k] = struct{}{}
	}
}

func (s StringSet) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s StringSet) IsPrefixOf(str string) bool {
	for k := range s {
		if strings.HasPrefix(str, k) {
			return true
		}
	}
	return false
}
