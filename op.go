package wsocket

type Op int

const (
	OpUnknow       Op = -1
	OpContinuation Op = 1
	OpText         Op = 2
	OpBinary       Op = 3
	OpClose        Op = 4
	OpPing         Op = 5
	OpPong         Op = 6
)
