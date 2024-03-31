package ex00

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func NewTreeNode(hasToy bool) *TreeNode {
	return &TreeNode{
		hasToy,
		nil,
		nil,
	}
}

func AreToysBalanced(root *TreeNode) bool {
	if root == nil {
		return false
	}
	var countR, countL uint
	// dfs
	// Вызываем рекурсивную функцию
	// которая подсчитывает количество игрушек для левой стороны
	// Затем вызываем такую же для правой и сравниваем результат
	// если равны то дерево сбалансированное

	countL = countToys(root.Left)
	countR = countToys(root.Right)

	// если вызвать эту функцию для root, то она подсчитает количество игрушек во всем дереве
	// allToys := countToys(root)

	return countR == countL
}

func countToys(node *TreeNode) uint {
	if node == nil {
		return 0
	}
	var res uint

	if node.HasToy {
		// инкрементируем счетчик игрушек
		res++
	}
	// вызываем подсчет игрушек пока не дойдем до nil
	res += countToys(node.Left)
	res += countToys(node.Right)

	return res
}
