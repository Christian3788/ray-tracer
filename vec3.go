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

func linearToGamma(c float64) float64 {
    if c > 0 {
        return math.Sqrt(c)
    }
    return 0
}

func colorToByte(c float64) int {
    intensity := Interval{0.0, 0.999}
    return int(256 * intensity.Clamp(linearToGamma(c)))
}

// WriteColor outputs a single pixel's colour in PPM byte format.
func (c Color) WriteColor(w io.Writer) {
    fmt.Fprintf(w, "%d %d %d\n",
        colorToByte(c.X),
        colorToByte(c.Y),
        colorToByte(c.Z),
    )
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

func RandomInUnitSphere() Vec3 {
    for {
        p := RandomVec3Range(-1, 1)
        if p.LengthSquared() >= 1 {
            continue
        }
        return p
    }
}

func RandomInUnitDisk() Vec3 {
    for {
        p := Vec3{rand.Float64()*2 - 1, rand.Float64()*2 - 1, 0}
        if p.LengthSquared() >= 1 {
            continue
        }
        return p
    }
}

func Reflect(v, n Vec3) Vec3 {
    return v.Sub(n.Scale(2 * Dot(v, n)))
}

func (v Vec3) NearZero() bool {
    const s = 1e-8
    return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}
