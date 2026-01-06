package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// BaseChord базовые аккорды
type BaseChord byte

const (
	BaseChordC BaseChord = iota + 1
	BaseChordDb
	BaseChordD
	BaseChordEb
	BaseChordE
	BaseChordF
	BaseChordGb
	BaseChordG
	BaseChordAb
	BaseChordA
	BaseChordBb
	BaseChordB
)

// Chord аккорды со всеми (нет) страшными значками
type Chord struct {
	Base    BaseChord
	IsMajor bool
}

// ChordChange изменение между двумя аккордами, описывает сам переход
type ChordChange struct {
	Steps       int8
	TishIsMajor bool
	NextIsMajor bool
}

var ErrUnknownChord = errors.New("unknown chord")

func NewChord(src string) (_ Chord, err error) {
	result := Chord{}

	clearer := regexp.MustCompile(`(maj|sus|dim|add|-|\+|\d|\/.*|\*|\(.*\))`)
	onlyBaseChord := clearer.ReplaceAllString(src, "")

	if !strings.HasSuffix(onlyBaseChord, "m") {
		result.IsMajor = true
	} else {
		onlyBaseChord = onlyBaseChord[:len(onlyBaseChord)-1]
	}

	result.Base, err = NewBaseChord(onlyBaseChord)
	if err != nil {
		return Chord{}, err
	}

	return result, nil
}

func NewBaseChord(src string) (BaseChord, error) {
	var result BaseChord
	switch src {
	case "C":
		result = BaseChordC
	case "G":
		result = BaseChordG
	case "D":
		result = BaseChordD
	case "A":
		result = BaseChordA
	case "E":
		result = BaseChordE
	case "Cb", "B", "H":
		result = BaseChordB
	case "F#", "Gb":
		result = BaseChordGb
	case "C#", "Db":
		result = BaseChordDb
	case "Ab", "G#":
		result = BaseChordAb
	case "Eb", "D#":
		result = BaseChordEb
	case "Bb", "A#", "Hb":
		result = BaseChordBb
	case "F":
		result = BaseChordF
	default:
		return BaseChord(0), fmt.Errorf("%w: unknown chord %s", ErrUnknownChord, src)
	}

	return result, nil
}

// ChordsChangeChain цепочка изменений аккордов, которая и представляет музыку
type ChordsChangeChain []ChordChange

// Song песня
type Song struct {
	Name              string
	ChordsChangeChain ChordsChangeChain
}

type SongKeys []SongKey

func (sk SongKeys) String() string {
	uniqSongs := map[string]map[string]struct{}{}
	for _, item := range sk {
		artistSongs, ok := uniqSongs[item.Artist]
		if !ok {
			artistSongs = map[string]struct{}{}
		}
		artistSongs[item.Name] = struct{}{}
		uniqSongs[item.Artist] = artistSongs
	}

	builder := strings.Builder{}
	for artist, songs := range uniqSongs {
		builder.WriteString(artist)
		builder.WriteString("\n")
		for song := range songs {
			builder.WriteString("\t")
			builder.WriteString(song)
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

type SongKey struct {
	Name   string
	Artist string
}

type ArtistSongs struct {
	Songs     []Song
	Artist    string
	ArtistKey string
}

// NewSong конструктор песни
func NewSong(
	name string,
	artist string,
	chords []Chord,
) Song {
	return Song{
		Name:              name,
		ChordsChangeChain: NewChordsChangeChainFromChords(chords),
	}
}

// NewChordsChangeChainFromChords создает представление песни из аккордов
func NewChordsChangeChainFromChords(chords []Chord) ChordsChangeChain {
	if len(chords) < 2 {
		return nil
	}

	result := make(ChordsChangeChain, 0, len(chords)-1)
	for i := 1; i < len(chords); i++ {
		result = append(result, NewChordChange(chords[i-1], chords[i]))
	}

	return result
}

// NewChordChange создает переход двух аккордов
func NewChordChange(chord1 Chord, chord2 Chord) ChordChange {
	// Получим разницу тонов
	change := chord2.level() - chord1.level()
	if change < -6 {
		change += 12
	}
	if change > 6 {
		change -= 12
	}

	return ChordChange{
		TishIsMajor: chord1.IsMajor,
		NextIsMajor: chord2.IsMajor,
		Steps:       change,
	}
}

func (c Chord) level() int8 {
	return int8(c.Base)
}
