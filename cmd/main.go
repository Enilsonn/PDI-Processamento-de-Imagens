package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/corr"
	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/hist"
	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/ioimg"
	"github.com/Enilsonn/PDI-Processamento-de-Imagens/internal/kernels"
)

func main() {
	// flags
	input := flag.String("input", "", "caminho d imagem de entrada (png ou jpeg)")
	output := flag.String("output", "saida.png", "caminho da imagem de saída")
	filtro := flag.String("filter", "", "caminho para txt do filtro")
	histograma := flag.String("histogram", "", "tipo do histograma: equalize, expand, local, equalize-local")
	localM := flag.Int("m", 3, "largura da janela local")
	localN := flag.Int("n", 3, "altura da janela local")
	flag.Parse()

	if *input == "" {
		log.Fatal("use -input para informar um caminho vádido da imagem")
	}

	// abrindo imagem
	img, err := ioimg.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	imgRGB := ioimg.ToRGB(img)

	// aplicando filtro (se houver)
	if *filtro != "" {
		kernel, err := kernels.LoadKernel(*filtro)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", kernel)

		imgRGB = corr.AplicarCorrelação(imgRGB, kernel)
	}

	// aplicando histograma (se houver)
	if *histograma != "" {
		switch *histograma {
		case "equalize":
			imgRGB = hist.Equalização(imgRGB)

		case "expand":
			imgRGB = hist.Expansão(imgRGB)

		case "local":
			imgRGB = hist.EqualizaçãoLocal(imgRGB, *localM, *localN)

		case "equalize-local":
			imgRGB = hist.Equalização(imgRGB)
			imgRGB = hist.EqualizaçãoLocal(imgRGB, *localM, *localN)

		default:
			log.Fatal("histograma inválido, use: -histogram qualize, expand, local ou equalize-local")
		}
	}
	if err := ioimg.Save(*output, imgRGB); err != nil {
		log.Fatalf("erro ao salvar imagem: %#v", err)
	}

	fmt.Printf("Processamento concluído e a imagem foi salva em: %v\n", *output)
}
