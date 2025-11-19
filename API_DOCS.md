# Documentação da API - Integração Frontend

## 📍 URL Base

**Desenvolvimento Local:**
```
http://localhost:3000
```

**Produção (Fly.io):**
```
https://top20-api.fly.dev
```

**Swagger UI:**
- Local: `http://localhost:3000/api/docs/`
- Produção: `https://top20-api.fly.dev/api/docs/`

---

## 🔌 Endpoints

### 1. Health Check

**GET** `/api/health`

Verifica se a API está funcionando.

**Response 200:**
```json
{
  "status": "healthy"
}
```

**Response 503:**
```json
{
  "status": "unhealthy",
  "error": "mensagem de erro"
}
```

---

### 2. Submeter Top 20 Jogadores

**POST** `/api/submissions`

Submete uma lista com **exatamente 20 jogadores**. Apenas **uma submissão por IP** é permitida.

**Request Body:**
```json
{
  "players": [
    { "position": 1, "name": "Cristiano Ronaldo" },
    { "position": 2, "name": "Lionel Messi" },
    { "position": 3, "name": "Neymar Jr" },
    { "position": 4, "name": "Kylian Mbappé" },
    { "position": 5, "name": "Erling Haaland" },
    { "position": 6, "name": "Kevin De Bruyne" },
    { "position": 7, "name": "Robert Lewandowski" },
    { "position": 8, "name": "Mohamed Salah" },
    { "position": 9, "name": "Harry Kane" },
    { "position": 10, "name": "Vinícius Jr" },
    { "position": 11, "name": "Luka Modrić" },
    { "position": 12, "name": "Karim Benzema" },
    { "position": 13, "name": "Virgil van Dijk" },
    { "position": 14, "name": "Thibaut Courtois" },
    { "position": 15, "name": "Sadio Mané" },
    { "position": 16, "name": "Joshua Kimmich" },
    { "position": 17, "name": "Casemiro" },
    { "position": 18, "name": "Alisson Becker" },
    { "position": 19, "name": "Toni Kroos" },
    { "position": 20, "name": "Son Heung-min" }
  ],
  "submittedBy": "João Silva"
}
```

**Validações:**
- ✅ Array `players` deve conter **exatamente 20 jogadores**
- ✅ Campo `submittedBy` é **obrigatório**
- ✅ Cada player deve ter `position` (number) e `name` (string)

**Response 201 - Sucesso:**
Sem corpo de resposta, apenas status `201 Created`

**Response 400 - Número incorreto de jogadores:**
```json
{
  "error": "Invalid number of players",
  "message": "Exactly 20 players are required, got 15"
}
```

**Response 409 - IP já submeteu:**
```json
{
  "error": "IP address already submitted",
  "message": "Only one submission per IP address is allowed"
}
```

**Response 500 - Erro no servidor:**
```json
{
  "error": "Error message",
  "message": "Detailed error description"
}
```

---

### 3. Listar Submissões

**GET** `/api/submissions`

Retorna todas as submissões ou filtra por quem submeteu.

**Query Parameters:**
- `submittedBy` (opcional): Filtra submissões por nome

**Exemplos:**
```
GET /api/submissions
GET /api/submissions?submittedBy=João Silva
```

**Response 200:**
```json
[
  {
    "id": 1,
    "players": [
      { "position": 1, "name": "Cristiano Ronaldo" },
      { "position": 2, "name": "Lionel Messi" },
      ...
    ],
    "submittedBy": "João Silva",
    "ipAddress": "192.168.1.1",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "players": [...],
    "submittedBy": "Maria Santos",
    "ipAddress": "192.168.1.2",
    "createdAt": "2024-01-15T11:00:00Z"
  }
]
```

**Response vazio:**
```json
[]
```

---

### 4. Estatísticas de Jogador

**GET** `/api/players/stats`

Retorna estatísticas de um jogador específico: quantas vezes ele apareceu e em quais posições.

**Query Parameters:**
- `name` (obrigatório): Nome do jogador

**Exemplos:**
```
GET /api/players/stats?name=Cristiano Ronaldo
GET /api/players/stats?name=Lionel Messi
```

