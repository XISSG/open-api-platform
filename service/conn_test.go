package service

import (
	"testing"
)

func TestConn(t *testing.T) {
	initRedis()
	initDBConn()
}
