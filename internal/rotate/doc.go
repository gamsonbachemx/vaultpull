// Package rotate implements secret rotation for vaultpull.
//
// It detects secrets that have exceeded their configured TTL and writes
// freshly generated values back to HashiCorp Vault, ensuring that local
// .env files remain in sync with up-to-date credentials.
//
// Basic usage:
//
//	cfg := &rotate.Config{
//		Paths:     []string{"secret/myapp"},
//		TTL:       24 * time.Hour,
//		Generator: func(key string) (string, error) {
//			// produce a new secret value for key
//			return generateNewValue(key)
//		},
//	}
//
//	r, err := rotate.New(cfg, vaultClient)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	n, err := r.Rotate(ctx, secrets)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("rotated %d secret(s)\n", n)
package rotate
