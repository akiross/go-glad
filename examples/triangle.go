package main

// Hello world in OpenGL: create a triangle on screen

import (
	"log"
	"runtime"

	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.5-core/gl"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting")

	win := glad.NewOGLWindow(800, 600, "Hello Triangle",
		glad.CoreProfile(true),
		glad.Resizable(false),
		glad.ContextVersion(4, 5),
		//glad.VSync(true),
	)
	defer glad.Terminate()
	// Enable VSync
	glad.SwapInterval(1)

	var (
		bgCol              = []float32{0.3, 0.3, 0.3, 1.0}
		vertexShaderSource = `#version 450 core
	in vec2 pos;
	in vec3 col;
	out vec3 vCol;
	void main() { gl_Position = vec4(pos, 0.0, 1.0); vCol = col; }`
		fragmentShaderSource = `#version 450 core
	in vec3 vCol;
	out vec3 color;
	void main() { color = vCol; }`
	)

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

	var bindPos uint32 = 0 // Bind position in the VAO

	vao := glad.NewVertexArrayObject()              // Create VAO
	vbo := glad.NewVertexBufferObject()             // Create VBO
	vbo.BufferData32(vertPosAndCol, gl.STATIC_DRAW) // Set VBO data
	vao.VertexBuffer32(bindPos, vbo, 0, 5)          // Bind VBO to position 0 in VAO, no offset, 5 elements stride

	attrPos := program.GetAttributeLocation("pos") // Find attribute position for "pos"
	vao.AttribFormat32(attrPos, 2, 0)              // Specify attribute size and relative offset
	vao.AttribBinding(bindPos, attrPos)            // Bind attributes to position 0
	vao.EnableAttrib(attrPos)

	attrCol := program.GetAttributeLocation("col") // Find position of "col"
	vao.AttribFormat32(attrCol, 3, 2)              // Specify 3 floats, after 2 floats
	vao.AttribBinding(bindPos, attrCol)            // Bind attribute to position 1
	vao.EnableAttrib(attrCol)

	program.Use()
	vao.Bind()

	for !win.ShouldClose() {
		gl.ClearBufferfv(gl.COLOR, 0, &bgCol[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		win.SwapBuffers()
		glad.PollEvents()
	}
}
