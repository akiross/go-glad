package main

// Two triangles in OpenGL, separated by VBO, using "Auto Build" utility

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

	win := glad.NewOGLWindow(800, 600, "Easy Builder",
		glad.CoreProfile(true),
		glad.Resizable(false),
		glad.ContextVersion(4, 5),
		//glad.VSync(true),
	)
	defer glad.Terminate()
	// Enable VSync
	glad.SwapInterval(1)

	var (
		vssTriangle = `#version 440 core
in vec2 pos;
in vec3 col;
out vec3 vCol;
void main() {
	gl_Position = vec4(pos, 0.0, 1.0);
	vCol = col;
}`
		fssTriangle = `#version 440 core
in vec3 vCol;
out vec3 color;
void main() { color = vCol; }`
		vssTexture = `#version 440 core
in vec2 pos;
in vec2 uv;
out vec2 vUV;
// in vec3 col;
// out vec3 vCol;
void main() {
	gl_Position = vec4(pos, 0.0, 1.0);
	vUV = uv;
	// vCol = col;
}`
		fssTexture = `#version 440 core
// in vec3 vCol;
in vec2 vUV;
out vec3 color;
uniform sampler2D sampler;
void main() {
	color = texture(sampler, vUV).rgb;
	// color = vCol;
}`
	)

	// Create a texture
	txrImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			txrImg.SetRGBA(x, y, color.RGBA{uint8(255.0 * float32(x%8) / 7.0), uint8(255 * float32(y%16) / 15.0), 0, 255})
		}
	}

	// First example: a triangle with a different color for each vertex.
	// We will use a single VBO and 2 attributes, "pos" and "col".
	// The data is interleaved
	autoTri := glad.AutoBuild(&glad.Config{
		Shaders: []glad.Shader{
			glad.NewShader(vssTriangle, gl.VERTEX_SHADER),
			glad.NewShader(fssTriangle, gl.FRAGMENT_SHADER),
		},
		// Attributes are read in order from each buffer
		Attributes: []glad.Attr{{0, "pos", 2}, {0, "col", 3}},
		Data: [][]float32{
			[]float32{ // Interleaved position and color data
				-1.0, -1.0, 1.0, 0.0, 0.0,
				0.0, 1.0, 0.0, 1.0, 0.0,
				1.0, -1.0, 0.0, 0.0, 1.0,
			},
		},
		DataUsages: []uint32{gl.STATIC_DRAW},
		Primitives: gl.TRIANGLES,
		Offscreen:  &glad.Rect{0, 0, 800, 600},
	})

	// Second example: two triangles that form a square.
	// We demonstrate how to use separated buffers, and one attribute per buffer.
	// Also, instead of using a TRIANGLE_STRIP, we use 2 triangles and one
	// element buffer to specify the indices.
	autoScr := glad.AutoBuild(&glad.Config{
		Shaders: []glad.Shader{
			glad.NewShader(vssTexture, gl.VERTEX_SHADER),
			glad.NewShader(fssTexture, gl.FRAGMENT_SHADER),
		},
		Attributes: []glad.Attr{{0, "pos", 2}, {1, "uv", 2}}, // First field specifies the Data array to use
		Data: [][]float32{
			[]float32{-0.9, -0.9, -0.9, 0.9, 0.9, -0.9, 0.9, 0.9}, // pos data
			[]float32{0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1.0},     // uv data
		},
		DataUsages: []uint32{gl.STATIC_DRAW, gl.STATIC_DRAW, gl.STATIC_DRAW}, // Last value is for EBO data usage
		//Primitives: gl.TRIANGLE_STRIP, // We could use TRIANGLE_STRIP and it would work
		// But here we use two TRIANGLES and use index drawing
		Primitives: gl.TRIANGLES,
		Elements:   []int16{0, 1, 2, 1, 3, 2},
		Textures:   []glad.Texture{autoTri.BgTxr},
		//Images:     []image.Image{txrImg},
		ClearColor: []float32{0.6, 0.6, 0.6, 1.0},
	})

	autoTri.AutoDraw() // Draw once

	for !win.ShouldClose() {
		autoScr.AutoDraw() // Draw always
		win.SwapBuffers()
		glad.PollEvents()
	}
}
