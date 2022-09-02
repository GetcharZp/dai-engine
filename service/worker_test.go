package service

import (
	"testing"
)

func TestNewWorker(t *testing.T) {
	w := NewWork([]string{"192.168.1.8:2379"}, "", "")
	w.Run()
	select {}
}
