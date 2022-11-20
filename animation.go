package rlottie

/*
#include <stdlib.h>
#include <rlottie_capi.h>

size_t u32t_sz = sizeof(uint32_t);
void _lottie_animation_property_override(Lottie_Animation *animation, const Lottie_Animation_Property type, const char *keypath, double v1, double v2, double v3)
{
	// this func uses va_list
	lottie_animation_property_override(animation, type, keypath, v1, v2, v3);
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type Animation struct {
	animation *Lottie_Animation_S
}

// AnimationFromData constructs an animation object from JSON string data.
// data, the JSON string data.
// key, the string that will be used to cache the JSON string data.
// resourcePath, the path that will be used to load external resource needed by the JSON data.
//
// returns Animation object that can build the contents of the
// Lottie resource represented by JSON string data.
func AnimationFromData(data, key, rescourcePath string) (r Animation, err error) {
	var (
		cData         = C.CString(data)
		cKey          = C.CString(key)
		cResourcePath = C.CString(rescourcePath)
	)
	defer C.free(unsafe.Pointer(cData))
	defer C.free(unsafe.Pointer(cKey))
	defer C.free(unsafe.Pointer(cResourcePath))

	animation := C.lottie_animation_from_data(cData, cKey, cResourcePath)
	if animation == nil {
		err = ErrLottieLoadFailed
		return
	}

	r = Animation{animation}
	return
}

// AnimationFromFile constructs an animation object from file path.
// path, Lottie resource file path
//
// returns Animation object that can build the contents of the
// Lottie resource represented by file path.
func AnimationFromFile(path string) (r Animation, err error) {
	var (
		cPath = C.CString(path)
	)
	defer C.free(unsafe.Pointer(cPath))

	animation := C.lottie_animation_from_file(cPath)
	if animation == nil {
		err = ErrLottieLoadFailed
		return
	}

	r = Animation{animation}
	return
}

// Destroy Animation object resource.
func (a Animation) Destroy() {
	C.lottie_animation_destroy(a.animation)
}

// GetSize returns default viewport size of the Lottie reource.
func (a Animation) GetSize() (d Dimension) {
	C.lottie_animation_get_size(a.animation,
		(*C.size_t)(unsafe.Pointer(&d.Width)),
		(*C.size_t)(unsafe.Pointer(&d.Height)))
	return
}

// GetDuration returns total animation duration of Lottie resource in second.
// it uses totalFrame() and frameRate() to calculate the duration.
// duration = totalFrame() / frameRate()
func (a Animation) GetDuration() float64 {
	return float64(C.lottie_animation_get_duration(a.animation))
}

// GetTotalFrame returns total number of frames present in the Lottie resource.
func (a Animation) GetTotalFrame() uint {
	return uint(C.lottie_animation_get_totalframe(a.animation))
}

// GetFrameRate returns default framerate of the Lottie resource.
func (a Animation) GetFrameRate() float64 {
	return float64(C.lottie_animation_get_framerate(a.animation))
}

// RenderTree get the render tree which contains the snapshot of the animation object
// at `frameNum`, the content of the animation in that frame number.
func (a Animation) RenderTree(frameNum, width, height uint) (r *LOTLayerNode) {
	renderTree := C.lottie_animation_render_tree(a.animation,
		C.size_t(frameNum),
		C.size_t(width),
		C.size_t(height))
	return &LOTLayerNode{renderTree}
}

// GetFrameAtPos maps position to frame number and returns it.
// pos, position in the range [ 0.0 .. 1.0 ]
func (a Animation) GetFrameAtPos(pos float64) uint {
	return uint(C.lottie_animation_get_frame_at_pos(a.animation, C.float(pos)))
}

// Render the content of the frame `frameNum` to `buf`
// frameNum, the frame number needs to be rendered.
// width, width of the surface
// height, height of the surface
// bytesPerLine, stride of the surface in bytes.
func (a Animation) Render(buf []uint32, frameNum, width, height, bytesPerLine uint) (rcast []uint8) {
	dim := width * height
	C.lottie_animation_render(a.animation,
		C.size_t(frameNum),
		(*C.uint32_t)(unsafe.Pointer(&buf[0])),
		C.size_t(width),
		C.size_t(height),
		C.size_t(bytesPerLine))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&rcast))
	sh.Cap = int(dim * (32 / 8)) // 32-bit casted to 8-bit
	sh.Len = sh.Cap
	sh.Data = uintptr(unsafe.Pointer(&buf[0]))
	return
}

// RenderAsync render the content of the frame `frameNum` to `buf` asynchronously.
// frameNum, frame number needs to be rendered.
// width, width of the surface
// height, height of the surface
// bytesPerLine, stride of the surface in bytes.
func (a Animation) RenderAsync(buf []uint32, frameNum, width, height, bytesPerLine uint) (rcast []uint8) {
	dim := width * height
	C.lottie_animation_render_async(a.animation,
		C.size_t(frameNum),
		(*C.uint32_t)(unsafe.Pointer(&buf[0])),
		C.size_t(width),
		C.size_t(height),
		C.size_t(bytesPerLine))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&rcast))
	sh.Cap = int(dim * (32 / 8))
	sh.Len = sh.Cap
	sh.Data = uintptr(unsafe.Pointer(&buf[0]))
	return
}

// RenderFlush request to finish the current async renderer job for this animation object.
// If render is finished the this call returns immediately.
// If not, it waits till render job finish and then return.
//
// warning: User must call `lottie_animation_render_async()` and `lottie_animation_render_flush()`
// in pair to get the benefit of async rendering.
//
// returns pixel buffer it finished rendering.
func (a Animation) RenderFlush() *C.uint {
	val := C.lottie_animation_render_flush(a.animation)
	return val
}

// PropertyOverride request to change properties of this animation object.
// type, property type. (Lottie_Animation_Property)
// keypath, specific content of target.
// props, ... property values.
func (a Animation) PropertyOverride(ptype int, keypath string, props ...float64) {
	var (
		cKeypath = C.CString(keypath)
	)
	defer C.free(unsafe.Pointer(cKeypath))

	var v1, v2, v3 float64
	switch {
	case len(props) >= 3:
		v1, v2, v3 = props[0], props[1], props[2]
	case len(props) >= 2:
		v1, v2 = props[0], props[1]
	case len(props) > 0:
		v1 = props[0]
	}

	C._lottie_animation_property_override(a.animation,
		Lottie_Animation_Property(ptype),
		cKeypath,
		C.double(v1), C.double(v2), C.double(v3))
}

// GetMarkerList returns list of markers in the Lottie resource
// LOTMarkerList has a `LOTMarker` list and size of list
// LOTMarker has the marker's name, start frame, and end frame.
func (a Animation) GetMarkerList() (r []LOTMarker) {
	markerList := C.lottie_animation_get_markerlist(a.animation)
	if markerList == nil {
		return nil
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	sh.Cap = int(markerList.size)
	sh.Len = int(markerList.size)
	sh.Data = uintptr(unsafe.Pointer(markerList.ptr))
	return
}
