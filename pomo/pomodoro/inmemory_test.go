package pomodoro_test

import (
	"testing"

	"github.com/rhysmeister/pomo/pomodoro"
	"github.com/rhysmeister/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
