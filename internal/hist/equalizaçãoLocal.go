package hist

import (
	"image"
	"image/color"

	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/utils"
)

func EqualizaçãoLocal(img *image.RGBA, m, n int) *image.RGBA {
	dimensões := img.Bounds()
	imgNova := image.NewRGBA(dimensões)

	altura := dimensões.Dy()
	largura := dimensões.Dx()

	// para equalização não estravazar a janela (offset)
	metadeM := m / 2
	metadeN := n / 2

	for y := dimensões.Min.Y; y < dimensões.Max.Y; y++ {
		for x := dimensões.Min.X; x < dimensões.Max.X; x++ {
			/*
				// definindo a janela ao redor do pixel central na largura
				minX := utils.Limiar(x-metadeM, 0, largura-1)
				maxX := utils.Limiar(x+metadeM, 0, largura-1)

				// definindo a janela ao redor do pixel central na altura
				minY := utils.Limiar(y-metadeN, 0, altura-1)
				maxY := utils.Limiar(y+metadeN, 0, altura-1)
			*/

			minX := utils.Max(0, x-metadeM)
			maxX := utils.Min(largura-1, x+metadeM)

			minY := utils.Max(0, y-metadeN)
			maxY := utils.Min(altura-1, y+metadeN)

			// calculando histograma dentro da janela
			histR := make([]int, 256)
			histG := make([]int, 256)
			histB := make([]int, 256)

			for j := minY; j < maxY; j++ {
				for i := minX; i < maxX; i++ {
					rgb := img.RGBAAt(int(i), int(j))
					histR[rgb.R]++
					histG[rgb.G]++
					histB[rgb.B]++
				}
			}

			// calculando a funcão acumulada da janela
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

			cdfMinR := 0
			cdfMinG := 0
			cdfMinB := 0

			for i := 0; i < 256; i++ {
				if cdfR[i] != 0 && cdfMinR == 0 {
					cdfMinR = cdfR[i]
				}
				if cdfG[i] != 0 && cdfMinG == 0 {
					cdfMinG = cdfG[i]
				}
				if cdfB[i] != 0 && cdfMinB == 0 {
					cdfMinB = cdfB[i]
				}

			}

			mapR := make([]uint8, 256)
			mapG := make([]uint8, 256)
			mapB := make([]uint8, 256)

			totalPixels := int((maxX - minX) * (maxY - minY))

			for i := 0; i < 256; i++ {
				mapR[i] = uint8(((cdfR[i] - cdfMinR) * 255) / (totalPixels - cdfMinR))
				mapG[i] = uint8(((cdfG[i] - cdfMinG) * 255) / (totalPixels - cdfMinG))
				mapB[i] = uint8(((cdfB[i] - cdfMinB) * 255) / (totalPixels - cdfMinB))

			}

			rgb := img.RGBAAt(x, y)
			imgNova.SetRGBA(x, y, color.RGBA{
				R: mapR[rgb.R],
				G: mapG[rgb.G],
				B: mapB[rgb.B],
				A: rgb.A,
			})

		}
	}

	return imgNova
}
