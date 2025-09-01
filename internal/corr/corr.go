package corr

import (
	"image"
	"image/color"

	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/kernels"
	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/utils"
)

// essa função é a alma do precessamento, pois enquanto os outras compoem
// os dados e os leem, esta aplica a máscara a imagem e gerando a nova imagem

func AplicarCorrelação(img *image.RGBA, kernel *kernels.Kernel) *image.RGBA {
	dimensoes := img.Bounds()
	imgComFiltro := image.NewRGBA(dimensoes)

	kerLargura := kernel.Largura
	kerAltura := kernel.Altura

	a := kerLargura / 2 // offset em x
	b := kerAltura / 2  // offset em y

	// agora vamos iterar pelos pixels da imagem
	for y := dimensoes.Min.Y; y < dimensoes.Max.Y; y++ {
		for x := dimensoes.Min.X; x < dimensoes.Max.X; x++ {

			var sumR, sumG, sumB int

			// agora para cada pixel (central) vamos aplicar o kernel
			for ky := 0; ky < kerAltura; ky++ {
				for kx := 0; kx < kerLargura; kx++ {

					iy := y + ky - b // pixel da imagem + "pixel" do kernel - offset do kernel em altura
					ix := x + kx - a // pixel da imagem + "pixel" do kernel - offset do kernel em largura

					// caso o ponto extravaze a imagem, vamos ignorá-lo
					if ix < 0 || ix >= dimensoes.Max.X || iy < 0 || iy >= dimensoes.Max.Y {
						continue
					}

					rgb := img.RGBAAt(ix, iy) // pegando valor RGB daquele pixel
					pesoDaMáscara := kernel.Máscara[ky][kx]

					sumR += int(rgb.R) * pesoDaMáscara
					sumG += int(rgb.G) * pesoDaMáscara
					sumB += int(rgb.B) * pesoDaMáscara
				}
			}

			//bias
			sumR += kernel.Bias
			sumG += kernel.Bias
			sumB += kernel.Bias

			// escala
			if kernel.Escala != 0 && kernel.Escala != 1 {
				sumR /= kernel.Escala
				sumG /= kernel.Escala
				sumB /= kernel.Escala
			}

			// funcao de ativação
			switch kernel.FunçãoDeAtivação {
			case "relu":
				sumR = utils.Max(0, sumR)
				sumG = utils.Max(0, sumG)
				sumB = utils.Max(0, sumB)

			case "identidade":
				// nao faz nada

			}

			// caso especial da o kernel sobel
			if kernel.Nome == "sobel_v" || kernel.Nome == "sobel_h" {
				sumR = utils.Abs(sumR)
				sumG = utils.Abs(sumG)
				sumB = utils.Abs(sumB)
			}

			imgComFiltro.SetRGBA(x, y, color.RGBA{
				R: utils.Limiar(sumR, 0, 255),
				G: utils.Limiar(sumG, 0, 255),
				B: utils.Limiar(sumB, 0, 255),
				A: img.RGBAAt(x, y).A,
			})
		}
	}
	return imgComFiltro
}
