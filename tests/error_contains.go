package tests

import (
	"strings"
	"testing"
)

func ErrorContains(t *testing.T, err error, substring string) bool {
	if err == nil {
		t.Errorf("Ожидалась ошибка, содержащая '%s', но получено значение nil", substring)
		return false
	}
	if !strings.Contains(err.Error(), substring) {
		t.Errorf("Ожидалась ошибка, содержащая '%s', но получено: '%s'", substring, err.Error())
		return false
	}
	return true
}
