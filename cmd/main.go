package main

import (
	"log"
	"os"
	"strconv"

	"gitlab.com/playing-possum-garbage/chords-matcher/entity"
	"gitlab.com/playing-possum-garbage/chords-matcher/service/scraper"
	"gitlab.com/playing-possum-garbage/chords-matcher/service/storage"
	"gitlab.com/playing-possum-garbage/chords-matcher/use_case/match_chords"
	"gitlab.com/playing-possum-garbage/chords-matcher/use_case/scrap_chords"
	"gitlab.com/playing-possum-garbage/chords-matcher/use_case/scrap_favorites"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]
	command := args[0]
	switch command {
	case "scrap":
		if len(args) < 2 {
			log.Panicln("need page number")
		}
		if args[1] == "F" {
			uc := scrap_favorites.New(scraper.New(), storage.New())
			err := uc.ScrapFavChords()
			if err != nil {
				log.Panicln(err)
			}
			return
		}
		uc := scrap_chords.New(
			scraper.New(),
			storage.New(),
		)

		page, err := strconv.Atoi(args[1])
		if err != nil {
			log.Panicln(err)
		}
		err = uc.ScrapChords(page)
		if err != nil {
			log.Fatalln(err)
		}
	case "match":
		uc := match_chords.New(storage.New())
		chords := make([]entity.Chord, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			chord, err := entity.NewChord(args[i])
			if err != nil {
				log.Fatalln(err)
			}
			chords = append(chords, chord)
		}
		songs, err := uc.MatchChords(chords)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(songs.String())
	default:
		help()
	}
}

func help() {
	log.Println(`
	Использование:
		Для получения песен на букву: scrap 28 (получить и сохранить все песни на 28 букву т.е. Я. Да, я - 28 буква см. сайт с которого скрейпим)
		Для поиска аккордов: match Em C G D
	`)
}
