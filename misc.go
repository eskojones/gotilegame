package main

import (
	"fmt"
	"strconv"
	"strings"
)

func tokenize(bytes []byte, open string, close string) []string {
	str := open + string(bytes) + close
	var tokens []string
	index := strings.LastIndex(str, open)
	for token := 0; index != -1; token++ {
		index_end := strings.Index(str[index+1:], close)
		if index_end == -1 {
			break
		}
		parens := str[index+1 : index+1+index_end]
		tokens = append(tokens, parens)
		fmt.Printf("%02d \"%s\"\n", token, strings.ReplaceAll(strings.Trim(strings.TrimLeft(parens, " "), " "), "\n", ""))
		str = fmt.Sprintf("%s@token%04x@%s", str[:index], token, str[index+index_end+2:])
		index = strings.LastIndex(str, open)
	}
	return tokens
}

func compile(tokens []string, open string, close string) string {
	var out = fmt.Sprintf("@token%04x@", len(tokens)-1)
	for strings.Contains(out, "@token") {
		start := strings.Index(out, "@token")
		token_id_str := out[start+6 : start+10]
		token_id, _ := strconv.ParseUint(token_id_str, 16, 0)
		// fmt.Printf("'%s' %d %s\n", token_id_str, token_id, tokens[token_id])
		out = fmt.Sprintf("%s%s%s%s%s", out[:start], open, tokens[token_id], close, out[start+11:])
	}
	return out
}
