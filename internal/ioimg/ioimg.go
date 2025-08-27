package ioimg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/tiff"
)

func Open(path string) (image.Image, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar acessar diretório: %v", err)
	}
	defer file.Close()

	// identificando o formato da imagem para decodificação
	ext := filepath.Ext(path)

	// decodificação
	switch ext {
	case ".png", ".PNG":
		img, err := png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("erro ao decodificar png: %v", err)
		}
		return img, nil
	case ".jpg", "jpeg", ".JPG", ".JPEG":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("erro ao decodificar jpeg: %v", err)
		}
		return img, nil

	case ".tif", ".tiff":
		img, err := tiff.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("erro ao decoficicar tiff: %v", err)
		}
		return img, nil
	default:
		// peng e jpeg sao os mais comuns, porém se o formato da imagem for
		// outro, podemos tentar converter de forma genérica com 'image'.Decode
		img, _, err := image.Decode(file) // o valor omitido é o nome do arquivo
		if err != nil {
			return nil, fmt.Errorf("erro ao processar imagem %s: %v", ext, err)
		}
		return img, nil
	}
}

func Save(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	// identificando o formato da imagem para codificação
	ext := filepath.Ext(path)

	// codificação
	switch ext {
	case ".png", ".PNG":
		if err := png.Encode(file, img); err != nil {
			return fmt.Errorf("erro ao codificar png: %v", err)
		}
	case ".jpg", "jpeg", ".JPG", ".JPEG":
		// a função Encode para jpeg da lib padrão pede que seja passada a
		// qualidade obrigatoriamente da imagem coloquei 90 (1 - 100), mas
		// é um magic number
		if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("erro ao codificar jpeg: %v", err)
		}

	case ".tif", ".tiff":
		if err := tiff.Encode(file, img, nil); err != nil {
			return fmt.Errorf("erro ao codificar tiff: %v", err)
		}
	default:
		return fmt.Errorf("extensão não suportada, escolha entre .png .PNG .jpg .jpeg .JPG .JPEG: %v", err)
	}
	return nil
}

// como o valor do pixel em uma banda é (-256, 255) usaremos uint8 (8bits)
// ademais o Deconde e Encode usados acima usam 16bits, por isso devemos
// dividir os pixels por 256(2ˆ8 ou melhor: (n >> 8)) para ficarmos no espaço desejado
func ToRGB(imgOriginal image.Image) *image.RGBA {
	dimensoes := imgOriginal.Bounds()
	largura := dimensoes.Dx()
	altura := dimensoes.Dy()

	imgEmRGB := image.NewRGBA(dimensoes)
	// percorrendo a imagem
	for y := 0; y < altura; y++ {
		for x := 0; x < largura; x++ {
			r, g, b, a := imgOriginal.At(x, y).RGBA()
			imgEmRGB.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a),
			})
		}
	}
	return imgEmRGB
}
