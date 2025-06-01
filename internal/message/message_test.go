package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	parameters := []struct {
		name             string
		content          string
		mentionsEveryone bool
	}{
		{
			name:             "message mentioning everyone",
			content:          "testMessage",
			mentionsEveryone: true,
		},
		{
			name:             "message not mentioning everyone",
			content:          "anotherTestMessage",
			mentionsEveryone: false,
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			msg := NewMessage(parameter.content, parameter.mentionsEveryone)

			assert.NotNil(t, msg)
			assert.Equal(t, parameter.content, msg.Content)
			assert.Equal(t, parameter.mentionsEveryone, msg.MentionsEveryone)
		})
	}
}
