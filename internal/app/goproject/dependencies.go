package goproject

import "example.com/example/goproject/internal/pkg/config"

type GoProjectDeps struct {
}

func NewGoProjectDeps(conf *config.Config) *GoProjectDeps {
	deps := &GoProjectDeps{}
	return deps
}
