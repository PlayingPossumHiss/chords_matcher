package storage

import "gitlab.com/playing-possum-garbage/chords-matcher/entity"

func artistSongFromEntity(src entity.ArtistSongs) artistSong {
	return artistSong{
		Artist: src.Artist,
		Songs:  songsFormEntity(src.Songs),
	}
}

func artistSongToEntity(src artistSong) entity.ArtistSongs {
	return entity.ArtistSongs{
		Artist: src.Artist,
		Songs:  songsToEntity(src.Songs),
	}
}

func songsFormEntity(src []entity.Song) []song {
	result := make([]song, 0, len(src))
	for _, song := range src {
		result = append(result, songFormEntity(song))
	}
	return result
}

func songsToEntity(src []song) []entity.Song {
	result := make([]entity.Song, 0, len(src))
	for _, song := range src {
		result = append(result, songToEntity(song))
	}
	return result
}

func songFormEntity(src entity.Song) song {
	return song{
		Name:              src.Name,
		ChordsChangeChain: chordsChangeChainFormEntity(src.ChordsChangeChain),
	}
}

func songToEntity(src song) entity.Song {
	return entity.Song{
		Name:              src.Name,
		ChordsChangeChain: chordsChangeChainToEntity(src.ChordsChangeChain),
	}
}

func chordsChangeChainFormEntity(src entity.ChordsChangeChain) chordsChangeChain {
	result := make(chordsChangeChain, 0, len(src))
	for _, chordChange := range src {
		result = append(result, chordChangeFormEntity(chordChange))
	}
	return result
}

func chordsChangeChainToEntity(src chordsChangeChain) entity.ChordsChangeChain {
	result := make(entity.ChordsChangeChain, 0, len(src))
	for _, chordChange := range src {
		result = append(result, chordChangeToEntity(chordChange))
	}
	return result
}

func chordChangeFormEntity(src entity.ChordChange) chordChange {
	return chordChange{
		Steps:       src.Steps,
		TishIsMajor: src.TishIsMajor,
		NextIsMajor: src.NextIsMajor,
	}
}

func chordChangeToEntity(src chordChange) entity.ChordChange {
	return entity.ChordChange{
		Steps:       src.Steps,
		TishIsMajor: src.TishIsMajor,
		NextIsMajor: src.NextIsMajor,
	}
}
