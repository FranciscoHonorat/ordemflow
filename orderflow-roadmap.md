# OrderFlow Analytics Platform
## Roadmap Completo — Go + DDD + Hexagonal Architecture + Docker + Kubernetes

---

## O Projeto

**OrderFlow** é uma plataforma de gestão de pedidos com pipeline analítico em tempo real.

- O lado **OLTP** gerencia pedidos, estoque e notificações com consistência transacional (PostgreSQL).
- O lado **OLAP** consome eventos via Kafka, agrega métricas e expõe relatórios analíticos (ClickHouse).
- Um **CDC connector** (Change Data Capture) captura mutações do OLTP e publica eventos no Kafka.
- Tudo containerizado com **Docker** e orquestrado com **Kubernetes**.

Este projeto cobre os capítulos 1–12 de *Designing Data-Intensive Applications* na prática.

---

## Stack Tecnológica

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.22+ |
| OLTP | PostgreSQL 16 |
| OLAP | ClickHouse 24 |
| Mensageria | Apache Kafka 3 + Zookeeper |
| Cache | Redis 7 |
| Containerização | Docker + Docker Compose |
| Orquestração | Kubernetes (k3s local ou minikube) |
| Observabilidade | Prometheus + Grafana + OpenTelemetry |
| API Gateway | NGINX Ingress Controller |
| Migrations | golang-migrate |
| Testes | testify + testcontainers-go |

---

## Bounded Contexts — DDD Subdomains

```
┌─────────────────────────────────────────────────────────────────┐
│                      OrderFlow Platform                         │
│                                                                 │
│  ┌──────────────────┐   ┌──────────────────┐                   │
│  │   Order Context  │   │ Inventory Context │                   │
│  │  (Core Domain)   │   │ (Supporting)      │                   │
│  │                  │   │                  │                   │
│  │  - Order         │   │  - Product       │                   │
│  │  - OrderItem     │   │  - Stock         │                   │
│  │  - OrderStatus   │   │  - StockMovement │                   │
│  └────────┬─────────┘   └────────┬─────────┘                   │
│           │ Domain Events        │ Domain Events               │
│           └──────────┬───────────┘                             │
│                      ▼                                         │
│              ┌───────────────┐                                  │
│              │     Kafka     │                                  │
│              │  Event Bus    │                                  │
│              └───────┬───────┘                                  │
│                      │                                         │
│           ┌──────────┴───────────┐                             │
│           ▼                      ▼                             │
│  ┌─────────────────┐   ┌──────────────────┐                   │
│  │Analytics Context│   │Notification Ctx  │                   │
│  │ (Supporting)    │   │ (Generic Domain) │                   │
│  │                 │   │                  │                   │
│  │ - OrderMetric   │   │  - Notification  │                   │
│  │ - SalesReport   │   │  - Channel       │                   │
│  │ - FunnelEvent   │   │  - Template      │                   │
│  └─────────────────┘   └──────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

### Por que esses subdomains?

- **Order (Core Domain)**: é o coração do negócio. Onde mais complexidade de domínio vive.
- **Inventory (Supporting Domain)**: suporta o core, tem regras próprias mas não é diferencial competitivo.
- **Analytics (Supporting Domain)**: lê eventos, não escreve no domínio principal — fronteira clara.
- **Notification (Generic Domain)**: poderia ser substituído por um SaaS (SendGrid, etc.) sem impacto no negócio.

---

## Arquitetura Hexagonal por Serviço

Cada bounded context é um serviço Go independente com a seguinte estrutura interna:

```
service/
├── domain/                  ← NÚCLEO — zero dependências externas
│   ├── entity/              ← Entidades e Aggregates
│   ├── valueobject/         ← Value Objects imutáveis
│   ├── event/               ← Domain Events
│   ├── repository/          ← Interfaces (Ports de saída)
│   └── service/             ← Domain Services
│
├── application/             ← CASOS DE USO — depende só do domain
│   ├── command/             ← Comandos (escrita)
│   ├── query/               ← Queries (leitura) — CQRS
│   └── port/
│       ├── input/           ← Driving Ports (interfaces de entrada)
│       └── output/          ← Driven Ports (interfaces de saída)
│
├── infrastructure/          ← ADAPTADORES — depende de tudo externo
│   ├── persistence/
│   │   ├── postgres/        ← Adapter PostgreSQL
│   │   └── clickhouse/      ← Adapter ClickHouse
│   ├── messaging/
│   │   └── kafka/           ← Adapter Kafka (producer/consumer)
│   ├── cache/
│   │   └── redis/           ← Adapter Redis
│   └── http/
│       ├── handler/         ← Handlers HTTP (chi/fiber)
│       └── middleware/      ← Auth, logging, tracing
│
├── cmd/
│   └── server/
│       └── main.go          ← Wiring — monta o hexágono
│
└── config/                  ← Configuração via env vars
```

### A regra de ouro da arquitetura hexagonal

```
domain ←── application ←── infrastructure
  ↑                              │
  └──────── (via interfaces) ────┘
