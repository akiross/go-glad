package main

// Hello world in OpenGL: create a triangle on screen

import (
	"image"
	"math/rand"

	glad "github.com/akiross/go-glad"
	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v4.5-core/gl"

	// "github.com/llgcode/draw2d"
	// "github.com/llgcode/draw2d/draw2dgl"
	// "github.com/llgcode/draw2d/draw2dkit"

	"log"
	"runtime"
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

	var bgCol = []float32{0.3, 0.3, 0.3, 1.0}

	var (
		vertexShaderSource = `
			#version 440 core
			in vec2 pos;

			in vec3 col;
			out vec3 vCol;

			in vec2 uv;
			out vec2 vUV;

			void main() {
				gl_Position = vec4(pos, 0.0, 1.0);
				vCol = col;
				vUV = uv;
			}`
		fragmentShaderSource = `
			#version 440 core
			in vec3 vCol;
			out vec4 color;

			in vec2 vUV;
			uniform sampler2D sampler;

			void main() {
				color = texture(sampler, vUV);
			}`
	)

	vertShader := glad.NewShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragShader := glad.NewShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	program := glad.NewProgram()
	program.AttachShaders(vertShader, fragShader)
	program.Link()

	vertShader.Delete()
	fragShader.Delete()

	// Data to be used when drawing
	// Format: X, Y, R, G, B, U, V
	vertPosAndCol := []float32{
		+0.9, +0.9, 1.0, 1.0, 0.0, 1.0, 0.0,
		-0.9, +0.9, 0.0, 1.0, 0.0, 0.0, 0.0,
		+0.9, -0.9, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.9, -0.9, 1.0, 0.0, 0.0, 0.0, 1.0,
	}

	vbo := glad.NewVertexBufferObject()             // Create VBO
	vbo.BufferData32(vertPosAndCol, gl.STATIC_DRAW) // Set VBO data

	var bindPos uint32                     // Bind position in the VAO
	vao := glad.NewVertexArrayObject()     // Create VAO
	vao.VertexBuffer32(bindPos, vbo, 0, 7) // Bind VBO to position 0 in VAO, no offset, 5 elements stride

	attrPos := program.GetAttributeLocation("pos") // Find attribute position for "pos"
	vao.AttribFormat32(attrPos, 2, 0)              // Specify attribute size and relative offset
	vao.AttribBinding(bindPos, attrPos)            // Bind attributes to position 0
	vao.EnableAttrib(attrPos)

	attrCol := program.GetAttributeLocation("col")
	vao.AttribFormat32(attrCol, 3, 2)
	vao.AttribBinding(bindPos, attrCol)
	vao.EnableAttrib(attrCol)

	attrUV := program.GetAttributeLocation("uv")
	vao.AttribFormat32(attrUV, 2, 5)
	vao.AttribBinding(bindPos, attrUV)
	vao.EnableAttrib(attrUV)

	tw, th := 100.0, 75.0
	txr := glad.NewTexture(gl.TEXTURE_2D)
	txr.Storage(1, gl.RGBA8, []int{int(tw), int(th)})
	txr.Bind(0)
	txr.SetFilters(gl.NEAREST, gl.NEAREST)

	img := image.NewRGBA(image.Rect(0, 0, int(tw), int(th)))

	// Create context and draw a bouncing ball
	ctx := gg.NewContextForRGBA(img)

	vao.Bind()

	r := 10.0
	x, y := r+rand.Float64()*(tw-2.0*r), r+rand.Float64()*(th-2.0*r)
	dx, dy := 2*rand.Float64()-0.5, 2*rand.Float64()-0.5

	for !win.ShouldClose() {
		gl.ClearBufferfv(gl.COLOR, 0, &bgCol[0])
		program.Use()
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		if x > tw-r {
			x = tw - r
			dx = -dx
		} else if x < r {
			x = r
			dx = -dx
		}
		if y > th-r {
			y = th - r
			dy = -dy
		} else if y < r {
			y = r
			dy = -dy
		}

		ctx.SetRGB(0, 0, 0)
		ctx.Clear()

		ctx.DrawCircle(x, y, r)
		ctx.SetRGB(1, 0, 0)
		ctx.Fill()

		x += dx
		y += dy

		txr.Image2D(img)

		win.SwapBuffers()
		glad.PollEvents()
	}
}
