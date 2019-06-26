package main

// Show a font with some text

import (
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"log"
	"runtime"

	glad "github.com/akiross/go-glad"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/golang/freetype"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting")

	win := glad.NewOGLWindow(800, 600, "Atlas",
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
in vec2 uv;
in vec3 col;
out vec2 vUV;
out vec3 vCol;
void main() {
	gl_Position = vec4(pos, 0.0, 1.0);
	vUV = uv;
	vCol = col;
}`
		fragmentShaderSource = `#version 450 core
in vec2 vUV;
in vec3 vCol;
out vec3 color;
uniform sampler2D sampler;
void main() {
	color = texture(sampler, vUV).rgb * vCol;
	// color = vCol;
}`
	)

	vertShader := glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader := glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	// program := glad.PrepareProgram(vertShader, fragShader)
	program := glad.NewProgram()
	program.AttachShaders(vertShader, fragShader)
	program.Link()

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

	// Create a framebuffer and a texture as render target
	txr := glad.NewTexture(gl.TEXTURE_2D)
	txr.Storage(1, gl.RGBA8, []int{512, 256})
	txr.Bind(0)
	txr.Image2D(rgba)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)
	// Attach texture to target
	// fbo.Texture(gl.COLOR_ATTACHMENT0, txr)

	// This quad will be used to render the texture
	// Format: X, Y, U, V, R, G, B
	quad := []float32{
		-0.9, -0.9, 0.0, 1.0, 1.0, 0.0, 0.0,
		-0.9, +0.9, 0.0, 0.0, 0.0, 1.0, 0.0,
		+0.9, -0.9, 1.0, 1.0, 0.0, 0.0, 1.0,
		+0.9, +0.9, 1.0, 0.0, 1.0, 1.0, 1.0,
	}

	// Create a buffer object to store the quad data above
	vbo := glad.NewVertexBufferObject()
	vbo.BufferData32(quad, gl.STATIC_DRAW) // FIXME this could be a BufferStorage32 as well

	// Create a VAO to associate shader attributes to the buffer
	var bindPos uint32 // Bind position for the VBO
	vao := glad.NewVertexArrayObject()

	vao.VertexBuffer32(bindPos, vbo, 0, 7) // Specify a buffer to be used, this is done once

	// Define attributes to draw the texture on quad
	attrPos := program.GetAttributeLocation("pos")
	vao.AttribFormat32(attrPos, 2, 0)   // Specify the format for the data of the attribute
	vao.AttribBinding(bindPos, attrPos) // Bind the attribute to the bind index
	vao.EnableAttrib(attrPos)

	attrUV := program.GetAttributeLocation("uv")
	vao.AttribFormat32(attrUV, 2, 2) // Same as above, but convenient for floats
	vao.AttribBinding(bindPos, attrUV)
	vao.EnableAttrib(attrUV)

	attrCol := program.GetAttributeLocation("col")
	vao.AttribFormat32(attrCol, 3, 4)
	vao.AttribBinding(bindPos, attrCol)
	vao.EnableAttrib(attrCol)

	program.Use()
	vao.Bind()

	for !win.ShouldClose() {
		gl.ClearBufferfv(gl.COLOR, 0, &bgCol[0])
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4) // Draw a quad
		win.SwapBuffers()
		glad.PollEvents()
	}
}
