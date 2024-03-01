package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "net"
    "os"
    "sync"
)

func worker(subdomains <-chan string, resolver *net.Resolver, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()
    for subdomain := range subdomains {
        txtRecords, err := resolver.LookupTXT(context.Background(), "default._bimi."+subdomain)
        if err != nil {
            continue
        }
        for _, txt := range txtRecords {
            results <- fmt.Sprintf("%s: %s", subdomain, txt)
        }
    }
}

var concurrency int
var customDNS string
var dnsPort string

func main() {
    flag.StringVar(&customDNS, "dns", "", "Custom DNS resolver address (ip:port)")
    flag.IntVar(&concurrency, "c", 20, "set the concurrency level")
    flag.StringVar(&dnsPort, "port", "53", "DNS server port")
    flag.Parse()

    subdomains := make(chan string, concurrency)
    results := make(chan string, concurrency)

    var wg sync.WaitGroup

    resolver := net.DefaultResolver
    if customDNS != "" {
        resolver = &net.Resolver{
            PreferGo: true,
            Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
                return net.Dial("udp", customDNS+":"+dnsPort)
            },
        }
    }

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go worker(subdomains, resolver, &wg, results)
    }

    go func() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            subdomains <- scanner.Text()
        }
        close(subdomains)
    }()

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        fmt.Println(result)
    }
}
