package service

import "errors"

var ErrInvalidInput = errors.New("invalid task input")
var ErrCycleDetected = errors.New("A cycle has been detected")
