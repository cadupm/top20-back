# Deploy no Fly.io - Guia Passo a Passo

## ⚠️ Importante sobre Docker Compose

O **Fly.io NÃO suporta docker-compose**. Você precisa:
- Usar o **Postgres gerenciado** do Fly.io
- Deployar a **API separadamente** usando o Dockerfile

## 🎉 Free Trial

O Fly.io oferece **7 dias grátis** para testar. Depois:
- **Free tier**: $5/mês de crédito grátis
- **Postgres**: ~$2-3/mês para desenvolvimento
- **API**: Gratuito com auto-stop (escala para 0 quando não está em uso)

## 📋 Passo 1: Instalar o Fly CLI

```bash
curl -L https://fly.io/install.sh | sh
```

**Adicionar ao PATH** (adicione ao seu `~/.bashrc` ou `~/.zshrc`):
```bash
export FLYCTL_INSTALL="/home/$USER/.fly"
export PATH="$FLYCTL_INSTALL/bin:$PATH"
```

Recarregar o shell:
```bash
source ~/.zshrc  # ou source ~/.bashrc
```

## 🚀 Passo 2: Login no Fly.io

```bash
fly auth login
```

Isso abrirá o navegador para você fazer login.

## 🗄️ Passo 3: Criar o Banco de Dados PostgreSQL

```bash
fly postgres create --name top20-db --region gru
```

**⚠️ IMPORTANTE:** Anote as credenciais que aparecerem!
- Username
- Password  
- Hostname
- Database name

## 📦 Passo 4: Launch da Aplicação

Na raiz do projeto, execute:

```bash
fly launch
```

O Fly.io vai:
1. Detectar o Dockerfile automaticamente
2. Sugerir um nome para a app (você pode mudar)
3. Perguntar sobre região (escolha `gru` - São Paulo)
4. Perguntar se quer fazer deploy agora - **responda NÃO** (ainda precisamos configurar o banco)

Isso vai criar um arquivo `fly.toml` automaticamente.

## 🔗 Passo 5: Conectar o Banco de Dados

Substitua `top20-api` pelo nome da sua app se for diferente:

```bash
fly postgres attach top20-db --app top20-api
```

Isso cria automaticamente a variável `DATABASE_URL` com a connection string do Postgres.

## ⚙️ Passo 6: Configurar Variáveis de Ambiente

Edite o arquivo `fly.toml` que foi criado e adicione:

```toml
[env]
  API_PORT = "8080"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0

  [[http_service.checks]]
    grace_period = "10s"
    interval = "30s"
    method = "GET"
    timeout = "5s"
    path = "/api/health"
```

## 🚀 Passo 7: Deploy!

```bash
fly deploy
```

Ou para build mais rápido (build remoto):
```bash
fly deploy --remote-only
```

## ✅ Passo 8: Verificar o Deploy

```bash
# Ver logs
fly logs

# Ver status
fly status

# Abrir no navegador
fly open
```

## 🌐 URL da Aplicação

Após o deploy, sua API estará disponível em:
```
https://[nome-da-sua-app].fly.dev
```

Endpoints disponíveis:
- `https://[nome-da-sua-app].fly.dev/api/health`
- `https://[nome-da-sua-app].fly.dev/api/submissions`
- `https://[nome-da-sua-app].fly.dev/api/players/stats?name=Cristiano%20Ronaldo`
- `https://[nome-da-sua-app].fly.dev/api/docs/`

## 🔧 Comandos Úteis

### Monitoramento
```bash
# Ver logs em tempo real
fly logs

# Ver métricas no dashboard
fly dashboard

# SSH na máquina
fly ssh console

# Ver informações da app
fly info
```

### Gerenciamento
```bash
# Fazer novo deploy
fly deploy

# Restart da aplicação
fly apps restart top20-api

# Listar todas as apps
fly apps list

# Listar bancos de dados
fly postgres list
```

### Variáveis de Ambiente
```bash
# Ver todas as variáveis (secrets)
fly secrets list

# Adicionar variável
fly secrets set NOME=valor

# Remover variável
fly secrets unset NOME
```

### Escalonamento
```bash
# Escalar memória
fly scale memory 512

# Escalar número de instâncias
fly scale count 2

# Ver regiões disponíveis
fly platform regions

# Adicionar região
fly regions add gru
```

### Banco de Dados
```bash
# Conectar via proxy
fly proxy 5432 -a top20-db

# Em outro terminal, conectar com psql
psql postgres://postgres:senha@localhost:5432/top20

# Ver databases no cluster
fly postgres db list -a top20-db

# Backup do banco
fly postgres backup -a top20-db
```

## 🐛 Troubleshooting

### Aplicação não inicia

```bash
# Ver logs para identificar o problema
fly logs

# Ver detalhes das máquinas
fly status --all

# SSH para debugar
fly ssh console
```

### Banco de dados não conecta

```bash
# Verificar se o attach foi feito corretamente
fly postgres list

# Ver a DATABASE_URL
fly secrets list

# Testar conexão
fly ssh console
# Dentro da máquina:
env | grep DATABASE_URL
```

### Erro de build

```bash
# Rebuild com logs mais detalhados
fly deploy --verbose

# Limpar cache e rebuildar
fly deploy --no-cache
```

### App não responde

```bash
# Verificar health check
fly checks list

# Ver se a porta está correta
fly status

# Restart
fly apps restart
```

## 🔄 CI/CD com GitHub Actions

Crie `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Fly.io

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: superfly/flyctl-actions/setup-flyctl@master
      
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

**Configurar o token:**

```bash
# Gerar token
fly tokens create deploy

# Adicionar no GitHub:
# Settings > Secrets and variables > Actions > New repository secret
# Name: FLY_API_TOKEN
# Value: [cole o token gerado]
```

## 🗑️ Destruir Recursos

**⚠️ CUIDADO:** Isso é irreversível!

```bash
# Destruir aplicação
fly apps destroy top20-api

# Destruir banco de dados
fly apps destroy top20-db
```

## 📝 Checklist Completo

- [ ] Instalar Fly CLI
- [ ] Fazer login (`fly auth login`)
- [ ] Criar banco PostgreSQL (`fly postgres create`)
- [ ] Anotar credenciais do banco
- [ ] Executar `fly launch` (escolher região `gru`)
- [ ] Conectar banco à app (`fly postgres attach`)
- [ ] Configurar `fly.toml` com porta e health check
- [ ] Deploy (`fly deploy`)
- [ ] Verificar logs (`fly logs`)
- [ ] Testar endpoints (`fly open`)
- [ ] Verificar Swagger em `/api/docs/`

## 🔐 Segurança

1. **Nunca commite** credenciais no Git
2. Use `fly secrets` para variáveis sensíveis
3. O Fly.io usa HTTPS automaticamente
4. O código já suporta `DATABASE_URL` (nativo do Fly.io)
5. Health check está configurado em `/api/health`

## 📚 Documentação Oficial

- [Fly.io Docs](https://fly.io/docs/)
- [Fly Postgres](https://fly.io/docs/postgres/)
- [Dockerfile Deploy](https://fly.io/docs/languages-and-frameworks/dockerfile/)
- [Fly CLI Reference](https://fly.io/docs/flyctl/)

