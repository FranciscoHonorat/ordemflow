# OrderFlow Analytics Platform

O **OrderFlow** é uma plataforma em desenvolvimento de gestão de pedidos com um pipeline analítico em tempo real. O objetivo principal deste projeto é aplicar na prática os conceitos de sistemas distribuídos de alta performance, cobrindo os pilares de **Clean Architecture**, **Domain-Driven Design (DDD)**, além de engenharia de dados (OLTP vs OLAP) e orquestração moderna.

---

## 🛠️ Stack Tecnológica Planejada

* **Linguagem:** Go 1.22+
* **Armazenamento Transacional (OLTP):** PostgreSQL 16
* **Armazenamento Analítico (OLAP):** ClickHouse 24
* **Mensageria & Eventos:** Apache Kafka 3
* **Cache:** Redis 7
* **Infraestrutura:** Docker & Kubernetes (k3s / minikube)

---

## 🔬 Onde Estou Agora: Estudo de Performance & Otimização de Memória

Antes de acoplar infraestrutura complexa, foquei em garantir que o núcleo do domínio fosse o mais performático possível. Durante o desenvolvimento do Value Object `Money` no **Order Context**, realizei um estudo prático de **Escape Analysis** utilizando as ferramentas nativas de benchmark do Go para analisar o impacto do Garbage Collector e alocações na Heap vs Stack.