```

O `domain` não importa nada de `infrastructure`. A `infrastructure` implementa as interfaces definidas no `domain`. O `cmd/server/main.go` é o único lugar que conhece tudo e faz o wiring (Dependency Injection manual).

---

## Estrutura de Repositórios

```
orderflow/
├── services/
│   ├── order-service/       ← OLTP: pedidos
│   ├── inventory-service/   ← OLTP: estoque
│   ├── analytics-service/   ← OLAP: métricas
│   ├── notification-service/← Genérico: alertas
│   └── cdc-connector/       ← Captura mudanças e publica no Kafka
│
├── infra/
│   ├── docker/              ← Dockerfiles por serviço
│   ├── k8s/                 ← Manifests Kubernetes
│   │   ├── base/            ← Configs compartilhadas
│   │   ├── order-service/
│   │   ├── analytics-service/
│   │   ├── databases/       ← StatefulSets: postgres, clickhouse
│   │   └── messaging/       ← StatefulSet: kafka
│   └── observability/       ← Prometheus, Grafana configs
│
├── shared/                  ← Pacotes Go compartilhados
│   ├── events/              ← Schemas de eventos (contratos)
│   ├── errors/              ← Erros de domínio padronizados
│   └── middleware/          ← Auth, tracing compartilhados
│
└── docker-compose.yml       ← Ambiente de desenvolvimento local
```

---

## Fluxo de Dados Completo

```
Cliente HTTP
     │
     ▼
[NGINX Ingress]
     │
     ├──► [order-service] ──WRITE──► [PostgreSQL]
     │         │                          │
     │         │ publica Domain Event      │ WAL
     │         ▼                          ▼
     │      [Kafka]        [cdc-connector polling]
     │         │                          │
     │         └──────────┬───────────────┘
     │                    ▼
     │           [analytics-service]
     │                    │
     │                    ├──WRITE──► [ClickHouse]
     │                    └──READ───► expõe /metrics, /reports
     │
     └──► [inventory-service] ──► [PostgreSQL]
               │
               └──► [Redis] (cache de stock disponível)
```

---

## Detalhamento por Serviço

### order-service

**Responsabilidade**: criar, atualizar e consultar pedidos.

```go
// domain/entity/order.go
type Order struct {
    id         OrderID
    customerID CustomerID
    items      []OrderItem
    status     OrderStatus
    total      Money
    createdAt  time.Time
    events     []DomainEvent   // eventos pendentes de publicação
}

// Aggregate root: toda mutação passa por aqui
func (o *Order) AddItem(product ProductID, qty int, price Money) error {
    if o.status != StatusDraft {
        return ErrOrderNotEditable
    }
    // regra de domínio: máximo 50 itens por pedido
    if len(o.items) >= 50 {
        return ErrMaxItemsExceeded
    }
    o.items = append(o.items, NewOrderItem(product, qty, price))
    o.recalculateTotal()
    o.events = append(o.events, OrderItemAdded{...})
    return nil
}

func (o *Order) Confirm() error {
    if len(o.items) == 0 {
        return ErrEmptyOrder
    }
    o.status = StatusConfirmed
    o.events = append(o.events, OrderConfirmed{
        OrderID:    o.id,
        CustomerID: o.customerID,
        Total:      o.total,
        OccurredAt: time.Now(),
    })
    return nil
}
```

```go
// application/command/confirm_order.go
type ConfirmOrderCommand struct {
    OrderID string
}

type ConfirmOrderHandler struct {
    repo      repository.OrderRepository     // driven port
    publisher port.output.EventPublisher     // driven port
    uow       port.output.UnitOfWork         // driven port
}

