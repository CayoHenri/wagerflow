# WagerFlow

Backend para processamento distribuído de apostas e movimentações financeiras, desenvolvido em Go.

O projeto implementa uma carteira financeira capaz de processar transações de apostas via HTTP e mensageria, mantendo consistência financeira, idempotência persistente, concorrência entre múltiplas instâncias e recuperação segura em cenários de falha.

> Projeto desenvolvido para o desafio **Backend — Processamento Distribuído de Apostas em Go**.

## Status do projeto

🚧 **Em desenvolvimento**

Implementado até o momento:

- [x] Inicialização do projeto com Go Modules
- [x] Composição da aplicação com Uber Fx
- [x] Lifecycle da aplicação
- [x] Configuração por variáveis de ambiente
- [x] Suporte a `.env` para desenvolvimento local
- [x] Testes da configuração com Testify
- [ ] Value Object `Money`
- [ ] Domínio de Wallet
- [ ] Domínio de Wager Transaction
- [ ] PostgreSQL
- [ ] Migrations
- [ ] Ledger financeiro
- [ ] Idempotência persistente
- [ ] API HTTP
- [ ] Autenticação OAuth2/OIDC
- [ ] Integração com SQS
- [ ] Inbox
- [ ] Outbox
- [ ] Processamento de referências pendentes
- [ ] Observabilidade
- [ ] Testes de integração
- [ ] Testes de concorrência e recuperação
- [ ] Docker Compose

## Objetivo

O WagerFlow tem como objetivo processar movimentações financeiras relacionadas a apostas garantindo que operações concorrentes ou repetidas não causem inconsistências no saldo das carteiras.

O sistema deverá suportar operações como:

- abertura de carteira;
- apostas (`BET`);
- ganhos (`WIN`);
- perdas (`LOSS`);
- reembolsos (`REFUND`);
- estornos (`ROLLBACK`);
- consulta do ledger;
- reconciliação financeira.

O mesmo processamento financeiro será utilizado independentemente da origem da requisição:

```text
                 ┌─────────────┐
                 │    HTTP     │
                 └──────┬──────┘
                        │
                        ▼
               ┌─────────────────┐
               │ Wager Use Case  │
               └────────┬────────┘
                        │
                ┌───────┴───────┐
                ▼               ▼
             Wallet           Ledger
                │
                ▼
            PostgreSQL

                 ▲
                 │
                 │
            ┌────┴────┐
            │   SQS   │
            └─────────┘
```

HTTP e SQS não possuirão implementações independentes das regras financeiras. Ambos utilizarão os mesmos casos de uso da aplicação.

## Stack

Tecnologias previstas para o projeto:

- **Go**
- **Go Modules**
- **Uber Fx**
- **PostgreSQL**
- **pgx**
- **AWS SQS**
- **LocalStack/MiniStack**
- **OAuth2 / OpenID Connect**
- **Keycloak**
- **Docker**
- **Docker Compose**
- **Testify**

A preferência para acesso ao PostgreSQL será `pgx` com SQL explícito, permitindo que transações, locks e garantias de concorrência permaneçam visíveis no código.

## Arquitetura

O projeto seguirá uma separação entre domínio, aplicação e detalhes de infraestrutura.

```text
┌──────────────────────────────────────┐
│           Infrastructure             │
│                                      │
│ HTTP • PostgreSQL • SQS • Keycloak   │
│ Fx • Logging • Metrics               │
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│             Application              │
│                                      │
│ Use Cases • Ports • Orchestration    │
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│                Domain                │
│                                      │
│ Money • Wallet • Wager • Ledger      │
│ Business Rules • Invariants          │
└──────────────────────────────────────┘
```

A direção das dependências será:

```text
Infrastructure → Application → Domain
```

O domínio não deverá depender de:

- Uber Fx;
- HTTP;
- PostgreSQL/pgx;
- AWS SQS;
- Keycloak;
- Docker;
- bibliotecas específicas de infraestrutura.

### Estrutura atual

