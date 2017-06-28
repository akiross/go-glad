package main

// Hello world in OpenGL: create a triangle on screen

import (
	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.4-core/gl"
	"image"
	"image/color"
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
	// Format: X, Y, U, V
	vertPosAndUV := []float32{
		-1.0, -1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0,
	}

	// Create a texture
	txrImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			txrImg.SetRGBA(x, y, color.RGBA{uint8(255.0 * float32(x%8) / 7.0), uint8(255 * float32(y%16) / 15.0), 0, 255})
		}
	}

	vao := glad.NewVertexArrayObject()
	vao.Bind()

	vbo := glad.NewVertexBufferObject()
	vbo.Bind()
	vbo.BufferData32(vertPosAndUV, gl.STATIC_DRAW)

	txr := glad.NewTexture()
	txr.Bind()
	txr.Image2D(txrImg)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)

	attrPos := program.GetAttributeLocation("pos")
	attrPos.PointerFloat32(2, false, 4, 0)
	attrPos.Enable()

	attrUV := program.GetAttributeLocation("uv")
	attrUV.PointerFloat32(2, false, 4, 2)
	attrUV.Enable()

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
in vec2 uv;
out vec2 vUV;
void main() { gl_Position = vec4(pos, 0.0, 1.0); vUV = uv; }`
	fragmentShaderSource = `#version 440 core
in vec2 vUV;
out vec3 color;
uniform sampler2D sampler;
void main() { color = texture(sampler, vUV).rgb; }`
)
