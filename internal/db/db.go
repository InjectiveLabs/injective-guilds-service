package db

import (
	"context"
)

type DBService interface {
	Disconnect(ctx context.Context) error
}
