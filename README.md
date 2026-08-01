# clickGo

CLI de integração com a API do [ClickUp](https://clickup.com) — tarefas, listas, espaços e time tracking, com autenticação OAuth 2.0 e interface estilizada no terminal.

## Instalação

Baixe o binário para o seu sistema operacional na página [Releases](https://github.com/LuizFelipeOliver/clickGo/releases), descompacte e execute:

```bash
./clickGo --help
```

Opcional: mova para um diretório no `PATH` para usar de qualquer lugar:

```bash
sudo mv clickGo /usr/local/bin/
```

## Primeiro uso

1. Crie um app OAuth no ClickUp (avatar → **Settings** → **Apps** → *Create new app*) com redirect `http://localhost:8080/callback`. Só Workspace owners/admins conseguem criar.
2. Autentique informando o `client_id` e o `client_secret`:

```bash
clickGo auth --client-id SEU_CLIENT_ID --client-secret SEU_CLIENT_SECRET
```

O comando abre o navegador, captura o token automaticamente e salva em `~/.config/clickGo/.config.yaml` (permissões `0600`).

Alternativas para as credenciais: variáveis de ambiente `CLICKUP_CLIENT_ID` / `CLICKUP_CLIENT_SECRET`, ou informar uma única vez — ficam salvas no config.

## Uso

```bash
clickGo auth              # fluxo OAuth completo
clickGo auth status       # mostra se já há token salvo
clickGo me                # exibe o usuário autenticado (testa o token)
```

## Autenticação

- **OAuth 2.0** (grant type *authorization code*): URL de autorização `https://app.clickup.com/api`, troca do code em `POST /api/v2/oauth/token`.
- Token de acesso salvo localmente; usado no header `Authorization: Bearer <token>`.
- Não expire (por padrão), mas você pode revogá-lo nas configurações do ClickUp.

## Desenvolvimento

```bash
go build ./...   # compila
go vet ./...     # análise estática
go test ./...    # testes
```

## Licença

MIT — veja [LICENSE](LICENSE).
