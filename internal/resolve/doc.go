// Package resolve implements Vault secret path resolution for vaultpull.
//
// It supports two resolution mechanisms:
//
//  1. Variable interpolation using ${VAR} syntax, where variables are
//     supplied at construction time via WithVars.
//
//  2. Path aliasing, where short human-readable names are mapped to full
//     Vault paths via WithAliases. Aliases are expanded before variable
//     interpolation, so aliases may themselves contain ${VAR} placeholders.
//
// Example usage:
//
//	r := resolve.New(
//		resolve.WithVars(map[string]string{"ENV": "production"}),
//		resolve.WithAliases(map[string]string{
//			"myapp": "secret/${ENV}/myapp/config",
//		}),
//	)
//	path, err := r.Resolve("myapp")
//	// path == "secret/production/myapp/config"
package resolve
