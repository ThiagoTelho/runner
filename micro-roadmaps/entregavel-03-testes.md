# Micro-roadmap — Entregável 3: **Testes**

**Referências:** `especificacao.md` (§7 item 3, §5 critérios de aceitação), `roadmap.md` (Fase 5).

**Critério de conclusão:** testes unitários, de integração, cenários de erro e de aceitação baseados nos critérios definidos, executáveis em CI quando possível.

---

## 1. Estratégia e infraestrutura

- [ ] Escolher frameworks por stack (ex.: JUnit/Mockito para Java; framework nativo do CLI).
- [ ] Configurar execução em CI (build + test em push/PR).
- [ ] Política de dados de teste (fixtures, arquivos temporários, portas efêmeras).

## 2. Testes unitários

### `assinador.jar`

- [ ] Validador de parâmetros: casos válidos e inválidos por campo.
- [ ] Simulação de criação: resposta esperada para conjunto representativo de entradas.
- [ ] Simulação de validação: ramos válido/inválido conforme regra definida.
- [ ] Utilitários PKCS#11 / stubs quando hardware não estiver disponível.

### CLI `assinatura`

- [ ] Parsing de argumentos e help.
- [ ] Lógica de escolha **local vs HTTP** (matriz de flags/defaults).
- [ ] Formatação de saída e mapeamento de códigos de erro.

### CLI `simulador`

- [ ] Lógica de cache “já tenho a versão mais recente”.
- [ ] Verificação de portas (mock de socket/bind).

## 3. Testes de integração

- [ ] `assinatura` → `assinador.jar` em **modo local** (subprocesso real ou container de teste).
- [ ] `assinatura` → **Assinador HTTP** (subir servidor em teste, cliente real).
- [ ] Fluxo **JDK ausente** simulado (mock de download ou ambiente controlado).
- [ ] Fluxo **download do simulador.jar** (mock de API GitHub Releases + happy path com fixture).

## 4. Cenários de erro

- [ ] Parâmetros inválidos em todas as camadas (CLI e JAR).
- [ ] Porta ocupada ao subir Assinador ou Simulador.
- [ ] Falha de rede / timeout no download (JDK ou simulador).
- [ ] Processo filho morto inesperadamente; servidor HTTP indisponível.
- [ ] Respostas malformadas do servidor HTTP (quando aplicável).

## 5. Testes de aceitação

- [ ] Matriz **US-01**: checklist da spec transformado em cenários automatizados ou BDD onde couber.
- [ ] Matriz **US-02**: validação + simulações + PKCS#11 conforme escopo.
- [ ] Matriz **US-03**: start/stop/status, portas, download condicional.
- [ ] Matriz **US-04**: detecção de JDK, download nas três plataformas (mínimo: job matrix no CI ou documentação de execução manual + evidências).
- [ ] Matriz **US-05**: verificação de presença/nomenclatura de binários e checksums em pipeline de release (cruza com Entregável 6).

## 6. Critérios de qualidade

- [ ] Definir meta mínima de cobertura ou “cobertura por áreas críticas” (validação + integração).
- [ ] Falha de teste bloqueia merge/release.
- [ ] Relatórios de teste anexados ou visíveis no CI.

---

## Dependências típicas

- Entregáveis 1 e 2 com APIs estáveis o suficiente para integração.
- Entregável 6: job que valide checksums e presença de `.sig`/`.pem` após release (aceitação US-05 + §9).
