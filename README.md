# Perfect World Sorteios

O Perfect World Sorteios é uma ferramenta desenvolvida em Go (Golang) para facilitar a realização de sorteios em servidores de Perfect World. Projetado para automatizar o processo de distribuição de recompensas, este programa oferece flexibilidade e personalização para atender às necessidades específicas de cada servidor e comunidade.

Com suporte para sorteios de itens, moedas e gold, além da possibilidade de configurar múltiplos sorteios consecutivos, o Perfect World Sorteios torna o processo de premiação mais justo, transparente e eficiente.

Ideal para administradores de servidores de Perfect World que buscam uma solução confiável e fácil de usar, o Perfect World Sorteios simplifica a organização e execução de sorteios, proporcionando uma experiência melhor para os jogadores e facilitando a gestão de eventos e promoções.


## Recursos e Funcionalidades

### Sorteios de Itens, Moedas e Gold

O programa permite realizar sorteios de itens, moedas e gold em servidores de Perfect World, proporcionando uma maneira fácil e eficiente de distribuir recompensas para os jogadores.

### Automatização de Sorteios

O sistema pode ser automatizado para realizar sorteios em intervalos predefinidos, oferecendo conveniência e regularidade nas distribuições de prêmios.

### Configuração de Múltiplos Sorteios

O programa pode ser configurado para realizar múltiplos sorteios consecutivos, permitindo que você defina quantos sorteios deseja realizar em uma única execução.

### Configuração Personalizada

- **Definição de Level Mínimo**: Possibilidade de configurar um level mínimo para participação no sorteio, garantindo que apenas jogadores com um nível mínimo estabelecido possam concorrer.

- **Definição de Cultivo Mínimo**: Opção de estabelecer um cultivo mínimo necessário para participação no sorteio, assegurando que apenas jogadores com o cultivo mínimo exigido possam ser elegíveis.

Estas configurações personalizadas permitem adaptar o sorteio às necessidades específicas do servidor e dos jogadores, garantindo uma distribuição justa de prêmios.

## Compilação

### 1. Compile o código

```bash
go build -o sorteio
```
### 2. Dê permissão de execução
```bash
chmod +x sorteio
```

O código foi projetado para ser executado periodicamente via crontab. No entanto, ele também oferece a flexibilidade de ser utilizado de forma isolada via linha de comando. Isso significa que você pode executar o projeto tanto para uma única vez quanto programá-lo para execuções automáticas e periódicas, conforme a sua necessidade.


## Utilização / Execução

### 1. Execução via linha de comando
Utilizada caso queira fazer um sorteio isolado sem a necessidade de agendamento periódico. Para executar o programa via linha de comando, basta executa-lo diretamente:

```bash
./sorteio
```

### 2. Execução via Crontab
Utilizado para agendar sorteios periódicos. Ao configurar o programa para rodar via crontab, ele será executado automaticamente em intervalos predefinidos. 

#### 2.1. Abra o crontab com o comando:

```bash
    crontab -e
```

#### 2.2. Adicione o programa ao crontab:

```bash
    0 */6 * * * /caminho/do/executavel/sorteio
```
No exemplo acima o sorteio será executado uma vez a cada 6 horas.

## Créditos

Este projeto foi inspirado e utiliza conhecimentos de diversas fontes. Agradeço a todos os desenvolvedores e comunidades que compartilham conhecimento e ferramentas que possibilitaram a criação deste projeto. Dentre eles vale destacar:

- [pwdev.ru](http://pwdev.ru/) - Material utilizado para consulta de opcodes e estrutura dos pacotes.
- [hrace009](https://github.com/hrace009/perfect-world-api) - Material utilizado para estudo de geração e descompactação de pacotes.

## Contato

Para mais informações, sugestões ou contribuições, sinta-se à vontade para entrar em contato comigo através das minhas redes sociais:

- [**GitHub**](https://github.com/frankduque)
- [**LinkedIn**](https://www.linkedin.com/in/frankduque3/)
-  [**Whatsapp**](https://api.whatsapp.com/send?phone=5562992844985&text=Ol%C3%A1%20Frank%20como%20vai?)
- **E-mail**: franklr2229@gmail.com

Estou sempre aberto a feedbacks e colaborações para tornar este projeto ainda melhor!
