package main

import (
	"image"
	"math"
	"os"
	"runtime/trace"
	"sync"
	"time"

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

type ImageMandelBrotSettings struct {
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

	outFile := `image\out\mandelbrot.jpg`
	startTime := time.Now()

	var size float64 = 4

	width := 1000 * size
	height := 1000 * size
	mandel := image.NewRGBA(image.Rect(0, 0, (int)(width), (int)(height)))

	img := &ImageMandelBrotSettings{
		x1: -0.08, x2: -0.07,
		y1: -0.825, y2: -0.835,
		height: (uint32)(height),
		width:  (uint32)(width),
		image:  mandel,
		speed:  3,
	}

	newImage := PixelLoop(img)

	newFile, err := os.Create(outFile)
	check(err)

	q := jpeg.Options{Quality: 90}
	err = jpeg.Encode(newFile, newImage.image, &q)
	check(err)

	println("Total time elapsed :", time.Since(startTime).Milliseconds(), "ms")

}

// PixelLoop Loops Over every Pixel and apply the function then set the color
func PixelLoop(mb *ImageMandelBrotSettings) *ImageMandelBrotSettings {
	colorLookup := colorGen(mb.speed)

	wg := &sync.WaitGroup{}
	for y := 0; y <= (int)(mb.height); y++ {

		wg.Add(1)
		normalizedY := mathutil.Normalize((float64)(y), 0.0, (float64)(mb.height), mb.y1, mb.y2)

		go func(y int, normalizedY float64, wg *sync.WaitGroup) {

			for x := 0; x <= (int)(mb.width); x++ {
				normalizedX := mathutil.Normalize((float64)(x), 0.0, (float64)(mb.width), mb.x1, mb.x2)

				z := cplx.NewComplex(0, 0)
				c := cplx.NewComplex(normalizedX, normalizedY)

				maxIter := 500

				/*for i := 0; i < maxIter; i++ {
					if cplx.Abs(z) > 2 {
						mb.image.Set(x, y, colorLookup[i])
						break
					}
					z = cplx.Add(cplx.Mult(z, z), c)
				}*/

				for i := 0; i < maxIter; i++ {
					if cplx.Sq(z) > 4 {
						mb.image.Set(x, y, colorLookup[i])
						break
					}

					z.MultBy(z)
					z.AddTo(c)
				}

			}
			wg.Done()
		}(y, normalizedY, wg)
	}
	wg.Wait()

	return mb
}

// Implement the Go color.Color interface.
func (col Color) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R*65535.0 + 0.5)
	g = uint32(col.G*65535.0 + 0.5)
	b = uint32(col.B*65535.0 + 0.5)
	a = 0xFFFF
	return
}

func colorGen(speed float64) []Color {
	var c []Color

	for i := 0; i < 1000; i++ {
		angle := (float64)((int)((float64)(i)*speed+360/3*2) % 360.0)
		c = append(c, Hsv(angle, 1.0, 1.0))
	}

	return c
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

// Averaging general runs
// Go   Par time --  1,786 ms
// Rust Par time --  1,625 ms
