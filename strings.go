package utils

type StringArray []string

func (s StringArray) Len() int { return len(s) }

func (s StringArray) Less(i, j int) bool { return s[i] < s[j] }

func (s StringArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func RemoveStrDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		if key != "" {
			result = append(result, key)
		}
	}
	return result
}

func removeSpecifiedStringFromSlice(elements []string, element string) []string {
	var clean StringArray
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedStringFromSlice(elements []string, toRemove ...string) []string {
	var clean StringArray
	if len(toRemove) > 0 {
		for i := range toRemove {
			clean = removeSpecifiedStringFromSlice(elements, toRemove[i])
		}
		return clean
	}

	return clean
}

func UnquoteString(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
