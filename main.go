package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cellsScale = 8
	cellsWidth = 80
	cellsHeight = 80	
	simStart = ebiten.KeySpace
	simStop = ebiten.KeyEscape
)

type Game struct{
	worldDisplay *image.RGBA
	cellsGrid [cellsHeight][cellsWidth]Cell
	simRunning bool
}

type Cell struct {
	alive bool
	placed bool
}

func (g *Game) Update() error {

	if !g.simRunning {
		cx, cy := getCellPosFromCursorPos(ebiten.CursorPosition())
		for i := range g.cellsGrid {
			for j := range g.cellsGrid[i] {
				if ( i == cy && j == cx ) {
					g.cellsGrid[i][j].alive = true
					if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
						g.cellsGrid[i][j].placed = true
					}
				} else {
					if !g.cellsGrid[i][j].placed {
						g.cellsGrid[i][j].alive = false
					}
				}
				g.updatePixel(i, j, &g.cellsGrid[i][j])
			}
		}
		if ebiten.IsKeyPressed(simStart) {
			g.simRunning = true
		}
	} else {
		for i := range g.cellsGrid {
			for j := range g.cellsGrid[i] {
				neighbours := g.getNeighbours(i, j)
				aliveNeighbours := 0
				for _, n := range neighbours {
					if n.alive {
						aliveNeighbours++
					}
				}
				if g.cellsGrid[i][j].alive {
					if aliveNeighbours < 2 || aliveNeighbours > 3 {
						g.cellsGrid[i][j].alive = false
					}
				} else {
					if aliveNeighbours == 3 {
						g.cellsGrid[i][j].alive = true
					}
				}
			}
		}

		for i := range g.cellsGrid {
			for j := range g.cellsGrid[i] {
				g.updatePixel(i, j, &g.cellsGrid[i][j])
			}
		}
		if ebiten.IsKeyPressed(simStop) {
			g.simRunning = false
		
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.worldDisplay.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return cellsWidth, cellsHeight
}

func main() {
	ebiten.SetWindowSize(cellsWidth * cellsScale, cellsHeight * cellsScale)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{
		worldDisplay: image.NewRGBA(image.Rect(0, 0, cellsWidth, cellsHeight)),
		cellsGrid: [cellsHeight][cellsWidth]Cell{},
	}); err != nil {
		log.Fatal(err)
	}
}

func getCellPosFromCursorPos(x, y int) (int, int) {
	 var cx, cy = x , y 
	 if ( cx < 0 ) {
		cx = 0
	 } else if ( cx > cellsWidth ) {
		cx = cellsWidth
	 }
	 if ( cy < 0 ) {
		cy = 0
	 } else if ( cy > cellsHeight ) {
		cy = cellsHeight
	 }
	 
	 return cx, cy
}

func (g *Game) updatePixel(row int, col int, cell *Cell) {
	pixelIndex := row * cellsWidth + col
	
	c := uint8(0x0)
	if cell.alive { 
		c = 0xff
	}

	//set rgba value of corresponding pixel in displayImage
	g.worldDisplay.Pix[pixelIndex*4] = c //R
	g.worldDisplay.Pix[pixelIndex*4+1] = c //G
	g.worldDisplay.Pix[pixelIndex*4+2] = c //B
	g.worldDisplay.Pix[pixelIndex*4+3] = 0xff //A
}

func (g *Game) getNeighbours(row int, col int) []Cell {
	neighbours := []Cell{}
	above := row - 1
	below := row + 1
	left := col - 1
	right := col + 1

	if above >= 0 {
		neighbours = append(neighbours, g.cellsGrid[above][col])
		if left >= 0 {
			neighbours = append(neighbours, g.cellsGrid[above][left])
		}
		if right < cellsWidth {
			neighbours = append(neighbours, g.cellsGrid[above][right])
		}
	}
	if below < cellsHeight {
		neighbours = append(neighbours, g.cellsGrid[below][col])
		if left >= 0 {
			neighbours = append(neighbours, g.cellsGrid[below][left])
		}
		if right < cellsWidth {
			neighbours = append(neighbours, g.cellsGrid[below][right])
		}
	}
	if left >= 0 {
		neighbours = append(neighbours, g.cellsGrid[row][left])
	}
	if right < cellsWidth {
		neighbours = append(neighbours, g.cellsGrid[row][right])
	}
	return neighbours
}



