package main

import (
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) moveOne() {
	//根据速率移动
	b.position = b.position.Add(b.velocity)

	//如果碰到屏幕的边界，反弹
	next := b.position.Add(b.velocity)
	//x轴上
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	//y轴
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
}

func (b *Boid) start() {
	for {
		b.moveOne()
		//睡5毫秒，给其他协程执行时间
		time.Sleep(5 * time.Millisecond)
	}
}
func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64()*2 - 1.0), (rand.Float64()*2 - 1.0)},
		id:       bid,
	}
	boids[bid] = &b
	go b.start()
}
