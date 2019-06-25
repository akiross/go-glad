package glad

//import "unicode/utf8"

// Function to define an atlas (for a given alphabet)
type Atlas struct {
	minX, minY, maxX, maxY []float32
	atlas                  Texture
}

func NewAtlas(abc string) *Atlas {
	var atl Atlas
	/*
		l := utf8.RuneCountInString(abc)
		atl.minX = make([]float32, l)
		atl.maxX = make([]float32, l)
		atl.minY = make([]float32, l)
		atl.maxY = make([]float32, l)

		a = NewTexture()
		atl.atlas = a
	*/
	return &atl
}
