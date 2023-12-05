package godns

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/root4loot/goutils/log"
	"github.com/root4loot/publicresolvers"
)

type Runner struct {
	Options Options     // options for the runner
	Results chan Result // channel to receive results
}

// Options contains options for the runner
type Options struct {
	Concurrency int      // number of concurrent requests
	Timeout     int      // timeout in seconds
	Delay       int      // delay in seconds
	DelayJitter int      // delay jitter in seconds
	Verbose     bool     // verbose logging
	Resolvers   []string // resolvers to use
	Protocol    string   // protocol to use
}

// Result contains the DNS resolution result for a domain.
type Result struct {
	Domain     string
	IPv4       []string
	IPv6       []string
	ResolvedBy string
}

func init() {
	log.Init("godns")
}

// DefaultOptions returns default options
func DefaultOptions() *Options {
	publicresolvers, err := publicresolvers.FetchResolversTrusted()
	if err != nil {
		// Use fallback resolvers
		publicresolvers = []string{"8.8.8.8", "8.8.4.4", "208.67.222.222", "208.67.220.220"}
	}

	return &Options{
		Concurrency: 10,
		Timeout:     5,
		Delay:       0,
		DelayJitter: 0,
		Resolvers:   publicresolvers,
		Protocol:    "udp",
	}
}

// NewRunner returns a new runner
func NewRunner() *Runner {
	options := DefaultOptions()
	return &Runner{
		Results: make(chan Result),
		Options: *options,
	}
}

// NewRunnerWithOptions returns a new runner with the specified options
func NewRunnerWithOptions(options Options) *Runner {
	return &Runner{
		Results: make(chan Result),
		Options: options,
	}
}

// Single resolves a single domain and returns the result
func Single(host string) (result Result) {
	r := NewRunner()
	r.Options.Concurrency = 1
	return r.worker(host)
}

// Multiple resolves multiple domains and returns the results
func (r *Runner) Multiple(hosts []string) (results []Result) {
	if r.Options.Concurrency > len(hosts) {
		r.Options.Concurrency = len(hosts)
	}

	sem := make(chan struct{}, r.Options.Concurrency)
	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)
		sem <- struct{}{}
		go func(h string) {
			defer func() { <-sem }()
			defer wg.Done()
			results = append(results, r.worker(h))
		}(host)
	}

	wg.Wait()
	return
}

// MultipleStream resolves multiple domains and streams the results using channels
func (r *Runner) MultipleStream(results chan<- Result, host ...string) {
	defer close(results)

	if r.Options.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	sem := make(chan struct{}, r.Options.Concurrency)
	var wg sync.WaitGroup
	for _, h := range host {
		sem <- struct{}{}
		wg.Add(1)
		go func(u string) {
			defer func() { <-sem }()
			defer wg.Done()
			results <- r.worker(u)
			time.Sleep(time.Millisecond * 100) // make room for processing results
		}(h)
		time.Sleep(r.getDelay() * time.Millisecond) // delay between requests
	}
	wg.Wait()
}

// worker is the worker that resolves a domain
func (r *Runner) worker(host string) Result {
	log.Debug("Resolving", host)

	var result Result

	dialer := &net.Dialer{
		Timeout: time.Duration(r.Options.Timeout) * time.Second,
	}

	c := &dns.Client{
		Net:     r.Options.Protocol,
		Dialer:  dialer,
		Timeout: time.Duration(r.Options.Timeout) * time.Second,
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)

	// Try each resolver until we get a response or exhaust all the resolvers
	for _, resolver := range r.Options.Resolvers {
		resolver = net.JoinHostPort(resolver, "53")
		// Make DNS request
		respV4, _, errV4 := c.Exchange(m.Copy(), resolver)
		if errV4 != nil {
			continue // Try the next resolver
		}

		// Clear the response for the next query
		m.Answer = nil

		// Make the DNS request with the current resolver for IPv6
		m.SetQuestion(dns.Fqdn(host), dns.TypeAAAA) // Query for IPv6
		respV6, _, errV6 := c.Exchange(m, resolver)
		if errV6 != nil {
			continue // Try the next resolver
		}

		// Extract IPv4 addresses from the response
		if len(respV4.Answer) > 0 {
			for _, ans := range respV4.Answer {
				if a, ok := ans.(*dns.A); ok {
					result.IPv4 = append(result.IPv4, a.A.String())
				}
			}
		}

		// Extract IPv6 addresses from the response
		if len(respV6.Answer) > 0 {
			for _, ans := range respV6.Answer {
				if aaaa, ok := ans.(*dns.AAAA); ok {
					result.IPv6 = append(result.IPv6, aaaa.AAAA.String())
				}
			}
		}

		// Set the domain in the result
		result.Domain = host
		// Set the resolver used in the result
		result.ResolvedBy = resolver

		break // Got successful responses, no need to try other resolvers
	}

	return result
}

// getDelay returns a random delay between Delay and Delay+DelayJitter
func (r *Runner) getDelay() time.Duration {
	if r.Options.DelayJitter != 0 {
		return time.Duration(r.Options.Delay + rand.Intn(r.Options.DelayJitter))
	}
	return time.Duration(r.Options.Delay)
}
