package spellers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewYandexSpeller(t *testing.T) {
	speller := NewYandexSpeller()

	speller_type := reflect.TypeOf(speller)
	if speller_type != reflect.TypeOf(&YandexSpeller{}) {
		t.Error("\"speller\" variable should be of type "+
			"\"YandexSpeller\", but is ", speller_type)
	}

	if len(speller.Config) != 3 {
		t.Errorf("Speller's config map does not contain the "+
			"right amount of keys: should have 3, but has %d",
			len(speller.Config))
	}

	for _, v := range []string{"options", "lang", "format"} {
		_, exists := speller.Config[v]
		if !exists {
			t.Errorf("Map key \"%s\" doesn't exist", v)
		}
	}

	if speller.Config["options"] != YA_OPT_IGNORE_DIGITS|
		YA_OPT_IGNORE_URLS {
		t.Errorf("\"options\" key doesn't have the "+
			"right values: must be %d, but is %S",
			YA_OPT_IGNORE_DIGITS|YA_OPT_IGNORE_URLS,
			speller.Config["options"])
	}

	for _, v := range []string{"lang", "format"} {
		if speller.Config["lang"] != "" {
			t.Errorf("\"%s\" key expected to be empty, "+
				"but contains \"%s\"", v,
				speller.Config["lang"])
		}
	}
}

func TestYandexSpeller_Check(t *testing.T) {
	speller := NewYandexSpeller()

	// nil
	_, err := speller.Check(nil, nil)
	if err == nil || err.Error() != "empty text field" {
		t.Fatalf("speller.Check(): expected \"empty text "+
			"field\", but got %s\n", err.Error())
	}

	// empty string
	empty := ""

	_, err = speller.Check(&empty, &empty)
	if err == nil || err.Error() != "empty text field" {
		t.Fatalf("speller.Check(): expected \"empty text "+
			"field\", but got %s\n", err.Error())
	}

	// correct data
	mock_handler := func(w http.ResponseWriter, r *http.Request) {
		response := `[[{"code":1,"pos":12,"row":0,"col":12,"len":6,"word":"\u0442\u0435\u043a\u043c\u0441\u0442","s":["\u0442\u0435\u043a\u0441\u0442","\u0442\u0435\u0441\u0442"]}],[{"code":1,"pos":4,"row":0,"col":4,"len":4,"word":"\u0434\u043e\u0438\u043d","s":["\u043e\u0434\u0438\u043d","\u0434\u043e\u0438\u043d","\u0434\u043e\u043d","\u0434\u0430\u0438\u043d"]},{"code":1,"pos":9,"row":0,"col":9,"len":6,"word":"\u0442\u0435\u043a\u0441\u0441\u0442","s":["\u0442\u0435\u043a\u0441\u0442"]}]]`
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}

	mock_server := httptest.NewServer(
		http.HandlerFunc(mock_handler),
	)
	defer mock_server.Close()

	speller.URL = mock_server.URL
	texts := []string{
		"проверочный текмст",
		"еще доин тексст",
	}

	expected := &[][]corrections{
		{
			{"текмст", []string{"текст", "тест"}, 1, 12, 0, 12},
		},
		{
			{"доин", []string{"один", "доин", "дон", "даин"},
				1, 4, 0, 4},
			{"тексст", []string{"текст"}, 1, 9, 0, 9},
		},
	}

	corrections, err := speller.Check(&texts[0], &texts[1])
	if err != nil {
		t.Fatalf("speller.Check(): %s\n", err)
	}
	if !reflect.DeepEqual(*corrections, *expected) {
		t.Error("speller.Check(): Corrections array isn't equal to expected array")
		t.Log("==== corrections ====")
		t.Log(corrections)
		t.Log()
		t.Log("==== expected ====")
		t.Log(expected)
		t.FailNow()
	}

	// empty URL
	speller.URL = ""
	corrections, err = speller.Check(&texts[0], &texts[1])
	if err == nil {
		t.Fatalf("speller.Check(): expected an error from " +
			"http.PostFrom(), but got nil\n")
	}

	// non-json response
	mock_handler = func(w http.ResponseWriter, r *http.Request) {
		response := `0000000000000000`
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}

	mock_server = httptest.NewServer(
		http.HandlerFunc(mock_handler),
	)
	speller.URL = mock_server.URL

	corrections, err = speller.Check(&texts[0], &texts[1])
	if err == nil {
		t.Fatalf("speller.Check(): Expected an error from " +
			"json.Unmarshal, got nil\n")
	}

	// wrong json response
	mock_handler = func(w http.ResponseWriter, r *http.Request) {
		response := `["data":""]`
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}

	mock_server = httptest.NewServer(
		http.HandlerFunc(mock_handler),
	)
	speller.URL = mock_server.URL

	corrections, err = speller.Check(&texts[0], &texts[1])
	if err == nil {
		t.Fatalf("speller.Check(): Expected an error from " +
			"json.Unmarshal, got nil\n")
	}
}

func TestYandexSpeller_Fix(t *testing.T) {
	speller := NewYandexSpeller()

	texts := []string{
		"проверочный текмст",
		"еще доин тексст",
	}

	valid_texts := []string{
		"проверочный текст",
		"еще один текст",
	}

	expected := &[][]corrections{
		{
			{"текмст", []string{"текст", "тест"}, 1, 12, 0, 12},
		},
		{
			{"доин", []string{"один", "доин", "дон", "даин"},
				1, 4, 0, 4},
			{"тексст", []string{"текст"}, 1, 9, 0, 9},
		},
	}

	// valid case
	err := speller.Fix(expected, &texts[0], &texts[1])
	if err != nil {
		t.Fatalf("speller.Fix(): got an error on a valid "+
			"case: %s\n", err)
	}
	for i, v := range texts {
		if v != valid_texts[i] {
			t.Fatalf("speller.Fix(): fixed string didn't "+
				"match: expected '%s', got '%s'",
				valid_texts[i], v)
		}
	}

	// incorrect number of strings
	err = speller.Fix(expected, &texts[0])
	if err == nil || err.Error() != "YandexSpeller.Fix(): "+
		"array with spellchecks is bigger than number of "+
		"supplied strings" {
		t.Fatalf("speller.Fix(): expected error due to"+
			"number of supplied strings, but got %s\n", err)
	}

}
