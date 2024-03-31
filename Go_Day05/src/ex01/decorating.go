package ex01

import (
	"day05/ex00"
)

// Нужно пройтись змейкой по дереву и вернуть все значения bool

func UnrollGarland(root *ex00.TreeNode) []bool {
	// создаем слайс указателей на ноды
	var row = []*ex00.TreeNode{root}
	var result []bool

	// флаг для того чтобы идти зиг загом
	var flagRow = true
	for len(row) > 0 {
		// слайс для следующего ряда
		var nextRow []*ex00.TreeNode
		var tmpResults []bool

		// идем в цикле по нодам
		for _, node := range row {
			// складываем значения во временный слайс
			tmpResults = append(tmpResults, node.HasToy)
			// добалвяем в слайс следующего ряда ноды
			if node.Left != nil {
				nextRow = append(nextRow, node.Left)
			}
			if node.Right != nil {
				nextRow = append(nextRow, node.Right)
			}
		}

		// если есть флаг значит нужно идти справа налево
		if flagRow {
			for i := len(tmpResults) - 1; i >= 0; i-- {
				result = append(result, tmpResults[i])
			}
		} else {
			// если нет просто append всех эелементов в res слайс
			result = append(result, tmpResults...)
		}
		// переходим на следующий ряд
		row = nextRow
		// меняем флаг
		flagRow = !flagRow
	}
	return result
}
