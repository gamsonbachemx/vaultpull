// Package hook implements pre/post sync lifecycle hooks for vaultpull.
//
// Hooks are shell commands configured by the user that run before or after
// a vault sync operation. They are useful for tasks such as:
//
//   - Restarting a service after secrets are refreshed
//   - Notifying a monitoring system of a sync event
//   - Validating environment state before pulling secrets
//
// # Usage
//
//	cfg := hook.Config{
//	    Pre:  []string{"echo starting sync"},
//	    Post: []string{"systemctl reload myapp", "echo done"},
//	}
//	runner, err := hook.New(cfg)
//	if err != nil { ... }
//
//	if err := runner.RunPre(ctx); err != nil { ... }
//	// ... perform sync ...
//	if err := runner.RunPost(ctx); err != nil { ... }
//
// Commands are split on whitespace; shell features like pipes and redirects
// are not supported. Each hook runs sequentially; if any hook fails the
// remaining hooks are skipped and an error is returned.
package hook