func (h *ConfirmOrderHandler) Handle(ctx context.Context, cmd ConfirmOrderCommand) error {
    return h.uow.Execute(ctx, func(ctx context.Context) error {
        order, err := h.repo.FindByID(ctx, OrderID(cmd.OrderID))
        if err != nil {
            return err
        }
        if err := order.Confirm(); err != nil {
            return err
        }
        if err := h.repo.Save(ctx, order); err != nil {
            return err
        }
        // publica eventos acumulados no aggregate
        return h.publisher.PublishAll(ctx, order.PullEvents())
    })
}
```

### analytics-service

**Responsabilidade**: consumir eventos do Kafka e projetar métricas no ClickHouse.

```go
// infrastructure/messaging/kafka/order_consumer.go
type OrderEventConsumer struct {
    reader  *kafka.Reader
    handler application.OrderEventHandler
}

func (c *OrderEventConsumer) Start(ctx context.Context) error {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            return err
        }
        event, err := deserialize(msg.Value)
        if err != nil {
            // dead letter queue — não bloqueia o consumer
            c.dlq.Send(ctx, msg)
            c.reader.CommitMessages(ctx, msg)
            continue
        }
        if err := c.handler.Handle(ctx, event); err != nil {
            return err  // retentativa — não faz commit
        }
        c.reader.CommitMessages(ctx, msg)
    }
}
```

```go
// infrastructure/persistence/clickhouse/metrics_repository.go
// ClickHouse é otimizado para INSERT em batch e SELECT analítico
func (r *MetricsRepository) InsertOrderConfirmed(ctx context.Context, e OrderConfirmedEvent) error {
    query := `
        INSERT INTO order_metrics 
            (event_time, order_id, customer_id, total_amount, items_count)
        VALUES (?, ?, ?, ?, ?)`
    return r.db.ExecContext(ctx, query,
        e.OccurredAt, e.OrderID, e.CustomerID, e.Total, e.ItemsCount)
}

// Query analítica — isso é o que o ClickHouse faz melhor que PostgreSQL
func (r *MetricsRepository) RevenueByDay(ctx context.Context, from, to time.Time) ([]DailyRevenue, error) {
    query := `
        SELECT 
            toDate(event_time)  AS day,
            sum(total_amount)   AS revenue,
            count()             AS orders,
            avg(total_amount)   AS avg_ticket
        FROM order_metrics
        WHERE event_time BETWEEN ? AND ?
        GROUP BY day
        ORDER BY day`
    // ...
}
```

---

## Docker — Uma Imagem por Serviço

### Dockerfile (multi-stage — padrão de produção)

```dockerfile
# infra/docker/order-service/Dockerfile

# Stage 1: build
FROM golang:1.22-alpine AS builder
WORKDIR /app

# cache de dependências separado do código
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /bin/order-service \
    ./services/order-service/cmd/server

# Stage 2: runtime mínimo
FROM gcr.io/distroless/static-debian12
COPY --from=builder /bin/order-service /order-service
EXPOSE 8080
ENTRYPOINT ["/order-service"]
```

Por que `distroless`? Imagem sem shell, sem package manager — superfície de ataque mínima. Tamanho final: ~10MB vs ~300MB de uma imagem ubuntu.

### docker-compose.yml (desenvolvimento local)

```yaml
version: "3.9"

services:
  # ─── OLTP ────────────────────────────────────────
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: orderflow
      POSTGRES_USER: orderflow
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports: ["5432:5432"]
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "orderflow"]
      interval: 5s

  # ─── OLAP ────────────────────────────────────────
  clickhouse:
    image: clickhouse/clickhouse-server:24
    volumes:
      - clickhouse_data:/var/lib/clickhouse
    ports: ["8123:8123", "9000:9000"]

  # ─── Mensageria ──────────────────────────────────
  zookeeper:
    image: confluentinc/cp-zookeeper:7.6.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.6.0
    depends_on: [zookeeper]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
    ports: ["9092:9092"]

  # ─── Cache ───────────────────────────────────────
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]

  # ─── Serviços ────────────────────────────────────
  order-service:
    build:
      context: .
      dockerfile: infra/docker/order-service/Dockerfile
    environment:
      DATABASE_URL: postgres://orderflow:secret@postgres:5432/orderflow
      KAFKA_BROKERS: kafka:9092
      PORT: 8080
    depends_on:
      postgres: { condition: service_healthy }
      kafka: { condition: service_started }
    ports: ["8080:8080"]

  analytics-service:
    build:
      context: .
      dockerfile: infra/docker/analytics-service/Dockerfile
    environment:
      CLICKHOUSE_URL: clickhouse://clickhouse:9000/orderflow
      KAFKA_BROKERS: kafka:9092
      KAFKA_GROUP_ID: analytics-service
    depends_on: [clickhouse, kafka]

  inventory-service:
    build:
      context: .
      dockerfile: infra/docker/inventory-service/Dockerfile
    environment:
      DATABASE_URL: postgres://orderflow:secret@postgres:5432/orderflow
      REDIS_URL: redis://redis:6379
      KAFKA_BROKERS: kafka:9092
    depends_on:
      postgres: { condition: service_healthy }
      redis: { condition: service_started }

