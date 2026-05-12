// Package prompt provides interactive terminal confirmation prompts
// for use in CLI workflows such as confirming destructive sync operations
// before overwriting existing .env files.
//
// Example usage:
//
//	c := prompt.New(nil, nil) // defaults to os.Stdin / os.Stdout
//	ok, err := c.Confirm("Overwrite existing .env?", false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !ok {
//		fmt.Println("Aborted.")
//		return
//	}
package prompt
