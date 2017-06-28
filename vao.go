package goglad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type VertexArrayObject uint32

func NewVertexArrayObject() VertexArrayObject {
	var vao uint32
	gl.CreateVertexArrays(1, &vao)
	return VertexArrayObject(vao)
}

func (vao VertexArrayObject) Bind() {
	gl.BindVertexArray(uint32(vao))
}

func (vao VertexArrayObject) Unbind() {
	gl.BindVertexArray(0)
}

func (vao VertexArrayObject) Delete() {
	var v = uint32(vao)
	gl.DeleteVertexArrays(1, &v)
}

func (vao VertexArrayObject) AttribBinding(bindIndex uint32, attr VertexAttrib) {
	gl.VertexArrayAttribBinding(uint32(vao), uint32(attr), bindIndex)
}

func (vao VertexArrayObject) VertexBuffer(bindIndex uint32, buffer VertexBufferObject, offset, stride int32) {
	gl.VertexArrayVertexBuffer(uint32(vao), bindIndex, uint32(buffer), int(offset), stride)
}

func (vao VertexArrayObject) VertexBuffer32(bindIndex uint32, buffer VertexBufferObject, offset, stride int32) {
	gl.VertexArrayVertexBuffer(uint32(vao), bindIndex, uint32(buffer), int(offset)*4, stride*4)
}

func (vao VertexArrayObject) AttribFormat(attr VertexAttrib, size int32, dataType uint32, normalize bool, offset uint32) {
	gl.VertexArrayAttribFormat(uint32(vao), uint32(attr), size, dataType, normalize, offset)
}

func (vao VertexArrayObject) AttribFormat32(attr VertexAttrib, size int32, offset uint32) {
	gl.VertexArrayAttribFormat(uint32(vao), uint32(attr), size, gl.FLOAT, false, offset*4)
}

func (vao VertexArrayObject) EnableAttrib(attr VertexAttrib) {
	gl.EnableVertexArrayAttrib(uint32(vao), uint32(attr))
}
