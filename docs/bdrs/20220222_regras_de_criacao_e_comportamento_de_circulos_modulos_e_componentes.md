# Criação e comportamento de módulos e componentes

- **Status:** aceito
- **Decisores:** [barbararochazup](https://github.com/barbararochazup), [nathannascimentozup](https://github.com/author:nathannascimentozup), [tiagoangelozup](https://github.com/tiagoangelozup), 
- **Date:** 2022-02-213

## Contexto

**O que é um módulo?** Um módulo é um conjunto de componente gerenciados pelo CharlesCD. É a representação de um repositório, sendo ele do Git ou Helm

**O que é um componente?** Um componente é a representação de uma ou mais aplicações gerenciados pelo CharlesCD..

**O que é um anel?** Um anel é uma versão de deploy de um componente com regras de segmentação especificas.

**O que é o mar aberto?** É o anel padrão, sem nenhuma regra de segmentação.

## Regras
- Um módulo possui somente a referência a um repositório 
- Componentes não são livremente cadastrados, disponibilidade de cadastro é de acordo com o que está mapeado no repositório informado no módulo.
- Ao realizar o cadastro do módulo é realizado deploy de todos os componentes em mar aberto; o status de mar aberto também passo a ser ativo caso esteja inativo
    - Essa opção pode ser desativada, caso seja, será aplicado todos os yamls menos os do tipo deployment, permanecendo mar aberto no status em que ele se encontra
- De um 1 em 1 min (valor padrão, mas pode ser configurável), um jobs irá validar se ocorreu alguma alteração no repositório, e se ocorreu irá atualizar as refs; fazendo inclusive atualização nos circulos já deployados.


