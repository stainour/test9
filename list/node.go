package list

type Node struct {
	Data string
	Next *Node
	Prev *Node
	Rand *Node // произвольный элемент внутри списка
}
