# Processamento Digital de Imagens em Go

Este projeto implementa técnicas de **Processamento Digital de Imagens (PDI)** em Go, incluindo filtros via convolução, equalização de histograma, expansão de contraste e equalização local. Ele permite aplicar filtros personalizados carregados de arquivos `.txt` e processar imagens diretamente pelo terminal.

---

## ⚡ Funcionalidades

- Aplicação de filtros via **convolução** (ex.: Sobel, Gaussiano).  
- **Equalização de histograma global** para melhorar contraste.  
- **Equalização de histograma local** para regiões específicas.  
- **Expansão de histograma** para aumentar contraste baseado na faixa real da imagem.  
- Suporte a filtros customizados via arquivos `.txt`.  

---

## 🖥️ Exemplo de linha de comando

Para processar uma imagem usando um filtro específico:

```bash
go run ./cmd -input ./Shapes.png -output ./outputs/Shapes_gaussiano5x5.png -filter ./filters/gaussiano5x5.txt
```

### Parâmetros:

- `-input` → Caminho da imagem de entrada.  
- `-output` → Caminho da imagem de saída.  
- `-filter` → Caminho do arquivo `.txt` do kernel/filtro.  
- Para equalização de histograma:
  - Global: `-histogram global`
  - Local: `-histogram local -m 7 -n 7`  
- Para expansão de histograma: `-histogram expansao`

---

## 🗂️ Estrutura do Projeto

```
.
├── Relatório-Projeto-1-PDI.pdf
├── Shapes.png
├── cmd
│   └── main.go
├── filters
│   ├── box10x1.txt
│   ├── box10x10.txt
│   ├── box1x10.txt
│   ├── gaussiano5x5.txt
│   ├── sobel_h.txt
│   └── sobel_v.txt
├── go.mod
├── go.sum
├── internal
│   ├── corr
│   │   └── corr.go
│   ├── hist
│   │   ├── equalização.go
│   │   ├── equalizaçãoLocal.go
│   │   └── expansão.go
│   ├── ioimg
│   │   └── ioimg.go
│   ├── kernels
│   │   └── kernels.go
│   └── utils
│       └── utils.go
├── outputs
│   ├── Shapes_box10x1.png
│   ├── Shapes_box10x10.png
│   ├── Shapes_box1x10.png
│   ├── Shapes_gaussiano5x5.png
│   ├── Shapes_sobel_h.png
│   ├── Shapes_sobel_v.png
│   ├── menina_box10x1.png
│   ├── menina_box10x10.png
│   ├── menina_box1x10.png
│   ├── menina_gaussiano5x5.png
│   ├── menina_sobel_h.png
│   └── menina_sobel_v.png
└── testpat.1k.color2.tif
```

---

## 📝 Formato dos filtros (`.txt`)

Um arquivo de filtro deve seguir o seguinte padrão:

```
nome: sobel_h
funçãoDeAtivação: identidade
bias: 0
tamanho: 3x3
escala: 1
máscara:
-1 -2 -1
0 0 0
1 2 1
```

- **nome:** Nome do filtro.  
- **funçãoDeAtivação:** Pode ser `identidade` ou `relu`.  
- **bias:** Valor adicionado após a correlação.  
- **tamanho:** Largura x Altura do kernel.  
- **escala:** Valor para normalizar o resultado da soma.  
- **máscara:** Matriz de pesos do filtro.

---

## 📈 Exemplos de uso

- Aplicando filtro Gaussiano 5x5:

```bash
go run ./cmd -input ./Shapes.png -output ./outputs/Shapes_gaussiano5x5.png -filter ./filters/gaussiano5x5.txt
```

- Equalização global:

```bash
go run ./cmd -input ./Shapes.png -output ./outputs/Shapes_equal_global.png -histogram global
```

- Equalização local 7x7:

```bash
go run ./cmd -input ./Shapes.png -output ./outputs/Shapes_equal_local.png -histogram local -m 3 -n 3
```

- Expansão de histograma:

```bash
go run ./cmd -input ./Shapes.png -output ./outputs/Shapes_exp.png -histogram expansao
```

---

## 🛠️ Requisitos

- Go 1.20+  
- Biblioteca padrão `image`  
- Não depende de bibliotecas externas para processar imagens.

---

## 🔧 Como contribuir

1. Fork este repositório.  
2. Crie sua branch: `git checkout -b feature/nome-da-feature`  
3. Faça suas alterações e commits.  
4. Abra um pull request explicando suas mudanças.  

---

## 📜 Licença

MIT License © Enilson Lima

