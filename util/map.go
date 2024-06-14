package util

func Map[In, Out any](inSlice []In, predicate func(In) Out) []Out {
	outSlice := make([]Out, len(inSlice))

	for i := range inSlice {
		outSlice[i] = predicate(inSlice[i])
	}

	return outSlice
}
