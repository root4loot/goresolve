![Go version](https://img.shields.io/badge/Go-v1.19-blue.svg) [![Contribute](https://img.shields.io/badge/Contribute-Welcome-green.svg)](CONTRIBUTING.md)

# goresolve

Quickly resolve domains using reliable resolvers ([publicresolvers](https://github.com/root4loot/publicresolvers))

## Example
```
go get github.com/root4loot/goresolve@master
```

```go
package main

import (
	"fmt"
	"strings"

	"github.com/root4loot/goresolve"
)

func main() {
	// Options
	options := goresolve.DefaultOptions()
	options.Concurrency = 5
	options.Timeout = 5
	options.Delay = 0
	options.DelayJitter = 0
	options.Resolvers = []string{"208.67.222.222", "208.67.220.220"}
	
	r := goresolve.NewRunnerWithOptions(*options)

	// Single domain
	fmt.Println("Single:")

	result := goresolve.Single("example.com")
	fmt.Printf("Domain: %s\n", result.Domain)

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

	// Multiple domains
	fmt.Println("\nMultiple:")

	results := r.Multiple([]string{"example.com", "google.com", "github.com"})
	for _, result := range results {
		fmt.Printf("Domain: %s\n", result.Domain)
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

	// Multiple domains using channels
	fmt.Println("\nMultipleStream")

	streamResults := make(chan goresolve.Result)
	go r.MultipleStream(streamResults, "example.com", "google.com", "github.com")
	for result := range streamResults {
		fmt.Printf("Domain: %s\n", result.Domain)
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
```

---

## Contributing

Contributions are welcome. If you find any bugs or have suggestions for improvements, feel free to open an issue or submit a pull request.
