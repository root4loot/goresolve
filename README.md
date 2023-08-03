![Go version](https://img.shields.io/badge/Go-v1.19-blue.svg) [![Contribute](https://img.shields.io/badge/Contribute-Welcome-green.svg)](CONTRIBUTING.md)

# GoDNS

GoDNS is a lightweight and easy-to-use Golang library that allows you to make DNS requests for single or multiple domains. It supports both IPv4 and IPv6 address resolutions and provides optional concurrency and timeout settings. GoDNS can also stream results using channels for efficient handling of large DNS resolution tasks.

## Example
```
go get github.com/root4loot/godns@master
```

```go
package main

import (
	"fmt"
	"strings"

	"github.com/root4loot/godns"
)

func main() {
	// Options
	options := godns.DefaultOptions()
	options.Concurrency = 5
	options.Timeout = 5
	options.Delay = 0
	options.DelayJitter = 0
	r := godns.NewRunnerWithOptions(*options)

	// Single domain
	fmt.Println("Single:")

	result := godns.Single("example.com")
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
	}

	// Multiple domains using channels
	fmt.Println("\nMultipleStream")

	streamResults := make(chan godns.Result)
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
	}
}
```

---

## Contributing

Contributions are welcome. If you find any bugs or have suggestions for improvements, feel free to open an issue or submit a pull request.