volumes:
  postgres_data:
  clickhouse_data:
```

---

## Kubernetes — Manifests por Serviço

### Estrutura dos manifests

```
infra/k8s/
├── namespace.yaml
├── base/
│   └── configmap.yaml          ← configs não-secretas compartilhadas
├── databases/
│   ├── postgres-statefulset.yaml
│   ├── postgres-service.yaml
│   ├── postgres-pvc.yaml
│   ├── clickhouse-statefulset.yaml
│   └── clickhouse-service.yaml
├── messaging/
│   ├── zookeeper-statefulset.yaml
│   ├── kafka-statefulset.yaml
│   └── kafka-service.yaml
├── order-service/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml                ← auto-scaling
│   └── configmap.yaml
├── analytics-service/
│   ├── deployment.yaml
│   └── service.yaml
└── ingress.yaml                ← roteamento externo
```

### namespace.yaml

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: orderflow
  labels:
    app.kubernetes.io/managed-by: kubectl
```

### PostgreSQL — StatefulSet (não Deployment)

```yaml
# infra/k8s/databases/postgres-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet           # ← banco de dados SEMPRE StatefulSet
metadata:
  name: postgres
  namespace: orderflow
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        env:
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:             # ← senha nunca em configmap
              name: postgres-secret
              key: password
        - name: POSTGRES_DB
          value: orderflow
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        readinessProbe:
          exec:
            command: ["pg_isready", "-U", "postgres"]
          initialDelaySeconds: 5
          periodSeconds: 5
  volumeClaimTemplates:                # ← PVC gerenciado pelo StatefulSet
  - metadata:
      name: postgres-data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 5Gi
```

### order-service — Deployment

```yaml
# infra/k8s/order-service/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: orderflow
spec:
  replicas: 2                          # ← mínimo 2 para HA
  selector:
    matchLabels:
      app: order-service
  strategy:
    type: RollingUpdate                # ← zero downtime deploy
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: order-service
        image: orderflow/order-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: order-service-secret
              key: database-url
        - name: KAFKA_BROKERS
          valueFrom:
            configMapKeyRef:
              name: orderflow-config
              key: kafka-brokers
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### HPA — Horizontal Pod Autoscaler

```yaml
# infra/k8s/order-service/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: order-service-hpa
  namespace: orderflow
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: order-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Ingress — roteamento externo

```yaml
# infra/k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: orderflow-ingress
  namespace: orderflow
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
spec:
  ingressClassName: nginx
  rules:
  - host: orderflow.local
    http:
      paths:
      - path: /orders(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: order-service
            port:
              number: 80
      - path: /analytics(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: analytics-service
            port:
              number: 80
```

---

## Roadmap — 12 Semanas

### Semana 1 — Fundamentos e Setup

**Objetivo**: ambiente funcionando, estrutura do projeto criada.

Diário:
- Dia 1: instalar Go 1.22, Docker Desktop, minikube ou k3s. Ler cap. 1 do Kleppmann.
- Dia 2: criar o monorepo, configurar Go workspaces (`go work`), escrever o primeiro `Makefile`.
- Dia 3: implementar o domain layer do `order-service` — entidades `Order`, `OrderItem`, `Money` (value object).
- Dia 4: escrever testes unitários para o domain — sem banco, sem HTTP, só regras de negócio.
- Dia 5: revisar, refatorar, documentar decisões em `docs/architecture-decisions/`.

**Critério de conclusão**: `go test ./...` verde no domain layer.

---

### Semana 2 — Domain Layer Completo

**Objetivo**: todos os 4 bounded contexts com domain layer implementado.

