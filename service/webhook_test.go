package service_test

import (
	"testing"

	"github.com/nbtca/saturday/service"
	"github.com/stretchr/testify/assert"
)

func TestGithubWebHook_ExtractCommand(t *testing.T) {
	comment := "@nbtca-bot accept"
	command, err := service.ExtractCommand(comment)

	// Assert the command is correctly extracted
	assert.NoError(t, err)
	assert.Equal(t, "accept", command)
}

func TestGithubWebHook_ExtractSizeLabel(t *testing.T) {
	label := "size:m"
	size, err := service.ExtractSizeLabel(label)

	// Assert the size label is correctly extracted
	assert.NoError(t, err)
	assert.Equal(t, "m", size)
}
