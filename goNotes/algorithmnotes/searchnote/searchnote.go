package main

import "fmt"

func binarysearch(target int, source []int) int {

	posStart := 0
	posEnd := len(source) - 1

	for posStart < posEnd {
		pos := (posStart + posEnd) / 2
		if target > source[pos] {
			posStart = pos + 1
		} else {
			posEnd = pos
		}
	}

	if target == source[posStart] {
		return posStart
	} else {
		return -1
	}
}

func binarysearch_test() {
	arr1 := []int{1}
	arr2 := []int{19, 41}
	arr3 := []int{34, 61, 93}
	arr4 := []int{20, 61, 66, 70, 84}
	arr5 := []int{3, 21, 37, 53, 56, 59, 60, 75, 82, 94}

	a1 := []int{0, 1, 2}
	for i := 0; i < len(a1); i++ {
		fmt.Println(binarysearch(a1[i], arr1))
	}
	fmt.Println("------------------")

	a2 := []int{18, 19, 20, 41, 42}
	for i := 0; i < len(a2); i++ {
		fmt.Println(binarysearch(a2[i], arr2))
	}
	fmt.Println("------------------")

	a3 := []int{18, 34, 35, 61, 62, 93, 94}
	for i := 0; i < len(a3); i++ {
		fmt.Println(binarysearch(a3[i], arr3))
	}
	fmt.Println("------------------")

	a4 := []int{19, 20, 21, 61, 62, 66, 67, 70, 71, 84, 85}
	for i := 0; i < len(a4); i++ {
		fmt.Println(binarysearch(a4[i], arr4))
	}

	fmt.Println("------------------")

	a5 := []int{2, 3, 4, 21, 22, 37, 38, 53, 54, 56, 57, 59, 60, 61, 75, 76, 82, 83, 94, 95}
	for i := 0; i < len(a5); i++ {
		fmt.Println(binarysearch(a5[i], arr5))
	}
}

func main() {

}