- Implementar `inventory-service` domain: `Product`, `Stock`, `StockMovement`.
- Implementar eventos de domínio: `OrderConfirmed`, `StockReserved`, `StockInsufficient`.
- Definir todas as interfaces de repositório (driven ports).
- Definir contratos de eventos no pacote `shared/events/`.
- Estudar: cap. 2 do Kleppmann (modelos de dados).

**Critério de conclusão**: interfaces definidas, testes de domínio completos.

---

### Semana 3 — Application Layer + PostgreSQL Adapter

**Objetivo**: primeiro caso de uso funcionando com banco real.

- Implementar `ConfirmOrderHandler` com Unit of Work.
- Implementar adapter PostgreSQL para `OrderRepository`.
- Configurar `golang-migrate` para migrations.
- Subir PostgreSQL com Docker Compose.
- Escrever testes de integração com `testcontainers-go` (sobe PostgreSQL real no teste).

```go
// Teste de integração — não usa mock
func TestOrderRepository_Integration(t *testing.T) {
    ctx := context.Background()
    pg, _ := testcontainers.RunContainer(ctx, postgres.RunContainer)
    defer pg.Terminate(ctx)

    repo := NewPostgresOrderRepository(pg.ConnectionString())
    order := domain.NewOrder(...)
    
    require.NoError(t, repo.Save(ctx, order))
    found, err := repo.FindByID(ctx, order.ID())
    require.NoError(t, err)
    assert.Equal(t, order.Status(), found.Status())
}
```

**Critério de conclusão**: criar e confirmar um pedido persiste no PostgreSQL.

---

### Semana 4 — HTTP Adapter + Docker

**Objetivo**: API REST funcionando, imagem Docker buildada.

- Implementar handlers HTTP com `chi` ou `fiber`.
- Implementar middleware: logging estruturado (`slog`), request ID, recovery.
- Criar `Dockerfile` multi-stage para `order-service`.
- Testar a imagem local: `docker build` → `docker run`.
- Estudar: cap. 3 do Kleppmann (storage engines e indexes).

```bash
# Deve funcionar ao final da semana
docker compose up order-service postgres
curl -X POST localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "abc", "items": [...]}'
```

**Critério de conclusão**: `docker compose up` funciona, API responde corretamente.

---

### Semana 5 — Kafka + Domain Events

**Objetivo**: eventos publicados no Kafka ao confirmar pedido.

- Subir Kafka no Docker Compose.
- Implementar `KafkaEventPublisher` (adapter de saída).
- Integrar publisher no `ConfirmOrderHandler`.
- Aprender sobre partições, offsets, consumer groups.
- Usar `kafka-ui` ou `kcat` para inspecionar mensagens.
- Estudar: cap. 11 do Kleppmann (stream processing).

**Critério de conclusão**: confirmar um pedido → mensagem aparece no Kafka topic `order.confirmed`.

---

### Semana 6 — inventory-service + CDC Connector

**Objetivo**: estoque reservado ao confirmar pedido; CDC capturando mudanças.

- Implementar `inventory-service` completo (domain → application → infrastructure).
- Implementar `StockReservationConsumer` que escuta `order.confirmed` e reserva estoque.
- Implementar `cdc-connector`: polling do PostgreSQL via `created_at`/`updated_at` ou WAL com `pglogical`.
- Estudar: cap. 10 e 11 do Kleppmann (batch e stream processing, CDC).

**Critério de conclusão**: confirmar pedido → estoque decrementado → evento publicado.

---

### Semana 7 — analytics-service + ClickHouse

**Objetivo**: pipeline OLAP funcionando end-to-end.

- Subir ClickHouse no Docker Compose.
- Criar schema analítico: tabela `order_metrics` com MergeTree engine.
- Implementar `analytics-service`: consome Kafka, insere no ClickHouse em batch.
- Implementar endpoint de relatório: receita por dia, ticket médio, top produtos.
- Comparar query no PostgreSQL vs ClickHouse na mesma massa de dados — sinta a diferença.
- Estudar: cap. 3 do Kleppmann (column-oriented storage).

**Critério de conclusão**: 1000 pedidos inseridos → query analítica retorna resultado correto.

---

### Semana 8 — Observabilidade

**Objetivo**: sistema observável com métricas, logs e traces.

