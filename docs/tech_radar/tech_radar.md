# Operator Tech Radar  


## Golang
**O que é:** Linguagem de programação criada pela Google e lançada em código livre em novembro de 2009. É uma linguagem compilada e focada em produtividade e programação concorrente.
**Como usamos:** Utilizamos para desenvolvimento do Operator. Para entender melhor o que porque do uso da linguagem veja a ADR sobre a decisão de uso [aqui](https://github.com/ZupIT/charlescd-adrs/blob/main/pt-br/20211213-usar-golang-como-linguagem-do-projeto.md)
**Tipo:** Linguagem
**Referência:** 
[go.dev](https://go.dev/)


## Docker
**O que é:** Plataforma aberta que permite através da containerização facilitar o desenvolvimento, implantação e execuçõa de aplicações em ambientes isolados.
**Como usamos:** Utilizamos para containerização do Operator a fim de facilitar sua implantação e desenvolvimento.
**Tipo:** Plataforma
**Referência:** 
[docker.com](https://docker.com/)


## Kubernetes
**O que é:** Plataforma aberta que permite facilita a ajuda na orquestração de container, como o docker.

**Como usamos:** É uma das tecnologias core do CharlesCD. Extendemos o comportamento do Kubernetes por meio de CRDs a fim de entregar Hypothesis Driven Development.
**Tipo:** Plataforma
**Referência:** 
[kubernetes.io](https://kubernetes.io/pt-br/)


## FluxCD Source Controller
**O que é:** Operador do Kubernetes especializado em buscar artefatos a serem usados pelo Kubernetes de origens externas como o Git, Helm e S3.
**Como usamos:** Utilizamos para buscar os artefatos armazenados no git ou helm a serem usados como manifestos pelo Kubernetes junto ao CharlesCD e disponibilizar dentro da rede do Kubernetes. O FluxCD já permite suporte ao GitOps.

**Tipo:** Ferramenta
**Referência:** 
[fluxcd/source-controller](https://github.com/fluxcd/source-controller)


## Manifestival
**O que é:** Biblioteca que permite e facilita a manipulação de manifestos de recursos do Kubernetes.
**Como usamos:** Utilizamos para ler, manipula e aplicar utilizados pelo CharlesCD.
**Tipo:** Biblioteca.
**Referência:** 
[manifestival](https://github.com/manifestival/manifestival)


## Helm
**O que é:** Gerenciador de pacotes Kubernetes que auxilia a instalar e gerenciar o ciclo de vida de aplicações Kubernetes. 
**Como usamos:** Utilizamos para instalar e gerenciar aplicações que utilizem o Helm como ferramenta.
**Tipo:** Ferramenta.
**Referência:** 
[helm.sh](https://helm.sh/)

## Kustomize
**O que é:** Gerenciador de pacotes Kubernetes que auxilia a instalar e gerenciar o ciclo de vida de aplicações Kubernetes.
**Como usamos:** Utilizamos para instalar e gerenciar aplicações que utilizem o Kustomize como ferramenta.
**Tipo:** Ferramenta.
**Referência:** 
[kustomize.io](https://kustomize.io/)

## Google/Wire
**O que é:** Biblioteca que possibilita injeção de dependência em tempo de compilação para código golang.
**Como usamos:** Utilizamos para realizar a injeção de depedências dentro do projeto Operator, facilitando a manutenção do código.
**Tipo:** Biblioteca.
**Referência:** 
[google/wire](https://github.com/google/wire)


## Ginkgo
**O que é:**  Framework de testes para Golang.
**Como usamos:** TODO
**Tipo:** Framework.
**Referência:** 
[onsi/ginkgo](https://github.com/onsi/ginkgo)

## Gomega
**O que é:** Biblioteca de apoio ao Ginkgo, que permite uso assertions no teste.
**Como usamos:** TODO
**Tipo:** Biblioteca.
**Referência:** 
[onsi/gomega](https://onsi.github.io/gomega/)


## GoGetter
**O que é:** Biblioteca que possibilita o download de arquivos e diretórios por meio de uma url. 
**Como usamos:** Utilizamos para realizar a injeção de depedências dentro do projeto Operator, facilitando a manutenção do código.
**Tipo:** Biblioteca.
**Referência:** 
[hashicorp/go-getter](https://github.com/hashicorp/go-getter)


## Makefile
**O que é:** Arquivo que contem conjunto de diretivas para automatizar um conjunto de processos (instalar, desinstalar, remover arquivos, etc)
**Como usamos:** Utilizamos para gerar uma série de comandos utilitários para manutenção e uso do projeto.
**Tipo:** Arquivo.
**Referência:** 
[gnu.org/makefiles](https://www.gnu.org/software/make/manual/make.html#Introduction)