package main

import "fmt"

func main() {
	var name string
	var requestCount int
	var latencyMs float64
	var isHealthy bool
	fmt.Println("zero values:", name, requestCount, latencyMs, isHealthy)

	var apiName string = "orders-api"
	version := 1
	rateLimit := 100.0
	published := true
	fmt.Println("declared values:", apiName, version, rateLimit, published)

	const maxRetries = 3
	fmt.Println("const:", maxRetries)
}
