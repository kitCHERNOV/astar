package main

import (
	"container/heap"
	"fmt"
	"math"
    "math/rand"
)


func HeuristicFunc(point1, point2 Point) float64 {
	return math.Abs(float64(point1.x - point2.x)) +  math.Abs(float64(point1.y - point2.y))
}

func ReconstructPath(node *Node) []*Node {
    path := make([]*Node, 0)
    current := node
    
    for current != nil {
        path = append([]*Node{current}, path...)
        current = current.Parent
    }
    
    return path
}

func AStar(grid *Grid, start, goal Point) ([]*Node, error) {
	// Проверка существования точек
	if !grid.IsValid(start) {
		return nil, fmt.Errorf("start point (%d,%d) isn't available", start.x, start.y)
	}
	if !grid.IsValid(goal) {
		return nil, fmt.Errorf("goal point (%d,%d) isn't available", goal.x, goal.y)
	}

	// Start 
	openList := &OpenList{}
	heap.Init(openList)
	closedList := NewClosedList()

	// Create start node
	startNode := &Node{
		Position: start,
		GCost: 0,
		HCost: HeuristicFunc(start, goal),
	}

	startNode.FCost = startNode.GCost + startNode.HCost
	
	heap.Push(openList, startNode)

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*Node)

		if current.Position == goal {
			return ReconstructPath(current), nil
		}

		// текущая точка уже пройдена
		closedList.Add(current)

		// поиск соседних точек
		neighbors := grid.GetNeighbors(current)

		for _, neighbor := range neighbors {
			// Не рассматриваем уже закрытые узлы
			if closedList.Contains(neighbor.Position) {
				continue
			}

			// Вычисляем значение G от рассматриваой точки (neighbor)
			tentativeG := current.GCost + 1.0 // тк растояние между точками 1
			
			// TODO: Написать функцию добавления единичного расстояния.

			// Есть ли сосед в открытом списке
			existingNode := openList.Contains(neighbor.Position)

			if existingNode == nil {
				// Добавим точки в список краевых точек
				neighbor.GCost = tentativeG
				neighbor.HCost = HeuristicFunc(neighbor.Position, goal)
				neighbor.FCost = neighbor.GCost + neighbor.HCost
				neighbor.Parent = current

				heap.Push(openList, neighbor)
			}else if tentativeG < existingNode.GCost {
				existingNode.GCost = tentativeG
				existingNode.FCost = existingNode.GCost + existingNode.HCost
				existingNode.Parent = current
				
				// Обновление позиции в куче
				heap.Fix(openList, existingNode.Index)
			}
		}
	}

	// Если открытый пуст, то путь не найден
	return nil, fmt.Errorf("путь от (%d,%d) до (%d,%d) отсутсвует", start.x, start.y, goal.x, goal.y)
}

// Обновленная функция main с графической визуализацией
func main() {
    // Создаем сетку 10x10
    // grid := NewGrid(10, 10)
    
    // // Добавляем препятствия
    // obstacles := [][2]int{
    //     {2, 2}, {2, 3}, {2, 4}, {2, 5},
    //     {3, 5}, {4, 5}, {5, 5},
    //     {7, 1}, {7, 2}, {7, 3}, {7, 4},
    // }
    
    // for _, obs := range obstacles {
    //     grid.AddObstacle(Point{obs[0], obs[1]})
    // }
    // ================================================== //

     // Создаем большую сетку 100x100
    grid := NewGrid(55, 55)
    
    // Стартовая и целевая точки
    startX, startY := 0, 0
    goalX, goalY := 99, 99
    
    // 1. Случайные препятствия (20% плотность)
    for x := 0; x < 100; x++ {
        for y := 0; y < 100; y++ {
            if (x != startX || y != startY) && (x != goalX || y != goalY) {
                if rand.Float64() < 0.2 {
                    grid.AddObstacle(Point{x, y})
                }
            }
        }
    }
    
    // 2. Диагональные стены
    for i := 10; i < 30; i++ {
        grid.AddObstacle(Point{i, i})
        grid.AddObstacle(Point{i, 90-i})
    }
    
    // 3. Вертикальные коридоры
    for y := 20; y < 80; y++ {
        grid.AddObstacle(Point{25, y})
        grid.AddObstacle(Point{50, y})
        grid.AddObstacle(Point{75, y})
    }
    
    // 4. Горизонтальные барьеры с проходами
    for x := 10; x < 90; x++ {
        if x%15 != 0 { // оставляем проходы каждые 15 клеток
            grid.AddObstacle(Point{x, 30})
            grid.AddObstacle(Point{x, 60})
        }
    }

    // createMaze(grid, 40, 40, 20, 20)

    // ================================================== //
    fmt.Println("Поиск пути с помощью алгоритма A*")
    fmt.Println("==================================")
    
    // Ищем путь от (0,0) до (9,9)
    // startX, startY := 0, 0
    // goalX, goalY := 9, 9
    start := Point{x: 0, y: 0,}
	goal := Point{x: 45, y: 30}
    path, err := AStar(grid, start, goal)
    
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    
    fmt.Printf("Путь найден! Длина пути: %d шагов\n", len(path)-1)
    fmt.Printf("Маршрут: ")
    for i, node := range path {
        if i > 0 {
            fmt.Print(" -> ")
        }
        fmt.Printf("(%d,%d)", node.Position.x, node.Position.y)
    }
    fmt.Println()
    
    // Создаем графические визуализации
    fmt.Println("\nСоздание графической визуализации...")
    
    // Детальная версия
    // err = PlotGridDetailed(grid, path, "astar_detailed.png")
    err = PlotGrid(grid, path, "astar_colored.png")
    if err != nil {
        fmt.Printf("Ошибка создания детального графика: %v\n", err)
    } else {
        fmt.Println("Детальный график сохранен как: astar_detailed.png")
    }
}


// func createMaze(grid *Grid, startX, startY, width, height int) {
//     for x := startX; x < startX+width; x += 3 {
//         for y := startY; y < startY+height; y += 3 {
//             // Создаем блоки с проходами
//             grid.AddObstacle(Point{x, y})
//             grid.AddObstacle(Point{x+1, y})
//             grid.AddObstacle(Point{x, y+1})
//         }
//     }
// }
