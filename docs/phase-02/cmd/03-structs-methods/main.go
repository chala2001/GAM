package main

import "fmt"

type RegisteredAPI struct {
	Name     string
	Upstream string
	Path     string
}

func (a RegisteredAPI) Describe() string {
	return fmt.Sprintf("%s -> %s%s", a.Name, a.Upstream, a.Path)
}

func main() {
	orders := RegisteredAPI{
		Name:     "orders-api",
		Upstream: "http://orders-service:8081",
		Path:     "/orders/*",
	}
	fmt.Println(orders.Describe())

	var empty RegisteredAPI
	fmt.Printf("zero value struct: %+v\n", empty)
}
