package services

import "example.com/kode-notes/spellers"

type Speller interface {
	spellers.Speller
}

func NewSpellersService() SpellerService {
	return SpellerService{
		Yandex: spellers.NewYandexSpeller(),
	}
}
