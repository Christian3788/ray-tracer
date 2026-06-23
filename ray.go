package main

func (r Ray) At(t float64) Point3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
