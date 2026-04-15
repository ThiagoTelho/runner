# Micro-roadmap — Entregável 1: Código-fonte da aplicação **assinatura**

**Referências:** `especificacao.md` (§4.1, §5 US-01, US-04, §6, §8.2), `roadmap.md` (Fases 0, 2, 3).

**Critério de conclusão:** implementação completa, compatível com Windows / Linux / macOS, código documentado (comentários e estrutura onde agregar valor).

> **Estado no repositório (2026-04):** comandos `criar` e `validar` invocam o Assinador via subprocesso local (`java -jar`); resolução de `java` (`JAVA_HOME` ou PATH) e do JAR (`RUNNER_ASSINADOR_JAR` ou `assinador.jar` junto ao binário). Modo servidor HTTP, detecção de instância e JDK automático (US-04) ainda pendentes.

---

## 1. Projeto e baseline multiplataforma

- [x] Inicializar módulo/pacote do CLI com convenções de build para **amd64** nas três plataformas.
- [x] Definir estratégia de paths (diretório de instalação, cache, `assinador.jar` embutido ou configurável). *(Feito: `RUNNER_ASSINADOR_JAR` ou `assinador.jar` ao lado do binário; cache/embed ainda não.)*
- [x] Garantir que operações de filesystem e subprocessos sejam portáveis (separadores, quoting, variáveis de ambiente).

## 2. Interface de linha de comandos

- [x] Estrutura de comandos/subcomandos para **criar** e **validar** assinatura (alinhado às referências FHIR da spec).
- [x] Help integrado (`--help` / subcomandos) com mensagens claras e consistentes.
- [ ] Validação inicial de entrada no CLI (antes de invocar o Assinador), quando aplicável.

## 3. Integração com `assinador.jar`

### Modo local (cold start)

- [x] Resolver executável `java` (sistema ou JDK provisionado — ver Entregável 1 × US-04).
- [x] Montar linha de comando `java -jar assinador.jar …` com argumentos corretos.
- [x] Capturar stdout/stderr e códigos de saída; propagar erros de forma legível.

### Modo servidor HTTP (warm start)

- [ ] Cliente HTTP com contrato alinhado ao Assinador em modo servidor.
- [ ] **Política padrão:** usar servidor quando o usuário não forçar modo local.
- [ ] **Porta padrão** do servidor; permitir override explícito.
- [ ] **Detectar** instância já em execução na porta e reutilizar quando não houver instrução contrária.
- [ ] **Subir** o Assinador em modo servidor quando necessário.
- [ ] **Parar** o processo na porta padrão ou na porta indicada.
- [ ] **Encerramento programado** após N minutos sem interação, quando o usuário solicitar.

## 4. Apresentação ao usuário

- [x] Formatar e exibir resultado das operações de forma legível (criação e validação).
- [ ] Tratamento unificado de falhas (rede, processo, timeouts, respostas inválidas).

## 5. Qualidade e documentação no código

- [x] Organização em pacotes/módulos coerentes com a arquitetura (CLI × integração × utilitários).
- [x] Documentação no código nos pontos não óbvios (modo local vs HTTP, resolução de JDK).
- [ ] Revisar mensagens de erro para serem acionáveis (o que corrigir e como).

## 6. Verificação antes do fechamento do entregável

- [ ] Smoke tests manuais nas três plataformas (ou em CI) para fluxos principais.
- [ ] Checklist US-01 e critérios US-04 relacionados ao uso do Java pelo CLI.

---

## Dependências típicas

- Contrato mínimo do **assinador.jar** (Entregável 2) para integração estável.
- Provisionamento de JDK (micro-roadmap embutido neste entregável via US-04; detalhado em tarefas da seção 3).
