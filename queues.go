package main

import (
	"container/heap"
	"fmt"
)

type Point struct {
	x, y int
}

type Node struct {
	Position Point
	GCost float64 // g(n) - фактическое расстояние от старта до текущей позиции
	HCost float64 // h(n) - эвристическая оценка от текущей позиции до цели
	FCost float64 // f(n) = g(n) + h(n) - общая оценка качества пути через эту позицию
	Parent *Node  // Указатель на родительский узел, чтобы востановить путь
	Index int 	  // Индекс в куче для heap.Fix
}

func (n *Node) String() string {
	return fmt.Sprintf("(%d, %d)", n.Position.x, n.Position.y)
}

// implemention priority queue
type OpenList []*Node
 
func (ol OpenList) Len() int {
	return len(ol)
} 

func (ol OpenList) Less(i, j int) bool {
	return ol[i].FCost < ol[j].FCost
}

func (ol OpenList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
	ol[i].Index = i
	ol[j].Index = j
}

func (ol *OpenList) Push(x interface{}) {
	n := len(*ol)
	node := x.(*Node)
	node.Index = n
	*ol = append(*ol, node)
}

func (ol *OpenList) Pop() interface{} {
	old := *ol
	n := old.Len() - 1
	node := old[n]
	node.Index = -1 // Элемент выбыл из очереди с приоритетом
	*ol = old[0:n]
	return node
}

func (ol *OpenList) Update(node *Node, g, h float64) {
	node.GCost = g
	node.HCost = h
	node.FCost = g + h
	heap.Fix(ol, node.Index) // Обновляем приоритет
}

func (ol *OpenList) Contains(point Point) *Node {
	for _, node := range *ol {
		if node.Position.x == point.x && node.Position.y == point.y {
			return node
		}
	}
	return nil
}


// ===================================================== // 
// Реализация списка обработанных узлов // 

type ClosedList struct {
	nodes map[string]*Node
}

func NewClosedList() *ClosedList {
	return &ClosedList{
		nodes: make(map[string]*Node),
	}
}

func (cl *ClosedList) Add(node *Node) {
	key := fmt.Sprintf("%d,%d", node.Position.x, node.Position.y)
	cl.nodes[key] = node
}

func (cl *ClosedList) Contains(point Point) bool {
	key := fmt.Sprintf("%d,%d", point.x, point.y)
	_, exists := cl.nodes[key]
	return exists
}




