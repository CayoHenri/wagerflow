# Domain Design Decisions

Este documento registra as principais decisões de domínio adotadas no WagerFlow até o momento.

O objetivo é manter o `README.md` enxuto e concentrar aqui as regras, invariantes e responsabilidades dos objetos financeiros já implementados.

## Princípios do domínio

O domínio deve permanecer independente de detalhes de infraestrutura.

Pacotes em `internal/domain` não devem depender de:

- Uber Fx;
- HTTP;
- PostgreSQL/pgx;
- AWS SQS;
- Keycloak;
- Docker;
- bibliotecas específicas de infraestrutura.

A direção de dependências adotada é:

```text
Infrastructure -> Application -> Domain
```

O domínio contém regras de negócio e invariantes. A camada de aplicação será responsável por orquestrar múltiplos objetos de domínio e coordenar persistência transacional.

## Money

Pacote:

```text
internal/domain/money
```

`Money` é um Value Object imutável para representação de valores monetários.

### Representação

Valores monetários não utilizam `float32` ou `float64`.

O valor é armazenado em unidades monetárias menores usando `int64`:

```text
10.25 BRL -> 1025
```

Estrutura conceitual:

```go
type Money struct {
    amount   int64
    currency Currency
}
```

### Criação

Há duas formas principais:

```go
money.Parse("10.25", money.BRL)
money.NewFromMinorUnits(1025, money.BRL)
```

`Parse` é destinado a valores externos em formato decimal.

`NewFromMinorUnits` é destinado a valores internos já convertidos para unidades menores.

A separação evita ambiguidades sobre o significado de um inteiro como `1025`.

### Parsing externo

`Parse` aceita valores decimais com no máximo duas casas:

```text
"10"    -> 10.00
"10.2"  -> 10.20
"10.25" -> 10.25
"0.01"  -> 0.01
```

São rejeitados, entre outros:

- valor vazio;
- valor negativo;
- notação científica;
- mais de duas casas decimais;
- vírgula decimal;
- `NaN`;
- `Infinity`;
- formato não numérico.

Valores financeiros externos negativos são rejeitados.

### Valores negativos internos

`Money` pode representar valores negativos internamente.

Essa decisão permite operações matemáticas como `Negate` e `Subtract`.

```text
Money(-10.00 BRL)
└── matematicamente válido

Wallet(balance=-10.00 BRL)
└── financeiramente inválido
```

Essa separação mantém responsabilidades claras:

```text
Money  -> regras matemáticas
Wallet -> regras financeiras da carteira
```

### Operações

Atualmente `Money` fornece:

- `Add`;
- `Subtract`;
- `Negate`;
- `Compare`;
- `Equal`;
- `IsZero`;
- `IsPositive`;
- `IsNegative`;
- `String`.

Operações entre moedas diferentes são rejeitadas.

Também existem proteções contra overflow e underflow em operações com `int64`.

### Currency

A moeda é representada por um Value Object próprio.

Atualmente o domínio reconhece o subconjunto de moedas necessário para o desenvolvimento inicial.

A expansão para ISO 4217 deve preservar a mesma regra: moedas inválidas ou incompatíveis não podem participar de operações financeiras.

## Wallet

Pacote:

```text
internal/domain/wallet
```

`Wallet` é um Aggregate Root responsável por proteger as invariantes do saldo de uma carteira.

Estrutura conceitual:

```go
type Wallet struct {
    id        string
    playerID  string
    currency  money.Currency
    balance   money.Money
    version   int64
    createdAt time.Time
    updatedAt time.Time
}
```

Os campos são privados para impedir mutações externas que contornem as regras do aggregate.

### Invariantes

Uma Wallet válida deve garantir:

```text
id obrigatório
playerId obrigatório
currency válida
balance válido
balance.currency == wallet.currency
balance >= 0
version >= 1
timestamps válidos
updatedAt >= createdAt
```

### New e Restore

Foram separados dois fluxos de construção.

`New` cria uma carteira nova:

```text
version = 1
createdAt = now
updatedAt = now
```

`Restore` reconstrói uma Wallet já persistida e preserva:

```text
balance
version
createdAt
updatedAt
```

Essa separação evita tratar estado vindo do banco como uma nova carteira.

