package glad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type VertexAttrib uint32

func (va VertexAttrib) Enable() {
	gl.EnableVertexAttribArray(uint32(va))
}

func (va VertexAttrib) Disable() {
	gl.DisableVertexAttribArray(uint32(va))
}

/*
// size: number of components per vertex (e.g. 3D vertices -> 3)
// dataType: gl.FLOAT, etc
// normalized: define if data have to be normalized
// stride: bytes between two vertices, 0 means they are tightly packed
// offset: bytes of offset to the first element in the array
func (va VertexAttrib) Pointer(size int32, dataType uint32, normalize bool, stride, offset int32) {
	gl.VertexAttribPointer(uint32(va), size, dataType, normalize, stride, gl.PtrOffset(int(offset)))
}

// Same as Pointer, but stride and offsets are count of elements, dataType is set to gl.FLOAT
func (va VertexAttrib) PointerFloat32(size int32, normalize bool, stride, offset int32) {
	gl.VertexAttribPointer(uint32(va), size, gl.FLOAT, normalize, stride*4, gl.PtrOffset(int(offset)*4))
}

// Same as PointerFloat32, but for Float64
func (va VertexAttrib) PointerFloat64(size int32, normalize bool, stride, offset int32) {
	gl.VertexAttribPointer(uint32(va), size, gl.DOUBLE, normalize, stride*8, gl.PtrOffset(int(offset)*8))
}
*/
