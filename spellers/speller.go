package spellers

type Speller interface {
	Check(...*string) (*[][]corrections, error)
	Fix(*[][]corrections, ...*string) error
}
