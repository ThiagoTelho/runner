# Contexto do projeto — Sistema Runner

Instruções para agentes (e colaboradores) que trabalham neste repositório. A fonte de verdade dos requisitos é `especificacao.md`.

## O que é

- Trabalho prático de **Implementação e Integração** (Eng. Software, 2026), alinhado ao ecossistema **HubSaúde** (SES-GO / UFG).
- Objetivo: facilitar execução de **aplicações Java** por **linha de comando**, sem exigir que o usuário configure Java manualmente.
- O produto é **CLI-only** (sem GUI). Integradores usam o Runner para assinatura **simulada** e para gerenciar o ciclo de vida do simulador HubSaúde.

## Componentes principais

- **`assinatura`**: CLI multiplataforma (Windows, Linux, macOS, **amd64**) que fala com o Assinador.
- **`assinador.jar`**: aplicação **Java** que valida parâmetros com rigor, **simula** criar/validar assinatura, suporta **PKCS#11** no escopo previsto, e pode rodar em modo **local** (subprocesso) ou **servidor HTTP** (warm start).
- **`simulador`**: CLI multiplataforma que **baixa** o `simulador.jar` mais recente (GitHub Releases da disciplina), faz cache local, inicia/para/mostra **status** e checa **portas** antes de subir.
- **JDK**: o sistema deve **detectar** JDK na versão exigida e, se faltar, **baixar** e usar de forma consistente nas três plataformas (US-04).

## Integração Assinador ↔ CLI `assinatura`

- **Local**: `assinatura` executa `java -jar assinador.jar …` (cold start).
- **HTTP**: `assinador.jar` como servidor; `assinatura` envia requisições (warm start).
- Padrão esperado: **preferir modo servidor** quando o usuário não forçar local; **detectar** instância já na porta padrão; permitir **parar** na porta padrão ou indicada; **encerrar após N minutos** sem interação quando solicitado (US-01).
- Erros devem ser capturados, propagados de forma compreensível e mostrados ao usuário com informação suficiente para correção (`especificacao.md` §6.3).

## Escopo: o que NÃO fazer

- Nenhuma assinatura ou validação **criptográfica real**; sem integração com **AC**; sem persistência de assinaturas; sem **GUI**; sem autenticação de usuários do Runner; sem geração de certificados (`especificacao.md` §4.2).
- O **`simulador.jar`** em si **não é desenvolvido** pelo Sistema Runner; só é **obtido** e **orquestrado** pelo CLI `simulador` (US-03). A spec também lista entregável de **código-fonte do Simulador HubSaúde** (§7 item 7): pode ser exigência acadêmica separada — ver `micro-roadmaps/entregavel-07-codigo-simulador-hubsaude.md` e alinhar com o orientador.

## Distribuição e segurança

- Binários pré-compilados para **Windows (.exe), Linux (.AppImage), macOS (.dmg)**, **amd64**, para `assinatura` e `simulador`, via **GitHub Releases**, **SemVer**, **SHA256** por artefato.
- Todos os artefatos de release devem ser assinados com **Cosign** (OIDC / Sigstore), com `<artefato>`, `<artefato>.sig`, `<artefato>.pem`, de preferência **no CI/CD** (`especificacao.md` §7 e §9).

## Referências de parâmetros e arquitetura

- Parâmetros de assinatura alinhados às páginas FHIR citadas em `especificacao.md` §10.
- Diagramas **C4** (contexto e contêineres): paths indicados na spec (`diagramas/`).

## Documentação interna deste repositório

- `especificacao.md`: requisitos, user stories, entregáveis.
- `roadmap.md`: fases e macro tarefas.
- `micro-roadmaps/`: detalhamento por entregável (checklists).

## Orientação para implementação (agentes)

- Ler trechos relevantes de `especificacao.md` antes de alterar comportamento visível ao usuário ou contratos entre CLI e JAR.
- Priorizar **validação de entrada** e **mensagens de erro úteis** no `assinador.jar` (maior esforço esperado na simulação).
- Manter **portabilidade** (paths, subprocessos, rede) nos CLIs.
- Não expandir escopo para funcionalidades fora da §4.1 sem atualizar a spec e o orientador.
