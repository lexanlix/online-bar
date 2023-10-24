package menu_db

import (
	"restapi/internal/domain/menu"
	"testing"
)

// REFRESH THIS
func Test_EncodeInsertValue(t *testing.T) {
	dr := map[string][]menu.NewDrinkDTO{}
	menu_id := ""

	s := EncodeInsertValue(&dr, menu_id)

	if s != "CAST(ARRAY[] AS DrinksGroup [])" {
		t.Fatalf("empty drink array get incorrect string: %s", s)
	}

}

func Test_DecodeMenuRequest(t *testing.T) {
	dr := "{\"\"}"

	drinks, err := DecodeMenuRequest(dr)
	if err != nil {
		t.Fatalf("parsing empty drinks menu request err: %v", err)
	}

	if drinks != nil {
		t.Fatalf("parsing empty drinks menu error")
	}

}
