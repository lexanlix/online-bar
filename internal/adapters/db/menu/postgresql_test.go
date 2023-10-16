package menu_db

import (
	"restapi/internal/domain/menu"
	"testing"
)

func Test_GetInsertValue(t *testing.T) {
	dr := map[string][]menu.Drink{}

	s := GetInsertValue(&dr)

	if s != "CAST(ARRAY[] AS DrinksGroup [])" {
		t.Fatalf("empty drink array get incorrect string: %s", s)
	}

}
