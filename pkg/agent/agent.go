package agent

import (
	"context"
)

type Agent interface {
	Run(ctx context.Context, input string) (string, error)
}
