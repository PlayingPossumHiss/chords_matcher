package entity

import (
	"reflect"
	"testing"
)

func TestNewSong(t *testing.T) {
	type args struct {
		name   string
		artist string
		chords []Chord
	}
	tests := []struct {
		name string
		args args
		want Song
	}{
		{
			name: "все идет по плану",
			args: args{
				name:   "Все идет по плану",
				artist: "ГрОб",
				chords: []Chord{
					{
						Base: BaseChordA,
					},
					{
						Base:    BaseChordF,
						IsMajor: true,
					},
					{
						Base:    BaseChordC,
						IsMajor: true,
					},
					{
						Base:    BaseChordE,
						IsMajor: true,
					},
				},
			},
			want: Song{
				Name: "Все идет по плану",
				ChordsChangeChain: ChordsChangeChain{
					ChordChange{
						NextIsMajor: true,
						Steps:       -4,
					},
					ChordChange{
						TishIsMajor: true,
						NextIsMajor: true,
						Steps:       1,
					},
					ChordChange{
						TishIsMajor: true,
						NextIsMajor: true,
						Steps:       4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSong(tt.args.name, tt.args.artist, tt.args.chords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChord(t *testing.T) {
	tests := []struct {
		name    string
		want    Chord
		wantErr bool
	}{
		{
			name: "Cm7-5",
			want: Chord{
				Base: BaseChordC,
			},
		},
		{
			name: "E7*",
			want: Chord{
				Base:    BaseChordE,
				IsMajor: true,
			},
		},
		{
			name: "Dm6add9",
			want: Chord{
				Base: BaseChordD,
			},
		},
		{
			name: "C+7/E",
			want: Chord{
				Base:    BaseChordC,
				IsMajor: true,
			},
		},
		{
			name: "Cm",
			want: Chord{
				Base: BaseChordC,
			},
		},
		{
			name: "Am(V)",
			want: Chord{
				Base: BaseChordA,
			},
		},
		{
			name: "E5",
			want: Chord{
				Base:    BaseChordE,
				IsMajor: true,
			},
		},
		{
			name: "F/C",
			want: Chord{
				Base:    BaseChordF,
				IsMajor: true,
			},
		},
		{
			name: "C6sus",
			want: Chord{
				Base:    BaseChordC,
				IsMajor: true,
			},
		},
		{
			name: "C#m7",
			want: Chord{
				Base: BaseChordDb,
			},
		},
		{
			name: "F#7m",
			want: Chord{
				Base: BaseChordGb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChord(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChord() = %v, want %v", got, tt.want)
			}
		})
	}
}
