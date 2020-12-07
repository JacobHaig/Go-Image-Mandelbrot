package main

import (
	"image"
	"math"
	"os"
	"runtime/trace"
	"sync"

	"image/jpeg"

	// Using this complex number impl doubles our speed
	cplx "mandelbrot/complex"
	"mandelbrot/mathutil"
)

/// go run main.go 2> trace.out
/// go tool trace trace.out

// Check is used just to panic at unexpected Errors.
// Use this if you want the program to panic during Errors.
func check(err error) {
	if err != nil {
		print("Error Occured: ")
		panic(err)
	}
}

// Other ImageMandelBrot settings
// x1 = -2.0,   x2 = 1.0,    y1 = 1.0,    y2 = -1.0
// x1 = -0.2,   x2 = 0.0,    y1 = -0.8,   y2 = -1.0
// x1 = -0.1,   x2 = -0.05,  y1 = -0.8,   y2 = -0.85

type MandelBrotSettings struct {
	x1, y1, x2, y2 float64
	height, width  uint32
	image          *image.RGBA
	speed          float64
}

type Color struct {
	R, G, B float64
}

// Start of the program, general setup and reading of directories
func main() {

	// Change to true and comment out prints to preform a trace
	if false {
		trace.Start(os.Stderr)
		defer trace.Stop()
	}

	//startTime := time.Now()

	// Set up Mandelbrot Settings
	outFile := `image\out\mandelbrot.jpg`
	var size float64 = 4
	width := 1000 * size
	height := 1000 * size

	img := &MandelBrotSettings{
		x1: -0.08, x2: -0.07,
		y1: -0.825, y2: -0.835,

		height: (uint32)(height),
		width:  (uint32)(width),
		image:  image.NewRGBA(image.Rect(0, 0, (int)(width), (int)(height))),
		speed:  3,
	}

	newImage := PixelLoop(img)

	// Save the image with the specified file name
	newFile, err := os.Create(outFile)
	check(err)

	q := jpeg.Options{Quality: 90}
	err = jpeg.Encode(newFile, newImage.image, &q)
	check(err)

	//println("Total time elapsed :", time.Since(startTime).Milliseconds(), "ms")
}

// PixelLoop preforms the computation of mandelbrot image with the given Mandelbrot settings
func PixelLoop(mb *MandelBrotSettings) *MandelBrotSettings {
	// Generate lookup arrays
	colorLookup := colorGen(mb.speed)
	normLookup := normalGen(mb) // normLookup[0] is X, normLookup[1] is Y

	wg := &sync.WaitGroup{}
	for y, normY := range normLookup[1] {

		wg.Add(1)
		go func(y int, normY float64, wg *sync.WaitGroup) {

			for x, normX := range normLookup[0] {
				z := cplx.NewComplex(0, 0)
				c := cplx.NewComplex(normX, normY)

				maxIter := 500

				// Complex Loop
				for i := 0; i < maxIter; i++ {
					if cplx.Sq(z) > 4 {
						// Apply color to everything that eventailly ends,
						// If a point does converge it will be colored.
						mb.image.Set(x, y, colorLookup[i])
						break
					}

					z.MultBy(z).AddTo(c)
				}
				// If the Complex Loop doesnt converge, the color is impliately Black
			}
			wg.Done()
		}(y, normY, wg)

	}
	wg.Wait()

	return mb
}

// normalGen generates all possible X and Y value that could be iterated over.
// This is faster than creating a 2d array of all the pairs as every value is
// used more than once. For instance, keep track of the X value in the 2d array:
// 0,0  1,0  2,0
// 0,1  1,1  2,1
// 0,2  1,2  2,2
// Since X and Y are the same for there column/row, we can just store them in an array
func normalGen(mb *MandelBrotSettings) [][]float64 {
	arr := make([][]float64, 2)
	arr[0] = make([]float64, mb.width)
	arr[1] = make([]float64, mb.height)

	for x := 0; x < (int)(mb.width); x++ {
		arr[0][x] = mathutil.Normalize((float64)(x), 0.0, (float64)(mb.width), mb.x1, mb.x2)
	}
	for y := 0; y < (int)(mb.height); y++ {
		arr[1][y] = mathutil.Normalize((float64)(y), 0.0, (float64)(mb.height), mb.y1, mb.y2)
	}

	return arr
}

// Implement the Go color.Color interface.
func (col Color) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R*65535.0 + 0.5)
	g = uint32(col.G*65535.0 + 0.5)
	b = uint32(col.B*65535.0 + 0.5)
	a = 0xFFFF
	return
}

func colorGen(speed float64) (c []Color) {
	for i := 0; i < 1000; i++ {
		angle := (float64)((int)((float64)(i)*speed+360/3*2) % 360.0)
		c = append(c, Hsv(angle, 1.0, 1.0))
	}
	return
}

// Hsv creates a new Color given a Hue in [0..360], a Saturation and a Value in [0..1]
func Hsv(H, S, V float64) Color {
	Hp := H / 60.0
	C := V * S
	X := C * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := V - C
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = C
		g = X
	case 1.0 <= Hp && Hp < 2.0:
		r = X
		g = C
	case 2.0 <= Hp && Hp < 3.0:
		g = C
		b = X
	case 3.0 <= Hp && Hp < 4.0:
		g = X
		b = C
	case 4.0 <= Hp && Hp < 5.0:
		r = X
		b = C
	case 5.0 <= Hp && Hp < 6.0:
		r = C
		b = X
	}

	return Color{m + r, m + g, m + b}
}
