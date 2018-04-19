package main

import "errors"

func example(code string) (int, error) {
	if code == "hoge" {
		return 1, nil
	}
	return 0, errors.New("Code must be hoge")
}

func Sum(a, b int) int {
	return a + b
}
