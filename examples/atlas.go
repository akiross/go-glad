package main

// Hello world in OpenGL: create a triangle on screen

import (
	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"log"
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
		attrCol    glad.VertexAttrib
	)

	vertShader = glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader = glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program = glad.PrepareProgram(vertShader, fragShader)

	vertShader.Delete()
	fragShader.Delete()

	// Load the font
	fontBytes, err := ioutil.ReadFile("font.ttf")
	if err != nil {
		panic(err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	text := []string{
		`The Quick Brown Fox`,
		`Jumps Over The Lazy God`,
		`This is a multiline text`,
		`... and we are here to rock!`,
	}

	var size, spacing float64 = 12, 1

	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 512, 256))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(196)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// This quad will be used to render the texture
	// Format: X, Y, U, V, R, G, B
	quad := []float32{
		-0.9, -0.9, 0.0, 1.0, 1.0, 0.0, 0.0,
		-0.9, 0.9, 0.0, 0.0, 0.0, 1.0, 0.0,
		0.9, -0.9, 1.0, 1.0, 0.0, 0.0, 1.0,
		0.9, 0.9, 1.0, 0.0, 1.0, 1.0, 1.0,
	}

	// Create a VAO and VBO for quad (Texture)
	vao = glad.NewVertexArrayObject()
	vao.Bind()
	vbo = glad.NewVertexBufferObject()
	vbo.Bind()
	vbo.BufferData32(quad, gl.STATIC_DRAW)
	// Define attributes to draw the texture on quad
	attrPos = program.GetAttributeLocation("pos")
	attrPos.PointerFloat32(2, false, 7, 0)
	attrUV = program.GetAttributeLocation("uv")
	attrUV.PointerFloat32(2, false, 7, 2)
	attrCol = program.GetAttributeLocation("col")
	attrCol.PointerFloat32(3, false, 7, 4)

	// Create a framebuffer and a texture as render target
	txr = glad.NewTexture()
	txr.Bind()
	txr.Image2D(rgba)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)
	// Attach texture to target
	//fbo.Texture(gl.COLOR_ATTACHMENT0, txr)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	program.Use()
	attrPos.Enable()
	attrUV.Enable()
	attrCol.Enable()

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
in vec3 col;
out vec2 vUV;
out vec3 vCol;
void main() { gl_Position = vec4(pos, 0.0, 1.0); vUV = uv; vCol = col; }`
	fragmentShaderSource = `#version 440 core
in vec2 vUV;
in vec3 vCol;
out vec3 color;
uniform sampler2D sampler;
void main() { color = texture(sampler, vUV).rgb * vCol; }`
)
