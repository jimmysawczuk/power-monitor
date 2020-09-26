package tmplfunc

func Seq(max int) []int {
	if max < 0 {
		return nil
	}

	v := make([]int, 0, max)
	for i := 0; i < max; i++ {
		v = append(v, i)
	}

	return v
}

func Add(a, b int) int {
	return a + b
}
