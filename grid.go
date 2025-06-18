package main

import "fmt"

type Grid struct {
	Width, Height int
	Obstacles map[string]bool
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Width: width,
		Height: height,
		Obstacles: make(map[string]bool),
	}
}

func (g *Grid) AddObstacle(point Point) {
	key := fmt.Sprintf("%d,%d", point.x, point.y)
	g.Obstacles[key] = true
}

func (g *Grid) IsObstacle(point Point) bool {
	key := fmt.Sprintf("%d,%d", point.x, point.y)
	return g.Obstacles[key]
}

func (g *Grid) IsValid(point Point) bool {
	return point.x >= 0 && point.x < g.Width && point.y >= 0 && point.y < g.Height && !g.IsObstacle(point)
}

// Получение соседей для текущей вершины
func (g *Grid) GetNeighbors(node *Node) []*Node {
	neighbors := make([]*Node, 0)

	directions := [][2]int{
        {0, 1},  // вверх
        {0, -1}, // вниз
        {1, 0},  // вправо
        {-1, 0}, // влево
    }

	for _, dir := range directions {
		newX := node.Position.x + dir[0]
		newY := node.Position.y + dir[1]
		
		if g.IsValid(Point{newX, newY}) {
			neighbor := &Node{
				Position: Point{newX, newY},
			}
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}