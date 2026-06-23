package main

import (
	"fmt"
	"os"
)

type Ray struct {
	Origin    Point3
	Direction Vec3
}

func rayColor(r Ray) Color {
	unitDir := r.Direction.Unit()
	a := 0.5 * (unitDir.Y + 1.0)
	white := NewVec3(1.0, 1.0, 1.0)
	blue := NewVec3(0.5, 0.7, 1.0)
	return white.Scale(1.0-a).Add(blue.Scale(a))
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

	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	viewportUL := cameraCenter.
		Sub(NewVec3(0, 0, focalLength)).
		Sub(viewportU.Div(2)).
		Sub(viewportV.Div(2))
	pixel00 := viewportUL.Add(pixelDeltaU.Add(pixelDeltaV).Scale(0.5))

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