**Response 200:**
```json
{
  "playerName": "Cristiano Ronaldo",
  "totalSubmissions": 15,
  "positionBreakdown": [
    { "position": 1, "count": 8 },
    { "position": 2, "count": 5 },
    { "position": 3, "count": 2 }
  ]
}
```

**Interpretação:**
- O jogador "Cristiano Ronaldo" apareceu em 15 submissões no total
- 8 pessoas colocaram ele na posição 1
- 5 pessoas colocaram ele na posição 2
- 2 pessoas colocaram ele na posição 3

**Response 400 - Nome não fornecido:**
```json
{
  "error": "Missing required parameter",
  "message": "Player name is required"
}
```

**Response 404 - Jogador não encontrado:**
```json
{
  "error": "Player not found",
  "message": "Player 'Nome do Jogador' was not found in any submission"
}
```

**Notas:**
- A busca é case-insensitive (não diferencia maiúsculas de minúsculas)
- Não mostra quem votou em cada posição, apenas o count
- As posições são ordenadas numericamente

---

## 💻 Exemplos de Código

### JavaScript (Fetch API)

#### 1. Submeter jogadores

```javascript
const submitPlayers = async (players, submittedBy) => {
  try {
    const response = await fetch('http://localhost:3000/api/submissions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        players,
        submittedBy,
      }),
    });

    if (response.status === 201) {
      console.log('Submissão criada com sucesso!');
      return { success: true };
    }

    // Tratar erros
    const error = await response.json();
    return { success: false, error };
  } catch (err) {
    console.error('Erro na requisição:', err);
    return { success: false, error: 'Erro de conexão' };
  }
};

// Uso:
const players = [
  { position: 1, name: 'Cristiano Ronaldo' },
  { position: 2, name: 'Lionel Messi' },
  // ... 18 jogadores restantes
];

const result = await submitPlayers(players, 'João Silva');
if (result.success) {
  alert('Top 20 submetido com sucesso!');
} else {
  alert(`Erro: ${result.error.message}`);
}
```

#### 2. Listar submissões

```javascript
const getSubmissions = async (submittedBy = null) => {
  try {
    const url = submittedBy
      ? `http://localhost:3000/api/submissions?submittedBy=${encodeURIComponent(submittedBy)}`
      : 'http://localhost:3000/api/submissions';

    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error('Erro ao buscar submissões');
    }

    const data = await response.json();
    return data;
  } catch (err) {
    console.error('Erro:', err);
    return [];
  }
};

// Uso:
const allSubmissions = await getSubmissions();
const joaoSubmissions = await getSubmissions('João Silva');
```

#### 3. Estatísticas de jogador

```javascript
const getPlayerStats = async (playerName) => {
  try {
    const url = `http://localhost:3000/api/players/stats?name=${encodeURIComponent(playerName)}`;
    
    const response = await fetch(url);
    
    if (response.status === 404) {
      return { found: false, message: 'Jogador não encontrado' };
    }
    
    if (!response.ok) {
      throw new Error('Erro ao buscar estatísticas');
    }

    const data = await response.json();
    return { found: true, data };
  } catch (err) {
    console.error('Erro:', err);
    return { found: false, error: err.message };
  }
};

// Uso:
const stats = await getPlayerStats('Cristiano Ronaldo');
if (stats.found) {
  console.log(`${stats.data.playerName} apareceu em ${stats.data.totalSubmissions} submissões`);
  stats.data.positionBreakdown.forEach(pos => {
    console.log(`Posição ${pos.position}: ${pos.count} vezes`);
  });
}
```

#### 4. Health check

```javascript
const checkHealth = async () => {
  try {
    const response = await fetch('http://localhost:3000/api/health');
    const data = await response.json();
    return data.status === 'healthy';
  } catch (err) {
    return false;
  }
};

// Uso:
const isHealthy = await checkHealth();
console.log('API está funcionando:', isHealthy);
```

---

### React / Next.js

#### Hook customizado

```typescript
// hooks/useSubmissions.ts
import { useState } from 'react';

