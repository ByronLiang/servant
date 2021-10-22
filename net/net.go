package net

type Server interface {
	Start() error
	Stop() error
}
