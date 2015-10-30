package help

import (
	"math"
)

// 三维向量
type Vec3 struct {
	X float32
	Y float32
	Z float32
}

// 是同一个点
func (t *Vec3) Eq(v *Vec3) bool {
	return t.X == v.X && t.Y == v.Y && t.Z == v.Z
}

// 交换点
func (t *Vec3) Swap(v *Vec3) {
	t.X, v.X = v.X, t.X
	t.Y, v.Y = v.Y, t.Y
	t.Z, v.Z = v.Z, t.Z
}

// 复制点
func (t *Vec3) Copy(v *Vec3) {
	t.X, t.Y, t.Z = v.X, v.Y, v.Z
}

// 点加
func (t *Vec3) Add(v *Vec3) {
	t.X, t.Y, t.Z = t.X+v.X, t.Y+v.Y, t.Z+v.Z
}

// 点减
func (t *Vec3) Sub(v *Vec3) {
	t.X, t.Y, t.Z = t.X-v.X, t.Y-v.Y, t.Z-v.Z
}

// 点乘
func (t *Vec3) Mult(v *Vec3) {
	t.X, t.Y, t.Z = t.X*v.X, t.Y*v.Y, t.Z*v.Z
}

// 缩放
func (t *Vec3) Scale(s float32) {
	t.X, t.Y, t.Z = t.X*s, t.Y*s, t.Z*s
}

// 逆缩放
func (t *Vec3) Div(s float32) {
	if s != 0 {
		inv := 1 / s
		t.X, t.Y, t.Z = t.X*inv, t.Y*inv, t.Z*inv
	}
}

// 点积
func (t *Vec3) Dot(v *Vec3) float64 {
	return float64(t.X)*float64(v.X) + float64(t.Y)*float64(v.Y) + float64(t.Z)*float64(v.Z)
}

// 这个点到原点的距离
func (t *Vec3) Len() float64 {
	return math.Sqrt(t.Dot(t))
}

// 这个点到原点距离的平方
func (t *Vec3) LenSqr() float64 {
	return t.Dot(t)
}

// 两点距离
func (t *Vec3) Dist(v *Vec3) float64 {
	return math.Sqrt(t.DistSqr(v))
}

// 两点距离的平方
func (t *Vec3) DistSqr(v *Vec3) float64 {
	x, y, z := float64(v.X-t.X), float64(v.Y-t.Y), float64(v.Z-t.Z)
	return x*x + y*y + z*z
}

// 求角度
func (t *Vec3) Ang(v *Vec3) float64 {
	magnitude := math.Sqrt(t.Dot(t) * v.Dot(v))
	if magnitude != 0 {
		return math.Acos(t.Dot(v) / magnitude)
	}
	// log.Printf("Dev error. vector.Vec3:Ang division by zero")
	return 0
}

// 叉积
func (t *Vec3) Cross(a, b *Vec3) {
	t.X, t.Y, t.Z = a.Y*b.Z-a.Z*b.Y, a.Z*b.X-a.X*b.Z, a.X*b.Y-a.Y*b.X
}
