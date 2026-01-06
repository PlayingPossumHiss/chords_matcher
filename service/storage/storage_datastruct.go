package storage

type artistSong struct {
	Artist string `json:"artist_name"`
	Songs  []song `json:"songs"`
}

type chordsChangeChain []chordChange

type song struct {
	Name              string            `json:"name"`
	ChordsChangeChain chordsChangeChain `json:"chords_changes"`
}

type chordChange struct {
	Steps       int8 `json:"steps"`
	TishIsMajor bool `json:"tish_is_major"`
	NextIsMajor bool `json:"next_is_major"`
}
