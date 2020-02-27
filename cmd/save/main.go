package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-internetarchive/wayback"
	"log"
	"sync"
)

func save(ctx context.Context, m *wayback.WaybackMachine, uris ...string) {

	wg := new(sync.WaitGroup)

	for _, u := range uris {

		wg.Add(1)

		go func(uri string) {

			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			err := m.Save(ctx, uri)

			if err != nil {
				log.Printf("Failed to save '%s', %v\n", uri, err)
			}

		}(u)
	}

	wg.Wait()
}

func main() {

	flag.Parse()

	opts, err := wayback.DefaultWaybackMachineOptions()

	if err != nil {
		log.Fatalf("Failed to create options, %v", err)
	}

	m, err := wayback.NewWaybackMachine(opts)

	if err != nil {
		log.Fatalf("Failed to create wayback machine, %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	uris := flag.Args()

	save(ctx, m, uris...)
}
