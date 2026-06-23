package main

import (
	"math/rand"
	"time"
)

type Ray struct {
	Origin    Point3
	Direction Vec3
}

func randomScene() HittableList {
	var world HittableList
	ground := Lambertian{Albedo: NewVec3(0.5, 0.5, 0.5)}
	world.Add(Sphere{Center: NewVec3(0, -1000, 0), Radius: 1000, Mat: ground})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Length() > 0.9 {
				if chooseMat < 0.8 {
					albedo := RandomVec3().Mul(RandomVec3())
					world.Add(Sphere{Center: center, Radius: 0.2, Mat: Lambertian{Albedo: albedo}})
				} else if chooseMat < 0.95 {
					albedo := NewVec3(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()))
					fuzz := 0.5 * rand.Float64()
					world.Add(Sphere{Center: center, Radius: 0.2, Mat: Metal{Albedo: albedo, Fuzz: fuzz}})
				} else {
					world.Add(Sphere{Center: center, Radius: 0.2, Mat: Dielectric{RefractionIndex: 1.5}})
				}
			}
		}
	}

	world.Add(Sphere{Center: NewVec3(0, 1, 0), Radius: 1.0, Mat: Dielectric{RefractionIndex: 1.5}})
	world.Add(Sphere{Center: NewVec3(-4, 1, 0), Radius: 1.0, Mat: Lambertian{Albedo: NewVec3(0.4, 0.2, 0.1)}})
	world.Add(Sphere{Center: NewVec3(4, 1, 0), Radius: 1.0, Mat: Metal{Albedo: NewVec3(0.7, 0.6, 0.5), Fuzz: 0.0}})

	return world
}

func main() {
	rand.Seed(time.Now().UnixNano())
	world := randomScene()

	cam := Camera{
		AspectRatio:     16.0 / 9.0,
		ImageWidth:      400,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		VFov:            20.0,
		LookFrom:        NewVec3(-2, 2, 1),
		LookAt:          NewVec3(0, 0, -1),
		VUp:             NewVec3(0, 1, 0),
		DefocusAngle:    0.6,
		FocusDist:       3.4,
	}
	cam.Render(&world)
}
