package utils

import (
	"math/rand"
)

// RandomPick 泛型函数，随机从输入列表中选择x个元素返回
func RandomPick[T any](list []T, x int) []T {
	n := len(list)
	if n == 0 {
		return []T{} // 空列表处理
	}

	result := make([]T, 0, x)

	if n >= x {
		// 当列表长度大于等于x时，随机取x个不同的元素
		indices := rand.Perm(n)[:x]
		for _, idx := range indices {
			result = append(result, list[idx])
		}
	} else {
		// 当列表长度小于x时，重复选择，尽量均匀分布
		repeatCount := x / n
		remainder := x % n

		// 每个元素至少出现repeatCount次
		for i := 0; i < repeatCount; i++ {
			result = append(result, list...)
		}

		// 处理剩余的部分，随机取remainder个
		indices := rand.Perm(n)[:remainder]
		for _, idx := range indices {
			result = append(result, list[idx])
		}
	}

	return result
}
