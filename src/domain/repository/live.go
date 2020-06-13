package repository

import (
	"context"

	"github.com/shuymn/nijisanji-db-collector/src/domain/entity"
)

type LiveRepository interface {
	RegisterLiver(ctx context.Context, l *entity.Liver) (*entity.Liver, error)
}
