package utils

func VerifyBitwiseComparison(check uint64, roles ...uint64) bool {
	if len(roles) > 0 {
		// loop through roles to find a match
		for i := range roles {
			if check&roles[i] > 0 {
				return true
			}
		}
		// if no result return an error
		return false
	} else {
		return true
	}
}
