package hist

import (
	"image"
	"image/color"
	"sync"

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

	var mu sync.Mutex
	var semaforoPixels sync.WaitGroup

	for y := 0; y < altura; y++ {
		for x := 0; x < largura; x++ {
			/*
				É tentador tornar essa execução concorrente, porém a imgNova está na sessão crítica
				da execução por isso não estou certo se isso gera respostas não completas. Por isso vou
				usar um mutex e torcer para não ter condição de corrida
			*/
			semaforoPixels.Add(1)
			go func(x, y int) {
				defer semaforoPixels.Done()
				var semaforo sync.WaitGroup

				// definindo a janela ao redor do pixel central na largura
				minX := utils.Limiar(x-metadeM, 0, largura-1)
				maxX := utils.Limiar(x+metadeM, 0, largura-1)

				// definindo a janela ao redor do pixel central na altura
				minY := utils.Limiar(y-metadeN, 0, altura-1)
				maxY := utils.Limiar(y+metadeN, 0, altura-1)

				// calculando histograma dentro da janela
				histR := make([]int, 256)
				histG := make([]int, 256)
				histB := make([]int, 256)

				for j := minY; j < maxY; j++ {
					for i := minX; i < maxX; i++ {
						rgb := img.RGBAAt(x, y)
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

				for i := 0; i < 256; i++ {
					cdfR[i] = cdfR[i-1] + histR[i]
					cdfG[i] = cdfG[i-1] + histG[i]
					cdfB[i] = cdfB[i-1] + histB[i]
				}

				cdfMinR := 0
				cdfMinG := 0
				cdfMinB := 0

				for i := 0; i < 256; i++ {
					// como cada variavel so é atualizada (excrita) uma vez e o resto é leitura, vamos tornar concorrente
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
				semaforo.Wait() // como vamos tornar a proxima ação também concorrente, precisamos esperar a execução acima terminar por completo primeiro

				mapR := make([]uint8, 256)
				mapG := make([]uint8, 256)
				mapB := make([]uint8, 256)

				totalPixels := int((maxX - minX + 1) * (maxY - minY + 1))
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

				rgb := img.RGBAAt(x, y)
				mu.Lock()
				imgNova.SetRGBA(x, y, color.RGBA{
					R: mapR[rgb.R],
					G: mapG[rgb.G],
					B: mapB[rgb.B],
					A: rgb.A,
				})
				mu.Unlock()

			}(x, y)
		}
	}
	semaforoPixels.Wait()
	return imgNova
}
