package rlottie

// #cgo LDFLAGS: -lrlottie
// #include <rlottie_capi.h>
import "C"

type (
	Lottie_Animation_S        = C.struct_Lottie_Animation_S
	Lottie_Animation_Property = C.Lottie_Animation_Property
)

type LOTLayerNode struct {
	n *C.struct_LOTLayerNode
}

type LOTMask struct {
	n *C.struct_LOTMask
}

type LOTNode struct {
	n *C.struct_LOTNode
}

type LOTMarker struct {
	n *C.struct_LOTMarker
}
