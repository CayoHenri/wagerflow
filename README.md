# WagerFlow

Backend para processamento distribuído de apostas e movimentações financeiras, desenvolvido em Go.

O projeto implementa uma carteira financeira capaz de processar transações de apostas mantendo consistência financeira, idempotência persistente, concorrência entre múltiplas instâncias e recuperação segura em cenários de falha.

> Projeto desenvolvido para o desafio **Backend — Processamento Distribuído de Apostas em Go**.

## Status do projeto

🚧 **Em desenvolvimento**

Implementado até o momento:

- [x] Go Modules
- [x] Uber Fx
- [x] Lifecycle da aplicação
- [x] Configuração por variáveis de ambiente
- [x] Suporte a `.env` em desenvolvimento
- [x] Testes da configuração com Testify
- [x] CI com formatação, testes, coverage, race detector, vet e build
- [x] Value Object `Money`
- [x] Aggregate Root `Wallet`
- [x] `WalletLedgerEntry`
- [ ] `WagerTransaction`
- [ ] PostgreSQL
- [ ] Migrations
- [ ] Idempotência persistente
- [ ] API HTTP
- [ ] OAuth2/OIDC
- [ ] AWS SQS
- [ ] Inbox
- [ ] Transactional Outbox
- [ ] Referências pendentes
- [ ] Observabilidade
- [ ] Testes de integração
- [ ] Testes distribuídos de concorrência e recuperação
- [ ] Docker Compose

## Objetivo

O WagerFlow tem como objetivo processar movimentações financeiras relacionadas a apostas sem permitir inconsistências causadas por concorrência, reprocessamento ou falhas parciais.

Operações previstas:

- abertura de carteira;
- `BET`;
- `WIN`;
- `LOSS`;
- `REFUND`;
- `ROLLBACK`;
- consulta de ledger;
- reconciliação financeira.

HTTP e SQS deverão utilizar os mesmos casos de uso:

```text
HTTP -----------┐
                v
         Application Use Case
                |
        +-------+-------+
        |               |
      Wallet          Ledger
        |
        v
    PostgreSQL
                ^
                |
SQS ------------┘
```

## Arquitetura

A direção de dependências adotada é:

```text
Infrastructure -> Application -> Domain
```

O domínio não depende de Fx, HTTP, PostgreSQL, SQS, Keycloak ou Docker.

Documentação detalhada:

- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — arquitetura, concorrência, atomicidade, CI e decisões estruturais.
- [`docs/DOMAIN.md`](docs/DOMAIN.md) — decisões de `Money`, `Wallet` e `WalletLedgerEntry`.

## Estrutura atual

```text
wagerflow/
├── .github/
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── api/
├── docs/
│   ├── ARCHITECTURE.md
│   └── DOMAIN.md
├── internal/
│   ├── bootstrap/
│   ├── config/
│   └── domain/
│       ├── ledger/
│       ├── money/
│       └── wallet/
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

A estrutura será expandida conforme responsabilidades concretas surgirem.

## Stack

Stack atual e planejada:

- Go
- Go Modules
- Uber Fx
- Testify
- PostgreSQL
- pgx
- AWS SQS
- LocalStack/MiniStack
- OAuth2 / OpenID Connect
- Keycloak
- Docker
- Docker Compose

O acesso ao PostgreSQL será preferencialmente feito com `pgx` e SQL explícito, deixando transações, locks e garantias financeiras visíveis no código.

## Domínio implementado

### Money

`Money` utiliza `int64` em unidades monetárias menores, sem `float32`/`float64`.

Exemplo:

```text
"100.00" BRL -> 10000
```

O Value Object suporta parsing decimal, soma, subtração, negação, comparação, serialização e proteção contra overflow.

### Wallet

`Wallet` é o Aggregate Root responsável por saldo, moeda, versão e invariantes financeiras.

Regras atuais:

- saldo nunca pode ser negativo;
- crédito e débito exigem valor positivo;
- moeda da operação deve coincidir com a moeda da carteira;
- `version` incrementa somente quando o saldo muda;
- uma operação inválida não altera parcialmente o estado.

### Ledger

`WalletLedgerEntry` representa uma movimentação financeira imutável.

Movimentos válidos:

```text
DEBIT
CREDIT
```

O ledger valida a própria transição:

```text
CREDIT: before + amount == after
DEBIT:  before - amount == after
```

Detalhes completos em [`docs/DOMAIN.md`](docs/DOMAIN.md).

## Concorrência

A Wallet protege invariantes em memória, mas a concorrência distribuída será coordenada no PostgreSQL.

A estratégia inicial é lock pessimista por carteira:

```sql
BEGIN;

SELECT ...
FROM wallets
WHERE id = $1
FOR UPDATE;

-- aplicar regras do domínio
-- atualizar wallet
-- inserir ledger
-- inserir outbox

COMMIT;
```

Não será utilizado `sync.Mutex` como garantia financeira entre instâncias.

Cenário obrigatório:

```text
saldo inicial = 100.00 BRL

BET 80 ─┐
        ├─ concorrentes
BET 80 ─┘

resultado esperado:
- 1 PROCESSED
- 1 REJECTED (insufficient funds)
- saldo final = 20.00 BRL
- 1 débito no ledger
```

## Garantias financeiras

O projeto seguirá estas garantias:

- nenhum valor monetário em ponto flutuante;
- saldo persistido não negativo;
- idempotência persistente;
- ledger append-only;
- eventos externos somente após commit;
- atomicidade entre transação, wallet, ledger, inbox/outbox quando aplicável;
- coordenação distribuída por carteira;
- replays não reaplicam movimentos financeiros.

## CI

O workflow em `.github/workflows/ci.yml` executa:

```text
Formatting ──────┐
                 ├── Build
Tests / Quality ─┘
```

O estágio de testes inclui:

- `go mod tidy` + validação de diff;
- `go vet ./...`;
- `go test ./...`;
- coverage;
- `go test -race ./...`;
- upload de `coverage.out`.

O build só executa após formatação e testes passarem.

## Configuração

Variáveis atuais:

| Variável | Padrão | Descrição |
| --- | --- | --- |
| `APP_NAME` | `wagerflow` | Nome da aplicação |
| `APP_ENV` | `development` | Ambiente |
| `HTTP_PORT` | `8080` | Porta HTTP |

O `.env` é opcional para desenvolvimento local e não deve ser versionado.

## Desenvolvimento local

Validação recomendada antes de push:

```bash
gofmt -w .
go mod tidy
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Para executar a aplicação:

```bash
go run ./cmd/api
```

## Próximos passos

A próxima etapa de domínio será modelar `WagerTransaction`, começando por:

```text
Kind
State
estrutura básica
invariantes
```

Depois disso serão introduzidas persistência PostgreSQL, idempotência, API HTTP, SQS/Inbox, Outbox e autenticação.

## Documentação

- [Arquitetura](docs/ARCHITECTURE.md)
- [Decisões de domínio](docs/DOMAIN.md)