```text
wagerflow/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── bootstrap/
│   │   ├── lifecycle.go
│   │   └── module.go
│   │
│   └── config/
│       ├── config.go
│       ├── config_test.go
│       └── module.go
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

A estrutura será expandida conforme os componentes forem implementados. Pacotes não serão criados antecipadamente sem uma responsabilidade concreta.

## Uber Fx

O Uber Fx é utilizado para composição e gerenciamento do ciclo de vida da aplicação.

O entrypoint permanece pequeno:

```go
func main() {
	fx.New(
		config.Module,
		bootstrap.Module,
	).Run()
}
```

Os módulos serão responsáveis por registrar suas próprias dependências.

Exemplo:

```go
var Module = fx.Module(
	"config",
	fx.Provide(
		Load,
	),
)
```

O lifecycle do Fx será utilizado posteriormente para controlar recursos como:

```text
Application
│
├── PostgreSQL Pool
├── HTTP Server
├── SQS Consumer
├── Outbox Publisher
└── Pending Reference Worker
```

Durante o shutdown, os componentes deverão ser encerrados de maneira ordenada sempre que possível.

## Configuração

A aplicação utiliza variáveis de ambiente para configuração.

Durante o desenvolvimento local, um arquivo `.env` pode ser utilizado.

Exemplo:

```env
APP_NAME=wagerflow
APP_ENV=development
HTTP_PORT=8080
```

O `.env` não deve ser versionado.

O arquivo `.env.example` será mantido no repositório com todas as configurações necessárias para execução da aplicação.

Variáveis atualmente disponíveis:

| Variável | Padrão | Descrição |
| --- | --- | --- |
| `APP_NAME` | `wagerflow` | Nome da aplicação |
| `APP_ENV` | `development` | Ambiente de execução |
| `HTTP_PORT` | `8080` | Porta HTTP |

Novas variáveis serão documentadas conforme PostgreSQL, SQS e OIDC forem implementados.

## Garantias financeiras

O sistema será desenvolvido considerando as seguintes garantias.

### Valores monetários sem ponto flutuante

Valores financeiros não utilizarão:

```go
float32
float64
```

O Value Object `Money` será responsável pela representação e manipulação dos valores monetários.

### Saldo não negativo

Uma operação de débito nunca poderá resultar em saldo negativo.

Essa garantia não dependerá apenas de validações em memória.

O PostgreSQL também participará da proteção das invariantes financeiras.

### Idempotência persistente

Requisições repetidas não poderão gerar movimentações financeiras duplicadas.

A idempotência será persistida no banco de dados e deverá funcionar mesmo entre múltiplas instâncias da aplicação.

### Ledger append-only

Cada alteração de saldo terá uma entrada correspondente no ledger.

Entradas financeiras já registradas não deverão ser alteradas.

### Atomicidade

Uma movimentação financeira deverá persistir atomicamente os dados relacionados à operação.

Conceitualmente:

```text
BEGIN

transaction
wallet
ledger
inbox (quando aplicável)
outbox

COMMIT
```

Caso alguma etapa falhe:

```text
ROLLBACK
```

### Eventos somente após commit

Eventos externos não serão publicados antes da confirmação da transação responsável pela alteração financeira.

Para isso será utilizado o padrão **Transactional Outbox**.

## Concorrência

O sistema deverá funcionar corretamente com múltiplas instâncias processando operações simultaneamente.

A coordenação financeira ocorrerá por carteira, evitando locks globais.

A estratégia definitiva será documentada em `ARCHITECTURE.md`.

Um dos cenários obrigatórios será:

```text
Saldo inicial: R$ 100,00

           Wallet
          R$ 100,00
          /       \
         /         \
        ▼           ▼
   BET 80,00    BET 80,00
        │           │
        ▼           ▼
   PROCESSED      REJECTED
                     │
              insufficient funds

