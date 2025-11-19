# Top20 Backend

API simples em Go com PostgreSQL para gerenciar submissões de top 20 jogadores.

## Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=top20

# API Configuration
API_PORT=3000
```

## Executar com Docker Compose

```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:3000`

**Documentação Swagger:** `http://localhost:3000/api/docs/`

**Nota:** O `DB_HOST` será automaticamente sobrescrito para `postgres` dentro do container (configurado no docker-compose.yml).

## Endpoints

### POST /api/submissions
Submete uma lista de jogadores. **Apenas uma submissão por IP.**

**Regras de Validação:**
- O array `players` deve conter **exatamente 20 jogadores**
- O campo `submittedBy` é obrigatório

**Body:**
```json
{
  "players": [
    {"position": 1, "name": "Jogador 1"},
    {"position": 2, "name": "Jogador 2"},
    {"position": 3, "name": "Jogador 3"},
    ...
    {"position": 20, "name": "Jogador 20"}
  ],
  "submittedBy": "Nome do usuário"
}
```

**Resposta de Sucesso (201):**
Sem corpo de resposta, apenas status 201 Created.

**Resposta de Erro - Número incorreto de jogadores (400):**
```json
{
  "error": "Invalid number of players",
  "message": "Exactly 20 players are required, got 15"
}
```

**Resposta de Erro - IP já submeteu (409):**
```json
{
  "error": "IP address already submitted",
  "message": "Only one submission per IP address is allowed"
}
```

### GET /api/submissions
Retorna todas as submissões ou filtra por `submittedBy`.

**Query Params:**
- `submittedBy` (opcional): Filtra por quem submeteu

**Exemplos:**
```bash
# Todas as submissões
curl http://localhost:3000/api/submissions

# Filtrado por submittedBy
curl "http://localhost:3000/api/submissions?submittedBy=João"
```

### GET /api/players/stats
Retorna estatísticas de um jogador específico: em quantas submissões ele aparece e em quais posições.

**Query Params:**
- `name` (obrigatório): Nome do jogador

**Exemplo:**
```bash
curl "http://localhost:3000/api/players/stats?name=Cristiano%20Ronaldo"
```

**Resposta 200:**
```json
{
  "playerName": "Cristiano Ronaldo",
  "totalSubmissions": 10,
  "positionBreakdown": [
    { "position": 1, "count": 5 },
    { "position": 2, "count": 3 },
    { "position": 5, "count": 2 }
  ]
}
```

**Resposta 404:**
```json
{
  "error": "Player not found",
  "message": "Player 'Nome do Jogador' was not found in any submission"
}
```

### GET /api/health
Verifica o status da API e conexão com o banco.

```bash
curl http://localhost:3000/api/health
```

**Resposta:**
```json
{
  "status": "healthy"
}
```

## Documentação Swagger

A documentação interativa da API está disponível em `/api/docs/`:

```
http://localhost:3000/api/docs/
```

## Executar localmente (sem Docker)

1. Certifique-se de ter PostgreSQL rodando localmente
2. Configure as variáveis de ambiente no `.env`
3. Instale a ferramenta Swagger (apenas uma vez):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Gere a documentação Swagger:

```bash
swag init
```

5. Execute:

```bash
go run main.go
```

## Estrutura do Banco de Dados

A tabela `submissions` é criada automaticamente com a seguinte estrutura:

```sql
CREATE TABLE submissions (
  id SERIAL PRIMARY KEY,
  players JSONB NOT NULL,
  submitted_by VARCHAR(255) NOT NULL,
  ip_address VARCHAR(45) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

A constraint `UNIQUE` no `ip_address` garante que cada IP possa submeter apenas uma vez.

## Tecnologias

- **Go 1.21** - Linguagem de programação
- **PostgreSQL 15** - Banco de dados
- **Docker & Docker Compose** - Containerização
- **Swagger/OpenAPI** - Documentação da API

