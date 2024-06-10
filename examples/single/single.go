package main

import (
	"fmt"
	"strings"

	"github.com/root4loot/goresolve"
)

func main() {
	result := goresolve.Single("example.com")
	fmt.Printf("Target Domain: %s\n", result.TargetDomain)

	if len(result.IPv4) > 0 {
		fmt.Printf("IPv4: %s\n", strings.Join(result.IPv4, ", "))
	} else {
		fmt.Println("IPV4: None")
	}
	if len(result.IPv6) > 0 {
		fmt.Printf("IPv6: %s\n", strings.Join(result.IPv6, ", "))
	} else {
		fmt.Println("IPv6: None")
	}
	fmt.Println("Resolver:", result.ResolvedBy)
}
