# API Routes - Quick Reference

Base URL: `https://top20-api.fly.dev`

---

## 1. Health Check

```
GET /api/health
```

**Response 200:**
```json
{"status": "healthy"}
```

---

## 2. Submit Top 20 Players

```
POST /api/submissions
```

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "players": [
    {"position": 1, "name": "Cristiano Ronaldo"},
    {"position": 2, "name": "Lionel Messi"},
    ...
    {"position": 20, "name": "Neymar Jr"}
  ],
  "submittedBy": "João Silva"
}
```

**Rules:**
- Exactly 20 players required
- `submittedBy` is required
- One submission per IP address

**Response 201:** (empty body)

**Response 409:**
```json
{
  "error": "IP address already submitted",
  "message": "Only one submission per IP address is allowed"
}
```

---

## 3. List Submissions

```
GET /api/submissions
GET /api/submissions?submittedBy=João Silva
```

**Query Params:**
- `submittedBy` (optional): Filter by submitter name

**Response 200:**
```json
[
  {
    "id": 1,
    "players": [
      {"position": 1, "name": "Cristiano Ronaldo"},
      ...
    ],
    "submittedBy": "João Silva",
    "ipAddress": "192.168.1.1",
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

---

## 4. Player Statistics

```
GET /api/players/stats
GET /api/players/stats?name=Cristiano Ronaldo
```

**Query Params:**
- `name` (optional): Player name (case-insensitive). If provided, returns array with one player. If omitted, returns all players.

**Response 200 (specific player):**
```json
[
  {
    "playerName": "Cristiano Ronaldo",
    "totalSubmissions": 15,
    "positionBreakdown": [
      {"position": 1, "count": 8},
      {"position": 2, "count": 5},
      {"position": 3, "count": 2}
    ]
  }
]
```

**Response 200 (all players):**
```json
[
  {
    "playerName": "Cristiano Ronaldo",
    "totalSubmissions": 15,
    "positionBreakdown": [...]
  },
  {
    "playerName": "Lionel Messi",
    "totalSubmissions": 12,
    "positionBreakdown": [...]
  },
  ...
]
```

**Response 404:**
```json
{
  "error": "Player not found",
  "message": "Player 'Nome' was not found in any submission"
}
```

---

## Quick Examples

### cURL

```bash
# Health check
curl https://top20-api.fly.dev/api/health

# Submit
curl -X POST https://top20-api.fly.dev/api/submissions \
  -H "Content-Type: application/json" \
  -d '{"players":[{"position":1,"name":"Messi"},...], "submittedBy":"User"}'

# List all
curl https://top20-api.fly.dev/api/submissions

# Filter by user
curl "https://top20-api.fly.dev/api/submissions?submittedBy=João"

# Player stats - specific player (returns array with 1 item)
curl "https://top20-api.fly.dev/api/players/stats?name=Cristiano%20Ronaldo"

# Player stats - all players (returns array with all players)
curl "https://top20-api.fly.dev/api/players/stats"
```

### JavaScript

```javascript
// Submit
const response = await fetch('https://top20-api.fly.dev/api/submissions', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    players: [/* 20 players */],
    submittedBy: 'User Name'
  })
});

// Get submissions
const submissions = await fetch('https://top20-api.fly.dev/api/submissions')
  .then(r => r.json());

// Player stats - specific player (returns array)
const stats = await fetch('https://top20-api.fly.dev/api/players/stats?name=Messi')
  .then(r => r.json());
// stats[0].playerName === "Messi"

// Player stats - all players (returns array)
const allStats = await fetch('https://top20-api.fly.dev/api/players/stats')
  .then(r => r.json());
// allStats = [{playerName: "Player1", ...}, {playerName: "Player2", ...}, ...]
```

---

## Swagger UI

Interactive documentation:
```
https://top20-api.fly.dev/api/docs/
```

