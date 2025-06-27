package api

import (
	"strconv"
)

func parseFloatParam(query map[string][]string, key string) float64 {
	if val, ok := query[key]; ok {
		if f, err := strconv.ParseFloat(val[0], 64); err == nil {
			return f
		}
	}
	return 0
}

func parseIntParam(query map[string][]string, key string) int {
	if val, ok := query[key]; ok {
		if i, err := strconv.Atoi(val[0]); err == nil {
			return i
		}
	}
	return 0
}
