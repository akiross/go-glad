package main

// Showing a texture

import (
	"image"
	"image/color"
	"log"
	"runtime"

	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.5-core/gl"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting")

	win := glad.NewOGLWindow(800, 600, "Texture",
		glad.CoreProfile(true),
		glad.Resizable(false),
		glad.ContextVersion(4, 4),
		//glad.VSync(true),
	)
	defer glad.Terminate()
	// Enable VSync
	glad.SwapInterval(1)

	bgCol := []float32{0.3, 0.3, 0.3, 1.0}
	var (
		vertexShaderSource = `#version 440 core
in vec2 pos;
in vec2 uv;
out vec2 vUV;
void main() {
	gl_Position = vec4(pos, 0.0, 1.0);
	vUV = uv;
}`
		fragmentShaderSource = `#version 440 core
in vec2 vUV;
out vec4 color;
uniform sampler2D sampler;
void main() {
	color = vec4(0.1, 0.1, 0.1, 1.0) + texture(sampler, vUV);
}`
	)

	vertShader := glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader := glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program := glad.NewProgram()
	program.AttachShaders(vertShader, fragShader)
	program.Link()

	vertShader.Delete()
	fragShader.Delete()

	// Create a texture
	txrImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			txrImg.SetRGBA(x, y, color.RGBA{uint8(255.0 * float32(x%8) / 7.0), uint8(255 * float32(y%16) / 15.0), 0, 255})
		}
	}

	txr := glad.NewTexture(gl.TEXTURE_2D)
	txr.Storage(1, gl.RGBA8, []int{64, 64}) //txrImg.Bounds().Dx(), txrImg.Bounds().Dy())
	txr.Bind(0)
	txr.Image2D(txrImg)
	//txr.Clear(255, 0, 0, 255)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)

	// Data to be used when drawing
	// Format: X, Y, U, V
	vertPosAndUV := []float32{
		-1.0, -1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0,
	}

	// Create the buffer object to store vertex data
	vbo := glad.NewVertexBufferObject()
	vbo.BufferData32(vertPosAndUV, gl.STATIC_DRAW)

	var bindPos uint32
	vao := glad.NewVertexArrayObject()

	vao.VertexBuffer32(bindPos, vbo, 0, 4) // Bind the VBO to the bind index

	attrPos := program.GetAttributeLocation("pos")
	vao.AttribFormat32(attrPos, 2, 0)   // Specify attribute stride and offset
	vao.AttribBinding(bindPos, attrPos) // Link the attribute to the bind pos

	attrUV := program.GetAttributeLocation("uv")
	vao.AttribFormat32(attrUV, 2, 2)
	vao.AttribBinding(bindPos, attrUV)

	vao.EnableAttrib(attrPos)
	vao.EnableAttrib(attrUV)

	program.Use()
	vao.Bind()

	for !win.ShouldClose() {
		gl.ClearBufferfv(gl.COLOR, 0, &bgCol[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		win.SwapBuffers()
		glad.PollEvents()
	}
}
