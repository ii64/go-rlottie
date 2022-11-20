package rlottie

/*
#include <stdio.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type Dimension struct {
	Width  uint
	Height uint
}

type Point struct {
	X, Y float64
}

type Color struct {
	R, G, B, A uint8
}

type Stroke struct {
	Enable     uint8
	Width      float64
	Cap        int // LOTCapStyle
	Join       int // LOTJoinStyle
	MiterLimit float64
	DashArray  []float64
}

type Gradient struct {
	Type                      int   // LOTGradientType
	Stop                      []int // LOTGradientStop
	Start, End, Center, Focal Point
	Cradius                   float64
	Fradius                   float64
}

type Matrix struct {
	M11, M12, M13 float64
	M21, M22, M23 float64
	M31, M32, M33 float64
}

type ImageInfo struct {
	Data          *uint8
	Width, Height uint
	Alpha         uint8
	Matrix        Matrix
}

type Path struct {
	Pt  []float64
	Elm []string
}

// LOTLayerNode
func (lln *LOTLayerNode) Mask() (r []LOTMask) {
	ins := lln.n.mMaskList
	if ins.ptr == nil {
		return nil
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	sh.Cap = int(ins.size)
	sh.Len = int(ins.size)
	sh.Data = uintptr(unsafe.Pointer(ins.ptr))
	return
}
func (lln *LOTLayerNode) ClipPath() (r Path) {
	ins := lln.n.mClipPath
	if ins.ptPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Pt))
		sh.Cap = int(ins.ptCount)
		sh.Len = int(ins.ptCount)
		sh.Data = uintptr(unsafe.Pointer(ins.ptPtr))
	}
	if ins.elmPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Elm))
		sh.Cap = int(ins.elmCount)
		sh.Len = int(ins.elmCount)
		sh.Data = uintptr(unsafe.Pointer(ins.elmPtr))
	}
	return
}
func (lln *LOTLayerNode) LayerList() (r []LOTLayerNode) {
	// implement setter later..
	ins := lln.n.mLayerList
	if ins.ptr == nil {
		return nil
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	sh.Cap = int(ins.size)
	sh.Len = int(ins.size)
	sh.Data = uintptr(unsafe.Pointer(ins.ptr))
	return
}
func (lln *LOTLayerNode) NodeList() (r []LOTNode) {
	ins := lln.n.mNodeList
	if ins.ptr == nil {
		return
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	sh.Cap = int(ins.size)
	sh.Len = int(ins.size)
	sh.Data = uintptr(unsafe.Pointer(ins.ptr))
	return
}
func (lln *LOTLayerNode) Matte() int {
	// LOTMatteType
	return int(lln.n.mMatte)
}
func (lln *LOTLayerNode) Visible() int {
	return int(lln.n.mVisible)
}
func (lln *LOTLayerNode) Alpha() uint8 {
	return uint8(lln.n.mAlpha)
}
func (lln *LOTLayerNode) KeyPath() string {
	return C.GoString(lln.n.keypath)
}

// LOTMask
func (lm *LOTMask) Mode() int {
	// LOTMaskType
	return int(lm.n.mMode)
}
func (lm *LOTMask) Alpha() uint8 {
	return uint8(lm.n.mAlpha)
}
func (lm *LOTMask) Path() (r Path) {
	ins := lm.n.mPath
	if ins.ptPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Pt))
		sh.Cap = int(ins.ptCount)
		sh.Len = int(ins.ptCount)
		sh.Data = uintptr(unsafe.Pointer(ins.ptPtr))
	}
	if ins.elmPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Elm))
		sh.Cap = int(ins.elmCount)
		sh.Len = int(ins.elmCount)
		sh.Data = uintptr(unsafe.Pointer(ins.elmPtr))
	}
	return
}

// LOTNode
func (ln *LOTNode) Path() (r Path) {
	ins := ln.n.mPath
	if ins.ptPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Pt))
		sh.Cap = int(ins.ptCount)
		sh.Len = int(ins.ptCount)
		sh.Data = uintptr(unsafe.Pointer(ins.ptPtr))
	}
	if ins.elmPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Elm))
		sh.Cap = int(ins.elmCount)
		sh.Len = int(ins.elmCount)
		sh.Data = uintptr(unsafe.Pointer(ins.elmPtr))
	}
	return
}
func (ln *LOTNode) Color() *Color {
	c := ln.n.mColor
	return &Color{
		R: uint8(c.r),
		G: uint8(c.g),
		B: uint8(c.b),
		A: uint8(c.a),
	}
}
func (ln *LOTNode) Stroke() (r *Stroke) {
	ins := ln.n.mStroke
	r = new(Stroke)
	r.Enable = uint8(ins.enable)
	r.Width = float64(ins.width)
	r.Cap = int(ins.cap)
	r.MiterLimit = float64(ins.miterLimit)
	if ins.dashArray != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.DashArray))
		sh.Len = int(ins.dashArraySize)
		sh.Cap = int(ins.dashArraySize)
		sh.Data = uintptr(unsafe.Pointer(ins.dashArray))
	}
	return
}
func (ln *LOTNode) Gradient() (r *Gradient) {
	ins := ln.n.mGradient
	r = new(Gradient)
	r.Type = int(ins._type)
	if ins.stopPtr != nil {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&r.Stop))
		sh.Cap = int(ins.stopCount)
		sh.Len = int(ins.stopCount)
		sh.Data = uintptr(unsafe.Pointer(ins.stopPtr))
	}
	start, end, center, focal := ins.start, ins.end, ins.center, ins.focal
	r.Start = Point{float64(start.x), float64(start.y)}
	r.End = Point{float64(end.x), float64(end.y)}
	r.Center = Point{float64(center.x), float64(center.y)}
	r.Focal = Point{float64(focal.x), float64(focal.y)}
	r.Cradius = float64(ins.cradius)
	r.Fradius = float64(ins.fradius)
	return
}
func (ln *LOTNode) ImageInfo() (r *ImageInfo) {
	ins := ln.n.mImageInfo
	r = new(ImageInfo)
	r.Data = (*uint8)(ins.data)
	r.Width = uint(ins.width)
	r.Height = uint(ins.height)
	r.Alpha = uint8(ins.mAlpha)
	matx := ins.mMatrix
	r.Matrix = Matrix{
		float64(matx.m11), float64(matx.m12), float64(matx.m13),
		float64(matx.m21), float64(matx.m22), float64(matx.m23),
		float64(matx.m31), float64(matx.m32), float64(matx.m33),
	}
	return
}
func (ln *LOTNode) Flag() int {
	return int(ln.n.mFlag)
}
func (ln *LOTNode) BrushType() int { // LOTBrushType
	return int(ln.n.mBrushType)
}
func (ln *LOTNode) FillRule() int { // LOTFillRule
	return int(ln.n.mFillRule)
}
func (ln *LOTNode) KeyPath() string {
	return C.GoString(ln.n.keypath)
}

// LOTMarker
func (lm *LOTMarker) Name() string {
	return C.GoString(lm.n.name)
}
func (lm *LOTMarker) StartFrame() uint {
	return uint(lm.n.startframe)
}
func (lm *LOTMarker) EndFrame() uint {
	return uint(lm.n.endframe)
}