### Credit

`Credit` exige:

```text
amount válido
amount > 0
mesma currency da Wallet
operação sem overflow
```

Quando bem-sucedido:

```text
balance aumenta
version incrementa
updatedAt é atualizado
```

### Debit

`Debit` exige:

```text
amount válido
amount > 0
mesma currency da Wallet
saldo suficiente
```

O cálculo pode produzir temporariamente um `Money` negativo, mas a Wallet rejeita o resultado com `ErrInsufficientFunds`.

### Falhas não alteram estado

As operações calculam e validam o novo estado antes da mutação.

Se `Credit` ou `Debit` falhar:

```text
balance não muda
version não muda
updatedAt não muda
```

### Version

A versão é incrementada somente quando o saldo muda.

Uma futura operação `LOSS`, que não altera saldo, não deve incrementar `version`.

### Concorrência

A Wallet protege invariantes dentro de uma operação em memória, mas não resolve concorrência distribuída.

Não será utilizado `sync.Mutex` como proteção financeira principal porque múltiplas instâncias não compartilham memória.

A estratégia inicial é coordenar operações por Wallet no PostgreSQL com lock pessimista:

```sql
BEGIN;

SELECT ...
FROM wallets
WHERE id = $1
FOR UPDATE;

-- aplicar domínio
-- atualizar wallet
-- inserir ledger
-- inserir outbox

COMMIT;
```

A estratégia será validada por testes concorrentes envolvendo múltiplos processos.

## Wallet Ledger

Pacote:

```text
internal/domain/ledger
```

`WalletLedgerEntry` representa uma movimentação financeira auditável.

Estrutura conceitual:

```go
type WalletLedgerEntry struct {
    id            string
    walletID      string
    transactionID string
    direction     Direction
    amount        money.Money
    balanceBefore money.Money
    balanceAfter  money.Money
    createdAt     time.Time
}
```

### Append-only

O ledger é modelado como imutável.

Não existem setters, métodos de atualização ou `updatedAt`.

Essa regra será posteriormente reforçada no PostgreSQL.

### Direction

Movimentos válidos:

```text
DEBIT
CREDIT
```

O valor de `amount` permanece positivo. A direção representa o sentido financeiro.

### Invariantes

Uma entrada válida exige:

```text
id obrigatório
walletId obrigatório
transactionId obrigatório
direction válida
amount > 0
balances válidos e não negativos
mesma currency em amount/before/after
timestamp válido
transição matemática consistente
```

### Validação da transição

Para `CREDIT`:

```text
balanceBefore + amount == balanceAfter
```

Para `DEBIT`:

```text
balanceBefore - amount == balanceAfter
```

Exemplo válido:

```text
before = 100.00
DEBIT  = 20.00
after  = 80.00
```

Um histórico matematicamente inconsistente é rejeitado.

### LOSS não gera ledger

Uma futura transação `LOSS` terá valor financeiro zero e não mudará saldo.

Ela deverá ser registrada como transação processada, mas não deve criar uma entrada fictícia como:

```text
DEBIT 0.00
```

Por isso `WalletLedgerEntry` exige `amount > 0`.

## Coordenação entre Wallet e Ledger

A Wallet não cria diretamente um `WalletLedgerEntry`.

Essa decisão evita acoplá-la a detalhes como `transactionID` e `ledgerID`.

A futura camada de aplicação fará a coordenação:

```text
Application Use Case
        |
        |-- captura balanceBefore
        |-- Wallet.Debit/Credit
        |-- captura balanceAfter
        |-- cria WalletLedgerEntry
        `-- persiste tudo atomicamente
```

Responsabilidades:

```text
Money       -> matemática monetária
Wallet      -> saldo e invariantes
Ledger      -> auditoria
Application -> orquestração
PostgreSQL  -> atomicidade e concorrência distribuída
```

## Próximas decisões de domínio

O próximo agregado a ser modelado será `WagerTransaction`.

A primeira etapa deve introduzir gradualmente:

```text
Kind
State
estrutura básica
provider/player/wallet
round/game
Money
idempotency key
payload hash
referências financeiras
```

Os tipos e transições serão definidos antes da integração com PostgreSQL, HTTP ou SQS.
