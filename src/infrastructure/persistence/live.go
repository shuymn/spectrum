package persistence

import (
	"context"

	"github.com/guregu/dynamo"
	"github.com/shuymn/nijisanji-db-collector/src/domain/entity"
	"github.com/shuymn/nijisanji-db-collector/src/domain/repository"
	"golang.org/x/sync/errgroup"
)

type liveRepository struct {
	table dynamo.Table
}

func NewLiveRepository(cli *dynamo.DB) repository.LiveRepository {
	table := cli.Table("Livers")
	return &liveRepository{
		table: table,
	}
}

type Platform string

const (
	PlatformYouTube   Platform = "YouTube"
	PlatformTwitter   Platform = "Twitter"
	PlatformBilibili  Platform = "Bilibili"
	PlatformTwitch    Platform = "Twitch"
	PlatformFacebook  Platform = "Facebook"
	PlatformInstagram Platform = "Instagram"
)

const (
	KeyLiverID    = "LiverID"
	KeyIdentifier = "Identifier"
	KeyPlatform   = "Platform"
	KeyChannelID  = "ChannelID"
	KeyScreenName = "ScreenName"
	KeySpaceID    = "SpaceID"
	KeyLoginID    = "LoginID"
	KeyURL        = "URL"
)

func (r *liveRepository) RegisterLiver(ctx context.Context, l *entity.Liver) (*entity.Liver, error) {
	eg, ctx := errgroup.WithContext(ctx)

	for _, yt := range l.YouTube {
		yt := yt
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, yt.ID, PlatformYouTube)
				return update.Set(KeyChannelID, yt.ID).RunWithContext(ctx)
			}
		})
	}

	for _, tw := range l.Twitter {
		tw := tw
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, tw.ID, PlatformTwitter)
				return update.Set(KeyScreenName, tw.ID).RunWithContext(ctx)
			}
		})
	}

	for _, bb := range l.Bilibili {
		bb := bb
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, bb.ID, PlatformBilibili)
				return update.Set(KeySpaceID, bb.ID).RunWithContext(ctx)
			}
		})
	}

	for _, tv := range l.Twitch {
		tv := tv
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, tv.ID, PlatformTwitch)
				return update.Set(KeyLoginID, tv.ID).RunWithContext(ctx)
			}
		})
	}

	for _, fb := range l.Facebook {
		fb := fb
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, fb.ID, PlatformFacebook)
				return update.Set(KeyURL, fb.URL).RunWithContext(ctx)
			}
		})
	}

	for _, ig := range l.Instagram {
		ig := ig
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				update := r.prepareUpdate(l.ID, ig.ID, PlatformInstagram)
				return update.Set(KeyURL, ig.URL).RunWithContext(ctx)
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return l, nil
}

func (r *liveRepository) prepareUpdate(liverID string, suffix string, p Platform) *dynamo.Update {
	id := r.createIdentifier(p, suffix)
	return r.table.Update(KeyLiverID, liverID).Range(KeyIdentifier, id).Set(KeyPlatform, p)
}

func (r *liveRepository) createIdentifier(p Platform, suffix string) string {
	return string(p) + "#" + suffix
}
