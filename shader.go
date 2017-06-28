package goglad

import (
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Shader uint32

func NewShader(source string, shaderType uint32) Shader {
	if source == "" {
		log.Fatalln("Unable to create shader from empty string")
	}
	var sh Shader
	sh = Shader(gl.CreateShader(shaderType))
	csrc, free := gl.Strs(source + "\x00")
	gl.ShaderSource(uint32(sh), 1, csrc, nil)
	free()
	gl.CompileShader(uint32(sh))

	if sh.GetParameter(gl.COMPILE_STATUS) == gl.FALSE {
		infoLog := sh.GetInfoLog()
		log.Fatalln("Unable to compile shader:\n", infoLog)
	}
	return sh
}

func (sh Shader) Delete() {
	gl.DeleteShader(uint32(sh))
}

func (sh Shader) GetParameter(pname uint32) int32 {
	var val int32
	gl.GetShaderiv(uint32(sh), pname, &val)
	return val
}

func (sh Shader) GetInfoLog() string {
	logLen := sh.GetParameter(gl.INFO_LOG_LENGTH)
	infoLog := string(make([]byte, int(logLen+1)))
	var savedLen int32
	gl.GetShaderInfoLog(uint32(sh), logLen, &savedLen, gl.Str(infoLog))
	if savedLen+1 != logLen {
		log.Println("Shader Info Log different lengths reported:", logLen, savedLen)
	}
	return infoLog
}
