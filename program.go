package glad

import (
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// Program represents a program in the OpenGL context
type Program uint32

// PrepareProgram creates a program, attaches shaders and links before returning
func PrepareProgram(shaders ...Shader) Program {
	pr := NewProgram()
	pr.AttachShaders(shaders...)
	pr.Link()
	return pr
}

// NewProgram creates a program with a new name
func NewProgram() Program {
	return Program(gl.CreateProgram())
}

// AttachShaders attaches one or more shaders to the program
func (pr Program) AttachShaders(shaders ...Shader) {
	for _, sh := range shaders {
		gl.AttachShader(uint32(pr), uint32(sh))
	}
}

// Link links the attached shaders belonging to the program
// logging an error if link did not succeed
func (pr Program) Link() {
	gl.LinkProgram(uint32(pr))

	if pr.GetParameter(gl.LINK_STATUS) == gl.FALSE {
		infoLog := pr.GetInfoLog()
		log.Fatalln("Unable to link program:\n", infoLog)
	}
}

func (pr Program) GetParameter(pname uint32) int32 {
	var val int32
	gl.GetProgramiv(uint32(pr), pname, &val)
	return val
}

func (pr Program) GetInfoLog() string {
	logLen := pr.GetParameter(gl.INFO_LOG_LENGTH)
	infoLog := string(make([]byte, int(logLen+1)))
	var savedLen int32
	gl.GetProgramInfoLog(uint32(pr), logLen, &savedLen, gl.Str(infoLog))
	if savedLen+1 != logLen {
		log.Println("Program Info Log different lengths reported:", logLen, savedLen)
	}
	return infoLog
}

// Call this before linking to set location of attributes
func (pr Program) BindAttributeLocation(index uint32, name string) {
	cname := gl.Str(name + "\x00")
	gl.BindAttribLocation(uint32(pr), index, cname)
}

// Call this after linking to get location of attributes (e.g. if they were not set)
func (pr Program) GetAttributeLocation(name string) VertexAttrib {
	return VertexAttrib(gl.GetAttribLocation(uint32(pr), gl.Str(name+"\x00")))
}

func (pr Program) Use() {
	gl.UseProgram(uint32(pr))
}

func (pr Program) Delete() {
	gl.DeleteProgram(uint32(pr))
}
