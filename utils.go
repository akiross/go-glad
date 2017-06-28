package goglad

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

// TODO We could create a tool that allows to easily specify data and attributes in the same place
// that would replace defining the data and binding it to attributes: it would understand automatically
// the size of data, location in the array (offset, stride, type)
/*
	Easy definition of meshes (VAO+VBO+ATTRs)

	DefineMesh(Attr("pos", 2), Attr("col", 3), Attr("uv", 2), data)
*/

type attrSpec struct {
	name string
	size int
}

type dataSetter struct {
	specs []attrSpec
}

func (ds dataSetter) Float32(data ...float32) {
}

func (ds dataSetter) Float64(data ...float64) {
}

func (ds dataSetter) Int32(data ...int32) {
}

func Attr(name string, size int) attrSpec {
	return attrSpec{name, size}
}

func DefineMesh(specs ...attrSpec) dataSetter {
	var ds dataSetter
	ds.specs = specs
	return ds
}
