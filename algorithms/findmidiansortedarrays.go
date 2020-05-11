package algorithm

func findMidianSortedArrays(nums1 []int, nums2 []int) float64 {

	nums1Len := len(nums1)
	nums2Len := len(nums2)
	totalLen := nums1Len + nums2Len

	if nums1Len == 0 {
		return findMidianSortedArray(nums2)
	}

	if nums2Len == 0 {
		return findMidianSortedArray(nums1)
	}

	if totalLen%2 != 0 {
		return float64(findKthSortedArrays(nums1, nums2, totalLen/2+1))
	}

	midNum1 := findKthSortedArrays(nums1, nums2, totalLen/2)
	midNum2 := findKthSortedArrays(nums1, nums2, totalLen/2+1)
	return float64(midNum1)/2 + float64(midNum2)/2
}

func findMidianSortedArray(nums []int) float64 {

	numsLen := len(nums)

	if numsLen == 0 {
		return -1
	}

	if numsLen%2 != 0 {
		return float64(nums[numsLen/2])
	}

	mid1 := nums[numsLen/2-1]
	mid2 := nums[numsLen/2]
	return float64(mid1)/2 + float64(mid2)/2
}

func findKthSortedArrays(nums1 []int, nums2 []int, k int) int {

	start1, end1 := 0, len(nums1)-1
	start2, end2 := 0, len(nums2)-1

	for k/2 != 0 && start1 <= end1 && start2 <= end2 {

		nextIdx1 := end1
		if k/2 < end1-start1+1 {
			nextIdx1 = start1 + k/2 - 1
		}

		nextIdx2 := end2
		if k/2 < end2-start2+1 {
			nextIdx2 = start2 + k/2 - 1
		}

		if nums1[nextIdx1] > nums2[nextIdx2] {
			start2 = start2 + k/2 - 1
		}
	}
}
