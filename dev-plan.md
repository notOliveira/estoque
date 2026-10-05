# Dev Plan — API de Estoque (Go + MongoDB)

> Documento vivo. Atualizado a cada sessão. O foco é **aprendizado**, não entrega rápida.

---

## 1. Filosofia de trabalho

- **Você programa, eu guio.** Atuo como professor/mentor.
- **Antes de cada mudança:** explico conceito, aponto armadilhas, mostro o "porquê". Só depois você codifica.
- **Depois de cada mudança:** reviso linha a linha ou apontando diff, explico conceitos novos que surgiram.
- **Nada de "colar resposta".** Quando travar, eu releio, aponto o problema e oriento o próximo passo.
- **Comentários em código:** PT-BR, mantendo o estilo didático atual.
- **Escopo atual:** CRUD completo da Fase 7. Depois disso, planejamos juntos.

---

## 2. Stack fixada

| Componente | Versão / Pacote | Por que |
|---|---|---|
| **Go** | 1.22+ (roda em 1.26.1) | `net/http` ganhou pattern matching nativo |
| **MongoDB** | Community Edition | Banco principal; já configurado localmente |
| **Driver Mongo** | `go.mongodb.org/mongo-driver/v2` v2.8.1 | v2 consolidou tipos (`bson.ObjectID`, sem precisar de `/primitive`) |
| **Env loader** | `github.com/joho/godotenv` v1.5.1 | Carrega `.env` em dev |
| **HTTP router** | `net/http` (stdlib) | Sem framework; padrão 1.22+ cobre o necessário |
| **Módulo** | `github.com/notoliveira/estoque` | Definido no `go.mod` |

---

## 3. Estrutura de pastas

```
estoque/
├── cmd/api/main.go               # entry point — wirar dependências
├── internal/
│   ├── config/   config.go       # carrega .env
│   ├── db/       mongo.go        # Connect / GetCollection
│   ├── model/    product.go      # DTOs (Create/Update/Response/Error)
│   ├── repository/ product.go    # CRUD no Mongo (interface + impl)
│   ├── service/  product.go      # regras de negócio + validação
│   └── handler/  product.go      # HTTP handlers (net/http)
├── domain/product.go             # struct Product pura
├── .env / .env.example / .gitignore
└── go.mod / go.sum
```

Princípio: dependência sempre aponta pra "dentro". `handler → service → repository → domain`. Camadas não se pulam.

---

## 4. Estado atual

| Fase | Escopo | Status | Observação |
|---|---|---|---|
| **0** | Limpeza do driver (remover v1) | ✅ FEITO | só v2 no `go.mod` |
| **1** | Consertar `internal/db/mongo.go` | ✅ FEITO | defer removido, panic → error |
| **2** | DTOs em `internal/model/product.go` | ✅ FEITO | `ToDomain` método + `ToResponse` função |
| **2A** | Fix de bugs (driver v2) | ✅ FEITO | `bson.ObjectID` em vez de `primitive.ObjectID` |
| **3** | `internal/repository/product.go` | 🔶 PARCIAL | `Create`, `FindByID`, `FindAll` OK. Faltam `Update` (ReplaceOne) e `Delete` (DeleteOne) |
| **4** | `internal/service/product.go` | ⏳ PENDENTE | arquivo só tem `package product`; precisa virar `package service` |
| **5** | `internal/handler/product.go` | ⏳ PENDENTE | arquivo só tem `package handler` |
| **6** | `cmd/api/main.go` wirar tudo | ⏳ PENDENTE | só `func main() {}` |
| **7** | `gofmt` + `go vet` + `go build` + curl | ⏳ PENDENTE | 3 arquivos pendem formatação hoje |

### 4.1 Validação feita nesta sessão
- ✅ `go build ./...` limpo
- ✅ `go vet ./...` limpo
- ⚠️ `gofmt -l .` aponta: `domain/product.go`, `internal/db/mongo.go`, `internal/model/product.go` (resolve na Fase 7)

