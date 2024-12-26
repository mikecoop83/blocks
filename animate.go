package main

import (
	"image/color"
	"time"
)

func lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}

func lerpColor(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := lerp(float64(r1), float64(r2), t)
	g := lerp(float64(g1), float64(g2), t)
	b := lerp(float64(b1), float64(b2), t)
	a := lerp(float64(a1), float64(a2), t)

	return color.RGBA{
		R: uint8(r / 256),
		G: uint8(g / 256),
		B: uint8(b / 256),
		A: uint8(a / 256),
	}
}

type animatedEntity struct {
	currentColor  color.Color
	targetColor   color.Color
	animationTime time.Duration

	progress float64 // Interpolation progress (0.0 to 1.0)
}

const ticksPerSecond = 60

func (block *animatedEntity) tick() bool {
	// delta is how much the progress should increase so that it takes `animationTime` to reach 1.0
	deltaSizeForOneSecond := 1 / float64(ticksPerSecond)
	delta := deltaSizeForOneSecond / block.animationTime.Seconds()
	if block.progress < 1.0 {
		block.progress += delta // `ticks` is the delta time (e.g., 1.0 / 60.0 for 60 FPS)
		if block.progress > 1.0 {
			block.progress = 1.0
		}
		newColor := lerpColor(block.currentColor, block.targetColor, block.progress)
		block.currentColor = newColor
		return false
	}
	return true
}
