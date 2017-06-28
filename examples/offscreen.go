package main

// Hello world in OpenGL: create a triangle on screen

import (
	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

var nextFrame = false

func advanceFrame(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, md glfw.ModifierKey) {
	if act == glfw.Press && key == glfw.KeySpace {
		nextFrame = true
	}
}

func checkError() {
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Println("GL ERROR:", err)
	}
}

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
	win.SetKeyCallback(advanceFrame)

	var (
		programCol    glad.Program
		programTxr    glad.Program
		vertShader    glad.Shader
		fragShaderCol glad.Shader
		fragShaderTxr glad.Shader
		vaoC          glad.VertexArrayObject
		vaoT          glad.VertexArrayObject
		vboC          glad.VertexBufferObject
		vboT          glad.VertexBufferObject
		fbo           glad.FramebufferObject
		txr           glad.Texture
		attrPosC      glad.VertexAttrib
		attrCol       glad.VertexAttrib
		attrPosT      glad.VertexAttrib
		attrUV        glad.VertexAttrib
	)

	vertShader = glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShaderCol = glad.NewShader(fragmentShaderColSource, gl.FRAGMENT_SHADER)
	fragShaderTxr = glad.NewShader(fragmentShaderTxrSource, gl.FRAGMENT_SHADER)

	programCol = glad.PrepareProgram(vertShader, fragShaderCol)
	programTxr = glad.PrepareProgram(vertShader, fragShaderTxr)

	vertShader.Delete()
	fragShaderCol.Delete()
	fragShaderTxr.Delete()

	checkError()

	var unbind func()

	// Triangles to draw on the texture
	// Format: X, Y, R, G, B
	tris := []float32{
		-1.0, -1.0, 0.0, 0.0, 0.0,
		0.0, -1.0, 1.0, 0.0, 0.0,
		-.75, 0.0, 0.5, 1.0, 1.0,

		0.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 0.0, 1.0, 1.0,
		0.75, 0.0, 0.5, 0.0, 0.0,

		-0.75, 0.0, 0.0, 0.0, 1.0,
		0.75, 0.0, 1.0, 1.0, 0.0,
		0.0, 1.0, 0.5, 0.0, 1.0,
	}

	// This quad will be used to render the texture
	// Format: X, Y, U, V
	quad := []float32{
		0.0, -0.9, 0.0, 0.0,
		-0.9, 0.9, 0.0, 1.0,
		0.9, -0.9, 1.0, 0.0,
		0.0, 0.9, 1.0, 1.0,
	}

	// Create a VAO and VBO for triangles (Color)
	vaoC = glad.NewVertexArrayObject()
	vboC = glad.NewVertexBufferObject()
	unbind = glad.BlockBind(vaoC, vboC)
	vboC.BufferData32(tris, gl.STATIC_DRAW) // Fill buffer with data
	// Define attributes to draw triangles
	attrPosC = programCol.GetAttributeLocation("pos")
	attrPosC.PointerFloat32(2, false, 5, 0)
	attrCol = programCol.GetAttributeLocation("col")
	attrCol.PointerFloat32(3, false, 5, 2)
	unbind()

	// Create a VAO and VBO for quad (Texture)
	vaoT = glad.NewVertexArrayObject()
	vboT = glad.NewVertexBufferObject()
	unbind = glad.BlockBind(vaoT, vboT)
	vboT.BufferData32(quad, gl.STATIC_DRAW)
	// Define attributes to draw the texture on quad
	attrPosT = programTxr.GetAttributeLocation("pos")
	attrPosT.PointerFloat32(2, false, 4, 0)
	attrUV = programTxr.GetAttributeLocation("uv")
	attrUV.PointerFloat32(2, false, 4, 2)
	vaoT.Unbind()

	// Create a framebuffer and a texture as render target
	fbo = glad.NewFramebuffer()
	txr = glad.NewTexture()
	unbind = glad.BlockBind(fbo, txr)
	txr.Empty2D(win.GetSize())
	txr.SetFilters(gl.NEAREST, gl.NEAREST)
	// Attach texture to target
	fbo.Texture(gl.COLOR_ATTACHMENT0, txr)
	unbind()

	// Draw onto texture the triangles
	unbind = glad.BlockBind(fbo, vaoC) //, vboC)
	programCol.Use()
	disable := glad.BlockEnable(attrPosC, attrCol)

	gl.ClearColor(0.8, 0.8, 0.8, 1.0)
	gl.Viewport(0, 0, 800, 600)
	// Clear the texture, filling it with white
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// Draw the triangles over the texture
	gl.DrawArrays(gl.TRIANGLES, 0, 9)

	disable()
	unbind()

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	vaoT.Bind()
	txr.Bind()
	programTxr.Use()
	attrPosT.Enable()
	attrUV.Enable()

	for !win.ShouldClose() {
		// Draw triangles with texture
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glad.PollEvents()
	}
}

var (
	vertexShaderSource = `#version 440 core
in vec2 pos;
in vec3 col;
in vec2 uv;
out vec2 vUV;
out vec3 vCol;
void main() { gl_Position = vec4(pos, 0.0, 1.0); vUV = uv; vCol = col; }`
	fragmentShaderTxrSource = `#version 440 core
in vec2 vUV;
in vec3 vCol;
out vec3 color;
uniform sampler2D sampler;
void main() { color = texture(sampler, vUV).rgb; }`
	fragmentShaderColSource = `#version 440 core
in vec2 vUV;
in vec3 vCol;
out vec3 color;
uniform sampler2D sampler;
void main() { color = vCol; }`
)
