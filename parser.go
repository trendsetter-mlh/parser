package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func Parse(fn string) []byte {

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	a := strings.FieldsFunc(string(b), func(r rune) bool {
		switch r {
		case '\n':
			return true
		case '\t':
			return true
		case ' ':
			return true
		case '.':
			return true
		case ',':
			return true
		case '[':
			return true
		case ']':
			return true
		case '(':
			return true
		case ')':
			return true
		default:
			return false
		}
	})

	varCount := 0
	funcCount := 0

	funcNames := make(map[string]string)
	varNames := make(map[string]string)
	funcNameMap := make(map[string]string)
	varNameMap := make(map[string]string)

	for i, v := range a {
		if i == 0 {
			continue
		}

		for k, n := range funcNames {
			if v == k {
				a[i] = n
				continue
			}
		}

		for k, n := range varNames {
			if v == k {
				a[i] = n
				continue
			}
		}

		if a[i-1] == "var" || a[i-1] == "const" || a[i-1] == "let" || a[i-1] == "this" {
			if varNames[v] == "" {
				newName := "TRENDSET_VAR_" + strconv.Itoa(varCount)
				varNames[v] = newName
				varNameMap[newName] = v
			}
			a[i] = varNames[v]
			varCount++
			continue
		}

		if a[i-1] == "function" && (i > 1 && a[i-2][0] != '=') {
			if funcNames[v] == "" {
				newName := "TRENDSET_FUNC_" + strconv.Itoa(funcCount)
				funcNames[v] = newName
				funcNameMap[newName] = v
			}

			a[i] = funcNames[v]
			funcCount++
			continue
		}

		if a[i-1] != "require" && (v[0] == '"' || v[0] == '\'') {
			a[i] = "TRENDSET_STRING"
		}
	}

	type Resp struct {
		TokenList []string          `json:"token_list"`
		VarNames  map[string]string `json:"var_names"`
		FuncNames map[string]string `json:"func_names"`
	}

	p, err := json.Marshal(Resp{a, varNameMap, funcNameMap})
	if err != nil {
		panic(err)
	}

	return p
}

func main() {
	fmt.Printf("%s", Parse("./markov.js"))
}
