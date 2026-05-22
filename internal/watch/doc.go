// Package watch implements a polling watcher that periodically re-syncs
// secrets from HashiCorp Vault into local .env files.
//
// # Usage
//
// Create a Watcher by providing any Syncer implementation and a Config:
//
//		w, err := watch.New(syncer, watch.Config{
//			Interval: 30 * time.Second,
//			OnError: func(err error) {
//				log.Printf("sync failed: %v", err)
//			},
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// Blocks until ctx is cancelled.
//		if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//			log.Fatal(err)
//		}
//
// The watcher performs an immediate sync on start, then re-syncs on every
// tick. Errors are forwarded to OnError without stopping the loop, so
// transient Vault unavailability does not terminate the process.
package watch
