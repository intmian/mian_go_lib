package ai

import "context"

type IAgent interface {
	Chat(ctx context.Context, content string) (string, error)
}
