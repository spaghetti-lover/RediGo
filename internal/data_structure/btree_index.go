package data_structure

// BTreeNode represents a node in the B+ Tree
type BTreeNode struct {
	Items    []*Item
	Children []*BTreeNode
	IsLeaf   bool
	Parent   *BTreeNode
	Next     *BTreeNode
}

// BTreeIndex implements OrderedIndex using B+ Tree
type BTreeIndex struct {
	Root   *BTreeNode
	Degree int
}

// NewBTreeIndex creates a new B+ Tree index
func NewBTreeIndex(degree int) OrderedIndex {
	return &BTreeIndex{
		Root:   &BTreeNode{IsLeaf: true},
		Degree: degree,
	}
}

func (t *BTreeIndex) Add(score float64, member string) int {
	item := &Item{Score: score, Member: member}

	if len(member) == 0 {
		return 0
	}

	// Find the correct leaf to insert into
	node := t.Root
	for !node.IsLeaf {
		i := 0
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// Check if the member already exists in the leaf node
	for i, existingItem := range node.Items {
		if existingItem.Member == member {
			node.Items[i].Score = score
			return 0 // Updated existing item
		}
	}

	// Member does not exist, insert it into the sorted position
	i := 0
	for i < len(node.Items) && item.CompareTo(node.Items[i]) > 0 {
		i++
	}
	node.Items = append(node.Items[:i], append([]*Item{item}, node.Items[i:]...)...)

	// Split the node if it's over capacity
	if len(node.Items) > t.Degree-1 {
		t.splitNode(node)
	}
	return 1 // Added new item
}

func (t *BTreeIndex) Remove(member string) int {
	// Implementation for remove - simplified for now
	// TODO: Implement proper B+ Tree deletion
	return 0
}

func (t *BTreeIndex) GetRank(member string) int {
	rank := 0

	// Find the first leaf node
	node := t.Root
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse all leaf nodes from the beginning
	for node != nil {
		for _, item := range node.Items {
			if item.Member == member {
				return rank
			}
			rank++
		}
		node = node.Next
	}

	return -1
}

func (t *BTreeIndex) GetByRank(rank int) *Item {
	if rank < 0 {
		return nil
	}

	currentRank := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes
	for node != nil {
		for _, item := range node.Items {
			if currentRank == rank {
				return item
			}
			currentRank++
		}
		node = node.Next
	}

	return nil
}

func (t *BTreeIndex) GetRange(min, max float64) []*Item {
	var result []*Item
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes and collect items in range
	for node != nil {
		for _, item := range node.Items {
			if item.Score >= min && item.Score <= max {
				result = append(result, item)
			}
			if item.Score > max {
				return result
			}
		}
		node = node.Next
	}

	return result
}

func (t *BTreeIndex) GetRangeByRank(start, end int) []*Item {
	var result []*Item
	if start < 0 || end < start {
		return result
	}

	currentRank := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes
	for node != nil {
		for _, item := range node.Items {
			if currentRank >= start && currentRank <= end {
				result = append(result, item)
			}
			if currentRank > end {
				return result
			}
			currentRank++
		}
		node = node.Next
	}

	return result
}

func (t *BTreeIndex) Count() int {
	count := 0
	node := t.Root

	// Find the first leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Count items in all leaf nodes
	for node != nil {
		count += len(node.Items)
		node = node.Next
	}

	return count
}

func (t *BTreeIndex) Clear() {
	t.Root = &BTreeNode{IsLeaf: true}
}

// Helper methods for B+ Tree operations
func (t *BTreeIndex) splitNode(node *BTreeNode) {
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BTreeIndex) splitLeaf(node *BTreeNode) {
	medianIndex := len(node.Items) / 2

	newLeaf := &BTreeNode{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
	}

	newLeaf.Items = append(newLeaf.Items, node.Items[medianIndex:]...)
	node.Items = node.Items[:medianIndex]
	node.Next = newLeaf

	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*BTreeNode{newLeaf}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BTreeIndex) splitInternal(node *BTreeNode) {
	medianIndex := len(node.Items) / 2

	newInternal := &BTreeNode{
		IsLeaf: false,
		Parent: node.Parent,
	}

	promotedItem := node.Items[medianIndex]

	newInternal.Items = append(newInternal.Items, node.Items[medianIndex+1:]...)
	newInternal.Children = append(newInternal.Children, node.Children[medianIndex+1:]...)

	node.Items = node.Items[:medianIndex]
	node.Children = node.Children[:medianIndex+1]

	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	parent := node.Parent
	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*BTreeNode{newInternal}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BTreeIndex) splitRoot() {
	oldRoot := t.Root
	newRoot := &BTreeNode{}

	t.Root = newRoot
	oldRoot.Parent = newRoot
	newRoot.Children = append(newRoot.Children, oldRoot)

	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}
