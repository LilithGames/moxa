package utils

func Tuple2Arg1[T1, T2 any](a1 T1, a2 T2) T1 {
	return a1
}

func Tuple2Arg2[T1, T2 any](a1 T1, a2 T2) T2 {
	return a2
}
