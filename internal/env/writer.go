package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write serializes the provided secrets map into a .env file.
// Keys are sorted alphabetically for deterministic output.
// Existing file contents are overwritten.
func (w *Writer) Write(secrets map[string]string) error {
	if len(secrets) == 0 {
		return nil
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		// Quote values that contain spaces or special characters.
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(w.filePath, []byte(sb.String()), 0600)
}

// FormatKey converts a Vault secret path + field into an env-style key.
// e.g. "myapp/database" + "password" -> "MYAPP_DATABASE_PASSWORD"
func FormatKey(namespace, field string) string {
	raw := namespace + "_" + field
	replacer := strings.NewReplacer("/", "_", "-", "_", ".", "_")
	return strings.ToUpper(replacer.Replace(raw))
}
