package algorithm

// find first index, OOXX
func binarySearch1(nums []int, target int) int {

	if len(nums) == 0 {
		return -1
	}

	left, right := 0, len(nums) - 1
	for left + 1 < right {

		mid := (right-left)/2 + left
		if nums[mid] >= target {
			right = mid
		} else {
			left = mid
		}
	}

	if nums[left] == target {
		return left
	}

	if nums[right] == target {
		return right
	}

	return -1
}

// search in rotated array, half half
func binarySearch2(nums []int, target int) int {


}

// binary search result
func sqrt(num int) int {

	if num == 0 {
		return 0
	}

	left, right := 1, num

}
