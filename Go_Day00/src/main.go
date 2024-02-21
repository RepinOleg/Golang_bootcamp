package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	var nums []int
	mp := make(map[int]int)
	nums, err := readInput(mp, nums)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(nums) > 0 {
		sort.Ints(nums)
		err = callFunctions(nums, mp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

}

func callFunctions(nums []int, mp map[int]int) error {
	if len(os.Args) == 1 {
		fmt.Printf("Mean: %.1f\n", mean(nums))
		median(nums)
		mode(mp)
		sDeviation(nums, mean(nums))
	} else {
		for _, arg := range os.Args[1:] {
			if arg == "mean" {
				fmt.Printf("Mean: %.1f\n", mean(nums))
			} else if arg == "median" {
				median(nums)
			} else if arg == "mode" {
				mode(mp)
			} else if arg == "sd" {
				sDeviation(nums, mean(nums))
			} else {
				return fmt.Errorf("wrong argument")
			}
		}
	}
	return nil
}

func median(nums []int) {
	lenSlice := len(nums)
	if lenSlice%2 != 0 {
		fmt.Printf("Median: %.1f\n", float64(nums[lenSlice/2]))
	} else {
		sumTwoMiddleElems := nums[lenSlice/2] + nums[lenSlice/2-1]
		fmt.Printf("Median: %.1f\n", float64(sumTwoMiddleElems/2))
	}
}

func mode(mp map[int]int) {
	minVal := 0
	minKey := math.MaxInt64
	for key, val := range mp {
		if val >= minVal && key < minKey {
			minVal = val
			minKey = key
		}
	}
	fmt.Printf("Mode: %d\n", minKey)
}

func mean(nums []int) float64 {
	var sum float64 = 0
	for _, num := range nums {
		sum += float64(num)
	}
	mean := sum / float64(len(nums))
	return mean
}

func sDeviation(nums []int, mean float64) {
	var dispercia float64
	for _, num := range nums {
		dispercia += math.Pow(float64(num)-mean, 2)
	}
	deviation := math.Sqrt(dispercia / float64(len(nums)-1))
	fmt.Printf("SD: %.2f\n", deviation)
}

func readInput(mp map[int]int, nums []int) ([]int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil || num > 100000 || num < -100000 {
			return nil, fmt.Errorf("wrong number")
		}
		mp[num]++
		nums = append(nums, num)
	}
	return nums, nil
}
