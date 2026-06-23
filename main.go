package main

type Ray struct {
	Origin    Point3
	Direction Vec3
}

func main() {
	var world HittableList
	world.Add(Sphere{Center: NewVec3(0, 0, -1), Radius: 0.5})
	world.Add(Sphere{Center: NewVec3(0, -100.5, -1), Radius: 100})

	cam := Camera{
		AspectRatio:     16.0 / 9.0,
		ImageWidth:      400,
		SamplesPerPixel: 100,
		MaxDepth:        50,
	}
	cam.Render(&world)
}