- Adicionar `prometheus/client_golang` em todos os serviços: latência, throughput, erros.
- Configurar Grafana com dashboards: RED (Rate, Errors, Duration) por serviço.
- Implementar health checks: `/health/live` e `/health/ready`.
- Adicionar OpenTelemetry para distributed tracing.
- Subir Prometheus + Grafana no Docker Compose.
- Estudar: o que é SLI, SLO, SLA.

**Critério de conclusão**: Grafana mostrando métricas reais dos serviços.

---

### Semana 9 — Kubernetes: Databases e Messaging

**Objetivo**: PostgreSQL, ClickHouse e Kafka rodando no Kubernetes.

- Instalar minikube ou k3s localmente.
- Criar namespace `orderflow`.
- Aplicar StatefulSets para PostgreSQL e ClickHouse com PersistentVolumeClaims.
- Aplicar StatefulSet para Kafka + Zookeeper.
- Criar Secrets para senhas (não commitar senhas no git).
- Entender: por que StatefulSet para bancos e não Deployment.

```bash
kubectl apply -f infra/k8s/namespace.yaml
kubectl apply -f infra/k8s/databases/
kubectl apply -f infra/k8s/messaging/
kubectl get pods -n orderflow  # todos Running
```

**Critério de conclusão**: databases e kafka rodando no cluster, dados persistindo após restart do pod.

---

### Semana 10 — Kubernetes: Serviços + Ingress + HPA

**Objetivo**: todos os serviços rodando no cluster com roteamento externo.

- Criar Deployments para `order-service`, `analytics-service`, `inventory-service`.
- Configurar Services (ClusterIP para comunicação interna).
- Instalar NGINX Ingress Controller.
- Configurar Ingress para roteamento externo.
- Configurar HPA para `order-service`.
- Fazer rolling update sem downtime: `kubectl set image`.
- Estudar: ConfigMaps vs Secrets, por que não usar `latest` em produção.

**Critério de conclusão**: `curl orderflow.local/orders` funciona via Ingress, HPA escala sob load.

---

### Semana 11 — Testes, Segurança e Refinamento

**Objetivo**: sistema pronto para produção (ou próximo disso).

- Escrever testes end-to-end contra o cluster Kubernetes com `testcontainers-go`.
- Implementar autenticação JWT no `order-service`.
- Adicionar rate limiting no Ingress.
- Revisar todas as migrations, adicionar índices faltantes.
- Fazer load test com `k6` ou `hey`: 100 req/s por 60 segundos.
- Analisar bottlenecks com Grafana — otimizar queries lentas.

**Critério de conclusão**: load test passa sem erros, latência p99 < 200ms.

---

### Semana 12 — Portfólio e Documentação

**Objetivo**: projeto pronto para mostrar em entrevistas.

- Escrever `README.md` profissional: arquitetura, decisões de design, como rodar localmente.
- Criar diagrama de arquitetura (excalidraw ou mermaid).
- Gravar um demo de 2-3 minutos mostrando o sistema funcionando.
- Escrever um post técnico no LinkedIn/Medium explicando uma decisão de design difícil.
- Preparar 10 perguntas de entrevista que você consegue responder sobre o projeto.
- Push para GitHub público com CI/CD (GitHub Actions buildando e testando).

---

## Perguntas de Entrevista que este Projeto Responde

1. Como você garante consistência em transações distribuídas?
2. Qual a diferença entre OLTP e OLAP e quando usar cada um?
3. O que é CDC (Change Data Capture) e por que não usar triggers?
4. Como você implementaria rollback de eventos Kafka mal processados?
5. Por que StatefulSet para banco de dados e não Deployment?
6. Qual a diferença entre Liveness e Readiness probe no Kubernetes?
7. Como funciona Hexagonal Architecture e quais problemas ela resolve?
8. O que é um Aggregate no DDD e por que ele importa?
9. Como você faria backpressure em um consumer Kafka?
10. Por que ClickHouse é mais rápido que PostgreSQL para queries analíticas?

---

## Critério Final de Senioridade neste Projeto

Você atingiu nível **pleno** quando consegue explicar cada decisão de design do projeto sem consultar nada.

Você atingiu nível **sênior** quando consegue identificar os trade-offs de cada decisão e dizer em qual contexto você faria diferente.

Por exemplo: *"Usei Kafka aqui, mas para um sistema com 100 req/dia usaria PostgreSQL LISTEN/NOTIFY — Kafka tem overhead operacional que não se justifica em baixo volume."*

Esse raciocínio é o que diferencia um sênior. O projeto é o veículo para você construir esse raciocínio.
