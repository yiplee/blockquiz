package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yiplee/blockquiz/core"
)

func TestParseCommand(t *testing.T) {
	ctx := context.Background()
	p := New()

	commands, err := p.Parse(ctx, "2 b")
	assert.Nil(t, err)
	cmd := commands[0]
	assert.Equal(t, core.ActionAnswerQuestion, cmd.Action)
	assert.Equal(t, 1, cmd.Question)
	assert.Equal(t, 1, cmd.Answer)

	input := p.Encode(ctx, commands...)
	assert.Equal(t, "2 B", input)
}
