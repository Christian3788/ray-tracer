package main

import (
	"fmt"
	"os"
)

type Ray struct {
	Origin    Point3
	Direction Vec3
}

// hitSphere returns true if ray r hits a sphere of radius `radius` centred at `center`.
func hitSphere(center Point3, radius float64, r Ray) bool {
	oc := center.Sub(r.Origin) // C - Q
	a := Dot(r.Direction, r.Direction)
	b := -2.0 * Dot(r.Direction, oc)
	c := Dot(oc, oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant >= 0
}

func rayColor(r Ray) Color {
	if hitSphere(NewVec3(0, 0, -1), 0.5, r) {
		return NewVec3(1, 0, 0) // red
	}
	unitDir := r.Direction.Unit()
	a := 0.5 * (unitDir.Y + 1.0)
	return NewVec3(1, 1, 1).Scale(1.0 - a).
		Add(NewVec3(0.5, 0.7, 1.0).Scale(a))
}

func main() {
	// Image
	const aspectRatio = 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	// Camera
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	cameraCenter := NewVec3(0, 0, 0)

	// Vectors across the viewport edges
	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0) // negative — Y is flipped

	// Per-pixel deltas
	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	// Top-left pixel location
	viewportUL := cameraCenter.
		Sub(NewVec3(0, 0, focalLength)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))
	pixel00 := viewportUL.Add(pixelDeltaU.Add(pixelDeltaV).Scale(0.5))

	// Render
	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)
	for j := 0; j < imageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", imageHeight-j)
		for i := 0; i < imageWidth; i++ {
			pixelCenter := pixel00.
				Add(pixelDeltaU.Scale(float64(i))).
				Add(pixelDeltaV.Scale(float64(j)))
			rayDir := pixelCenter.Sub(cameraCenter)
			r := Ray{Origin: cameraCenter, Direction: rayDir}

			rayColor(r).WriteColor(os.Stdout)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