Saldo final: R$ 20,00
```

O resultado esperado será:

- exatamente uma aposta processada;
- exatamente uma aposta rejeitada;
- saldo final de `20.00 BRL`;
- exatamente um débito no ledger;
- replays sem novas movimentações.

## Money

O Value Object `Money` será responsável por valores financeiros.

Requisitos previstos:

- precisão fixa de duas casas decimais;
- moeda ISO 4217;
- parsing de valores decimais em string;
- nenhuma utilização de ponto flutuante;
- soma;
- subtração;
- comparação;
- negação;
- serialização;
- detecção de overflow;
- rejeição de moedas incompatíveis;
- rejeição de notação científica;
- rejeição de `NaN`;
- rejeição de `Infinity`;
- rejeição de escala superior a duas casas decimais.

A implementação será baseada em unidades monetárias menores (`int64`), por exemplo:

```text
"100.00" BRL
       │
       ▼
10000 centavos
```

## Operações financeiras

O domínio deverá suportar:

| Operação | Movimento |
| --- | --- |
| `OPENING` | Crédito |
| `BET` | Débito |
| `WIN` | Crédito |
| `LOSS` | Nenhum |
| `REFUND` | Crédito |
| `ROLLBACK` | Movimento inverso |

`LOSS` deverá registrar a transação, mas não deverá alterar:

- saldo;
- ledger;
- versão da carteira.

`REFUND` e `ROLLBACK` utilizarão referências a transações previamente processadas e terão proteção contra reversões financeiras duplicadas.

## API HTTP

Endpoints previstos pelo desafio:

```text
POST /wallets
GET  /wallets/:walletId
GET  /wallets/:walletId/ledger

POST /wagering/transactions

GET /wagering/transactions/:transactionId

GET /providers/:providerId/wagering/transactions/:externalTransactionId

POST /wallets/:walletId/reconciliation

GET /health/live
GET /health/ready
```

As operações financeiras externas utilizarão:

```http
Idempotency-Key: <key>
```

O comportamento deverá distinguir:

```text
mesma chave + mesmo payload
        ↓
idempotent replay

mesma chave + payload diferente
        ↓
conflict
```

## Mensageria

A entrada assíncrona utilizará uma fila FIFO.

```text
wager-transactions.fifo
```

DLQ:

```text
wager-transactions-dlq.fifo
```

O processamento utilizará um padrão Inbox persistente.

```text
SQS
 │
 ▼
Inbox
 │
 ▼
Wager Use Case
 │
 ▼
COMMIT
 │
 ▼
Delete SQS Message
```

A mensagem somente será removida da fila depois que seu processamento estiver duravelmente persistido.

## Transactional Outbox

Eventos serão registrados no banco dentro da mesma transação da operação financeira.

```text
Financial Transaction
        │
        ├── Wallet
        ├── Ledger
        └── Outbox
              │
            COMMIT
              │
              ▼
       Outbox Publisher
              │
              ▼
             SQS
```

O publisher será independente e deverá suportar:

- múltiplas instâncias;
- retries;
- exponential backoff;
- recuperação de trabalhos abandonados;
- event ID estável durante republicações.

Eventos previstos:

- `WagerTransactionProcessed`
- `WagerTransactionRejected`
- `WalletBalanceChanged`
- `WagerTransactionPendingReference`

## Autenticação

A API utilizará autenticação externa através de OAuth2/OpenID Connect.

O ambiente local deverá disponibilizar automaticamente o provedor de identidade e identidades necessárias para execução dos testes.

A implementação prevista utilizará Keycloak.

A autorização deverá impedir acesso indevido a recursos pertencentes a outros providers.

## Observabilidade

A aplicação utilizará logs estruturados em JSON.

Campos relevantes poderão incluir:

```text
correlationId
messageId
transactionId
walletId
providerId
```

Credenciais e payloads financeiros completos não deverão ser registrados.

Métricas previstas incluem:

- transações por status;
- duplicidades;
- retries;
- mensagens enviadas para DLQ;
- conflitos de concorrência;
- latência;
- atraso do outbox;
- divergências de reconciliação.

## Testes

O projeto utilizará o package `testing` do Go e Testify.

### Unitários

Os testes unitários deverão cobrir principalmente:

- `Money`;
- regras de Wallet;
- transições de estado;
- regras das operações;
- idempotência;
- conflito entre payloads.

### Integração

Os testes de integração utilizarão dependências reais containerizadas quando exigido pelo desafio, incluindo:

- PostgreSQL;
- provedor OIDC;
- SQS local.

Serão testados:

- migrations;
- constraints;
- atomicidade;
- ledger;
- Inbox;
- Outbox;
- retries;
- DLQ;
- autenticação;
- isolamento entre providers;
- reinicialização de processos.

### Concorrência

Também serão criados testes específicos para:

- mesma aposta enviada dezenas de vezes;
- duas apostas concorrentes maiores que o saldo disponível;
- carteiras independentes sendo processadas em paralelo;
- múltiplas instâncias da aplicação;
- concorrência entre publishers do Outbox;
- falhas e reinicializações durante o processamento.

## Executando os testes

```bash
go test ./...
```

Com detector de race conditions:

```bash
go test -race ./...
```

Análise estática:

```bash
go vet ./...
```

Formatação:

```bash
gofmt -w .
```

Durante o desenvolvimento, utilizamos o fluxo:

```text
gofmt -w .
      ↓
