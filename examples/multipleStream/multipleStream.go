package main

import (
	"fmt"
	"strings"

	"github.com/root4loot/goresolve"
)

func main() {
	options := goresolve.DefaultOptions()
	options.Concurrency = 5
	options.Timeout = 5
	options.Delay = 0
	options.DelayJitter = 0
	options.Resolvers = []string{"208.67.222.222", "208.67.220.220"}

	r := goresolve.NewRunnerWithOptions(*options)

	streamResults := make(chan goresolve.Result)
	go r.MultipleStream(streamResults, "example.com", "google.com", "github.com")
	for result := range streamResults {
		fmt.Printf("Target Domain: %s\n", result.TargetDomain)
		if len(result.IPv4) > 0 {
			fmt.Printf("IPv4: %s\n", strings.Join(result.IPv4, ", "))
		} else {
			fmt.Println("IPv4: None")
		}
		if len(result.IPv6) > 0 {
			fmt.Printf("IPv6: %s\n", strings.Join(result.IPv6, ", "))
		} else {
			fmt.Println("IPv6: None")
		}
		fmt.Println("Resolver:", result.ResolvedBy)
	}
}
