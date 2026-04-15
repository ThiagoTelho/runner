# Micro-roadmap — Entregável 6: **Artefatos executáveis**

**Referências:** `especificacao.md` (§7 item 6, §9 Cosign), `roadmap.md` (Fase 7).

**Critério de conclusão:** binários pré-compilados para Windows / Linux / macOS (**amd64**) para `assinatura` e `simulador`, publicados em **GitHub Releases**, com **SHA256** e trinca **Cosign** (`.sig`, `.pem`) por artefato, versionamento **SemVer**.

> **Estado no repositório (2026-04):** versões SemVer definidas no código (`0.1.0` nos CLIs); pipeline de release, GitHub Actions e publicação de artefatos **ainda não** constam no repositório.

---

## 1. Nomenclatura e versionamento

- [ ] Adotar padrão: `assinatura-<semver>-<os>-amd64.<ext>` e `simulador-<semver>-<os>-amd64.<ext>`.
- [ ] Extensões: `.exe` (Windows), `.AppImage` (Linux), `.dmg` (macOS).
- [ ] Processo de bump de versão (tags git ↔ SemVer).

## 2. Build multiplataforma

### CLI `assinatura`

- [ ] Pipeline que produza os três binários a partir do mesmo commit.
- [ ] Teste de smoke pós-build (executar `--version` ou equivalente).

### CLI `simulador`

- [ ] Idem: três binários com mesma política de versionamento.
- [ ] Garantir que o binário embute ou resolve corretamente o download do `simulador.jar` em runtime.

## 3. Conteúdo da release no GitHub

- [ ] Upload de todos os binários listados na spec.
- [ ] Arquivo(s) de **checksums SHA256** (ex.: `SHA256SUMS`) cobrindo cada artefato.
- [ ] Release notes: mudanças da versão, instruções breves de verificação.

## 4. Assinatura com Cosign (obrigatório)

Para **cada** artefato publicado:

- [ ] `<artefato>`
- [ ] `<artefato>.sig`
- [ ] `<artefato>.pem`

Tarefas:

- [ ] Configurar identidade **OIDC** no GitHub Actions (ou provedor aceito) para Cosign.
- [ ] Registrar assinatura no transparency log do **Sigstore**, conforme §9.
- [ ] Job de CI/CD que falha se algum arquivo da trinca estiver ausente.
- [ ] Documentar comando `cosign verify-blob` (cruza com Entregável 4).

## 5. Verificação e segurança

- [ ] Job opcional que baixa a release e verifica checksums + Cosign em ambiente limpo.
- [ ] Proteção de branches/tags conforme política do repositório (evitar sobrescrita de release).

## 6. Inclusão do `assinador.jar`

- [ ] Decisão: JAR na mesma release, release separada ou apenas no código-fonte — **documentar** na spec/README.
- [ ] Se o JAR for distribuído: aplicar mesma política de checksum e Cosign, se exigido pelo curso.

## 7. Checklist final de publicação

- [ ] Todos os binários baixáveis e com tamanho esperado.
- [ ] SemVer da tag = versão nos nomes dos arquivos.
- [ ] Links da release comunicados na documentação.

---

## Dependências típicas

- Entregáveis 1 e 2 com código estável e testes passando (Entregável 3).
- Conta e permissões no GitHub para OIDC/Cosign.
- Entregável 4 atualizado com instruções de verificação para integradores.
