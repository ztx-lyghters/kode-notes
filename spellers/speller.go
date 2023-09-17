package spellers

type Speller interface {
	Check([]string) (*[][]corrections, error)
	Prepare(map[string]interface{})
	Fix(*[][]corrections, ...*string)
}
