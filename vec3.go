package main

import (
    "fmt"
    "io"
    "math"
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

// WriteColor outputs a single pixel's colour in PPM byte format.
func (c Color) WriteColor(w io.Writer) {
    ir := int(255.999 * c.X)
    ig := int(255.999 * c.Y)
    ib := int(255.999 * c.Z)
    fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
}

func (v Vec3) Unit() Vec3 { return v.Div(v.Length()) }