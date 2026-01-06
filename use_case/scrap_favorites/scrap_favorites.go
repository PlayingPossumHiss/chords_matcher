package scrap_favorites

import (
	"context"
	"fmt"

	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
)

type Scraper interface {
	GetArtistSongs(
		ctx context.Context,
		groupID int,
		artistLinks []string,
	) (chan entity.ArtistSongs, chan error)
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

var favSongs = map[int][]string{
	0: {
		"https://181.amdm.ru/akkordi/5nizza/",
	},
	1: {
		"https://amdm.ru/akkordi/arefeva_olga/",
		"https://181.amdm.ru/akkordi/alla_pugacheva/",
		"https://758.amdm.ru/akkordi/ariya/",
		"https://758.amdm.ru/akkordi/akvarium/",
	}, // А
	2: {
		"https://amdm.ru/akkordi/bahyt_kompot/",
		"https://amdm.ru/akkordi/bashnya_rowan/",
		"https://758.amdm.ru/akkordi/bi_2/",
	}, // Б
	3: {
		"https://758.amdm.ru/akkordi/valentin_strykalo/",
	}, // В
	4: {
		"https://758.amdm.ru/akkordi/grazhdanskaya_oborona/",
	}, // Г
	5: {
		"https://amdm.ru/akkordi/dyagileva_yanka/",
		"https://758.amdm.ru/akkordi/ddt/",
		"https://758.amdm.ru/akkordi/dayte_tank/",
	}, // Д
	8: {
		"https://amdm.ru/akkordi/zemfira/",
		"https://amdm.ru/akkordi/zahar_may/",
	}, // З
	10: {
		"https://amdm.ru/akkordi/komsomolsk/",
		"https://amdm.ru/akkordi/kis_kis/",
		"https://758.amdm.ru/akkordi/korol_i_shut/",
	}, // К
	11: {
		"https://181.amdm.ru/akkordi/leningrad/",
		"https://758.amdm.ru/akkordi/lyubeh/",
		"https://758.amdm.ru/akkordi/lyapis_trubetskoy/",
	}, // Л
	12: {
		"https://181.amdm.ru/akkordi/melnitsa/",
		"https://758.amdm.ru/akkordi/mihail_krug/",
	}, // М
	14: {
		"https://amdm.ru/akkordi/operatsiya_plastilin/",
	},
	15: {
		"https://amdm.ru/akkordi/pornofilmy/",
		"https://amdm.ru/akkordi/pesni_iz_multikov/",
	}, // П
	17: {
		"https://amdm.ru/akkordi/semen_slepakov/",
		"https://amdm.ru/akkordi/pnevmoslon/",
		"https://758.amdm.ru/akkordi/sektor_gaza/",
		"https://758.amdm.ru/akkordi/splin/",
	}, // С
	19: {
		"https://amdm.ru/akkordi/umka_i_bronevichok/",
	}, // У
	21: {
		"https://amdm.ru/akkordi/hadn_dadn/",
	}, // X
	23: {
		"https://amdm.ru/akkordi/chayf/",
	}, // Ч
	34: {
		"https://amdm.ru/akkordi/flyour/",
	}, // F
}

func (uc *UseCase) ScrapFavChords() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("UseCase.ScrapFavChords: %w", err)
		}
	}()
	ctx, cancel := uc.scraper.GetContext()
	defer cancel()

	for i, links := range favSongs {
		songsCh, errCh := uc.scraper.GetArtistSongs(ctx, i, links)
		for song := range songsCh {
			err = uc.storage.SaveArtistSongs(song)
			if err != nil {
				return err
			}
		}
		for err = range errCh {
			if err != nil {
				return err
			}
		}
	}

	return nil
}
