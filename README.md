# clickGo

CLI de integração com a API do [ClickUp](https://clickup.com) — focada em visualizar rapidamente suas **tasks**, **subtasks** e **tasks relacionadas**, direto pelo terminal.

## Instalação

Baixe o binário para o seu sistema operacional na página [Releases](https://github.com/LuizFelipeOliver/clickGo/releases), descompacte e execute:

```bash
./clickGo --help
```

Opcional: mova para um diretório no `PATH` para usar de qualquer lugar:

```bash
sudo mv clickGo /usr/local/bin/
```

## Primeiro uso (Autenticação)

Execute o comando abaixo para abrir a página de configurações do ClickUp no navegador, copiar seu **Personal API Token (`pk_...`)** e colar no terminal:

```bash
clickGo auth
```

O token e o seu Workspace principal ficam salvos em `~/.config/clickGo/config.yaml` (permissões `0600`).

Alternativa: passar o token direto por flag:

```bash
clickGo auth --token pk_SEU_TOKEN
```

## Comandos

```bash
# Autenticação
clickGo auth                  # abre o navegador e salva seu Personal API Token
clickGo auth status           # verifica a conta e o Workspace autenticados
clickGo auth logout           # remove o token salvo localmente

# Tasks, Subtasks e Tasks Relacionadas
clickGo tasks                 # lista todas as tasks e subtasks onde você é responsável
clickGo tasks -a              # inclui tasks concluídas/fechadas
clickGo tasks -s "progress"   # filtra por status
clickGo tasks view <id>       # exibe detalhes da task, subtasks e tasks relacionadas
clickGo tasks view <id> -o    # exibe os detalhes e abre a task no navegador
```

## Desenvolvimento

```bash
go build ./...   # compila
go vet ./...     # análise estática
```

## Licença

MIT — veja [LICENSE](LICENSE).