go mod tidy
      ↓
go test ./...
      ↓
go test -race ./...
      ↓
go vet ./...
```

## Executando localmente

### Requisitos

No estado atual do projeto:

- Go instalado.

Conforme a infraestrutura for adicionada, também serão necessários:

- Docker;
- Docker Compose.

### Instalação

Clone o repositório:

```bash
git clone <repository-url>
```

Entre no diretório:

```bash
cd wagerflow
```

Instale/organize as dependências:

```bash
go mod tidy
```

Crie seu `.env`:

```bash
cp .env.example .env
```

No PowerShell:

```powershell
Copy-Item .env.example .env
```

Execute:

```bash
go run ./cmd/api
```

A aplicação utiliza atualmente:

```text
http://localhost:8080
```

> O servidor HTTP ainda está em desenvolvimento.

## Docker

🚧 **Em desenvolvimento**

A entrega final deverá permitir iniciar toda a infraestrutura através de:

```bash
docker compose up --build
```

O ambiente deverá incluir os serviços necessários para execução local e dos testes de integração.

## Migrations

🚧 **Em desenvolvimento**

O banco será criado através de migrations versionadas.

As migrations também serão responsáveis por constraints necessárias para garantir invariantes financeiras independentemente da aplicação.

## Documentação de arquitetura

Decisões arquiteturais e detalhes das estratégias de concorrência, idempotência e recuperação serão documentados separadamente em:

```text
ARCHITECTURE.md
```

Esse documento deverá explicar, entre outros pontos:

- limites arquiteturais;
- modelo transacional;
- estratégia de concorrência;
- prevenção de lost updates;
- idempotência;
- Inbox;
- Outbox;
- recuperação após falhas;
- processamento de referências pendentes;
- garantias fornecidas pelo PostgreSQL.

## Comandos da entrega

Ao final, os seguintes comandos deverão ser reproduzíveis:

```bash
docker compose up --build
```

```bash
go test ./...
```

```bash
go test -race ./...
```

```bash
go vet ./...
```

## Princípios do projeto

As decisões de implementação serão guiadas principalmente pelos seguintes princípios:

1. **Corretude financeira antes de performance.**
2. **Nenhum valor financeiro representado com ponto flutuante.**
3. **Garantias críticas persistidas no PostgreSQL.**
4. **Idempotência independente da instância da aplicação.**
5. **Eventos externos somente após commit.**
6. **Ledger financeiro auditável e imutável.**
7. **Concorrência coordenada por carteira, sem lock global.**
8. **HTTP e SQS compartilhando as mesmas regras de negócio.**
9. **Domínio independente de detalhes de infraestrutura.**
10. **Falhas e reinicializações tratadas como parte normal do sistema.**

---

## Licença

Este projeto foi desenvolvido como parte de um desafio técnico.