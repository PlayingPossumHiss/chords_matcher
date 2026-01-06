package scrap_chords

import (
	"context"
	"fmt"

	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
)

type Scraper interface {
	GetSongs(ctx context.Context, page int) (chan entity.ArtistSongs, chan error)
	GetAllSongs(ctx context.Context) (chan entity.ArtistSongs, chan error)
	GetContext() (context.Context, context.CancelFunc)
}

type Storage interface {
	SaveArtistSongs(entity.ArtistSongs) error
}

type UseCase struct {
	scraper Scraper
	storage Storage
}

func New(
	scraper Scraper,
	storage Storage,
) *UseCase {
	return &UseCase{
		scraper: scraper,
		storage: storage,
	}
}

func (uc *UseCase) ScrapChords(page int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("UseCase.ScrapChords: %w", err)
		}
	}()

	ctx, cancel := uc.scraper.GetContext()
	defer cancel()

	var (
		songsCh chan entity.ArtistSongs
		errCh   chan error
	)

	if page < 0 {
		songsCh, errCh = uc.scraper.GetAllSongs(ctx)
		if err != nil {
			return err
		}
	} else {
		songsCh, errCh = uc.scraper.GetSongs(ctx, page)
		if err != nil {
			return err
		}
	}

	for artistSongs := range songsCh {
		err := uc.storage.SaveArtistSongs(artistSongs)
		if err != nil {
			return err
		}
	}
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
