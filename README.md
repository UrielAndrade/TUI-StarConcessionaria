# Star Concessionária — Terminal User Interface (TUI)

Aplicação de terminal (TUI) corporativa, moderna e multiplataforma (Linux e Windows 10/11) desenvolvida em **Go (Golang)** com o ecossistema reativo **Bubble Tea**, conectada à API oficial da Star Concessionária.

---

## 🌟 Principais Recursos

- **Executável Único e Portátil:** Compilação estática (`CGO_ENABLED=0`) gerando um binário sem dependências externas (~14 MB).
- **Paridade Multiplataforma Real:** Suporte nativo ao Linux e Windows (Windows Terminal, PowerShell e CMD legado com habilitação automática de UTF-8 e Virtual Terminal Processing via Win32 API).
- **Dashboard Gerencial em Tempo Real:** Indicadores operacionais consolidados de estoque, faturamento, peças mais vendidas e alertas de frota.
- **Catálogo de Veículos e Estoque:** Navegação fluida por teclado com paginação, visualização prévia da ficha técnica e fotos.
- **Gestão de Peças e Vendas:** Tabela de estoque de autopeças e módulo de baixa por venda.
- **Autenticação Segura JWT:** Formulário interativo para login com tratamento nativo de RFC 7807 (`ProblemDetails`) e persistência segura de sessão.
- **Resiliência a Cold Start:** Tolerância ao ciclo de hibernação de serviços em nuvem (Render Free Tier) com retentativas e avisos textuais.

---

## 🚀 Como Executar

### Pré-requisitos
- [Go 1.23+](https://go.dev/) (apenas para compilação a partir do código-fonte).

### Compilação e Execução Direta

```bash
# Executar diretamente via Go
go run ./cmd/star-tui

# Ou construir e rodar o binário Linux
make run
```

### Compilação Cruzada (Linux e Windows)

```bash
# Gerar binários para Linux (star-tui) e Windows (star-tui.exe)
make build
```

---

## ⌨️ Atalhos de Teclado Globais

| Tecla | Ação |
| :---: | :--- |
| `1` | Acessar aba **Dashboard** |
| `2` | Acessar aba **Veículos** (Catálogo) |
| `3` | Acessar aba **Peças** (Estoque) |
| `4` | Acessar aba **Login** (Autenticação) |
| `Tab` | Alternar sequencialmente entre as abas |
| `↑` / `↓` | Navegar pelas linhas das tabelas |
| `r` | Recarregar dados da tela ativa a partir da API |
| `q` ou `Ctrl+C` | Encerrar a aplicação |

---

## 🏗️ Arquitetura do Projeto

O projeto segue os princípios de *Clean Architecture* e *The Elm Architecture* (TEA):

```
TUI-StarConcessionaria/
├── cmd/star-tui/          # Ponto de entrada e inicialização de plataforma (Win32 / POSIX)
├── internal/
│   ├── api/               # Cliente HTTP Resty, DTOs e tratamento RFC 7807
│   ├── config/            # Carregamento de YAML, ENV e caminhos cross-platform
│   ├── model/             # Estruturas e entidades de domínio
│   └── ui/                # Motor Bubble Tea, temas Lip Gloss, componentes e views
│       ├── components/    # Header, footer, cartões e modais
│       ├── theme/         # Paleta corporativa e estilos
│       └── views/         # Dashboard, Veículos, Peças e Login
└── Makefile               # Scripts de automação de testes e build
```

---

## 🧪 Testes Automatizados

O projeto conta com testes de integração cobrindo a comunicação com a API remota em produção:

```bash
make test
```
