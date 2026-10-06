# Architecture Decisions

Este documento registra as decisões arquiteturais atuais do WagerFlow.

Ele descreve decisões já adotadas e separa essas escolhas de componentes que ainda serão implementados.

## Visão geral

O sistema será organizado em três áreas principais:

```text
+--------------------------------------+
|            Infrastructure            |
| HTTP | PostgreSQL | SQS | OIDC | Fx |
+------------------+-------------------+
                   |
                   v
+--------------------------------------+
|             Application              |
| Use Cases | Ports | Orchestration    |
+------------------+-------------------+
                   |
                   v
+--------------------------------------+
|                Domain                |
| Money | Wallet | Ledger | Wager      |
+--------------------------------------+
```

A direção das dependências é:

```text
Infrastructure -> Application -> Domain
```

O domínio não conhece infraestrutura.

## Composição com Uber Fx

Uber Fx é usado para composição de dependências e lifecycle da aplicação.

O entrypoint deve permanecer pequeno e delegar composição aos módulos.

Exemplo conceitual:

```go
fx.New(
    config.Module,
    bootstrap.Module,
).Run()
```

Os módulos concretos serão adicionados apenas quando suas responsabilidades existirem.

Módulos futuros previstos:

```text
Database
Auth
Repositories
Application
HTTP
SQS
Outbox
Observability
```

A criação antecipada de pacotes vazios é evitada.

## Lifecycle

`fx.Lifecycle` será usado para controlar inicialização e encerramento ordenado de recursos como:

```text
PostgreSQL pool
HTTP server
SQS consumer
Outbox publisher
Pending-reference worker
```

Lifecycle não substitui garantias transacionais. Um processo pode morrer sem executar `OnStop`, portanto consistência e recuperação devem depender de estado durável.

## Configuração

A configuração atual é carregada por variáveis de ambiente.

Em desenvolvimento local, `.env` é opcional.

Variáveis atuais:

```text
APP_NAME
APP_ENV
HTTP_PORT
```

O `.env` não é versionado; `.env.example` deve acompanhar novas configurações.

## Persistência

A preferência é PostgreSQL com `pgx` e SQL explícito.

Motivos:

- transações visíveis;
- locks visíveis;
- constraints visíveis;
- comportamento de concorrência auditável;
- menor ocultação das garantias financeiras.

Migrations serão versionadas.

## Atomicidade financeira

Uma operação financeira deverá persistir tudo que pertence à decisão em uma única transação.

Modelo conceitual:

```text
BEGIN

WagerTransaction
Wallet
WalletLedgerEntry
Inbox      (quando origem for SQS)
Outbox

COMMIT
```

Se qualquer etapa falhar:

```text
ROLLBACK
```

Nenhum evento externo deve ser publicado antes do commit da transação de origem.

## Transactional Outbox

Eventos externos serão persistidos em Outbox dentro da mesma transação financeira.

```text
Financial Transaction
        |
        |-- Wallet
        |-- Ledger
        `-- Outbox
              |
            COMMIT
              |
              v
       Outbox Publisher
              |
              v
             SQS
```

O publisher deverá suportar múltiplas instâncias, retries, backoff e recuperação de trabalho abandonado.

## Inbox

O consumidor SQS terá Inbox persistente.

Fluxo previsto:

```text
SQS Message
    |
    v
Inbox check
    |
    v
Application Use Case
    |
    v
COMMIT
    |
    v
Delete SQS Message
```

A mensagem só deve ser removida depois que o resultado estiver duravelmente persistido.

## Concorrência distribuída

Locks locais de processo não serão a garantia principal.

A coordenação será feita por Wallet no PostgreSQL.

Estratégia inicial:

```sql
BEGIN;

SELECT ...
FROM wallets
WHERE id = $1
FOR UPDATE;

-- processar operação

COMMIT;
```

Isso serializa movimentos conflitantes da mesma Wallet sem bloquear carteiras independentes.

O desenho deve permitir:

```text
Wallet A ---- processamento independente ---->
Wallet B ---- processamento independente ---->
```

sem lock global.

## Idempotência

Idempotência será persistente.

A aplicação deverá distinguir:

```text
mesma key + mesmo payload
-> replay do resultado persistido

mesma key + payload diferente
-> conflict
```

Replays não podem reaplicar débito ou crédito.

HTTP e SQS devem convergir para o mesmo caso de uso e para a mesma regra de idempotência.

## Ledger

O ledger será append-only.

Cada mudança de saldo deve ter exatamente uma entrada financeira correspondente.

A imutabilidade deve ser protegida tanto no domínio quanto no banco.

`LOSS` não cria ledger porque não altera saldo.

## Eventos após commit

A aplicação não publica diretamente um evento financeiro durante a transação de domínio.

Em vez disso:

```text
DB transaction
  |
  `-- insert outbox
        |
      COMMIT
        |
        v
publisher assíncrono
```

Essa decisão evita eventos representando mudanças que posteriormente sofreram rollback.

## CI

O workflow atual possui três jobs principais:

```text
Formatting ──────┐
                 ├── Build
Tests / Quality ─┘
```

`Formatting` e `Tests / Quality` executam em paralelo.

O build possui dependência explícita dos dois.

### Formatting

Verifica `gofmt` sem modificar arquivos.

### Tests / Quality

Executa:

```text
go mod tidy + diff
go vet ./...
go test ./... com coverage
go test -race ./...
```

O `coverage.out` é armazenado como artifact.

Não há threshold mínimo de coverage neste estágio. Coverage é utilizado como indicador, não como objetivo isolado.

### Build

Executa somente após os jobs anteriores passarem:

```bash
go build ./...
```

## Testes

A estratégia de testes será dividida em:

```text
Unit
Integration
Concurrency
Failure/Recovery
```

### Unit

Valida regras puras de domínio.

Atualmente inclui `Money`, `Wallet`, `Ledger` e configuração.

### Integration

Futuramente utilizará serviços reais containerizados:

```text
PostgreSQL
OIDC provider
LocalStack/MiniStack
```

### Concurrency

Cenários obrigatórios incluem duas apostas concorrentes de `80` sobre saldo `100`, múltiplas instâncias e carteiras independentes.

### Failure/Recovery

Serão simuladas falhas em pontos críticos, como commit concluído antes do delete da mensagem e publishers concorrentes do Outbox.

## Observabilidade

Planejada:

- logs JSON;
- `correlationId`;
- `messageId`;
- `transactionId`;
- `walletId`;
- `providerId`;
- métricas de processamento, retries, duplicidade, DLQ, conflitos, outbox lag e reconciliação.

Credenciais e payloads financeiros completos não devem ser registrados.

## Decisões ainda abertas

Ainda serão detalhadas durante a implementação:

- schema PostgreSQL final;
- migrations;
- modelo completo de `WagerTransaction`;
- contratos dos repositories;
- API HTTP;
- configuração OIDC/Keycloak;
- formato das mensagens SQS;
- schema de Inbox e Outbox;
- política de retry;
- reconciliação;
- observabilidade concreta.

Essas decisões serão documentadas à medida que forem implementadas, evitando documentação especulativa excessiva.
