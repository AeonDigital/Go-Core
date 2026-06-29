# CLI_SPECS.md

## 1. Sintaxe e Regras de Escopo
* **Comandos:** Estrutura em árvore (`maincmd subcmd subsubcmd`). Cada nó é um escopo isolado.
* **Flags (Long):** Começam com `--` (ex: `--output`).
* **Flags (Short):** Começam com `-` e possuem até 3 caracteres (ex: `-out`). Sem aglutinação (ex: `-abc` é uma flag única chamada `abc`).
* **Atribuição:** Suporte exclusivo para `--flag valor` e `--flag=valor` (sem espaços ao redor do `=`).
* **Isolamento de Contexto:** Subcomandos NÃO herdam flags dos comandos pais. Cada comando ou subcomando processa estritamente suas próprias flags registradas.
* **Argumentos Posicionais:** Proibidos. Todo dado passado na linha de comando deve estar associado a uma flag (ex: `--file=relatorio.pdf`). Qualquer argumento solto ou sem flag correspondente resultará em erro imediato.

## 2. Interface de Ajuda (Help)
Todo comando ou subcomando gera uma interface de ajuda automática ao receber os argumentos `help`, `--help` ou `-h`.
* Se a descrição longa (`LongDescription`) não for expressamente definida no registro, o gerador de help usará a descrição curta (`ShortDescription`) no lugar.

## 3. Especificações de Registro

### 3.1 Registro de Comandos e Subcomandos
Para registrar um comando ou subcomando no roteador, as seguintes informações são obrigatórias:
* **Nome:** O termo disparador na linha de comando (ex: `build`).
* **ShortDescription:** Descrição breve usada em listagens gerais de ajuda.
* **LongDescription:** Descrição detalhada exibida na ajuda específica do comando (opcional, assume a curta se vazia).
* **Flags:** Lista de flags específicas e exclusivas que este comando aceita.
* **Subcommands:** Lista de comandos filhos que pertencem a este nó da árvore.

### 3.2 Registro de Flags
Cada flag registrada em um comando deve conter:
* **Long Name:** O nome completo da flag (sem os traços no registro, ex: `output`).
* **Short Name:** O nome reduzido de até 3 caracteres (sem o traço no registro, ex: `out`).
* **ShortDescription:** Descrição breve para o help.
* **LongDescription:** Descrição detalhada (opcional, assume a curta se vazia).
* **Required:** Booleano indicando se a presença da flag é obrigatória.
* **Default Value:** Valor padrão utilizado caso a flag não seja informada (ignorado se `Required` for true).
* **Type:** O tipo do dado para conversão e validação automática.

## 4. Sistema de Tipos e Validações Avançadas

### 4.1 Tipos Suportados
* **Primitivos:** `string`, `int`, `float`, `bool` (de ativação automática/sem valor expresso obrigatório).
* **Estruturados:** `json`, `duration`, `date`, `time`, `datetime`, `email`.
* **Sistema/Rede:** `path`, `filename`, `filepath`, `dirname`, `dirpath`, `url`, `fullurl`, `relativeurl`.
* **Arrays/Fatias:** Qualquer um dos tipos acima pode ser declarado como um array (ex: `[]int`, `[]string`) utilizando a sintaxe de array do JSON (ex: `--ids=""` ou `--nomes='["A","B"]'`).

### 4.2 Regras de Validação por Tipo

#### Tipos Quantificáveis (`int`, `float`, `duration`, `date`, `time`, `datetime`)
* **min:** Valor mínimo admissível. Se não definido, não há limite inferior.
* **max:** Valor máximo admissível. Se não definido, não há limite superior.

#### Tipo Texto (`string`)
* **min:** Número mínimo de caracteres que a string deve conter.
* **max:** Número máximo de caracteres que a string pode conter.
* **regex:** Expressão regular que o valor bruto da string deve passar para ser considerado correto.

#### Tipos de Array (`[]tipo`)
* **minlength:** Quantidade mínima de itens que a lista/fatia deve possuir.
* **maxlength:** Quantidade máxima de itens que a lista/fatia pode possuir.

#### Tipos de Sistema de Arquivos (`filepath`, `dirpath`)
As diretivas abaixo controlam de forma estrita o comportamento do validador em relação ao disco do sistema operacional:
* **must_exist:** Booleano que define se o recurso já deve obrigatoriamente existir em disco. 
  * Se `true`, a validação falha caso o caminho não seja encontrado.
  * Se `false`, o roteador valida apenas a sintaxe e as permissões do diretório pai.
* **access:** Define o nível de permissão exigido para o recurso. Aceita os valores:
  * `read`: Garante que a aplicação consegue abrir e ler o arquivo ou diretório.
  * `write`: Garante que a aplicação consegue modificar o recurso (se `must_exist: true`) ou que possui permissão de escrita no diretório pai para criá-lo (se `must_exist: false`).
  * `readwrite`: Exige que ambas as capacidades de leitura e modificação/criação estejam disponíveis.





