package algorithm

/*
给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那 两个 整数，并返回他们的数组下标。

你可以假设每种输入只会对应一个答案。但是，你不能重复利用这个数组中同样的元素
*/

// O(n)
func twoSum1(nums []int, target int) []int {

	hash := make(map[int]int, len(nums))
	for idx, num := range nums {
		if leftIdx, has := hash[target-num]; has {
			return []int{leftIdx, idx}
		}

		hash[num] = idx
	}

	return nil
}
