package utils

type IntArray []int
type Int8Array []int8
type Int16Array []int16
type Int32Array []int32
type Int64Array []int64

type UintArray []uint
type Uint8Array []uint8
type Uint16Array []uint16
type Uint32Array []uint32
type Uint64Array []uint64

type Float32Array []float32
type Float64Array []float64

func (s IntArray) Len() int           { return len(s) }
func (s IntArray) Less(i, j int) bool { return s[i] > s[j] }
func (s IntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Int8Array) Len() int           { return len(s) }
func (s Int8Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Int8Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Int16Array) Len() int           { return len(s) }
func (s Int16Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Int16Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Int32Array) Len() int           { return len(s) }
func (s Int32Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Int32Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Int64Array) Len() int           { return len(s) }
func (s Int64Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Int64Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s UintArray) Len() int           { return len(s) }
func (s UintArray) Less(i, j int) bool { return s[i] > s[j] }
func (s UintArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Uint8Array) Len() int           { return len(s) }
func (s Uint8Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Uint8Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Uint16Array) Len() int           { return len(s) }
func (s Uint16Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Uint16Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Uint32Array) Len() int           { return len(s) }
func (s Uint32Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Uint32Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Uint64Array) Len() int           { return len(s) }
func (s Uint64Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Uint64Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Float32Array) Len() int           { return len(s) }
func (s Float32Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Float32Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Float64Array) Len() int           { return len(s) }
func (s Float64Array) Less(i, j int) bool { return s[i] > s[j] }
func (s Float64Array) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func RemoveIntDuplicatesUnordered(elements []int) []int {
	encountered := map[int]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []int
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func RemoveInt8DuplicatesUnordered(elements []int8) []int8 {
	encountered := map[int8]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []int8
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func RemoveInt16DuplicatesUnordered(elements []int16) []int16 {
	encountered := map[int16]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []int16
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func RemoveInt32DuplicatesUnordered(elements []int32) []int32 {
	encountered := map[int32]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []int32
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func RemoveInt64DuplicatesUnordered(elements []int64) []int64 {
	encountered := map[int64]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []int64
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func removeSpecifiedIntFromSlice(elements []int, element int) []int {
	var clean IntArray
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedIntFromSlice(elements []int, toRemove int, another ...int) []int {
	if len(another) > 0 {
		var clean IntArray
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedIntFromSlice(elements, another[i])
		}
		return clean
	} else {
		return removeSpecifiedIntFromSlice(elements, toRemove)
	}
}

func removeSpecifiedInt8FromSlice(elements []int8, element int8) []int8 {
	var clean Int8Array
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedInt8FromSlice(elements []int8, toRemove int8, another ...int8) []int8 {
	if len(another) > 0 {
		var clean Int8Array
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedInt8FromSlice(elements, another[i])
		}
		return clean
	} else {
		return removeSpecifiedInt8FromSlice(elements, toRemove)
	}
}

func removeSpecifiedInt16FromSlice(elements []int16, element int16) []int16 {
	var clean Int16Array
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedInt16FromSlice(elements []int16, toRemove int16, another ...int16) []int16 {
	if len(another) > 0 {
		var clean Int16Array
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedInt16FromSlice(elements, another[i])
		}
		return clean
	} else {
		return removeSpecifiedInt16FromSlice(elements, toRemove)
	}
}

func removeSpecifiedInt32FromSlice(elements []int32, element int32) []int32 {
	var clean Int32Array
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedInt32FromSlice(elements []int32, toRemove int32, another ...int32) []int32 {
	if len(another) > 0 {
		var clean Int32Array
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedInt32FromSlice(elements, another[i])
		}
		return clean
	} else {
		return removeSpecifiedInt32FromSlice(elements, toRemove)
	}
}

func removeSpecifiedInt64FromSlice(elements []int64, element int64) []int64 {
	var clean Int64Array
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedInt64FromSlice(elements []int64, toRemove int64, another ...int64) []int64 {
	if len(another) > 0 {
		var clean Int64Array
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedInt64FromSlice(elements, another[i])
		}
		return clean
	} else {
		return removeSpecifiedInt64FromSlice(elements, toRemove)
	}
}
