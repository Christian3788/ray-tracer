package main

import (
	"math"
	"math/rand"
)

type Lambertian struct {
	Albedo Color
}

func (l Lambertian) Scatter(rIn Ray, rec *HitRecord) (Color, Ray, bool) {
	scatterDirection := rec.Normal.Add(RandomUnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}
	scattered := Ray{Origin: rec.P, Direction: scatterDirection}
	return l.Albedo, scattered, true
}

type Metal struct {
	Albedo Color
	Fuzz   float64
}

func (m Metal) Scatter(rIn Ray, rec *HitRecord) (Color, Ray, bool) {
	reflected := Reflect(rIn.Direction.Unit(), rec.Normal)
	scattered := Ray{Origin: rec.P, Direction: reflected.Add(RandomVec3Range(-1, 1).Scale(m.Fuzz))}
	return m.Albedo, scattered, Dot(scattered.Direction, rec.Normal) > 0
}

type Dielectric struct {
	RefractionIndex float64
}

func Refract(uv, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := math.Min(Dot(uv.Neg(), n), 1.0)
	rOutPerp := uv.Add(n.Scale(cosTheta)).Scale(etaiOverEtat)
	rOutParallel := n.Scale(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (d Dielectric) Scatter(rIn Ray, rec *HitRecord) (Color, Ray, bool) {
	attenuation := NewVec3(1, 1, 1)
	ri := d.RefractionIndex
	if rec.FrontFace {
		ri = 1.0 / d.RefractionIndex
	}
	unitDir := rIn.Direction.Unit()
	cosTheta := math.Min(Dot(unitDir.Neg(), rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1.0
	var direction Vec3
	if cannotRefract || Schlick(cosTheta, ri) > rand.Float64() {
		direction = Reflect(unitDir, rec.Normal)
	} else {
		direction = Refract(unitDir, rec.Normal, ri)
	}
	return attenuation, Ray{Origin: rec.P, Direction: direction}, true
}

func Schlick(cosine, refIdx float64) float64 {
	// Schlick's approximation for reflectance
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
