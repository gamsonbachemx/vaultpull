package env

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing values when a key already exists.
	StrategyKeepExisting
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged   map[string]string
	Added    []string
	Updated  []string
	Skipped  []string
}

// Merge combines incoming secrets into existing env entries using the given strategy.
// existing may be nil, in which case all incoming keys are treated as added.
func Merge(existing, incoming map[string]string, strategy MergeStrategy) MergeResult {
	if existing == nil {
		existing = make(map[string]string)
	}

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}

	result := MergeResult{Merged: merged}

	for k, v := range incoming {
		oldVal, exists := merged[k]
		switch {
		case !exists:
			merged[k] = v
			result.Added = append(result.Added, k)
		case strategy == StrategyOverwrite && oldVal != v:
			merged[k] = v
			result.Updated = append(result.Updated, k)
		default:
			result.Skipped = append(result.Skipped, k)
		}
	}

	return result
}