==========================================================================


Entrar na lógica de validação é a parte mais crítica, pois é aqui que garantimos que o código do desenvolvedor só receberá dados perfeitamente limpos e higienizados.
Como o seu roteador roda em duas fases, a validação acontece estritamente na Fase 2 (após sabermos qual comando executar e termos isolado suas flags). Para construir isso de forma profissional e reutilizável, o validador deve disparar erros claros para o usuário final do CLI.
Vamos detalhar a estratégia e os critérios para cada grupo de tipos que você definiu:
------------------------------
## 1. Sistema e Rede (O diferencial da sua CLI)
Para validar esses tipos sem criar dependências externas, usaremos os pacotes nativos do Go (os, path/filepath e net/url).

* path / dirname / filename (Validação Semântica/Texto):
* Aqui validamos apenas a sintaxe do nome, sem checar se ele existe no disco.
   * Estratégia: Usamos filepath.Clean(). Para filename, verificamos se o resultado não contém barras (/ ou \), garantindo que é apenas o nome de um arquivo isolado. Para dirname, garantimos que ele segue a estrutura de pastas permitida pelo sistema operacional.
* filepath / dirpath (Validação de Existência no Sistema):
* Aqui nós exigimos que o caminho realmente exista no computador de quem está rodando a CLI.
   * Estratégia: Usamos os.Stat(valor).
   * Se retornar um erro onde os.IsNotExist(err) é verdadeiro, a validação falha (o arquivo/pasta não existe).
      * Para dirpath, além de existir, checamos se info.IsDir() é verdadeiro. Se for um arquivo físico, falha.
      * Para filepath, garantimos que !info.IsDir() (é um arquivo, não uma pasta).
   * url / fullurl / relativeurl:
* Estratégia: Usamos url.Parse(valor). Se o parser do Go falhar, a URL é inválida.
   * Para fullurl: Exigimos que os campos Scheme (http/https) e Host estejam preenchidos.
   * Para relativeurl: Exigimos que o Scheme e o Host estejam vazios (ex: /api/v1/users).
   * Para url: Aceita qualquer uma das duas opções acima, desde que seja interpretável.

------------------------------
## 2. Dados Estruturados e Temporais

* email:
* Estratégia: Como não queremos regex gigantes e complexas que falham, a recomendação profissional é usar o validador nativo do Go presente no pacote net/mail. O método mail.ParseAddress(valor) valida e-mails seguindo estritamente a RFC 5322.
* json:
* Estratégia: Usamos json.Valid([]byte(valor)). Se o formato do JSON estiver quebrado, rejeita imediatamente.
* date / time / datetime / duration:
* Estratégia: Usamos time.Parse(layout, valor) ou time.ParseDuration(valor).
   * Definiremos layouts padrão para a sua CLI (ex: 2006-01-02 para data, 15:04:05 para hora, e 2006-01-02 15:04:05 para datetime). Se a string não encaixar no layout, o validador falha e avisa o formato correto esperado.

------------------------------
## 3. A Lógica dos Limites (min, max, regex)
O maior segredo aqui é a ordem de execução: a validação de limites só acontece depois que o validador de tipo deu sinal verde e converteu a string bruta para o tipo correto.

* Se for string: O validador conta os caracteres usando utf8.RuneCountInString(valor) (para apoiar caracteres especiais corretamente). Em seguida, compara com MinLength e MaxLength. Por fim, se houver uma string no campo Regex, roda o regexp.MustCompile(f.Regex).MatchString(valor).
* Se for quantificável (int, float, duration, date): O validador faz uma conversão rápida e checa se o valor está matematicamente entre Min e Max.
* Se for um Array ([]tipo):
1. O roteador quebra o JSON usando json.Unmarshal.
   2. Conta quantos itens o array possui.
   3. Valida contra MinItems e MaxItems.
   4. O laço interno: Roda uma repetição (for) e passa cada item individualmente pelas regras de tipo descritas acima (ex: se for []email, testa cada item na regra de e-mail).

------------------------------
## Como quer prosseguir?
Essa estratégia blinda completamente a sua CLI contra dados corrompidos. O que achou da abordagem de validação física para filepath/dirpath e semântica para os outros?
Para continuarmos desenhando nossa especificação:

* Fazem sentido os formatos padrão que propus para as datas e horas, ou sua empresa costuma usar padrões diferentes (como o formato brasileiro DD/MM/AAAA)?
* Se estiver confortável com as validações, podemos passar para o design do Help Automático, definindo como o terminal vai exibir visualmente os subcomandos e as flags disponíveis!


