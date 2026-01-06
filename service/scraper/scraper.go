package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
)

type Scraper struct {
}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) GetAllSongs(ctx context.Context) (chan entity.ArtistSongs, chan error) {
	const procName = "Scraper.GetAllSongs"

	errCh := make(chan error)
	resultCh := make(chan entity.ArtistSongs)

	go func() {
		for i := 1; i <= 54; i++ {
			subResultCh, subErrCh := s.GetSongs(ctx, i)
			for result := range subResultCh {
				resultCh <- result
			}
			for err := range subErrCh {
				errCh <- fmt.Errorf("%s: chromedp.Run: %w", procName, err)
				close(errCh)
				close(resultCh)
				break
			}
		}

		close(errCh)
		close(resultCh)
	}()

	return resultCh, errCh
}

func (s *Scraper) GetSongs(ctx context.Context, i int) (chan entity.ArtistSongs, chan error) {
	const procName = "Scraper.GetSongs"

	errCh := make(chan error)
	resultCh := make(chan entity.ArtistSongs)

	go func() {
		err := chromedp.Run(
			ctx,
			chromedp.Navigate(fmt.Sprintf("https://amdm.ru/chords/%d/", i)),
			chromedp.Sleep(time.Second*2),
			chromedp.ActionFunc(func(ctx context.Context) error {
				artistLinks, err := getAttributesFromDom(
					ctx,
					"a.artist",
					"href",
				)
				if err != nil {
					return err
				}

				subResultCh, subErrCh := s.GetArtistSongs(ctx, i, artistLinks)
				for result := range subResultCh {
					resultCh <- result
				}
				for err := range subErrCh {
					return err
				}

				return nil
			}),
		)

		if err != nil {
			errCh <- fmt.Errorf("%s: chromedp.Run: %w", procName, err)
			close(errCh)
			close(resultCh)
		}

		close(errCh)
		close(resultCh)
	}()

	return resultCh, errCh
}

func (s *Scraper) GetContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(
		context.Background(),
	)
}

func (s *Scraper) GetArtistSongs(
	ctx context.Context,
	groupID int,
	artistLinks []string,
) (chan entity.ArtistSongs, chan error) {
	const procName = "Scraper.getArtistSongs"

	errCh := make(chan error)
	resultCh := make(chan entity.ArtistSongs)

	go func() {
		for _, artistLink := range artistLinks {
			err := chromedp.Run(
				ctx,
				chromedp.Navigate(artistLink),
				chromedp.Sleep(time.Second*2),
				chromedp.ActionFunc(func(ctx context.Context) error {
					songLinks, err := getAttributesFromDom(
						ctx,
						"#tablesort a.g-link",
						"href",
					)
					if err != nil {
						return err
					}

					pathParts := strings.Split(artistLink, "/")
					atrist := "undefined"
					if len(pathParts) != 0 {
						atrist = pathParts[len(pathParts)-2]
					}
					artistKey := fmt.Sprintf(
						"%d-%s",
						groupID, atrist,
					)

					artistSongs, err := s.getSongs(
						ctx,
						artistKey,
						songLinks,
					)
					if err != nil {
						return err
					}
					if len(artistSongs.Songs) > 0 {
						resultCh <- artistSongs
					}

					return nil
				}),
			)

			if err != nil {
				errCh <- fmt.Errorf("%s: chromedp.Run: %w", procName, err)
				close(errCh)
				close(resultCh)
				break
			}
		}

		close(errCh)
		close(resultCh)
	}()

	return resultCh, errCh
}

func (s *Scraper) getSongs(
	ctx context.Context,
	artistKey string,
	songLinks []string,
) (_ entity.ArtistSongs, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Scraper.getSongs: %w", err)
		}
	}()

	result := entity.ArtistSongs{
		ArtistKey: artistKey,
	}
	for _, songLink := range songLinks {
		song := entity.Song{}

		err = chromedp.Run(
			ctx,
			chromedp.Navigate(songLink),
			chromedp.Sleep(750*time.Millisecond),
			chromedp.Text("[itemprop='name']", &song.Name, chromedp.ByQuery),
			chromedp.Text("[itemprop='byArtist']", &result.Artist, chromedp.ByQuery),
			chromedp.ActionFunc(func(ctx context.Context) error {
				rootNode, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}

				chordNodes, err := dom.QuerySelectorAll(rootNode.NodeID, ".podbor__chord").Do(ctx)
				if err != nil {
					return err
				}
				chords := make([]entity.Chord, 0, len(chordNodes))
				for _, chordNode := range chordNodes {
					attributesRaw, err := dom.GetAttributes(chordNode).Do(ctx)
					if err != nil {
						return err
					}
					chordRaw := attributesToMap(attributesRaw)["data-chord"]
					chord, err := entity.NewChord(chordRaw)
					if err != nil {
						log.Println(err)
						return nil
					}
					chords = append(chords, chord)
				}
				song.ChordsChangeChain = entity.NewChordsChangeChainFromChords(chords)

				result.Songs = append(result.Songs, song)
				return nil
			}),
		)
	}
	if err != nil {
		return entity.ArtistSongs{}, fmt.Errorf("chromedp.Run: %w", err)
	}

	return result, nil
}