interface Player {
  position: number;
  name: string;
}

interface Submission {
  id: number;
  players: Player[];
  submittedBy: string;
  ipAddress: string;
  createdAt: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';

export const useSubmissions = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submitPlayers = async (players: Player[], submittedBy: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/api/submissions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ players, submittedBy }),
      });

      if (response.status === 201) {
        setLoading(false);
        return { success: true };
      }

      const errorData = await response.json();
      setError(errorData.message);
      setLoading(false);
      return { success: false, error: errorData.message };
    } catch (err) {
      const errorMsg = 'Erro ao conectar com a API';
      setError(errorMsg);
      setLoading(false);
      return { success: false, error: errorMsg };
    }
  };

  const getSubmissions = async (submittedBy?: string): Promise<Submission[]> => {
    setLoading(true);
    setError(null);

    try {
      const url = submittedBy
        ? `${API_URL}/api/submissions?submittedBy=${encodeURIComponent(submittedBy)}`
        : `${API_URL}/api/submissions`;

      const response = await fetch(url);
      
      if (!response.ok) {
        throw new Error('Erro ao buscar submissões');
      }

      const data = await response.json();
      setLoading(false);
      return data;
    } catch (err) {
      setError('Erro ao buscar submissões');
      setLoading(false);
      return [];
    }
  };

  return {
    submitPlayers,
    getSubmissions,
    loading,
    error,
  };
};
```

#### Componente de exemplo

```typescript
// components/Top20Form.tsx
'use client';

import { useState } from 'react';
import { useSubmissions } from '@/hooks/useSubmissions';

export default function Top20Form() {
  const [players, setPlayers] = useState(
    Array.from({ length: 20 }, (_, i) => ({
      position: i + 1,
      name: '',
    }))
  );
  const [submittedBy, setSubmittedBy] = useState('');
  const { submitPlayers, loading, error } = useSubmissions();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const result = await submitPlayers(players, submittedBy);

    if (result.success) {
      alert('Top 20 submetido com sucesso!');
      // Reset form
      setPlayers(Array.from({ length: 20 }, (_, i) => ({
        position: i + 1,
        name: '',
      })));
      setSubmittedBy('');
    } else {
      alert(`Erro: ${result.error}`);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Seu nome:</label>
        <input
          type="text"
          value={submittedBy}
          onChange={(e) => setSubmittedBy(e.target.value)}
          required
        />
      </div>

      {players.map((player, index) => (
        <div key={index}>
          <span>{player.position}.</span>
          <input
            type="text"
            placeholder={`Jogador ${player.position}`}
            value={player.name}
            onChange={(e) => {
              const newPlayers = [...players];
              newPlayers[index].name = e.target.value;
              setPlayers(newPlayers);
            }}
            required
          />
        </div>
      ))}

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <button type="submit" disabled={loading}>
        {loading ? 'Enviando...' : 'Submeter Top 20'}
      </button>
    </form>
  );
}
```

---

### Axios

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Submeter jogadores
export const submitPlayers = async (players, submittedBy) => {
  try {
    await api.post('/api/submissions', { players, submittedBy });
    return { success: true };
  } catch (error) {
    if (error.response) {
      return { success: false, error: error.response.data };
    }
    return { success: false, error: 'Erro de conexão' };
  }
};

// Listar submissões
export const getSubmissions = async (submittedBy = null) => {
  try {
    const params = submittedBy ? { submittedBy } : {};
    const response = await api.get('/api/submissions', { params });
    return response.data;
  } catch (error) {
    console.error('Erro ao buscar submissões:', error);
    return [];
  }
};

// Estatísticas de jogador
export const getPlayerStats = async (playerName) => {
  try {
    const response = await api.get('/api/players/stats', {
      params: { name: playerName }
    });
    return response.data;
  } catch (error) {
    if (error.response?.status === 404) {
      return null; // Jogador não encontrado
    }
    console.error('Erro ao buscar estatísticas:', error);
    return null;
  }
};

// Health check
export const checkHealth = async () => {
  try {
    const response = await api.get('/api/health');
    return response.data.status === 'healthy';
  } catch (error) {
    return false;
  }
};
```

