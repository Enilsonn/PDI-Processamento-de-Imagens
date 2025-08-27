package hist

import (
	"image"
	"image/color"
	"sync"
)

// https://github.com/hybridgroup/gocv
// lib GoCV similar ao que o OpenCV faz em Python

// go get -u -d gocv.io/x/gocv

// esuqalização de histograma global

func Equalização(img *image.RGBA) *image.RGBA {
	dimensões := img.Bounds()
	largura := dimensões.Dx()
	altura := dimensões.Dy()
	var semaforo sync.WaitGroup

	// as frequências
	histR := make([]int, 256)
	histG := make([]int, 256)
	histB := make([]int, 256)

	// calculando histograma
	for y := 0; y < altura; y++ {
		for x := 0; y < largura; x++ {
			// como nesse caso varias gorrotines podem acessar a mesma variavel, histR[i] por exemplo,
			// tornar esse processamento concorrente não dá grandes ganhos por causa da quantidade de
			// condição de corrida
			rgb := img.RGBAAt(x, y)
			histR[rgb.R]++
			histG[rgb.G]++
			histB[rgb.B]++
		}
	}

	// acumu;lada
	cdfR := make([]int, 256)
	cdfG := make([]int, 256)
	cdfB := make([]int, 256)

	cdfR[0] = histR[0]
	cdfG[0] = histG[0]
	cdfB[0] = histB[0]

	for i := 1; i < 256; i++ {
		cdfR[i] = cdfR[i-1] + histR[i]
		cdfG[i] = cdfG[i-1] + histG[i]
		cdfB[i] = cdfB[i-1] + histB[i]
	}

	totalPixels := altura * largura

	// normalizando a acumulada
	mapR := make([]uint8, 256)
	mapG := make([]uint8, 256)
	mapB := make([]uint8, 256)

	// achar o menor valor > 0 para a acumulada de R G e B
	cdfMinR := 0
	cdfMinG := 0
	cdfMinB := 0
	for i := 0; i < 256; i++ {
		semaforo.Add(1)
		go func(i int) {
			defer semaforo.Done()
			if cdfR[i] != 0 && cdfMinR == 0 {
				cdfMinR = cdfR[i]
			}
			if cdfG[i] != 0 && cdfMinG == 0 {
				cdfMinG = cdfG[i]
			}
			if cdfB[i] != 0 && cdfMinB == 0 {
				cdfMinB = cdfB[i]
			}
		}(i)
	}
	semaforo.Wait()

	for i := 0; i < 256; i++ {
		semaforo.Add(1)
		go func(i int) {
			defer semaforo.Done()
			mapR[i] = uint8(((cdfR[i] - cdfMinR) * 255) / (totalPixels - cdfMinR))
			mapG[i] = uint8(((cdfG[i] - cdfMinG) * 255) / (totalPixels - cdfMinG))
			mapB[i] = uint8(((cdfB[i] - cdfMinB) * 255) / (totalPixels - cdfMinB))
		}(i)
	}
	semaforo.Wait()

	// aplicando mapeamento nos pixels
	imgNova := image.NewRGBA(dimensões)
	for y := 0; y < altura; y++ {
		for x := 0; x < largura; x++ {
			semaforo.Add(1)
			// na sessão crrítica desse código há apenas condição de corrida de leitura, o que não é um problema
			go func(x, y int) {
				defer semaforo.Done()
				rgb := img.RGBAAt(x, y)
				imgNova.SetRGBA(x, y, color.RGBA{
					R: mapR[rgb.R],
					G: mapG[rgb.G],
					B: mapB[rgb.B],
					A: rgb.A,
				})
			}(x, y)
		}
	}
	semaforo.Wait()

	return imgNova
}
