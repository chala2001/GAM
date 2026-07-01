package main

import "fmt"

func main() {
	var routes []string
	routes = append(routes, "/orders/*")
	routes = append(routes, "/payments/*")
	fmt.Println("routes:", routes, "len:", len(routes))

	rateLimits := map[string]int{
		"orders-api":   100,
		"payments-api": 50,
	}
	rateLimits["catalog-api"] = 200

	limit, ok := rateLimits["orders-api"]
	fmt.Println("orders-api limit:", limit, "found:", ok)

	missing, ok := rateLimits["shipping-api"]
	fmt.Println("shipping-api limit:", missing, "found:", ok)

	for api, limit := range rateLimits {
		fmt.Println(api, "->", limit)
	}
}