---

## 🔒 CORS

A API aceita requisições de qualquer origem. Se precisar restringir, configure no backend:

```go
// Adicionar no main.go antes dos handlers
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
```

---

## 📊 Estrutura de Dados

### Player
```typescript
interface Player {
  position: number;  // 1-20
  name: string;
}
```

### Submission (Request)
```typescript
interface SubmissionRequest {
  players: Player[];      // Exatamente 20 players
  submittedBy: string;    // Nome de quem está submetendo
}
```

### Submission (Response)
```typescript
interface SubmissionResponse {
  id: number;
  players: Player[];
  submittedBy: string;
  ipAddress: string;
  createdAt: string;      // ISO 8601 format
}
```

### PlayerStats
```typescript
interface PlayerStats {
  playerName: string;
  totalSubmissions: number;
  positionBreakdown: PositionCount[];
}

interface PositionCount {
  position: number;       // 1-20
  count: number;          // Quantas vezes apareceu nessa posição
}
```

### Error Response
```typescript
interface ErrorResponse {
  error: string;
  message: string;
}
```

---

## 🎯 Fluxo Completo

```javascript
// 1. Verificar se a API está online
const isOnline = await checkHealth();
if (!isOnline) {
  alert('API offline');
  return;
}

// 2. Preparar dados
const players = [
  { position: 1, name: 'Cristiano Ronaldo' },
  // ... 19 jogadores restantes
];
const submittedBy = 'João Silva';

// 3. Submeter
const result = await submitPlayers(players, submittedBy);

// 4. Tratar resultado
if (result.success) {
  // Sucesso - redirecionar ou mostrar mensagem
  window.location.href = '/obrigado';
} else {
  // Erro - mostrar mensagem específica
  if (result.error.error === 'IP address already submitted') {
    alert('Você já submeteu seu Top 20!');
  } else if (result.error.error === 'Invalid number of players') {
    alert('Você precisa selecionar exatamente 20 jogadores');
  } else {
    alert('Erro ao submeter: ' + result.error.message);
  }
}

// 5. Buscar submissões para exibir
const submissions = await getSubmissions();
console.log('Total de submissões:', submissions.length);
```

---

## 🐛 Tratamento de Erros

```javascript
const handleApiError = (error) => {
  if (!error) return 'Erro desconhecido';

  switch (error.error) {
    case 'IP address already submitted':
      return 'Você já votou! Apenas uma submissão por pessoa é permitida.';
    
    case 'Invalid number of players':
      return 'Você deve selecionar exatamente 20 jogadores.';
    
    case 'Missing required field':
      return 'Por favor, preencha seu nome.';
    
    default:
      return error.message || 'Erro ao processar sua requisição.';
  }
};

// Uso:
const result = await submitPlayers(players, submittedBy);
if (!result.success) {
  const userMessage = handleApiError(result.error);
  showToast(userMessage, 'error');
}
```

---

## 🔗 Variáveis de Ambiente

### `.env.local` (Next.js)
```env
NEXT_PUBLIC_API_URL=http://localhost:3000
```

### `.env` (React/Vite)
```env
VITE_API_URL=http://localhost:3000
```

### Produção
```env
NEXT_PUBLIC_API_URL=https://top20-api.fly.dev
```

---

## 📝 Notas Importantes

1. **Limite de submissões**: Apenas 1 submissão por IP
2. **Validação obrigatória**: Exatamente 20 jogadores
3. **Campo obrigatório**: `submittedBy` não pode estar vazio
4. **Data format**: `createdAt` está em formato ISO 8601 (UTC)
5. **HTTPS**: Em produção, sempre use HTTPS

---

## 🧪 Testar no Swagger

Acesse `http://localhost:3000/api/docs/` para testar os endpoints diretamente no navegador com interface interativa.

---

## 📞 Suporte

Problemas com a integração? Verifique:
1. URL da API está correta
2. CORS está permitido
3. Body da requisição está no formato correto
4. Todos os 20 jogadores foram enviados
5. Console do navegador para erros de CORS/rede

