package main

// Example of loading an image and using it as a texture

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"runtime"

	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.5-core/gl"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting")

	win := glad.NewOGLWindow(800, 600, "Texture image",
		glad.CoreProfile(true),
		glad.Resizable(false),
		glad.ContextVersion(4, 4),
		//glad.VSync(true),
	)
	defer glad.Terminate()
	// Enable VSync
	glad.SwapInterval(1)

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
	bindPos := uint32(0)
	vao = glad.NewVertexArrayObject()
	vbo = glad.NewVertexBufferObject()
	vbo.BufferData32(quad, gl.STATIC_DRAW)
	vao.VertexBuffer32(bindPos, vbo, 0, 4)

	// Define attributes to draw the texture on quad
	attrPos = program.GetAttributeLocation("pos")
	vao.AttribFormat32(attrPos, 2, 0)
	vao.AttribBinding(bindPos, attrPos)

	attrUV = program.GetAttributeLocation("uv")
	vao.AttribFormat32(attrUV, 2, 2)
	vao.AttribBinding(bindPos, attrUV)

	// Create a texture as render target
	txr = glad.NewTexture(gl.TEXTURE_2D)
	txr.Storage(1, gl.RGBA8, []int{img.Bounds().Dx(), img.Bounds().Dy()})
	txr.Image2D(img)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	txr.Bind(0)
	vao.Bind()
	program.Use()
	vao.EnableAttrib(attrPos)
	vao.EnableAttrib(attrUV)

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
