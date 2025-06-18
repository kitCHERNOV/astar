package main

import (
	// "container/heap"
	"fmt"
	"image/color"

	// "math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	// "gonum.org/v1/plot/plotter/moreland"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// ... (все предыдущие структуры Node, OpenList, ClosedList, Grid остаются без изменений)

// GridData представляет данные для отображения сетки в виде тепловой карты
type GridData struct {
    grid *Grid
    path []*Node
}

// Dims возвращает размеры сетки для HeatMap
func (gd GridData) Dims() (c, r int) {
    return gd.grid.Width, gd.grid.Height
}

// Z возвращает значение для каждой ячейки сетки
func (gd GridData) Z(c, r int) float64 {
    // Инвертируем Y координату для правильного отображения
    y := gd.grid.Height - 1 - r
    
    // Проверяем, является ли клетка частью пути
    for _, node := range gd.path {
        if node.Position.x == c && node.Position.y == y {
            return 0.5 // Путь - серый цвет
        }
    }
    
    // Проверяем препятствия
    if gd.grid.IsObstacle(Point{c, y}) {
        return 1.0 // Препятствие - черный цвет
    }
    
    return 0.0 // Свободная клетка - белый цвет
}

// X возвращает X координату для ячейки
func (gd GridData) X(c int) float64 {
    return float64(c)
}

// Y возвращает Y координату для ячейки
func (gd GridData) Y(r int) float64 {
    return float64(r)
}

// PlotGrid создает графическое отображение сетки с найденным путем
func PlotGrid(grid *Grid, path []*Node, filename string) error {
    // Создаем новый график
    p := plot.New()
    p.Title.Text = "A* Pathfinding Visualization"
    p.X.Label.Text = "X Coordinate"
    p.Y.Label.Text = "Y Coordinate"
    
    // Настройка размера графика
    p.X.Min = -0.5
    p.X.Max = float64(grid.Width) - 0.5
    p.Y.Min = -0.5
    p.Y.Max = float64(grid.Height) - 0.5
    
    // Создаем тепловую карту для основы сетки
    gridData := GridData{grid: grid, path: path}
    hm := plotter.NewHeatMap(gridData, nil)
    
    // Настройка цветовой схемы
    // colors := []color.Color{
    //     color.RGBA{255, 255, 255, 255}, // Белый для свободных клеток
    //     color.RGBA{128, 128, 128, 255}, // Серый для пути
    //     color.RGBA{0, 0, 0, 255},       // Черный для препятствий
    // }
	hm.Palette = palette.Heat(256, 1.0)
    
    p.Add(hm)
    
    // Добавляем линию пути если он существует
    if len(path) > 0 {
        pathPoints := make(plotter.XYs, len(path))
        for i, node := range path {
            pathPoints[i].X = float64(node.Position.x)
            pathPoints[i].Y = float64(grid.Height - 1 - node.Position.y) // Инвертируем Y
        }
        
        line, err := plotter.NewLine(pathPoints)
        if err != nil {
            return err
        }
        line.Color = color.RGBA{0, 0, 255, 255} // Синий цвет для пути
        line.Width = vg.Points(3)
        p.Add(line)
        
        // Добавляем маркер старта
        startPoint := plotter.XYs{{
            X: float64(path[0].Position.x),
            Y: float64(grid.Height - 1 - path[0].Position.y),
        }}
        startScatter, err := plotter.NewScatter(startPoint)
        if err != nil {
            return err
        }
        startScatter.GlyphStyle.Color = color.RGBA{0, 255, 0, 255} // Зеленый для старта
        startScatter.GlyphStyle.Radius = vg.Points(8)
        startScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[1]
        p.Add(startScatter)
        
        // Добавляем маркер цели
        goalPoint := plotter.XYs{{
            X: float64(path[len(path)-1].Position.x),
            Y: float64(grid.Height - 1 - path[len(path)-1].Position.y),
        }}
        goalScatter, err := plotter.NewScatter(goalPoint)
        if err != nil {
            return err
        }
        goalScatter.GlyphStyle.Color = color.RGBA{255, 0, 0, 255} // Красный для цели
        goalScatter.GlyphStyle.Radius = vg.Points(8)
        goalScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[2]
        p.Add(goalScatter)
        
        // Добавляем легенду
        p.Legend.Add("Path", line)
        p.Legend.Add("Start", startScatter)
        p.Legend.Add("Goal", goalScatter)
        p.Legend.Top = true
    }
    
    // Добавляем сетку для лучшей видимости
    p.Add(plotter.NewGrid())
    
    // Сохраняем график
    return p.Save(8*vg.Inch, 8*vg.Inch, filename)
}

// Альтернативная версия с более детальным отображением
func PlotGridDetailed(grid *Grid, path []*Node, filename string) error {
    p := plot.New()
    p.Title.Text = "A* Pathfinding - Detailed View"
    p.X.Label.Text = "X Coordinate"
    p.Y.Label.Text = "Y Coordinate"
    
    // Настройка размера графика
    p.X.Min = -0.5
    p.X.Max = float64(grid.Width) - 0.5
    p.Y.Min = -0.5
    p.Y.Max = float64(grid.Height) - 0.5
    
    // Создаем отдельные scatter plots для разных типов клеток
    var obstacles, freeCells, pathCells plotter.XYs
    
    // Создаем карту пути для быстрого поиска
    pathMap := make(map[string]bool)
    for _, node := range path {
        key := fmt.Sprintf("%d,%d", node.Position.x, node.Position.y)
        pathMap[key] = true
    }
    
    // Заполняем точки для каждого типа клеток
    for x := 0; x < grid.Width; x++ {
        for y := 0; y < grid.Height; y++ {
            plotY := float64(grid.Height - 1 - y) // Инвертируем Y
            key := fmt.Sprintf("%d,%d", x, y)
            
            if pathMap[key] {
                pathCells = append(pathCells, plotter.XY{
                    X: float64(x),
                    Y: plotY,
                })
            } else if grid.IsObstacle(Point{x, y}) {
                obstacles = append(obstacles, plotter.XY{
                    X: float64(x),
                    Y: plotY,
                })
            } else {
                freeCells = append(freeCells, plotter.XY{
                    X: float64(x),
                    Y: plotY,
                })
            }
        }
    }
    
    // Добавляем свободные клетки
    if len(freeCells) > 0 {
        freeScatter, err := plotter.NewScatter(freeCells)
        if err != nil {
            return err
        }
        freeScatter.GlyphStyle.Color = color.RGBA{240, 240, 240, 255}
        freeScatter.GlyphStyle.Radius = vg.Points(15)
        freeScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[5]
        p.Add(freeScatter)
    }
    
    // Добавляем препятствия
    if len(obstacles) > 0 {
        obstacleScatter, err := plotter.NewScatter(obstacles)
        if err != nil {
            return err
        }
        obstacleScatter.GlyphStyle.Color = color.RGBA{0, 0, 0, 255}
        obstacleScatter.GlyphStyle.Radius = vg.Points(20)
        obstacleScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[4]
        p.Add(obstacleScatter)
        p.Legend.Add("Obstacles", obstacleScatter)
    }
    
    // Добавляем путь
    if len(pathCells) > 0 {
        pathScatter, err := plotter.NewScatter(pathCells)
        if err != nil {
            return err
        }
        pathScatter.GlyphStyle.Color = color.RGBA{0, 0, 255, 200}
        pathScatter.GlyphStyle.Radius = vg.Points(12)
        pathScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[0]
        p.Add(pathScatter)
        p.Legend.Add("Path", pathScatter)
    }
    
    // Добавляем соединяющую линию для пути
    if len(path) > 0 {
        pathLine := make(plotter.XYs, len(path))
        for i, node := range path {
            pathLine[i].X = float64(node.Position.x)
            pathLine[i].Y = float64(grid.Height - 1 - node.Position.y)
        }
        
        line, err := plotter.NewLine(pathLine)
        if err != nil {
            return err
        }
        line.Color = color.RGBA{0, 0, 255, 150}
        line.Width = vg.Points(2)
        p.Add(line)
        
        // Специальные маркеры для старта и цели
        startPoint := plotter.XYs{{
            X: float64(path[0].Position.x),
            Y: float64(grid.Height - 1 - path[0].Position.y),
        }}
        startScatter, err := plotter.NewScatter(startPoint)
        if err != nil {
            return err
        }
        startScatter.GlyphStyle.Color = color.RGBA{0, 255, 0, 255}
        startScatter.GlyphStyle.Radius = vg.Points(15)
        startScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[1]
        p.Add(startScatter)
        p.Legend.Add("Start", startScatter)
        
        goalPoint := plotter.XYs{{
            X: float64(path[len(path)-1].Position.x),
            Y: float64(grid.Height - 1 - path[len(path)-1].Position.y),
        }}
        goalScatter, err := plotter.NewScatter(goalPoint)
        if err != nil {
            return err
        }
        goalScatter.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
        goalScatter.GlyphStyle.Radius = vg.Points(15)
        goalScatter.GlyphStyle.Shape = plotutil.DefaultGlyphShapes[2]
        p.Add(goalScatter)
        p.Legend.Add("Goal", goalScatter)
    }
    
    // Добавляем сетку
    p.Add(plotter.NewGrid())
    p.Legend.Top = true
    
    return p.Save(10*vg.Inch, 10*vg.Inch, filename)
}