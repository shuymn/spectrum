package repository

import "github.com/shuymn/nijisanji-db-collector/src/domain/entity"

type IchikaraRepository interface {
	FetchLivers() ([]*entity.Liver, error)
}
