package logger

import (
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	Info("test info:%v", errors.New("INFO"))
	Error("test error:%v", errors.New("ERROR"))
}
