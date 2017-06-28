package main

// Hello world in OpenGL: create a triangle on screen

import (
	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.4-core/gl"
	"log"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting")

	win := glad.NewOGLWindow(800, 600, "Gex",
		glad.CoreProfile(true),
		glad.Resizable(false),
		glad.ContextVersion(4, 4),
		//glad.VSync(true),
	)
	defer glad.Terminate()
	// Enable VSync
	glad.SwapInterval(1)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	vertShader := glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader := glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program := glad.NewProgram()
	program.AttachShaders(vertShader, fragShader)
	program.Link()

	vertShader.Delete()
	fragShader.Delete()

	// Data to be used when drawing
	// Format: X, Y, R, G, B,
	vertPosAndCol := []float32{
		-1.0, -1.0, 1.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 1.0, 0.0,
		1.0, -1.0, 0.0, 0.0, 1.0,
	}

	vao := glad.NewVertexArrayObject()
	vao.Bind()

	vbo := glad.NewVertexBufferObject()
	vbo.Bind()
	vbo.BufferData32(vertPosAndCol, gl.STATIC_DRAW)

	attrPos := program.GetAttributeLocation("pos")
	attrPos.PointerFloat32(2, false, 5, 0)
	attrPos.Enable()

	attrCol := program.GetAttributeLocation("col")
	attrCol.PointerFloat32(3, false, 5, 2)
	attrCol.Enable()

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		program.Use()
		vao.Bind()
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		win.SwapBuffers()
		glad.PollEvents()
	}
}

var (
	vertexShaderSource = `#version 440 core
in vec2 pos;
in vec3 col;
out vec3 vCol;
void main() { gl_Position = vec4(pos, 0.0, 1.0); vCol = col; }`
	fragmentShaderSource = `#version 440 core
in vec3 vCol;
out vec3 color;
void main() { color = vCol; }`
)
