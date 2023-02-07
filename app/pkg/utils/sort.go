package utils

import "sort"

func Sort[T string | int](a []T) {
	var list []T

	allKeys := make(map[T]bool)
	for _, item := range a {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}
