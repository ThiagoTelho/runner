# Micro-roadmap — Entregável 2: Código-fonte da aplicação **assinador.jar**

**Referências:** `especificacao.md` (§4.1, §5 US-02, §6, §8.1–8.2, §10 links FHIR), `roadmap.md` (Fase 1).

**Critério de conclusão:** Java com validação completa de parâmetros, simulação das operações (sem criptografia real), tratamento de erros claro.

> **Estado no repositório (2026-04):** núcleo funcional do Assinador implementado (Gradle, picocli, validadores com testes, simuladores, servidor HTTP). PKCS#11 ainda não há integração com provedor/token real (apenas validação de parâmetros de material). O binário `assinador.jar` não está versionado no Git — apenas build via `./gradlew`.

---

## 1. Projeto Java

- [x] Estrutura Maven ou Gradle (escolha única para o repositório).
- [x] Configuração de `jar` executável (manifest / `Main-Class`).
- [x] Versionamento alinhado ao **SemVer** do projeto.

## 2. Interface CLI do JAR

- [x] Parsing de argumentos para operações **criar assinatura** e **validar assinatura** (paridade com o que o CLI `assinatura` envia).
- [x] Modo de execução **one-shot** adequado à invocação local.

## 3. Validação rigorosa de parâmetros (foco principal)

- [x] Mapear todos os parâmetros exigidos pelas **especificações FHIR** referenciadas na spec.
- [x] Validar tipos, obrigatoriedade, formatos e combinações inválidas.
- [x] Mensagens de erro **específicas** (qual parâmetro falhou e por quê).
- [x] Testes unitários dedicados à camada de validação (cobre também Entregável 3).

## 4. Simulação de operações

### Criação

- [x] Para entradas válidas, retornar **assinatura simulada pré-construída** (exemplos estáveis e documentados).

### Validação

- [x] Para entradas válidas, retornar **resultado simulado pré-determinado** (regra simples documentada, ex.: prefixo ou flag de teste).

## 5. PKCS#11

- [ ] Expor fluxo que **interaja com a interface PKCS#11** (token/smart card) conforme critérios de aceitação.
- [x] Manter escopo: **sem** assinatura digital real — apenas encaixe arquitetural e validação de parâmetros de dispositivo quando aplicável.

## 6. Modo servidor HTTP

- [x] Servidor HTTP escutando porta configurável (default alinhado ao CLI).
- [x] Endpoints (ou contrato único) espelhando as operações de criar/validar.
- [x] Serialização de erros de validação de forma consumível pelo CLI.
- [x] Encerramento gracioso (SIGTERM / endpoint de shutdown se definido na arquitetura).

## 7. Tratamento de exceções e robustez

- [x] Capturar falhas de I/O, parsing HTTP, PKCS#11 não disponível, etc.
- [x] Não vazar stack traces brutos ao usuário final; log interno opcional para depuração.
- [x] Códigos de saída CLI e códigos HTTP consistentes com a gravidade do erro.

## 8. Documentação no código

- [x] Documentar contrato público (parâmetros, exemplos de chamada local e HTTP).
- [x] Comentar limitações explícitas: simulação, sem AC, sem persistência.

## 9. Artefato `assinador.jar`

- [x] Pipeline de build reproduzível gerando o JAR.
- [ ] Publicar/incluir o JAR no repositório ou como artefato de release conforme decisão do grupo (compatível com Entregável 6).

---

## Dependências típicas

- Referências FHIR (URLs na spec) como fonte da verdade para parâmetros.
- Alinhamento com Entregável 1 nos formatos de chamada local e HTTP.
