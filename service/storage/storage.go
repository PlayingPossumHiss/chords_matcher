package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
)

type Storage struct {
	root string
}

func New() *Storage {
	return &Storage{
		root: "./.storage",
	}
}

func (s *Storage) SaveArtistSongs(src entity.ArtistSongs) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Storage.SaveArtistSongs: %w", err)
		}
	}()

	data, err := json.Marshal(artistSongFromEntity(src))
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s.json", s.root, src.ArtistKey), data, 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	return nil
}

func (s *Storage) GetAllSongs() (_ []entity.ArtistSongs, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Storage.GetAllSongs: %w", err)
		}
	}()

	files, err := os.ReadDir(s.root)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}

	result := make([]entity.ArtistSongs, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		rawContent, err := os.ReadFile(fmt.Sprintf("%s/%s", s.root, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("os.ReadFile: %w", err)
		}
		artistSongs := artistSong{}
		err = json.Unmarshal(rawContent, &artistSongs)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		result = append(result, artistSongToEntity(artistSongs))
	}

	return result, nil
}
