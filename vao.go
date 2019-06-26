// Package glad allows to use OpenGL in Go
// It uses DSA-style OpenGL, to make it easier to learn and understand
// https://www.khronos.org/opengl/wiki/Direct_State_Access
// Generally, you want to:
// 1. create buffer objects and fill it with data
// 2. create a VAO NewVertexArrayObject
// 3. bind the VBO to some VAOs bind index with VertexBuffer
// 4. specify the format of some attributes with AttribFormat
// 5. bind the buffer to the attribute with AttribBinding
package glad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// VertexArrayObject stores data relative to some vertices stored in VBOs
// VertexArrayObject contain information on data location and layout in memory
// VertexArrayObject represents the vertex attribute state
// In the core profile, at least one VAO is mandatory
// VAOs have "attribute lists", where you can store data about your mesh
// each attribute you might store vertex positions, colors, normals, texture coords, etc
// The data themselves are stored as VBOs (plain data)
type VertexArrayObject uint32

// VertexAttrib represents a vertex attribute
// This can be obtained for example by using Program.GetAttributeLocation
type VertexAttrib uint32

// NewVertexArrayObject allocates the name for a VAO
func NewVertexArrayObject() VertexArrayObject {
	var vao uint32
	gl.CreateVertexArrays(1, &vao)
	return VertexArrayObject(vao)
}

// Bind the VAO maing it active
// Use this to select the vertex data to be used in the draw calls
func (vao VertexArrayObject) Bind() {
	gl.BindVertexArray(uint32(vao))
}

// Unbind any VAO currently bound
func (vao VertexArrayObject) Unbind() {
	gl.BindVertexArray(0)
}

// Delete the VAO freeing the name
func (vao VertexArrayObject) Delete() {
	var v = uint32(vao)
	gl.DeleteVertexArrays(1, &v)
}

// EnableAttrib the vertex attribute in the VAO, storing the state in the VAO
// This tells OpenGL to read the attribute data from the buffer
func (vao VertexArrayObject) EnableAttrib(attr VertexAttrib) {
	gl.EnableVertexArrayAttrib(uint32(vao), uint32(attr))
}

// VertexBuffer binds the buffer object to bindIndex
// bindIndex can be any value less than gl.MAX_VERTEX_ATTRIB_BINDINGS, but usually
// it is given the same index of the attribute
// offset and stride are in bytes
func (vao VertexArrayObject) VertexBuffer(bindIndex uint32, buffer VertexBufferObject, offset, stride int32) {
	gl.VertexArrayVertexBuffer(uint32(vao), bindIndex, uint32(buffer), int(offset), stride)
}

// VertexBuffer32 is like VertexBuffer but assuming we are working with 32 bit data
func (vao VertexArrayObject) VertexBuffer32(bindIndex uint32, buffer VertexBufferObject, offset, stride int32) {
	gl.VertexArrayVertexBuffer(uint32(vao), bindIndex, uint32(buffer), int(offset)*4, stride*4)
}

// AttribFormat specifies the format of the data associated to the attribute
// The state of the attribute is stored in the VAO
// size: number of components per vertex, 1, 2, 3 or 4 (e.g. 3D vertices -> 3)
// dataType: gl.FLOAT, etc
// normalized: define if data have to be normalized
// stride: bytes between two vertices, 0 means they are tightly packed
// offset: bytes of offset to the first element in the array
func (vao VertexArrayObject) AttribFormat(attr VertexAttrib, size int32, dataType uint32, normalize bool, relativeOffset uint32) {
	gl.VertexArrayAttribFormat(uint32(vao), uint32(attr), size, dataType, normalize, relativeOffset)
}

func (vao VertexArrayObject) AttribFormat32(attr VertexAttrib, size int32, offset uint32) {
	gl.VertexArrayAttribFormat(uint32(vao), uint32(attr), size, gl.FLOAT, false, offset*4)
}

// AttribBinding associates the attribute to the bind index
// By using the same bindIndex in VertexBuffer, AttribFormat and AttribBinding,
// the user can create a corrispondence between the attribute and the buffer.
func (vao VertexArrayObject) AttribBinding(bindIndex uint32, attr VertexAttrib) {
	gl.VertexArrayAttribBinding(uint32(vao), uint32(attr), bindIndex)
}

// TODO http://docs.gl/gl4/glVertexArrayElementBuffer
