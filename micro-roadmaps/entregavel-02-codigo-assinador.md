# Micro-roadmap — Entregável 2: Código-fonte da aplicação **assinador.jar**

**Referências:** `especificacao.md` (§4.1, §5 US-02, §6, §8.1–8.2, §10 links FHIR), `roadmap.md` (Fase 1).

**Critério de conclusão:** Java com validação completa de parâmetros, simulação das operações (sem criptografia real), tratamento de erros claro.

---

## 1. Projeto Java

- [ ] Estrutura Maven ou Gradle (escolha única para o repositório).
- [ ] Configuração de `jar` executável (manifest / `Main-Class`).
- [ ] Versionamento alinhado ao **SemVer** do projeto.

## 2. Interface CLI do JAR

- [ ] Parsing de argumentos para operações **criar assinatura** e **validar assinatura** (paridade com o que o CLI `assinatura` envia).
- [ ] Modo de execução **one-shot** adequado à invocação local.

## 3. Validação rigorosa de parâmetros (foco principal)

- [ ] Mapear todos os parâmetros exigidos pelas **especificações FHIR** referenciadas na spec.
- [ ] Validar tipos, obrigatoriedade, formatos e combinações inválidas.
- [ ] Mensagens de erro **específicas** (qual parâmetro falhou e por quê).
- [ ] Testes unitários dedicados à camada de validação (cobre também Entregável 3).

## 4. Simulação de operações

### Criação

- [ ] Para entradas válidas, retornar **assinatura simulada pré-construída** (exemplos estáveis e documentados).

### Validação

- [ ] Para entradas válidas, retornar **resultado simulado pré-determinado** (regra simples documentada, ex.: prefixo ou flag de teste).

## 5. PKCS#11

- [ ] Expor fluxo que **interaja com a interface PKCS#11** (token/smart card) conforme critérios de aceitação.
- [ ] Manter escopo: **sem** assinatura digital real — apenas encaixe arquitetural e validação de parâmetros de dispositivo quando aplicável.

## 6. Modo servidor HTTP

- [ ] Servidor HTTP escutando porta configurável (default alinhado ao CLI).
- [ ] Endpoints (ou contrato único) espelhando as operações de criar/validar.
- [ ] Serialização de erros de validação de forma consumível pelo CLI.
- [ ] Encerramento gracioso (SIGTERM / endpoint de shutdown se definido na arquitetura).

## 7. Tratamento de exceções e robustez

- [ ] Capturar falhas de I/O, parsing HTTP, PKCS#11 não disponível, etc.
- [ ] Não vazar stack traces brutos ao usuário final; log interno opcional para depuração.
- [ ] Códigos de saída CLI e códigos HTTP consistentes com a gravidade do erro.

## 8. Documentação no código

- [ ] Documentar contrato público (parâmetros, exemplos de chamada local e HTTP).
- [ ] Comentar limitações explícitas: simulação, sem AC, sem persistência.

## 9. Artefato `assinador.jar`

- [ ] Pipeline de build reproduzível gerando o JAR.
- [ ] Publicar/incluir o JAR no repositório ou como artefato de release conforme decisão do grupo (compatível com Entregável 6).

---

## Dependências típicas

- Referências FHIR (URLs na spec) como fonte da verdade para parâmetros.
- Alinhamento com Entregável 1 nos formatos de chamada local e HTTP.
