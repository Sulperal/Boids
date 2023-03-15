package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

// 根据视野范围的的其他boid计算移动速率(平均速度)
func (b *Boid) calcAcceleration() Vector2D {
	//lower是圆形视野(以自身的的坐标为圆心)的外切正方形的左下角，upper是右上角
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)

	//平均速率,和平均位置（视野内的boid都向这个点聚集），分离
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}

	//视野中有多少的boid
	count := 0.0

	//为boidMap 加上读写锁，这里是读锁，一个线程拿到读锁后，其他协程也可以拿到读锁，
	rWlock.RLock()
	//这里的Max和Min的作用视野范围再屏幕的边缘和四个角落做出限制
	for i := math.Max(lower.x, 0); i < math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j < math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				//圆形的视野范围
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					//累加平均速率
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					//fmt.Println("1")
					//从群体中分离
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{b.borderBounce(b.position.x, screenWidth),
		b.borderBounce(b.position.y, screenHeight)}

	//视野范围内存在其他的boids
	if count > 0 {
		//
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		//
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)

		//向avgPosition 以avgVelocity的速度移动
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}
	return accel
}

// 反弹,
func (b *Boid) borderBounce(pos, maxBorderos float64) float64 {
	//如果boid的位置靠经屏幕的左下角，反方向移动
	if pos < viewRadius {
		return 1 / pos
		//boid的位置再屏幕的右上角，反弹
	} else if pos > maxBorderos-viewRadius {
		return 1 / (pos - maxBorderos)
	}
	return 0

}
func (b *Boid) moveOne() {
	//因为calcAcceleration中也使用了锁，所以不能再放到下面的锁中
	acceleration := b.calcAcceleration()

	//另一种是写锁，一协程拿到写锁后，其他的协程不能在拿到写锁和读锁,写锁就是一般意义上的锁
	rWlock.Lock()
	//根据相邻的boid进行velocity 的改变，但是有个限制,让移动看着更加平滑
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)

	//移动,先修改boidMap为-1，然后更新boid移动后的boidMap
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	//根据速率移动
	b.position = b.position.Add(b.velocity)

	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	//如果碰到屏幕的边界，反弹
	/* 	next := b.position.Add(b.velocity)
	   	//x轴上
	   	if next.x >= screenWidth || next.x < 0 {
	   		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	   	}
	   	//y轴
	   	if next.y >= screenHeight || next.y < 0 {
	   		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	   	} */

	rWlock.Unlock()
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
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b

	//将boidMap值更新为对应的boid的id
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	////每一个boid都会通过协程移动，这时boidMap 这块共享内存会发生条件竞争，需要使用锁
	go b.start()
}
