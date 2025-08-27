package kernels

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// esse modulo é respon'sabvel por ler os modulos txt e montar os filtros
// portanto segue a mesma estrutura de ioimg, porém para retornar o objeto do tipo Kernel

/*
	O txt estará estruturado da seguinte forma:
		nome: SobelHorizontal
		funçãoDeAtivaçao: identity
		bias: 0
		tamanho: 3x3
		máscara:
		-1 -2 -1
		0  0  0
		1  2  1
*/
// maáscara deve ser o ultimo, pois suas linhas não têm prefixos

type Kernel struct {
	Nome             string
	Largura          int
	Altura           int
	Máscara          [][]int
	Bias             int
	FunçãoDeAtivação string
	Escala           int
}

func LoadKernel(path string) (*Kernel, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar abrir %s: %v", path, err)
	}
	defer file.Close()

	// na especificação do projeto Escala tem um valor opcional, então vamos inicializar
	// o objeto com ele setado como default (1)
	kernel := &Kernel{Escala: 1}

	scanner := bufio.NewScanner(file)
	var flagMascara bool
	var linhasDaMascara []string

	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha == "" {
			continue
		}

		// vamos separar a leitura pelas linhas que tem cabeçalho e pela máscara

		//linhas com cabeçalho:
		if strings.HasPrefix(linha, "nome:") {
			nome := strings.TrimPrefix(linha, "nome:")
			kernel.Nome = strings.TrimSpace(nome)
			continue

		} else if strings.HasPrefix(linha, "funçãoDeAtivação") {
			função := strings.TrimPrefix(linha, "funçãoDeAtivação")
			função = strings.TrimSpace(função)
			kernel.FunçãoDeAtivação = strings.ToLower(função)
			continue

		} else if strings.HasPrefix(linha, "bias:") {
			fmt.Println("aqui bias")
			bias := strings.TrimPrefix(linha, "bias:")
			bias = strings.TrimSpace(bias)
			kernel.Bias, err = strconv.Atoi(bias)
			if err != nil {
				return nil, fmt.Errorf("formato do bias é inválido")
			}
			continue

		} else if strings.HasPrefix(linha, "tamanho:") {
			tamanho := strings.TrimPrefix(linha, "tamanho:")
			tamanho = strings.TrimSpace(tamanho)
			MN := strings.Split(tamanho, "x")
			if len(MN) != 2 {
				return nil, fmt.Errorf("formato do tamanho é inválido")
			}
			kernel.Altura, err = strconv.Atoi(MN[0])
			if err != nil {
				return nil, fmt.Errorf("formato da altura é inválido")
			}
			kernel.Largura, err = strconv.Atoi(MN[1])
			if err != nil {
				return nil, fmt.Errorf("formato da largula é inválido")
			}
			continue

		} else if strings.HasPrefix(linha, "escala:") {
			escala := strings.TrimPrefix(linha, "escala:")
			escala = strings.TrimSpace(escala)
			kernel.Escala, err = strconv.Atoi(escala)
			if err != nil {
				return nil, fmt.Errorf("formato da escala é inválido")
			}
			continue

		} else if strings.HasPrefix(linha, "máscara:") {
			flagMascara = true
			continue

		} else if flagMascara {
			linhasDaMascara = append(linhasDaMascara, linha)
			continue
		}
	}

	// se houve algum erro de leitura acima, ele é ignorado pelo scanner.Scan()
	// porém é armazenado em um buffer que verificaremos agora
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("o arquivo está mal formatado")
	}

	// linhas da máscara:
	for _, linha := range linhasDaMascara {
		valores := strings.Fields(linha)

		if len(valores) != kernel.Largura {
			return nil, fmt.Errorf("a largura da máscara não está de acordo com a largura informada")
		}

		var linhaInt []int
		for _, valor := range valores {
			num, err := strconv.Atoi(valor)
			if err != nil {
				return nil, fmt.Errorf("o valor %s da máscara é inválido: %v", valor, err)
			}
			linhaInt = append(linhaInt, num)
		}

		kernel.Máscara = append(kernel.Máscara, linhaInt)
	}

	if len(kernel.Máscara) != kernel.Altura {
		return nil, fmt.Errorf("a altura da máscara não está de acordo com a altura informada")
	}

	return kernel, nil

}
