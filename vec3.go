package main

import (
    "fmt"
    "io"
    "math"
    "math/rand"
)

// Vec3 is a 3D vector. We use it for points, directions, AND colours.
type Vec3 struct {
    X, Y, Z float64
}

// Aliases — no real type safety, but they help readability.
type Point3 = Vec3
type Color  = Vec3

// Constructor
func NewVec3(x, y, z float64) Vec3 { return Vec3{x, y, z} }

// --- arithmetic ---
func (v Vec3) Add(u Vec3) Vec3 { return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z} }
func (v Vec3) Sub(u Vec3) Vec3 { return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z} }
func (v Vec3) Mul(u Vec3) Vec3 { return Vec3{v.X * u.X, v.Y * u.Y, v.Z * u.Z} }
func (v Vec3) Scale(t float64) Vec3 { return Vec3{v.X * t, v.Y * t, v.Z * t} }
func (v Vec3) Div(t float64) Vec3   { return v.Scale(1 / t) }
func (v Vec3) Neg() Vec3            { return Vec3{-v.X, -v.Y, -v.Z} }

// --- length / dot / cross ---
func (v Vec3) LengthSquared() float64 { return v.X*v.X + v.Y*v.Y + v.Z*v.Z }
func (v Vec3) Length() float64        { return math.Sqrt(v.LengthSquared()) }

func Dot(a, b Vec3) float64 {
    return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross(a, b Vec3) Vec3 {
    return Vec3{
        a.Y*b.Z - a.Z*b.Y,
        a.Z*b.X - a.X*b.Z,
        a.X*b.Y - a.Y*b.X,
    }
}

// colorToByte converts a colour component from [0,1] to a byte [0,255].
func colorToByte(c float64) int {
    return int(255.999 * c)
}

// WriteColor outputs a single pixel's colour in PPM byte format.
func (c Color) WriteColor(w io.Writer) {
    intensity := Interval{0.000, 0.999}
    ir := int(256 * intensity.Clamp(c.X))
    ig := int(256 * intensity.Clamp(c.Y))
    ib := int(256 * intensity.Clamp(c.Z))
    fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
}

func (v Vec3) Unit() Vec3 { return v.Div(v.Length()) }

func RandomVec3() Vec3 {
    return Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func RandomVec3Range(min, max float64) Vec3 {
    return Vec3{
        min + (max-min)*rand.Float64(),
        min + (max-min)*rand.Float64(),
        min + (max-min)*rand.Float64(),
    }
}

func RandomUnitVector() Vec3 {
    for {
        p := RandomVec3Range(-1, 1)
        lensq := p.LengthSquared()
        if lensq < 1e-8 || lensq > 1 {
            continue
        }
        return p.Unit()
    }
}
