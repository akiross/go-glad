package glad

import (
	"image"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Binder interface {
	Bind()
	Unbind()
}

type Enabler interface {
	Enable()
	Disable()
}

// Binds the objects passed as arguments and return
// a function to unbind them in reverse order
func BlockBind(objs ...Binder) func() {
	for i := 0; i < len(objs); i++ {
		objs[i].Bind()
	}
	return func() {
		for i := len(objs) - 1; i >= 0; i-- {
			objs[i].Unbind()
		}
	}
}

// Same as BlockBind but with Enable/Disable
func BlockEnable(objs ...Enabler) func() {
	for i := 0; i < len(objs); i++ {
		objs[i].Enable()
	}
	return func() {
		for i := len(objs) - 1; i >= 0; i-- {
			objs[i].Disable()
		}
	}
}

// Create
func MakeProgram(shaders ...Shader) Program {
	program := NewProgram()
	program.AttachShaders(shaders...)
	program.Link()
	for i := range shaders {
		shaders[i].Delete()
	}
	return program
}

// CheckError checks for OpenGL errors and print them if any
// Returns true if any error was found
func CheckError() bool {
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Println("GL ERROR:", err)
		return true
	}
	return false
}

// TODO We could create a tool that allows to easily specify data and attributes in the same place
// that would replace defining the data and binding it to attributes: it would understand automatically
// the size of data, location in the array (offset, stride, type)
/*
	Easy definition of meshes (VAO+VBO+ATTRs)

	DefineMesh(Attr("pos", 2), Attr("col", 3), Attr("uv", 2), data)
*/

// Attr describes an attribute
type Attr struct {
	Buff int    // Which of the Data buffers will be used
	Name string // Name of the attribute in the shader
	Size int32  // Number of elements for this attribute
}

// TxrSpec describes a texture to be loaded
// type TxrSpec struct {
// 	Width, Height int        // Empty image of given size
// 	Path          string     // Load texture from file
// 	Image         image.RGBA // Load texture from data
// }

// Rect describes a rectangle with integer coordinates
type Rect struct {
	X, Y, W, H int
}

// Config contains the specifications for automated build
type Config struct {
	Shaders    []Shader
	Attributes []Attr
	Data       [][]float32   // Multiple buffers of data (one for each VBO)
	Elements   []int16       // Indices of elements to use. If not nil, gl.DrawElements will be used instead of gl.DrawArrays
	DataUsages []uint32      // gl.STATIC_DRAW, etc. one for each slice in Data. If Elements is not nil, the last DataUsage is used for the EBO
	Primitives uint32        // gl.TRIANGLES, gl.POINTS, etc.
	ClearColor []float32     // Clear color before drawing
	Textures   []Texture     // List of pre-existing textures to use (attached before images)
	Images     []image.Image // Images to use to create new textures (attached after textures)
	Offscreen  *Rect         // If not nil, will create and render to FBO setting Viewport
}

type AutoConfig struct {
	BgTxr    Texture
	Textures []Texture // List of pre-existing textures to use (attached before images)
	Cfg      *Config
	Prog     Program
	FBO      FramebufferObject
	VAO      VertexArrayObject
	VBOs     []VertexBufferObject
	NumVert  int32
	//bp      uint32 // Binding point
}

