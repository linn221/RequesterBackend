package validate

// returns slice removing duplicate elements
func UniqueSlice[T comparable](slice []T) []T {
	inResult := make(map[T]struct{})
	var result []T
	for _, elm := range slice {
		if _, ok := inResult[elm]; !ok {
			// if not exists in map, append it, otherwise do nothing
			inResult[elm] = struct{}{}
			result = append(result, elm)
		}
	}
	return result
}
