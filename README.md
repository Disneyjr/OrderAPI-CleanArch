# Clean Architecture: Listagem de Orders (REST, gRPC e GraphQL)

Aplicação em Go utilizando Clean Architecture para criação e listagem de pedidos (`Orders`) através de três interfaces simultâneas: **REST**, **gRPC** e **GraphQL**.

## Como Rodar o Projeto

1. Certifique-se de ter o **Docker** e o **Docker Compose** instalados e em execução.
2. Na raiz do projeto, execute o comando único:
   ```bash
   docker compose up
   ```
   *(Opcional: use `docker compose up -d --build` para rodar em background e forçar o build da imagem).*

O Docker Compose irá automaticamente:
- Subir o banco de dados **MySQL** e o **RabbitMQ**;
- Executar as **migrações** criando a tabela `orders`;
- Compilar e iniciar a **aplicação Go** expondo os servidores Web (REST), gRPC e GraphQL.

---

## Portas dos Serviços

| Serviço | Porta | Endereço / Endpoint |
| :--- | :--- | :--- |
| **Web / REST API** | `8000` | `http://localhost:8000/order` |
| **gRPC Server** | `50051` | `localhost:50051` |
| **GraphQL Playground** | `8080` | `http://localhost:8080/` (Query endpoint: `/query`) |
| **MySQL** | `3306` | `localhost:3306` |
| **RabbitMQ** | `5672` / `15672` | `http://localhost:15672/` (`guest` / `guest`) |

---

## Como Testar

- **REST:** Utilize as requisições prontas no arquivo `api.http` na raiz do projeto (`POST http://localhost:8000/order` para criar e `GET http://localhost:8000/order` para listar).
- **GraphQL:** Acesse `http://localhost:8080` no navegador e execute a mutation `createOrder` e a query `ListOrders`:
  ```graphql
  mutation {
    createOrder(input: { id: "a", Price: 100.5, Tax: 0.5 }) {
      id Price Tax FinalPrice
    }
  }

  query {
    ListOrders {
      id Price Tax FinalPrice
    }
  }
  ```
- **gRPC:** Conecte na porta `50051` via **Evans** (`evans -r repl -p 50051`) ou **grpcurl** e chame os métodos `CreateOrder` e `ListOrders` do serviço `pb.OrderService`.
