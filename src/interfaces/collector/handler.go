package collector

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shuymn/nijisanji-db-collector/src/domain/service"
)

type Dependency struct {
	CollectService service.CollectService
}

func (d *Dependency) CollectLiversHandler(ctx context.Context) {
	if err := d.CollectService.CollectLiver(ctx); err != nil {
		log.Error().Err(err).Send()
	}
}
