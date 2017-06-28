package main

// Hello world in OpenGL: create a triangle on screen

import (
	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
	_ "image/png"
	"log"
	"os"
	"runtime"
)

var nextFrame = false

func advanceFrame(w *glfw.Window, key glfw.Key, sc int, act glfw.Action, md glfw.ModifierKey) {
	if act == glfw.Press && key == glfw.KeySpace {
		nextFrame = true
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
		program    glad.Program
		vertShader glad.Shader
		fragShader glad.Shader
		vao        glad.VertexArrayObject
		vbo        glad.VertexBufferObject
		txr        glad.Texture
		attrPos    glad.VertexAttrib
		attrUV     glad.VertexAttrib
	)

	vertShader = glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader = glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program = glad.PrepareProgram(vertShader, fragShader)

	vertShader.Delete()
	fragShader.Delete()

	// Load the image
	file, err := os.Open("image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// This quad will be used to render the texture
	// Format: X, Y, U, V
	quad := []float32{
		-0.9, -0.9, 0.0, 1.0,
		-0.9, 0.9, 0.0, 0.0,
		0.9, -0.9, 1.0, 1.0,
		0.9, 0.9, 1.0, 0.0,
	}

	// Create a VAO and VBO for quad (Texture)
	vao = glad.NewVertexArrayObject()
	vao.Bind()
	vbo = glad.NewVertexBufferObject()
	vbo.Bind()
	vbo.BufferData32(quad, gl.STATIC_DRAW)
	// Define attributes to draw the texture on quad
	attrPos = program.GetAttributeLocation("pos")
	attrPos.PointerFloat32(2, false, 4, 0)
	attrUV = program.GetAttributeLocation("uv")
	attrUV.PointerFloat32(2, false, 4, 2)

	// Create a framebuffer and a texture as render target
	txr = glad.NewTexture()
	txr.Bind()
	txr.Image2D(img)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)
	// Attach texture to target
	//fbo.Texture(gl.COLOR_ATTACHMENT0, txr)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	program.Use()
	attrPos.Enable()
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
in vec2 uv;
out vec2 vUV;
void main() { gl_Position = vec4(pos, 0.0, 1.0); vUV = uv; }`
	fragmentShaderSource = `#version 440 core
in vec2 vUV;
out vec3 color;
uniform sampler2D sampler;
void main() { color = texture(sampler, vUV).rgb; }`
)
