package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Camera struct {
	AspectRatio     float64
	ImageWidth      int
	SamplesPerPixel int
	MaxDepth        int

	imageHeight  int
	pixelSamples float64
	center       Point3
	pixel00      Point3
	pixelDeltaU  Vec3
	pixelDeltaV  Vec3
}

func (c *Camera) initialize() {
	c.imageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.imageHeight < 1 {
		c.imageHeight = 1
	}

	c.pixelSamples = float64(c.ImageWidth * c.imageHeight)
	c.center = NewVec3(0, 0, 0)

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(c.ImageWidth) / float64(c.imageHeight)

	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

	viewportUL := c.center.
		Sub(NewVec3(0, 0, focalLength)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))

	c.pixelDeltaU = viewportU.Div(float64(c.ImageWidth))
	c.pixelDeltaV = viewportV.Div(float64(c.imageHeight))
	c.pixel00 = viewportUL.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Scale(0.5))
}

func (c Camera) GetRay(u, v float64) Ray {
	pixelCenter := c.pixel00.
		Add(c.pixelDeltaU.Scale(u)).
		Add(c.pixelDeltaV.Scale(v))
	return Ray{Origin: c.center, Direction: pixelCenter.Sub(c.center)}
}

func (c Camera) rayColor(r Ray, depth int, world Hittable) Color {
	if depth <= 0 {
		return NewVec3(0, 0, 0)
	}

	var rec HitRecord
	if world.Hit(r, Interval{0, math.MaxFloat64}, &rec) {
		return rec.Normal.Add(NewVec3(1, 1, 1)).Scale(0.5)
	}

	unitDir := r.Direction.Unit()
	a := 0.5 * (unitDir.Y + 1.0)
	return NewVec3(1, 1, 1).Scale(1.0-a).
		Add(NewVec3(0.5, 0.7, 1.0).Scale(a))
}

func (c Camera) Render(world Hittable) {
	c.initialize()

	if c.SamplesPerPixel < 1 {
		c.SamplesPerPixel = 1
	}

	fmt.Printf("P3\n%d %d\n255\n", c.ImageWidth, c.imageHeight)
	for j := 0; j < c.imageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", c.imageHeight-j)
		for i := 0; i < c.ImageWidth; i++ {
			pixelColor := NewVec3(0, 0, 0)
			for s := 0; s < c.SamplesPerPixel; s++ {
				u := float64(i) + rand.Float64()
				v := float64(j) + rand.Float64()
				ray := c.GetRay(u, v)
				pixelColor = pixelColor.Add(c.rayColor(ray, c.MaxDepth, world))
			}
			pixelColor = pixelColor.Div(float64(c.SamplesPerPixel))
			pixelColor.WriteColor(os.Stdout)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
