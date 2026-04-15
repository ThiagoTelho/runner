# Micro-roadmap — Entregável 7: **Código-fonte do Simulador do HubSaúde**

**Referências:** `especificacao.md` (§7 item 7; §5 US-03 — o `simulador.jar` em si não é desenvolvido pelo Runner, mas este entregável pede código-fonte do Simulador), `roadmap.md` (R7.6).

**Critério de conclusão (conforme spec):** implementação completa, código bem documentado, compatível com Windows / Linux / macOS.

> **Nota de alinhamento:** Na mesma especificação, a US-03 afirma que o **Simulador (`simulador.jar`) não faz parte do escopo de desenvolvimento do Sistema Runner** e deve ser obtido via GitHub Releases. O Entregável 7 pode ser exigência acadêmica separada (outro repositório ou módulo). **Confirmar com o orientador** se este entregável é o JAR fornecido pela disciplina (apenas empacotamento) ou um projeto Java independente que vocês implementam.

> **Estado no repositório (2026-04):** não há código-fonte do Simulador HubSaúde neste mono-repo; o CLI `simulador` existe apenas como esqueleto (comandos e pacotes com TODO), sem download/cache/integração implementados.

---

## 0. Esclarecimento de escopo (obrigatório antes de codificar)

- [ ] Confirmar com o professor: o Simulador é **implementado pela equipe** ou **fornecido** e vocês apenas integram?
- [ ] Se for implementação própria: obter requisitos funcionais mínimos (portas, protocolos, compatibilidade HubSaúde).
- [ ] Definir repositório: mesmo mono-repo do Runner ou repositório dedicado com release linkada.

## 1. Projeto base

- [ ] Estrutura Java (ou stack acordada) com build reproduzível.
- [ ] Empacotamento como `simulador.jar` executável.
- [ ] Política de versionamento **SemVer** alinhada às releases usadas pelo CLI (Entregável 1).

## 2. Funcionalidade do simulador

- [ ] Comportamento esperado pelo HubSaúde / ambiente de laboratório (conforme material da disciplina).
- [ ] Configuração de **portas** e validação de conflito (complementa o que o CLI já verifica — US-03).
- [ ] Logs claros para depuração durante integração.

## 3. Multiplataforma

- [ ] Testes manuais ou automatizados em **Windows, Linux e macOS** (amd64).
- [ ] Tratar diferenças de filesystem, encoding e firewall local em documentação.

## 4. Documentação no código e para usuários

- [ ] README do módulo/repositório: como compilar, executar, parâmetros.
- [ ] Javadoc ou equivalente nas APIs públicas relevantes.
- [ ] Exemplos de execução direta com `java -jar simulador.jar`.

## 5. Integração com o Runner

- [ ] Publicar releases no GitHub esperado pelo CLI `simulador` (URL da disciplina ou do grupo).
- [ ] Garantir que o CLI detecte “versão mais recente” conforme metadados da API de Releases.
- [ ] Checklist: download condicional, cache local, start/stop/status (aceitação US-03).

## 6. (Se distribuído como binários)

- [ ] Alinhar com Entregável 6: empacotamento `.exe` / `.AppImage` / `.dmg` do **CLI simulador** que gerencia o JAR — distinto do código deste entregável, mas parte da mesma entrega global.

---

## Dependências típicas

- Requisitos do Simulador definidos pelo corpo docente ou por especificação HubSaúde adicional.
- Entregável 1: CLI que consome o artefato gerado aqui.
