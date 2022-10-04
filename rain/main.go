package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	title     = "Rain"
	winWidth  = 1024
	winHeight = 612

	maxFPS uint32 = 60

	raindropSpeedMin, raindropSpeedMax float32 = 10, 24
	raindropWidth  int32 = 10
	raindropHeight int32 = 70

	gravityAcceleration, friction float32 = 0.3, 0.98
)

var raindropsCount *int

var (
	window   *sdl.Window
	renderer *sdl.Renderer

	randGen *rand.Rand

	raindrops []Raindrop

	quit = false
)

type Particle struct {
	x, y float32
	size int32

	xVel, yVel, dist float32

	active bool
}

func newParticle(p_x, p_y float32, p_size int32, p_xVel, p_yVel, p_dist float32) Particle {
	return Particle{x: p_x, y: p_y,
	                size: p_size,

	                xVel: p_xVel, yVel: p_yVel,
	                dist: p_dist,

	                active: true}
}

func (p_particle *Particle) render(p_renderer *sdl.Renderer) {
	if !p_particle.active {
		return
	}

	size := float32(p_particle.size) * p_particle.dist

	rect := sdl.Rect{X: int32(p_particle.x - size / 2),
	                 Y: int32(p_particle.y - size / 2),
	                 W: int32(size),
	                 H: int32(size)}

	p_renderer.SetDrawColor(199, 107, 219, uint8(p_particle.dist * 255))
	p_renderer.FillRect(&rect)
}

func (p_particle *Particle) update(p_winHeight int32) {
	if !p_particle.active {
		return
	}

	p_particle.x += p_particle.xVel
	p_particle.y += p_particle.yVel

	p_particle.xVel *= friction
	p_particle.yVel *= friction

	p_particle.yVel += gravityAcceleration

	if int32(p_particle.y) >= p_winHeight {
		p_particle.active = false
	}
}

var particles [2048]Particle

func emitSplashParticles(p_x, p_y float32, p_size int32,
                         p_maxXVel, p_yVel, p_dist float32, p_count int) {
	for i, _ := range particles {
		if particles[i].active {
			continue
		}

		xVel := p_maxXVel - float32(randGen.Intn(int(p_maxXVel * 2)))
		yVel := p_yVel - float32(2 - randGen.Intn(4))

		particles[i] = newParticle(p_x, p_y, p_size, xVel, yVel, p_dist)

		p_count --
		if p_count <= 0 {
			break
		}
	}
}

type Raindrop struct {
	x, y float32
	w, h int32

	yVel, dist float32
}

func newRaindrop(p_winWidth, p_winHeight, p_w, p_h int32, p_yVel, p_dist float32) Raindrop {
	return Raindrop{x:  float32(randGen.Intn(int(p_winWidth))),
	                y: -float32(randGen.Intn(int(p_winHeight))) - float32(p_h) / 2,

	                w: p_w, h: p_h,
	                yVel: p_yVel,
	                dist: p_dist}
}

func (p_raindrop *Raindrop) render(p_renderer *sdl.Renderer) {
	width  := float32(p_raindrop.w) * (p_raindrop.dist)
	height := float32(p_raindrop.h) * (p_raindrop.dist)

	rect := sdl.Rect{X: int32(p_raindrop.x - float32(p_raindrop.w) / 2),
	                 Y: int32(p_raindrop.y - float32(p_raindrop.h) / 2),
	                 W: int32(width),
	                 H: int32(height)}

	p_renderer.SetDrawColor(199, 107, 219, uint8(p_raindrop.dist * 255))
	p_renderer.FillRect(&rect)
}

func (p_raindrop *Raindrop) update(p_winWidth, p_winHeight int32) {
	p_raindrop.y += p_raindrop.yVel

	if int32(p_raindrop.y - float32(p_raindrop.h) / 2) > p_winHeight {
		emitSplashParticles(p_raindrop.x, float32(winHeight - p_raindrop.w),
		                    p_raindrop.w, 10, -4 * p_raindrop.dist, p_raindrop.dist, 10)

		p_raindrop.y = -float32(p_raindrop.h) / 2
		p_raindrop.x = float32(randGen.Intn(int(p_winWidth)))
	}
}

func init() {
	if raindropSpeedMin >= raindropSpeedMax {
		panic("min raindrop speed >= max raindrop speed")
	}

	source := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(source)

	raindropsCount := flag.Int("raindrops_count", 35, "The count of all raindrops")

	flag.Parse()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	var err error
	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
	                               winWidth, winHeight, sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "nearest")
	if err := renderer.SetLogicalSize(winWidth, winHeight); err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	if err := renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		sdl.Log("%s", sdl.GetError())
		panic(err)
	}

	for i := 0; i < *raindropsCount; i ++ {
		dist  := float32(randGen.Intn(1000)) / 1000
		speed := raindropSpeedMin + dist * (raindropSpeedMax - raindropSpeedMin)

		raindrops = append(raindrops, newRaindrop(winWidth, winHeight,
		                                          raindropWidth, raindropHeight, speed, dist))
	}
}

func finish() {
	renderer.Destroy()
	window.Destroy()

	sdl.Quit()
}

func render() {
	renderer.SetDrawColor(7, 17, 36, sdl.ALPHA_OPAQUE)
	renderer.Clear()

	for _, raindrop := range raindrops {
		raindrop.render(renderer)
	}

	for _, particle := range particles {
		particle.render(renderer)
	}

	renderer.Present()
}

func input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent: quit = true
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE: quit = true
			}
		}
	}
}

func update() {
	for i, _ := range raindrops {
		raindrops[i].update(winWidth, winHeight)
	}

	for i, _ := range particles {
		particles[i].update(winHeight)
	}
}

func main() {
	for !quit {
		render()
		input()
		update()

		sdl.Delay(1000 / maxFPS)
	}

	finish()
}
