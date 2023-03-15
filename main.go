package main

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 300, 200
	boidCount                 = 300

	viewRadius = 13
	adjRate    = 0.015
)

var (
	green = color.RGBA{10, 255, 50, 255}
	//指针数组，元素是指针
	boids [boidCount]*Boid

	//协程共享区域,二维数组
	boidMap [screenWidth + 1][screenHeight + 1]int
	//锁
	// lock = sync.Mutex{}
	//读写锁
	rWlock = sync.RWMutex{}
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		screen.Set(int(boid.position.x+1), int(boid.position.y), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)

	}

}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	//地图的初始值为-1 ，如果在上面生成一个boid, 将-1改为 boid 的id
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}

	for i := 0; i < boidCount; i++ {

		createBoid(i)
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
