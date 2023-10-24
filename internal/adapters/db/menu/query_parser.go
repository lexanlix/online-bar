package menu_db

import (
	"fmt"
	"go/scanner"
	"go/token"
	"reflect"
	"strconv"
)

type lexer struct {
	scan  scanner.Scanner
	token token.Token
	pos   token.Pos
	lit   string
}

func (lex *lexer) next() {
	lex.pos, lex.token, lex.lit = lex.scan.Scan()
}

func (lex *lexer) consume(want token.Token) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, required %q", lex.lit, want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case token.IDENT:
		if lex.lit == "" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}

		var str string

		for i := 0; !endString(lex); i++ {
			if lex.token == token.SUB {
				lex.next()
				str += fmt.Sprintf("-%s", lex.lit)
				lex.next()
				continue
			}

			if i > 0 {
				str += fmt.Sprintf(" %s", lex.lit)
			} else {
				str += lex.lit
			}
			lex.next()
		}

		v.SetString(str)

		if lex.token == token.COMMA {
			lex.consume(token.COMMA)
		}
		return
	case token.INT:
		i, _ := strconv.Atoi(lex.lit)
		v.SetUint(uint64(i))
		lex.next()
		if lex.token == token.COMMA {
			lex.consume(token.COMMA)
		}
		return
	case token.LPAREN:
		lex.next()
		readList(lex, v)
		lex.next()
		if lex.token == token.COMMA {
			lex.consume(token.COMMA)
		}
		return
	case token.LBRACE:
		lex.next()
		readList(lex, v)
		lex.next()
		if lex.token == token.COMMA {
			lex.consume(token.COMMA)
		}
		return
	}
	panic(fmt.Sprintf("unexpected lex: %q", lex.lit))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Slice:
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))

			if lex.token == token.COMMA {
				lex.consume(token.COMMA)
			}
		}
	case reflect.Struct:
		for i := 0; !endStruct(lex); i++ {
			read(lex, v.FieldByIndex([]int{i}))
		}
	default:
		panic(fmt.Sprintf("cannot decode list to %v", v.Type()))
	}
}

func readStruct(lex *lexer, v reflect.Value) {
	lex.next()
	for i := 0; !endStruct(lex); i++ {
		read(lex, v.Field(i))
	}
}

func endString(lex *lexer) bool {
	switch lex.token {
	case token.COMMA:
		return true
	}
	return false
}

func endStruct(lex *lexer) bool {
	switch lex.token {
	case token.EOF:
		panic("end of file")
	case token.RPAREN:
		return true
	}
	return false
}

func endList(lex *lexer) bool {
	switch lex.token {
	case token.EOF:
		panic("end of file")
	case token.RBRACE:
		return true
	}
	return false
}

func UnmarshalQueryRow(data string, out interface{}) (err error) {
	lex := &lexer{
		scan: scanner.Scanner{},
	}

	dataBytes := []byte(data)
	file := token.NewFileSet().AddFile("tokenFile", 1, len(dataBytes))

	lex.scan.Init(file, dataBytes, nil, 2)

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error in %d: %v", lex.pos, x)
		}
	}()

	lex.next()
	readStruct(lex, reflect.ValueOf(out).Elem())
	lex.next()

	return nil
}
