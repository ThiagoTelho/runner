# Micro-roadmap — Entregável 4: **Documentação**

**Referências:** `especificacao.md` (§7 item 4), `roadmap.md` (Fase 6).

**Critério de conclusão:** manual de usuário do `assinatura`, doc técnica de integração, exemplos de uso, guia de instalação — tudo revisado e coerente com o produto entregue.

> **Estado no repositório (2026-04):** `especificacao.md` e `AGENTS.md` cobrem contexto e stack; ainda não há manuais de usuário por CLI nem guia de instalação/release como documentos dedicados.

---

## 1. Manual de usuário — CLI **assinatura**

- [ ] Visão geral do que o Runner faz e o que **não** faz (simulação, sem AC real).
- [ ] Requisitos: SO suportados, arquitetura **amd64**, espaço em disco, rede (downloads).
- [ ] Instalação por plataforma (referência ao guia de instalação ou seção consolidada).
- [ ] Comandos: criar assinatura, validar assinatura, opções de modo **local** vs **servidor HTTP**.
- [ ] Porta padrão, como mudar porta, como parar o servidor, timeout por inatividade.
- [ ] Mensagens de erro frequentes e como resolver.
- [ ] FAQ curto (JDK automático, onde ficam caches).

## 2. Manual de usuário — CLI **simulador**

- [ ] Comandos: iniciar, parar, status.
- [ ] Comportamento do download (GitHub Releases da disciplina) e cache.
- [ ] Verificação de portas e o que fazer se estiverem ocupadas.
- [ ] Coerência visual e terminológica com o manual do `assinatura`.

## 3. Documentação técnica da integração

- [ ] Diagrama ou descrição textual dos fluxos §6.1 e §6.2 (criar / validar).
- [ ] Contrato **local**: argumentos do `assinador.jar`, exemplos de linha de comando.
- [ ] Contrato **HTTP**: endpoints, métodos, corpos, códigos de status, formato de erro.
- [ ] Variáveis de ambiente, arquivos de configuração, prioridade de resolução do `java`/JDK.
- [ ] Estratégia de versionamento entre CLI e JAR (compatibilidade).

## 4. Exemplos de uso

- [ ] Scripts ou snippets copy-paste para cada plataforma (PowerShell, bash, zsh).
- [ ] Exemplo “primeira execução” com JDK baixado automaticamente.
- [ ] Exemplo fluxo completo: subir simulador + assinar + validar (simulado).
- [ ] Exemplo verificação de binário com **Cosign** (ligação com §9 e Entregável 6).

## 5. Guia de instalação

- [ ] Onde baixar releases oficiais (GitHub Releases).
- [ ] Verificação de **SHA256** dos artefatos.
- [ ] Verificação opcional/recomendada com **Cosign** (`verify-blob`).
- [ ] Instalação no PATH ou local fixo; permissões necessárias.

## 6. Revisão e publicação

- [ ] Ortografia e consistência de termos (Assinador, Runner, HubSaúde).
- [ ] Links internos (README → docs) e links externos (FHIR, C4) atualizados.
- [ ] Versão da documentação alinhada à versão do software (SemVer).

---

## Dependências típicas

- Produto estável o suficiente para não documentar comportamento provisório incorreto.
- Entregável 5: diagramas C4 e spec atualizados como fonte para a doc técnica.
- Entregável 6: nomes exatos dos artefatos e processo de release para o guia de instalação.
