package main

import (
	"fmt"
)

type Stack struct {
	s     []string
	limit int
}

func (s *Stack) Push(ss string) {
	if len(s.s) < s.limit {
		s.s = append(s.s, ss)
	} else {
		s.s = s.s[len(s.s)-1:]
		s.s = append(s.s, ss)
	}

}

func (s *Stack) Pop() string {
	d := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return d
}

func main() {
	var s *Stack = &Stack{limit: 2}
	s.Push("data1")
	s.Push("data2")
	s.Push("data3")
	s.Push("data4")
	s.Pop()
	fmt.Println(s)
}