### 4.2 Pendências menores (cosméticas)
- `internal/db/mongo.go` linha 25: `return client, err` → `return client, nil` (valor stale, não quebra)
- `gofmt -w .` na Fase 7 resolve o alinhamento dos 3 arquivos

---

## 5. Decisões já tomadas (não voltar atrás)

- **Driver v2 apenas.** Não usar `primitive.*` (subpacote foi consolidado em `bson`).
- **PUT = full replace**, sem ponteiros no Update DTO. PATCH parcial foi descartado nesta sessão.
- **`ProductRepository` interface definida no pacote `repository`.** Simples e funcional.
- **`ToDomain` é método** (nos DTOs); **`ToResponse` é função** no pacote `model`. Aceita-se a assimetria — Go não permite método em tipo de outro pacote.
- **Model não importa driver.** Quem fala de `bson.ObjectID` é o repository.
- **Service é dono de ID e timestamps.** DTOs não setam `ID`, `CreatedAt`, `UpdatedAt`.

---

## 6. Endpoints do CRUD (contrato HTTP)

| Método | Rota | Body entrada | Resposta | Códigos |
|---|---|---|---|---|
| `POST` | `/products` | `CreateProductRequest` | `ProductResponse` | 201, 400 |
| `GET` | `/products` | — | `[]ProductResponse` | 200 |
| `GET` | `/products/{id}` | — | `ProductResponse` | 200, 400, 404 |
| `PUT` | `/products/{id}` | `UpdateProductRequest` | `ProductResponse` | 200, 400, 404 |
| `DELETE` | `/products/{id}` | — | — | 204, 400, 404 |

Mapeamento de erros (após Fase 5):
- `service.ErrInvalidInput` / `service.ErrInvalidID` → **400**
- `repository.ErrNotFound` (propagado) → **404**
- outros → **500**

---

## 7. Próximos passos (curto prazo)

### Imediato — terminar a Fase 3
- **`Update`** com `ReplaceOne` + checagem de `MatchedCount == 0` → `ErrNotFound`
- **`Delete`** com `DeleteOne` + checagem de `DeletedCount == 0` → `ErrNotFound`

### Depois — Fases 4 → 7 (ordem do `fases.txt`)
1. **Fase 4 — service**: trocar `package product` por `package service`; implementar `Create`/`List`/`Get`/`Update`/`Delete` com validações e tradução de erros
2. **Fase 5 — handler**: 5 endpoints, helper `writeError`, mapeamento `err → status`, `r.PathValue("id")`
3. **Fase 6 — main**: wiring completo (config → db → repo → service → handler → mux), `defer client.Disconnect` no main
4. **Fase 7 — validação**: `gofmt -w .`, `go vet ./...`, `go build ./...`, teste com `curl`

---

## 8. Roadmap de aprendizado (pós-CRUD)

Cada item aqui é um **capítulo** a ser estudado quando chegar a hora. Não agora.

### 8.1 Endpoints e regras de estoque
- **`PATCH /products/{id}/stock`** — entrada/saída de unidades, com histórico.
- **`GET /products/low-stock?threshold=N`** — query parametrizada.

### 8.2 Middlewares
- **Logger middleware** — `log/slog` (stdlib Go 1.21+).
- **Recover middleware** — `defer recover()` evita derrubar o server num panic.
- **Request ID** — propagar um `X-Request-Id` por toda a cadeia (`context.WithValue`).
- **CORS** — pré-voo para futuro front.

### 8.3 Persistência e modelagem
- **Índices Mongo** — `slug` único, índice em `quantity` para `low-stock`.
- **Soft delete** — campo `deleted_at` com queries filtradas.
- **Validação estruturada** — `github.com/go-playground/validator/v10` com tags nos DTOs.
- **Decimal128** — preço em decimal real (não `float64`!), evita bug de ponto flutuante.

### 8.4 Autenticação e autorização
- **JWT** — `github.com/golang-jwt/jwt/v5`.
- **Middleware de auth** — extrair claims do header.
- **Roles** — `admin` vs `user` (quem pode criar/deletar?).

