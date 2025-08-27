package hist

import (
	"image"
	"image/color"
	"sync"
)

func Expansão(img *image.RGBA) *image.RGBA {
	dimensões := img.Bounds()

	minR, maxR := 255, 0
	minG, maxG := 255, 0
	minB, maxB := 255, 0

	// a condição de corrida nas variais min e max não compensa tornar o precessamento concorrente
	for y := dimensões.Min.Y; y < dimensões.Max.Y; y++ {
		for x := dimensões.Min.X; x < dimensões.Max.X; x++ {
			rgb := img.RGBAAt(x, y)
			r, g, b := int(rgb.R), int(rgb.G), int(rgb.B)

			// min e max em R
			if r < minR {
				minR = r
			}
			if r > maxR {
				maxR = r
			}

			// min e max em G
			if g < minG {
				minG = g
			}
			if g > maxG {
				maxG = g
			}

			// min e max em B
			if b < minB {
				minB = b
			}
			if b > maxB {
				maxB = b
			}
		}
	}

	rangeR := maxR - minR
	rangeG := maxG - minG
	rangeB := maxB - minB

	//eliminando a possibilidade de divisão por zero
	if rangeR == 0 {
		rangeR = 1
	}
	if rangeG == 0 {
		rangeG = 1
	}
	if rangeB == 0 {
		rangeB = 1
	}

	// aplicando a fórmula
	var semaforo sync.WaitGroup
	imagemNova := image.NewRGBA(dimensões)
	for y := dimensões.Min.Y; y < dimensões.Max.Y; y++ {
		for x := dimensões.Min.X; x < dimensões.Max.X; x++ {
			semaforo.Add(1)
			// agora podemos tornas a aplicação da fórmula concorrente
			go func(x, y int) {
				defer semaforo.Done()
				rgb := img.RGBAAt(x, y)

				newR := uint8((int(rgb.R) - minR) * 255 / rangeR)
				newG := uint8((int(rgb.G) - minG) * 255 / rangeG)
				newB := uint8((int(rgb.B) - minB) * 255 / rangeB)

				imagemNova.SetRGBA(x, y, color.RGBA{
					R: newR,
					G: newG,
					B: newB,
					A: rgb.A,
				})

			}(x, y)
		}
	}
	semaforo.Wait()
	return imagemNova
}
