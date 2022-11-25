package moxa

import (
	"fmt"
	"testing"
)

func TestParallelMap(t *testing.T) {
	args := []int{1, 2, 3, 4, 5}
	results, err := ParallelMap(args, func(i int) (int, error) {
		if i == 2 {
			return -1, fmt.Errorf("error2")
		} else if i == 4 {
			return -1, fmt.Errorf("error4")
		}
		return i, nil
	})
	fmt.Printf("%+v\n", results)
	fmt.Printf("%+v\n", err)
}
