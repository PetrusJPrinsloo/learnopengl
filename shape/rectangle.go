package shape

var (
	Rectangle = []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}
	Indices = []uint32{
		0, 1, 3,
		1, 2, 3,
	}
)
