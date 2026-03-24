# Micro-roadmap — Entregável 5: **Especificação** (documento de requisitos e arquitetura)

**Referências:** `especificacao.md` (§7 item 5, §10), `roadmap.md` (Fase 0).

**Critério de conclusão:** contexto e escopo definidos, diagramas C4, requisitos documentados e **alinhados à implementação** ao longo do projeto.

---

## 1. Contexto e escopo

- [ ] Manter visão geral, objetivos e **objetivos específicos** refletindo o produto final.
- [ ] Revisar §4.1 / §4.2 (dentro/fora do escopo) se houver mudanças aprovadas pelo orientador.
- [ ] Registrar decisões de escopo com data e motivo (changelog curto no doc ou seção “Histórico”).

## 2. Requisitos funcionais (user stories)

- [ ] Garantir que **US-01 a US-05** descrevem o comportamento real ou marcar deltas como “planejado”.
- [ ] Atualizar **critérios de aceitação** com checkboxes conforme forem cumpridos (ou manter matriz de rastreabilidade separada).
- [ ] Referências FHIR e links externos validados (sem links quebrados).

## 3. Diagramas C4

### Nível 1 — Contexto

- [ ] Atores (usuário, sistemas externos relevantes: HubSaúde, repositório de releases).
- [ ] Sistema Runner como caixa central e relações nomeadas.

### Nível 2 — Contêineres

- [ ] Contêineres: CLI `assinatura`, CLI `simulador`, `assinador.jar` (modo CLI e modo HTTP se forem visualmente distintos).
- [ ] Protocolos: subprocesso, HTTP, download HTTPS.
- [ ] Atualizar artefatos em `diagramas/` (fonte e export SVG/PNG conforme padrão do curso).

### (Opcional) Níveis 3+

- [ ] Somente se exigido pelo professor: componentes internos do Assinador ou do CLI.

## 4. Integração e fluxos

- [ ] Texto ou diagramas alinhados às seções **6.1**, **6.2** e **6.3** (incluindo tratamento de erros).
- [ ] Documentar portas padrão, formatos de mensagem e códigos relevantes.

## 5. Segurança e distribuição (referência cruzada)

- [ ] Resumo da obrigatoriedade **Cosign** / **Sigstore** (§9) e o que o usuário deve verificar.
- [ ] Apontar para pipeline CI/CD que assina artefatos (Entregável 6).

## 6. Entregáveis e rastreabilidade

- [ ] Tabela **requisito → implementação → teste** (pode ser anexo) para facilitar correção acadêmica.
- [ ] Garantir que nomes de artefatos na spec batem com o que é publicado nas releases.

## 7. Revisão final antes da entrega

- [ ] Leitura por segundo membro da equipe (revisão cruzada).
- [ ] PDF ou formato exigido pela disciplina, se houver.

---

## Dependências típicas

- Evolução dos Entregáveis 1, 2 e 4 (para texto e diagramas refletirem a realidade).
- Feedback do orientador sobre profundidade dos diagramas C4.