### 8.5 Testes
- **Unit no service** — com mock do repository (interface já facilita!).
- **Integration no handler** — `httptest.NewServer` + repositório real apontando pra Mongo de teste.
- **Table-driven tests** — convenção do Go para múltiplos cenários.
- **Test fixtures** — popular banco antes de cada teste.
- **Coverage** — `go test -cover ./...`.

### 8.6 Observabilidade e produção
- **`log/slog` estruturado** — JSON em prod, texto em dev.
- **Métricas** — `expvar` ou Prometheus.
- **Graceful shutdown** — `signal.NotifyContext`, `srv.Shutdown(ctx)`.
- **Health check** — `GET /healthz` e `GET /readyz`.

### 8.7 DevOps e documentação
- **Dockerfile** — `golang:1.22-alpine` multi-stage.
- **docker-compose** — sobe app + Mongo juntos.
- **OpenAPI/Swagger** — `github.com/swaggo/swag` ou `oapi-codegen`.
- **CI** — GitHub Actions: `gofmt`, `go vet`, `go test`, `go build`.

### 8.8 Performance e escala
- **Cache com Redis** — `GET /products/{id}` cacheado por TTL.
- **Paginação** — `?page=N&size=M` ou cursor-based.
- **Connection pooling** — tuning no driver Mongo.

### 8.9 Conceitos Go que vão aparecer pelo caminho

Marque conforme for estudando:

- [ ] `sentinel errors` + `errors.Is` / `errors.As` — ✅ já usamos
- [ ] `context.Context` (cancelamento, deadline, valores) — ✅ já usamos
- [ ] Pointer receiver vs value receiver — ✅ já usamos
- [ ] Satisfação implícita de interface (não precisa de `implements`)
- [ ] DI manual via construtores — virá na Fase 6
- [ ] Pattern matching do `ServeMux` 1.22+ (`r.PathValue`) — virá na Fase 5
- [ ] `json.Decoder` / `json.Encoder` — virá na Fase 5
- [ ] `bson.M` / `bson.D` / `bson.A` — ✅ já usamos parcialmente
- [ ] `time.Time` e fusos — virá na Fase 4
- [ ] `http.Handler` interface e composição
- [ ] `table-driven tests` — virá nos testes
- [ ] `httptest` — virá nos testes
- [ ] `mock` repository via interface — virá nos testes
- [ ] `slog` (log estruturado) — virá nos middlewares
- [ ] `signal.NotifyContext` (graceful shutdown) — virá no main
- [ ] `go-playground/validator` (validação com tags) — virá na evolução
- [ ] `golang-jwt/jwt` (auth) — virá na evolução
- [ ] `go:generate` e geração de código

---

## 9. Convenções do projeto (registrar aqui conforme decidir)

- **Idioma dos comentários:** PT-BR.
- **Idioma dos erros sentinel:** inglês (`product not found`).
- **Idioma dos logs:** inglês (padrão de mercado).
- **Formato:** `gofmt` é lei. Roda antes de cada commit.
- **Erros nunca ignorados:** nada de `_ =`. Sempre checar ou `defer log`.
- **Nada de magic strings:** constantes pra rotas, nomes de collection, etc.

---

## 10. Comandos úteis (cola)

```bash
go run ./cmd/api              # subir a API
go build -o estoque ./cmd/api # gerar binário
go build ./...                # checar compilação
go vet ./...                  # análise estática
gofmt -w .                    # formatar tudo
go mod tidy                   # sincronizar deps
```

---

## 11. Onde paramos da última vez

- **Fase 3 está em ~75%.**
- Feito: 3.1 (sentinel), 3.2 (`FindByID` propagando erro), 3.3 (`FindAll`).
- Resta: 3.4 (`Update`) e 3.5 (`Delete`).
- Quando concluir 3.4 e 3.5, **próximo capítulo = Fase 4 (service)**.
