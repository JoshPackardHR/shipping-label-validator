package models

import "context"

type Manager interface {
	CheckLabel(ctx context.Context, image string) (bool, error)
}
