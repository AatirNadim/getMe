package utils

func ConvertStringToBytes(str string) []byte {
	return []byte(str)
}

func ConvertBytesToString(b []byte) string {
	return string(b)
}

func DeleteDuplicateKeys(keys []string) []string {
	// Pre-allocate map with expected capacity to prevent rehashing
	uniqueKeys := make(map[string]struct{}, len(keys))

	// Pre-allocate slice with capacity to avoid reallocations during append
	result := make([]string, 0, len(keys))

	for _, key := range keys {
		if _, exists := uniqueKeys[key]; !exists {
			uniqueKeys[key] = struct{}{}
			result = append(result, key) // Appends without reallocation
		}
	}

	return result
}
