package services

import "github.com/ztx-lyghters/kode-notes/spellers"

type Speller interface {
	spellers.Speller
}

func NewSpellersService() SpellerService {
	return SpellerService{
		Yandex: spellers.NewYandexSpeller(),
	}
}
