package application

import (
	"context"

	"github.com/shuymn/nijisanji-db-collector/src/domain/repository"
	"github.com/shuymn/nijisanji-db-collector/src/domain/service"
	"golang.org/x/sync/errgroup"
)

type collectService struct {
	irepo repository.IchikaraRepository
	lrepo repository.LiveRepository
}

func NewCollectService(irepo repository.IchikaraRepository, lrepo repository.LiveRepository) service.CollectService {
	return &collectService{
		irepo: irepo,
		lrepo: lrepo,
	}
}

func (s *collectService) CollectLiver(ctx context.Context) error {
	ls, err := s.irepo.FetchLivers()
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, l := range ls {
		l := l
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				_, err = s.lrepo.RegisterLiver(ctx, l)
				if err != nil {
					return err
				}
				return nil
			}
		})
	}

	return eg.Wait()
}
