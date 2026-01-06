package match_chords

import (
	"fmt"
	"sync"

	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
)

type Storage interface {
	GetAllSongs() ([]entity.ArtistSongs, error)
}

type UseCase struct {
	storage Storage
}

func New(
	storage Storage,
) *UseCase {
	return &UseCase{
		storage: storage,
	}
}

func (uc *UseCase) MatchChords(sample []entity.Chord) (_ entity.SongKeys, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("UseCase.MatchChords: %w", err)
		}
	}()

	// Получаем мелодию
	sampleMelody := entity.NewChordsChangeChainFromChords(sample)

	// Получаем все мелодии, что у нас есть
	artistSongs, err := uc.storage.GetAllSongs()
	if err != nil {
		return nil, err
	}

	// Обходим каждую мелодию попереходно и проверям есть ли совпадение
	var result entity.SongKeys
	mx := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(artistSongs))
	for _, songs := range artistSongs {
		go func() {
			defer wg.Done()

			for _, song := range songs.Songs {
				for songStartPos := 0; songStartPos < len(song.ChordsChangeChain)-len(sample)-1; songStartPos++ {
					match := true
					for samplePos := range sampleMelody {
						if sampleMelody[samplePos] != song.ChordsChangeChain[samplePos+songStartPos] {
							match = false
							break
						}
					}
					if match {
						mx.Lock()
						result = append(result, entity.SongKey{
							Name:   song.Name,
							Artist: songs.Artist,
						})
						mx.Unlock()
						break
					}
				}
			}
		}()
	}
	wg.Wait()

	return result, nil
}
