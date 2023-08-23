package cmd

import (
	"github.com/rhysmeister/pomo/pomodoro"
	"github.com/rhysmeister/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
