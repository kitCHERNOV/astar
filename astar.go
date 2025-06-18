package main

import (
	// "container/heap"
	"container/heap"
	"fmt"
	"math"
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

			// Проверяем есть ли сосед в открытом списке
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
    grid := NewGrid(10, 10)
    
    // Добавляем препятствия
    obstacles := [][2]int{
        {2, 2}, {2, 3}, {2, 4}, {2, 5},
        {3, 5}, {4, 5}, {5, 5},
        {7, 1}, {7, 2}, {7, 3}, {7, 4},
    }
    
    for _, obs := range obstacles {
        grid.AddObstacle(Point{obs[0], obs[1]})
    }
    
    fmt.Println("Поиск пути с помощью алгоритма A*")
    fmt.Println("==================================")
    
    // Ищем путь от (0,0) до (9,9)
    // startX, startY := 0, 0
    // goalX, goalY := 9, 9
    start := Point{x: 0, y: 0,}
	goal := Point{x: 9, y: 9}
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
    
    // Простая версия с тепловой картой
    err = PlotGrid(grid, path, "astar_heatmap.png")
    if err != nil {
        fmt.Printf("Ошибка создания тепловой карты: %v\n", err)
    } else {
        fmt.Println("Тепловая карта сохранена как: astar_heatmap.png")
    }
    
    // Детальная версия
    err = PlotGridDetailed(grid, path, "astar_detailed.png")
    if err != nil {
        fmt.Printf("Ошибка создания детального графика: %v\n", err)
    } else {
        fmt.Println("Детальный график сохранен как: astar_detailed.png")
    }
    
    // Также сохраняем текстовую версию для сравнения
    fmt.Println("\nТекстовая версия:")
    PrintGridText(grid, path)
}

// Переименованная оригинальная функция для сравнения
func PrintGridText(grid *Grid, path []*Node) {
    pathMap := make(map[string]bool)
    for _, node := range path {
        key := fmt.Sprintf("%d,%d", node.Position.x, node.Position.y)
        pathMap[key] = true
    }
    
    fmt.Println("Сетка с найденным путем (текст):")
    for y := grid.Height - 1; y >= 0; y-- {
        for x := 0; x < grid.Width; x++ {
            key := fmt.Sprintf("%d,%d", x, y)
            
            if pathMap[key] {
                if x == path[0].Position.x && y == path[0].Position.y {
                    fmt.Print("S ") // Старт
                } else if x == path[len(path)-1].Position.x && y == path[len(path)-1].Position.y {
                    fmt.Print("G ") // Цель
                } else {
                    fmt.Print("* ") // Путь
                }
            } else if grid.IsObstacle(Point{x, y}) {
                fmt.Print("# ") // Препятствие
            } else {
                fmt.Print(". ") // Свободная клетка
            }
        }
        fmt.Println()
    }
}


