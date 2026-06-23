package main

import "math"

type Interval struct {
	Min, Max float64
}

var EmptyInterval = Interval{math.Inf(1), math.Inf(-1)}
var UniverseInterval = Interval{math.Inf(-1), math.Inf(1)}

func (i Interval) Contains(x float64) bool {
	return x >= i.Min && x <= i.Max
}

func (i Interval) Surrounds(x float64) bool {
	return x >= i.Min && x <= i.Max
}

func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

type Material interface {
	Scatter(rIn Ray, rec *HitRecord) (attenuation Color, scattered Ray, ok bool)
}

type HitRecord struct {
	P         Point3
	Normal    Vec3
	T         float64
	FrontFace bool
	Mat       Material
}

func (rec *HitRecord) SetFaceNormal(r Ray, outwardNormal Vec3) {
	rec.FrontFace = Dot(r.Direction, outwardNormal) < 0
	if rec.FrontFace {
		rec.Normal = outwardNormal
	} else {
		rec.Normal = outwardNormal.Neg()
	}
}

type Hittable interface {
	Hit(r Ray, rayT Interval, rec *HitRecord) bool
}

type Sphere struct {
	Center Point3
	Radius float64
	Mat    Material
}

func (s Sphere) Hit(r Ray, rayT Interval, rec *HitRecord) bool {
	oc := s.Center.Sub(r.Origin)
	a := r.Direction.LengthSquared()
	h := Dot(r.Direction, oc)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}
	sqrtd := math.Sqrt(discriminant)

	root := (-h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (-h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = r.At(rec.T)
	outwardNormal := rec.P.Sub(s.Center).Div(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat
	return true
}

type HittableList struct {
	Objects []Hittable
}

func (hl *HittableList) Add(obj Hittable) {
	hl.Objects = append(hl.Objects, obj)
}

func (hl *HittableList) Hit(r Ray, rayT Interval, rec *HitRecord) bool {
	tempRec := HitRecord{}
	hitAnything := false
	closestSoFar := rayT.Max

	for _, object := range hl.Objects {
		if object.Hit(r, Interval{rayT.Min, closestSoFar}, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}

	return hitAnything
}
