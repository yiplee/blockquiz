package parser

import (
	"context"
	"strings"

	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
)

type commandParser struct {
	// 一条消息里面最多处理多少 cmd
	partLimit int
}

func New() core.CommandParser {
	return &commandParser{partLimit: 3}
}

func (c *commandParser) Parse(ctx context.Context, input string) ([]*core.Command, error) {
	log := logger.FromContext(ctx)
	var commands []*core.Command

	input = strings.ToLower(input)
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ';'
	})

	if len(parts) > c.partLimit {
		parts = parts[:c.partLimit]
	}

	for _, part := range parts {
		if args := newArgs(part); len(args) > 0 {
			cmd, err := parseCommand(args)
			if err != nil {
				log.WithError(err).Errorf("parse command: %s", part)
				continue
			}

			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func (c *commandParser) Encode(ctx context.Context, cmds ...*core.Command) string {
	var parts []string
	for _, cmd := range cmds {
		if args := encodeCommand(cmd); len(args) > 0 {
			parts = append(parts, args.Encode())
		}
	}

	result := strings.Join(parts, ";")
	return result
}
