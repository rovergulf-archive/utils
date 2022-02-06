package utils

func VerifyBitwiseComparison(toCheck uint64, values ...uint64) bool {
	for i := range values {
		if toCheck&values[i] > 0 {
			return true
		}
	}
	// if no result return false
	return false
}