func AutoBuild(cfg *Config) *AutoConfig {
	var mo AutoConfig
	mo.Cfg = cfg

	// Setup shaders and program
	mo.Prog = NewProgram()
	mo.Prog.AttachShaders(cfg.Shaders...)
	mo.Prog.Link()

	for i := range cfg.Shaders {
		cfg.Shaders[i].Delete()
	}

	if cfg.Offscreen != nil {
		mo.FBO = NewFramebuffer()
		mo.BgTxr = NewTexture(gl.TEXTURE_2D)
		mo.BgTxr.Storage(1, gl.RGBA8, []int{cfg.Offscreen.W, cfg.Offscreen.H})
		mo.BgTxr.SetFilters(gl.NEAREST, gl.NEAREST) // FIXME use a setting
		mo.FBO.Texture(gl.COLOR_ATTACHMENT0, mo.BgTxr)
	}

	if cfg.ClearColor == nil {
		cfg.ClearColor = []float32{0.0, 0.0, 0.0, 1.0}
	}

	mo.VBOs = make([]VertexBufferObject, len(cfg.Data))
	for i := range cfg.Data {
		mo.VBOs[i] = NewVertexBufferObject()
		mo.VBOs[i].BufferData32(cfg.Data[i], cfg.DataUsages[i])
	}

	mo.VAO = NewVertexArrayObject()
	var offsets = make([]uint32, len(cfg.Data))
	// Prepare attributes
	for i := range cfg.Attributes {
		b := cfg.Attributes[i].Buff
		// Get attribute by name
		at := mo.Prog.GetAttributeLocation(cfg.Attributes[i].Name)
		// Specify format for attrib
		mo.VAO.AttribFormat32(at, cfg.Attributes[i].Size, offsets[b])
		// Next attribute starts where this ends: build relative offset
		offsets[b] += uint32(cfg.Attributes[i].Size)
		// Set binding
		mo.VAO.AttribBinding(uint32(b), at)
		// Enable attribute when using VAO
		mo.VAO.EnableAttrib(at)
	}

	// Set VBO specifiying the total stride (= sum of relative offsets)
	for i := range cfg.Data {
		mo.VAO.VertexBuffer32(uint32(i), mo.VBOs[i], 0, int32(offsets[i]))
	}

	// If elements are specified, create a VBO for that
	if cfg.Elements != nil {
		// The number of vertices to draw is given by elements array
		mo.NumVert = int32(len(cfg.Elements))
		// Now create a new Element Buffer Object
		var ebo uint32
		mo.VAO.Bind()
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		usage := cfg.DataUsages[len(cfg.DataUsages)-1]
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cfg.Elements)*2, unsafe.Pointer(&cfg.Elements[0]), usage)
		mo.VAO.Unbind()
	} else {
		// Compute number of vertices to draw
		mo.NumVert = int32(len(cfg.Data[0]) / int(offsets[0]))
		for i := 1; i < len(cfg.Data); i++ {
			nv := int32(len(cfg.Data[i]) / int(offsets[i]))
			if nv != mo.NumVert {
				panic("Inferred number of vertices not matching")
			}
		}
	}

	// Load images as textures
	if cfg.Images != nil {
		mo.Textures = make([]Texture, len(cfg.Images))
	}
	for i := range cfg.Images {
		txr := NewTexture(gl.TEXTURE_2D)
		txr.Storage(1, gl.RGBA8, []int{cfg.Images[i].Bounds().Dx(), cfg.Images[i].Bounds().Dy()})
		txr.Image2D(cfg.Images[i])
		txr.SetFilters(gl.NEAREST, gl.NEAREST)
		mo.Textures[i] = txr
	}

	return &mo
}

func (mo *AutoConfig) AutoDraw() {
	var bindUnit uint32
	if mo.Cfg.Offscreen != nil {
		mo.FBO.Bind()
		mo.BgTxr.Bind(bindUnit)
		bindUnit++
		gl.Viewport(
			int32(mo.Cfg.Offscreen.X),
			int32(mo.Cfg.Offscreen.Y),
			int32(mo.Cfg.Offscreen.W),
			int32(mo.Cfg.Offscreen.H))
		gl.ClearNamedFramebufferfv(uint32(mo.FBO), gl.COLOR, 0, &mo.Cfg.ClearColor[0])
	} else {
		gl.ClearNamedFramebufferfv(0, gl.COLOR, 0, &mo.Cfg.ClearColor[0])
	}

	// Bind pre-allocated textures
	for i := range mo.Cfg.Textures {
		mo.Cfg.Textures[i].Bind(bindUnit) // FIXME how does the shader know what is the texture to use?
		bindUnit++
	}
	// Bind image-loaded textures
	for i := range mo.Textures {
		mo.Textures[i].Bind(bindUnit)
		bindUnit++
	}

	mo.Prog.Use()
	mo.VAO.Bind()
	if mo.Cfg.Elements == nil {
		gl.DrawArrays(mo.Cfg.Primitives, 0, mo.NumVert)
	} else {
		gl.DrawElements(mo.Cfg.Primitives, mo.NumVert, gl.UNSIGNED_SHORT, nil)
	}
	mo.VAO.Unbind()

	for i := range mo.Textures {
		bindUnit--
		mo.Textures[i].Unbind(bindUnit)
	}

	for i := range mo.Cfg.Textures {
		bindUnit--
		mo.Cfg.Textures[i].Unbind(bindUnit)
	}
	bindUnit--
	if mo.Cfg.Offscreen != nil {
		mo.BgTxr.Unbind(bindUnit)
		mo.FBO.Unbind()
	}
}

// UpdateImage reloads the data of the i-th image into i-th texture
func (mo *AutoConfig) UpdateImage(i int) {
	mo.Textures[i].Image2D(mo.Cfg.Images[i])
}
