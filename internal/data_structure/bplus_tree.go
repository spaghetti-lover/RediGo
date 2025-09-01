package data_structure

func (i *Item) CompareTo(other *Item) int {
	if i.Score < other.Score {
		return -1
	}
	if i.Score > other.Score {
		return 1
	}

	//If Score is equal. Use Member
	if i.Member < other.Member {
		return -1
	}
	if i.Member > other.Member {
		return 1
	}
	return 0
}

type Node struct {
	Items    []*Item // A list of key-value pairs
	Children []*Node // A list of children nodes
	IsLeaf   bool    //True if it is a leaf node
	Parent   *Node   //Pointer to the parent node
	Next     *Node   //For leaf nodes, a pointer to the next leaf in the sequence
}

type BPlusTree struct {
	Root   *Node
	Degree int
}

func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		Root:   &Node{IsLeaf: true},
		Degree: degree,
	}
}

func (t *BPlusTree) Add(score float64, member string) int {
	item := &Item{Score: score, Member: member}

	if len(member) == 0 {
		return 0
	}

	//Find the correct leaf to insert into
	node := t.Root
	for !node.IsLeaf {
		// Find the correct child based on the score
		i := 0
		for i < len(node.Items) && score >= node.Items[i].Score {
			i += 1
		}
		node = node.Children[i]
	}

	// Check if the member already exists in the leaf node.
	for i, existingItem := range node.Items {
		if existingItem.Member == member {
			node.Items[i].Score = score
			return 1
		}
	}

	// Member does not exist, insert it into the sorted position.
	i := 0
	for i < len(node.Items) && score >= node.Items[i].Score {
		i += 1
	}
	node.Items = append(node.Items[:i], append([]*Item{item}, node.Items[i:]...)...)

	// Split the node if it's over capacity
	if len(node.Items) > t.Degree-1 {
		t.splitNode(node)
	}
	return 1
}

func (t *BPlusTree) splitNode(node *Node) {
	// If the node is the root, we need to create a new root.
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	// Split based on whether the node is a leaf or an internal node.
	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BPlusTree) GetScore(member string) (float64, bool) {
	node := t.Root
	// Traverse to the first leaf node
	// We have to search all leaf nodes since we don't know the score.
	for !node.IsLeaf {
		node = node.Children[0] // Always go to the leftmost child
	}

	// Iterate through all leaf nodes using the 'Next' pointer
	for node != nil {
		// Iterate through all items in the current leaf node
		for _, item := range node.Items {
			if item.Member == member {
				return item.Score, true
			}
		}
		node = node.Next
	}

	return 0, false // Member not found
}

func (t *BPlusTree) GetRank(member string) int {
	rank := 0

	// Find the first leaf node
	node := t.Root
	for !node.IsLeaf {
		node = node.Children[0] //Always go to the leftmost child
	}

	// Traverse all leaf nodes from the beginning
	for node != nil {
		for _, item := range node.Items {
			// Check if we found the member
			if item.Member == member {
				return rank
			}
			rank += 1
		}
		node = node.Next
	}

	return -1
}

func (t *BPlusTree) splitLeaf(node *Node) {
	medianIndex := len(node.Items) / 2

	// Create a new sibling node.
	newLeaf := &Node{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
	}

	// Move the second half of the times to the new leaf
	newLeaf.Items = append(newLeaf.Items, node.Items[medianIndex:]...)
	node.Items = node.Items[:medianIndex]

	// Update the 'Next' pointer for sequential traversal
	node.Next = newLeaf

	// Promote the first key of the new leaf to the parent.
	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	// Find the insertion point in the parent node.
	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex += 1
	}

	// Insert the promoted key and the new child node into the parent
	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newLeaf}, parent.Children[childIndex+1:]...)...)

	// If the parent now overflows, split it too.
	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitInternal(node *Node) {
	medianIndex := len(node.Items) / 2

	// Create a new sibling internal node.
	newInternal := &Node{
		IsLeaf: false,
		Parent: node.Parent,
	}

	// Promote the median key to the parent.
	promotedItem := node.Items[medianIndex]

	// Move the second half of the items and children to the new node.
	newInternal.Items = append(newInternal.Items, node.Items[medianIndex+1:]...)
	newInternal.Children = append(newInternal.Children, node.Children[medianIndex+1:]...)

	// Trim the original node.
	node.Items = node.Items[:medianIndex]
	node.Children = node.Children[:medianIndex+1]

	// Update parent pointers for the new children.
	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	// Now, insert the promoted key and the new node into the parent.
	parent := node.Parent

	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex += 1
	}

	// Insert the promoted key and the new child node into the parent.
	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newInternal}, parent.Children[childIndex+1:]...)...)

	// If the parent now overflows, split it too.
	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitRoot() {
	oldRoot := t.Root
	newRoot := &Node{}

	// Create a new root and set the old root as its first child.
	t.Root = newRoot
	oldRoot.Parent = newRoot
	newRoot.Children = append(newRoot.Children, oldRoot)

	// Split the old root.
	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}
