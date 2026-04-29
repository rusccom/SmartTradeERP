package search

import "strconv"

func Position(args []any) string {
	return strconv.Itoa(len(args) + 1)
}
