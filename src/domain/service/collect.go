package service

import "context"

type CollectService interface {
	CollectLiver(ctx context.Context) error
}
