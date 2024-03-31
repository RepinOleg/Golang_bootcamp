package test

import (
	ex00 "day05/ex00"
	"day05/ex01"
	"day05/ex02"
	"day05/ex03"
	"reflect"
	"testing"
)

func TestAreToysBalanced(t *testing.T) {
	//    0
	//   / \
	//  0   1
	// / \
	// 0   1
	tree := ex00.NewTreeNode(false)
	tree.Left = ex00.NewTreeNode(false)
	tree.Left.Left = ex00.NewTreeNode(false)
	tree.Left.Right = ex00.NewTreeNode(true)
	tree.Right = ex00.NewTreeNode(true)

	if !ex00.AreToysBalanced(tree) {
		t.Errorf("Expected true for a balanced tree")
	}
	//		 1
	//		/ \
	//     0   1
	//     	  /
	//        1
	tree = ex00.NewTreeNode(true)
	tree.Left = ex00.NewTreeNode(false)
	tree.Right = ex00.NewTreeNode(true)
	tree.Right.Left = ex00.NewTreeNode(true)

	if ex00.AreToysBalanced(tree) {
		t.Error("Expected false for a unbalanced tree")
	}

	//  1
	// / \
	//1   0
	tree = ex00.NewTreeNode(true)
	tree.Left = ex00.NewTreeNode(true)
	tree.Right = ex00.NewTreeNode(false)

	if ex00.AreToysBalanced(tree) {
		t.Error("Expected false for a unbalanced tree")
	}

	//   	 0
	// 		/ \
	//		1  0
	// 		\   \
	//  	1   1

	tree = ex00.NewTreeNode(false)
	tree.Left = ex00.NewTreeNode(true)

	tree.Left.Right = ex00.NewTreeNode(true)
	tree.Right = ex00.NewTreeNode(false)
	tree.Right.Right = ex00.NewTreeNode(true)

	if ex00.AreToysBalanced(tree) {
		t.Errorf("Expected false for a unbalanced tree")
	}
}

func TestUnrollGarland(t *testing.T) {
	root := &ex00.TreeNode{
		HasToy: true,
		Left: &ex00.TreeNode{
			HasToy: true,
			Left:   &ex00.TreeNode{HasToy: true},
			Right:  &ex00.TreeNode{HasToy: false},
		},
		Right: &ex00.TreeNode{
			HasToy: false,
			Left:   &ex00.TreeNode{HasToy: true},
			Right:  &ex00.TreeNode{HasToy: true},
		},
	}

	expected := []bool{true, true, false, true, true, false, true}
	result := ex01.UnrollGarland(root)

	if len(result) != len(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
		return
	}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected, result)
			return
		}
	}
}

func TestGetNCoolestPresents(t *testing.T) {
	h := make([]ex02.Present, 4)
	h[0].Value = 5
	h[0].Size = 1
	h[1].Value = 4
	h[1].Size = 5
	h[2].Value = 3
	h[2].Size = 1
	h[3].Value = 5
	h[3].Size = 2
	var res []ex02.Present

	t.Run("GetNCoolestPresents_1", func(t *testing.T) {
		expected := make([]ex02.Present, 4)
		expected[0].Value = 5
		expected[0].Size = 1
		expected[1].Value = 5
		expected[1].Size = 2
		expected[2].Value = 4
		expected[2].Size = 5
		expected[3].Value = 3
		expected[3].Size = 1
		res, _ = ex02.GetNCoolestPresents(h, 4)
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Incorrect results: %v != %v", expected, res)
		}
	})

	t.Run("GetNCoolestPresents_2", func(t *testing.T) {
		expected := make([]ex02.Present, 1)
		expected[0].Value = 5
		expected[0].Size = 1
		res, _ = ex02.GetNCoolestPresents(h, 1)
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Incorrect results: %v != %v", expected, res)
		}
	})

	t.Run("GetNCoolestPresents_3", func(t *testing.T) {
		presents := make([]ex02.Present, 3)
		presents[0].Value = 1
		presents[1].Value = 1
		presents[2].Value = 1
		presents[0].Size = 1
		presents[1].Size = 2
		presents[2].Size = 3

		expected := make([]ex02.Present, 3)
		expected[0].Value = 1
		expected[0].Size = 1
		expected[1].Value = 1
		expected[1].Size = 2
		expected[2].Value = 1
		expected[2].Size = 3
		res, _ = ex02.GetNCoolestPresents(presents, 3)
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Incorrect results: %v != %v", res, expected)
		}
	})

	t.Run("GetNCoolestPresents_3", func(t *testing.T) {
		presents := make([]ex02.Present, 1)
		presents[0].Value = 1
		presents[0].Size = 1

		res, err := ex02.GetNCoolestPresents(presents, 3)
		if res != nil || err == nil {
			t.Error("Expected error")
		}
	})
}

func TestGrabPresents(t *testing.T) {
	t.Run("GrabPresents_1", func(t *testing.T) {
		presents := []ex02.Present{
			{Value: 5, Size: 1},
			{Value: 4, Size: 5},
			{Value: 3, Size: 1},
			{Value: 5, Size: 2},
		}
		capacity := 7

		expected := []ex02.Present{
			{Value: 5, Size: 2},
			{Value: 3, Size: 1},
			{Value: 5, Size: 1},
		}
		res := ex03.GrabPresents(presents, capacity)
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Arrays are not equal: %v - %v", res, expected)
		}
	})

	t.Run("GrabPresents_2", func(t *testing.T) {
		presents := []ex02.Present{
			{Value: 5, Size: 1},
			{Value: 4, Size: 5},
			{Value: 3, Size: 1},
			{Value: 5, Size: 2},
		}
		capacity := 8
		res := ex03.GrabPresents(presents, capacity)
		expected := []ex02.Present{
			{Value: 5, Size: 2},
			{Value: 4, Size: 5},
			{Value: 5, Size: 1},
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Arrays are not equal: %v - %v", res, expected)
		}
	})

}
