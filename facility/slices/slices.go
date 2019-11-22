package slices

func DistinctInt64(inputs []int64) []int64 {
	d := map[int64]bool{}
	for _, input := range inputs {
		d[input] = true
	}
	results := []int64{}
	for k := range d {
		results = append(results, k)
	}
	return results
}
