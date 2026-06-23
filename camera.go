package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

type Camera struct {
	AspectRatio     float64
	ImageWidth      int
	SamplesPerPixel int
	MaxDepth        int
	VFov            float64 // vertical field of view, degrees
	LookFrom        Point3
	LookAt          Point3
	VUp             Vec3
	DefocusAngle   float64
	FocusDist      float64

	imageHeight   int
	pixelSamples  float64
	center        Point3
	pixel00       Point3
	pixelDeltaU   Vec3
	pixelDeltaV   Vec3
	defocusDiskU  Vec3
	defocusDiskV  Vec3
}

func (c *Camera) initialize() {
	c.imageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.imageHeight < 1 {
		c.imageHeight = 1
	}

	c.pixelSamples = float64(c.ImageWidth * c.imageHeight)
	c.center = c.LookFrom

	theta := c.VFov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2)
	viewportHeight := 2.0 * halfHeight
	viewportWidth := viewportHeight * float64(c.ImageWidth) / float64(c.imageHeight)

	w := c.LookFrom.Sub(c.LookAt).Unit()
	u := Cross(c.VUp, w).Unit()
	v := Cross(w, u)

	horizontal := u.Scale(viewportWidth)
	viewportV := v.Scale(-viewportHeight)

	viewportUL := c.center.
		Sub(w.Scale(c.FocusDist)).
		Sub(horizontal.Div(2)).
		Sub(viewportV.Div(2))

	c.pixelDeltaU = horizontal.Div(float64(c.ImageWidth))
	c.pixelDeltaV = viewportV.Div(float64(c.imageHeight))
	c.pixel00 = viewportUL.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Scale(0.5))

	defocusRadius := c.FocusDist * math.Tan((c.DefocusAngle/2) * math.Pi / 180)
	c.defocusDiskU = u.Scale(defocusRadius)
	c.defocusDiskV = v.Scale(defocusRadius)
}

func (c *Camera) GetRay(i, j int) Ray {
	offsetX := rand.Float64() - 0.5
	offsetY := rand.Float64() - 0.5
	pixelSample := c.pixel00.
		Add(c.pixelDeltaU.Scale(float64(i) + offsetX)).
		Add(c.pixelDeltaV.Scale(float64(j) + offsetY))

	origin := c.center
	if c.DefocusAngle > 0 {
		origin = origin.
			Add(c.defocusDiskU.Scale(offsetX)).
			Add(c.defocusDiskV.Scale(offsetY))
	}

	return Ray{Origin: origin, Direction: pixelSample.Sub(origin)}
}

func (c *Camera) rayColor(r Ray, depth int, world Hittable) Color {
	if depth <= 0 {
		return NewVec3(0, 0, 0)
	}

	var rec HitRecord
	if world.Hit(r, Interval{0, math.MaxFloat64}, &rec) {
		if rec.Mat != nil {
			attenuation, scattered, ok := rec.Mat.Scatter(r, &rec)
			if ok {
				return attenuation.Mul(c.rayColor(scattered, depth-1, world))
			}
			return NewVec3(0, 0, 0)
		}
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

	pixels := make([]Color, c.ImageWidth*c.imageHeight)

	var wg sync.WaitGroup
	workers := runtime.NumCPU()
	rowsPerWorker := (c.imageHeight + workers - 1) / workers
	var done atomic.Int64

	for w := 0; w < workers; w++ {
		startJ := w * rowsPerWorker
		if startJ >= c.imageHeight {
			break
		}
		endJ := startJ + rowsPerWorker
		if endJ > c.imageHeight {
			endJ = c.imageHeight
		}

		wg.Add(1)
		go func(startJ, endJ int) {
			defer wg.Done()
			for j := startJ; j < endJ; j++ {
				for i := 0; i < c.ImageWidth; i++ {
					pixelColor := NewVec3(0, 0, 0)
					for s := 0; s < c.SamplesPerPixel; s++ {
						ray := c.GetRay(i, j)
						pixelColor = pixelColor.Add(c.rayColor(ray, c.MaxDepth, world))
					}
					pixelColor = pixelColor.Div(float64(c.SamplesPerPixel))
					pixels[j*c.ImageWidth+i] = pixelColor
				}
				done.Add(1)
				fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", c.imageHeight-int(done.Load()))
			}
		}(startJ, endJ)
	}

	wg.Wait()

	fmt.Printf("P3\n%d %d\n255\n", c.ImageWidth, c.imageHeight)
	for j := 0; j < c.imageHeight; j++ {
		for i := 0; i < c.ImageWidth; i++ {
			pixels[j*c.ImageWidth+i].WriteColor(os.Stdout)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
