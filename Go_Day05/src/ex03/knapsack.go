package ex03

import (
	"day05/ex02"
)

// функция получает слайс подарков и вместимость рюкзака

func GrabPresents(presents []ex02.Present, capacity int) []ex02.Present {
	n := len(presents) // количество подарков
	// двумерный массив
	dp := make([][]int, n+1)

	// инициализируем внутренний массив
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= capacity; j++ {
			// Если размер текущего подарка превышает текущую вместимость j,
			// то значение ценности остается таким же, как у предыдущего подарка.
			if presents[i-1].Size > j {
				dp[i][j] = dp[i-1][j]
			} else {
				// Иначе выбираем максимум между ценностью, которую можно получить,
				// включив текущий подарок, и ценностью без него.
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-presents[i-1].Size]+presents[i-1].Value)
			}
		}
	}

	selectedPresents := make([]ex02.Present, 0)
	i, j := n, capacity
	for i > 0 && j > 0 {
		// Если значение ценности изменилось при добавлении текущего подарка,
		// то добавляем его в список выбранных подарков и обновляем вместимость.
		if dp[i][j] != dp[i-1][j] {
			selectedPresents = append(selectedPresents, presents[i-1])
			j -= presents[i-1].Size
		}
		i--
	}

	return selectedPresents
}
