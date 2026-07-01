package main

import "fmt"

type Quota struct {
	Remaining int
}

func (q Quota) TryConsumeByValue() {
	q.Remaining--
}

func (q *Quota) TryConsumeByPointer() {
	q.Remaining--
}

func main() {
	q := Quota{Remaining: 100}

	q.TryConsumeByValue()
	fmt.Println("after value-receiver call:", q.Remaining)

	q.TryConsumeByPointer()
	fmt.Println("after pointer-receiver call:", q.Remaining)

	qPtr := &q
	fmt.Println("dereferenced:", *qPtr)
}
