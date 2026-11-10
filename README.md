# Star Concessionária — Terminal User Interface (TUI)

Aplicação de terminal (TUI) corporativa, moderna e multiplataforma (Linux e Windows 10/11) desenvolvida em **Go (Golang)** com o ecossistema reativo **Bubble Tea**, conectada à API oficial da Star Concessionária (`https://full-cycle-web-api.onrender.com/`).

---

## 🌟 Principais Recursos

- **Executável Único e Portátil:** Compilação estática (`CGO_ENABLED=0`) gerando um binário sem dependências externas (~11 MB).
- **Paridade Multiplataforma Real:** Suporte nativo ao Linux e Windows (Windows Terminal, PowerShell e CMD legado com habilitação automática de UTF-8 Code Page 65001 e Virtual Terminal Processing via Win32 API).
- **10 Telas Completas:** Cobertura de 100% dos recursos da concessionária (Dashboard, Veículos, Peças, Vendas, Agendamentos, Reservas, Propostas, SAC, Lojas e Perfil/Login).
- **Formulários Interativos em Modal:** Suporte a cadastros completos, movimentações de estoque, upload de fotos, pareceres comerciais e respostas de chamados diretamente pelo terminal.
- **Autenticação Segura JWT:** Formulário interativo para login com tratamento nativo de RFC 7807 (`ProblemDetails`) e persistência segura de sessão.
- **Resiliência a Cold Start:** Tolerância ao ciclo de hibernação de serviços em nuvem (Render Free Tier) com retentativas e avisos visuais.

---

## 📱 Módulos e Telas do Sistema

| Aba | Tecla | Módulo | Funcionalidades e Cadastros |
| :---: | :---: | :--- | :--- |
| **1** | `1` | **Dashboard** | Métricas operacionais em tempo real: faturamento consolidado, frota disponível/reservada/vendida, estoque por loja e ranking de peças mais vendidas. |
| **2** | `2` | **Veículos** | Catálogo paginado com ficha técnica completa. Atalhos: `[a]` Cadastrar Novo Veículo, `[s]` Alterar Situação Operacional, `[i]` Adicionar Fotos à Galeria. |
| **3** | `3` | **Peças** | Almoxarifado com alerta de estoque baixo. Atalhos: `[a]` Cadastrar Nova Peça, `[v]` Registrar Venda/Baixa de Peça. |
| **4** | `4` | **Vendas** | Histórico consolidado de vendas de veículos. Atalhos: `[a]` Emitir/Registrar Venda de Veículo. |
| **5** | `5` | **Agenda** | Grade de visitas e test-drives. Atalhos: `[v]` Agendar Visita Presencial, `[t]` Agendar Test-Drive com CNH, `[s]` Alterar Status da Agenda. |
| **6** | `6` | **Reservas** | Gestão de intenções de compra e bloqueios de veículos. Atalhos: `[a]` Solicitar Reserva, `[s]` Atualizar Status com Motivo. |
| **7** | `7` | **Propostas** | Negociações comerciais com contrapropostas e veículos usados de entrada. Atalhos: `[a]` Enviar Nova Proposta, `[s]` Registrar Parecer Comercial. |
| **8** | `8` | **SAC** | Atendimento ao cliente e suporte a leads. Atalhos: `[a]` Registrar Novo Chamado/Dúvida, `[r]` Registrar Atendimento/Resposta. |
| **9** | `9` | **Lojas** | Rede de concessionárias físicas com endereços e geolocalização. Atalhos: `[a]` Cadastrar Nova Loja/Filial. |
| **0** | `0` | **Perfil / Login** | Autenticação por credenciais (Email/Senha), verificação de cargo e `[L]` Logout. |

---

## ⌨️ Atalhos de Teclado Globais

| Tecla | Ação |
| :---: | :--- |
| `1` até `9` | Acessar diretamente as abas 1 a 9 |
| `0` | Acessar aba de **Perfil / Login** |
| `Tab` | Alternar sequencialmente para a próxima aba |
| `Shift+Tab` | Alternar para a aba anterior |
| `↑` / `↓` | Navegar pelas linhas das tabelas |
| `r` | Recarregar dados da tela ativa diretamente da API |
| `Esc` | Fechar modal interativo ou retornar ao Dashboard |
| `q` ou `Ctrl+C` | Encerrar a aplicação |

---

## 📦 Releases Oficiais para Download

As versões compiladas prontas para uso estão disponíveis na aba de Releases do GitHub:

- **Linux (x86_64):** [Download da Release Linux](https://github.com/UrielAndrade/TUI-StarConcessionaria/releases/tag/v1.0.0-linux)
- **Windows (10 e 11 x86_64):** [Download da Release Windows](https://github.com/UrielAndrade/TUI-StarConcessionaria/releases/tag/v1.0.0-windows)

---

## 🚀 Como Compilar a Partir do Código-Fonte

### Pré-requisitos
- [Go 1.23+](https://go.dev/)

```bash
# Executar diretamente
go run ./cmd/star-tui

# Executar testes automatizados com API real
make test

# Compilar para Linux e Windows simultaneamente
make build
```

---

## 🏗️ Estrutura do Projeto

```
TUI-StarConcessionaria/
├── cmd/star-tui/          # Ponto de entrada e inicialização Win32 / POSIX
├── internal/
│   ├── api/               # Cliente HTTP Resty cobrindo todos os 42 endpoints
│   ├── config/            # Persistência de token e configurações
│   ├── model/             # Modelos de domínio e DTOs tipados
│   └── ui/                # Motor reativo Bubble Tea
│       ├── components/    # Header responsivo, footer, cartões e helpers
│       ├── forms/         # Formulários interativos para todos os cadastros
│       ├── theme/         # Paleta visual corporativa Lip Gloss
│       └── views/         # 10 módulos de telas completos
└── Makefile               # Scripts de automação
```
