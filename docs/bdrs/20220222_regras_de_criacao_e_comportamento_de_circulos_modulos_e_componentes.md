# Criação e comportamento de módulos e componentes

- **Status:** proposto
- **Decisores:** [barbararochazup](https://github.com/barbararochazup), [nathannascimentozup](https://github.com/author:nathannascimentozup), [tiagoangelozup](https://github.com/tiagoangelozup), 
- **Date:** 2022-02-213

## Contexto

**O que é um módulo?** Um módulo é um conjunto de componente gerenciados pelo CharlesCD. É a a representação de um repositório de chart do git ou helm.

**O que é um componente?** Um componente é um item deploiável descrito no charts referenciado no módulo. É a representação de uma aplicação.

**O que é um anel?** Um anel é uma versão de deploy de um componente com regras de segmentação especificas.

**O que é o mar aberto?** É o anel padrão, sem nenhuma regra de segmentação.
## Regras

- Um workspace pode ter mais de um módulo
- Um módulo possui somente a referência a um repositório de charts
- Componentes não são livremente cadastrados, disponibilidade de cadastro é de acodo com o que está mapeado no charts.
- Ao realizar o cadastro do módulo realiza deploy de todos os componentes em mar aberto e ativa o mar aberto caso esteja inativo
    - Essa opção pode ser desativada, caso seja, será aplicado todos os yamls menos os do tipo deployment, permanecendo mar aberto no status em que ele se encontra
- De um 1 em 1 min (valor padrão, mas pode ser configurável), um jobs irá validar se ocorreu alguma alteração nos charts, e se ocorreu irá atualizar as refs; fazendo inclusive atualização nos circulos que já deployados.



