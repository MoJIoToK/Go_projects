// Пакет tree создает и сортирует дерево комментариев.

package tree

import (
	"Comments/internal/model"
	"errors"
	"fmt"
)

// Root - структура дерева комментариев.
type Root struct {
	Comments []*Node
}

// Node - структура узла дерева комментариев.
type Node struct {
	model.Comment
	Childs []*Node `json:"Childs"`
}

// Ошибки при работе с пакетом.
var ErrEmptySlice = errors.New("empty comment array")

// Build - функция позволяет строить дерево комментариев. На вход принимается слайс комментариев.
func Build(comments []model.Comment) (Root, error) {
	const operation = "goComments.tree.Build"

	m := make(map[string]*Node)
	root := Root{}

	if len(comments) == 0 {
		return root, fmt.Errorf("%s: %w", operation, ErrEmptySlice)
	}

	for _, comment := range comments {
		node := &Node{Comment: comment}

		n, ok := m[comment.ID]
		if ok {
			node.Childs = n.Childs
		}
		m[comment.ID] = node

		if comment.ParentID == "" {
			root.Comments = append(root.Comments, node)
			continue
		}

		_, ok = m[comment.ParentID]
		if !ok {
			parent := &Node{}
			parent.Childs = append(parent.Childs, node)
			m[comment.ParentID] = parent
			continue
		}
		m[comment.ParentID].Childs = append(m[comment.ParentID].Childs, node)
	}

	return root, nil
}
