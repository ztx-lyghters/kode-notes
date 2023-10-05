package spellers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	YA_OPT_IGNORE_DIGITS         = 2
	YA_OPT_IGNORE_URLS           = 4
	YA_OPT_FIND_REPEAT_WORDS     = 8
	YA_OPT_IGNORE_CAPITALIZATION = 512
	YA_SPELLER_URL               = "https://speller.yandex.net/services/spellservice.json/checkTexts?"
)

type corrections struct {
	Word string   `json:"word"`
	Sgst []string `json:"s"`
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
}

type YandexSpeller struct {
	Config map[string]interface{}
	URL    string
}

func NewYandexSpeller() *YandexSpeller {
	return &YandexSpeller{
		Config: map[string]interface{}{
			"options": YA_OPT_IGNORE_DIGITS |
				YA_OPT_IGNORE_URLS,
			"lang":   "",
			"format": "",
		},
		URL: YA_SPELLER_URL,
	}
}

func (s *YandexSpeller) Check(texts ...*string) (*[][]corrections, error) {
	var result [][]corrections

	for _, str := range texts {
		if str == nil || strings.TrimSpace(*texts[0]) == "" {
			return nil, errors.New("empty text field")
		}
	}

	data := url.Values{}
	for k, v := range s.Config {
		data.Add(k, fmt.Sprint(v))
	}
	for _, v := range texts {
		if strings.TrimSpace(*v) != "" {
			data.Add("text", *v)
		}
	}

	r, err := http.PostForm(s.URL, data)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rbody, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *YandexSpeller) Fix(c_arr *[][]corrections, str ...*string) error {
	if len(*c_arr) > len(str) {
		return errors.New("YandexSpeller.Fix(): array with " +
			"spellchecks is bigger than number of " +
			"supplied strings")
	}

	for i, corrects := range *c_arr {
		text := *str[i]

		for j := len(corrects) - 1; j >= 0; j-- {
			c := corrects[j]
			if len(c.Sgst) > 0 {
				text = applySuggestion(text, c)
			}
		}

		*str[i] = text
	}

	return nil
}

// Dirty cyrillic problem solution via runes
func applySuggestion(text string, c corrections) string {
	runes := []rune(text)
	suggestion := []rune(c.Sgst[0])

	start := c.Pos
	end := start + len([]rune(c.Word))

	left := runes[:start]
	right := runes[end:]

	fixed := append(left, append(suggestion, right...)...)
	final := string(fixed)
	return final
}
