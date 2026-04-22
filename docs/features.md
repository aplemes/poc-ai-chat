# Demand Assistant — Guia de Funcionalidades

## Visão geral

O Demand Assistant é um chat com IA integrado ao formulário de demandas da Adeo. O assistente conduz uma conversa em linguagem natural, coleta as informações necessárias e preenche o formulário automaticamente — eliminando a necessidade de saber quais IDs e opções usar em cada campo.

---

## Funcionalidade 1 — Confirmação antes de preencher

### O que é

Antes de preencher o formulário, a IA apresenta um **card de revisão** com todos os valores que ela propõe. O usuário revisa e decide confirmar ou corrigir.

### Fluxo

```
Conversa com IA → IA coleta todos os campos → Card de revisão aparece no chat
                                                        ↓
                                          [Confirmar e preencher] → Formulário preenchido
                                          [Corrigir pelo chat]    → Usuário digita correção
```

### Como usar para extrair o máximo

- **Revise os campos âmbares primeiro** — são os que a IA inferiu do contexto e têm maior chance de erro (ver Feature 3)
- **Use o card como checklist** — cada campo é listado com seu nome e valor proposto antes de qualquer dado ir para o formulário
- **Clique em "Corrigir pelo chat"** se qualquer valor estiver errado — o cursor vai direto para a caixa de texto

### Exemplo de card de revisão

```
┌─────────────────────────────────────────────────┐
│ Revisar antes de preencher                      │
│                                                 │
│ title          Add payment method X on website  │
│ businessLine   Digital Data Tech                │
│ requesterBU    Leroy Merlin Brazil               │
│ whyDemand      Clientes não conseguem pagar...  │
│ whoIsImpacted  Clientes online: ~50.000 ❓       │  ← inferido
│ ...                                             │
│                                                 │
│ [Confirmar e preencher]  [Corrigir pelo chat]   │
└─────────────────────────────────────────────────┘
```

---

## Funcionalidade 2 — Correção via chat após preenchimento

### O que é

Após o formulário ser preenchido, o chat continua ativo como canal de correção. O usuário pode pedir ajustes diretamente pela conversa sem editar campo a campo manualmente.

### Fluxo

```
Formulário preenchido → IA exibe hint: "Precisa ajustar algo? É só me dizer pelo chat."
                                                        ↓
                               Usuário digita correção → IA aplica e preenche novamente
```

### Como usar para extrair o máximo

- **Mencione o campo e o novo valor explicitamente**
  - Bom: *"Muda o título para 'Enable loyalty program on the app'"*
  - Bom: *"A BU solicitante é Leroy Merlin Spain, não Brazil"*
  - Evite: *"Muda o primeiro campo"* — a IA pode não saber qual campo você quer dizer

- **Agrupe correções em uma mensagem**
  - *"Muda o título e o benefitCategory para Innovation"*
  - Isso evita múltiplas rodadas de confirmação

- **O histórico é preservado** — você não precisa repetir informações já dadas. A IA sabe tudo o que foi conversado antes.

- **Quantas correções quiser** — cada mensagem de correção atualiza o formulário com os campos ajustados, mantendo os demais intactos.

### Exemplo de conversa de correção

```
IA:       Formulário preenchido! Revise os campos.
          Precisa ajustar algo? É só me dizer pelo chat.

Usuário:  Na verdade a BU interessada é Leroy Merlin Poland, não Brazil.
          E muda o título para "Integrate payment gateway X into checkout"

IA:       [preenche formulário com as duas correções]
          Feito! Alguma outra alteração?
```

---

## Funcionalidade 3 — Transparência sobre incerteza

### O que é

A IA distingue campos que o usuário **declarou explicitamente** de campos que ela **inferiu do contexto**. Campos inferidos recebem indicadores visuais de incerteza para que o usuário saiba onde revisar com mais atenção.

### Indicadores visuais

| Local | Campo confirmado | Campo inferido |
|---|---|---|
| Card de revisão (chat) | — | Badge **❓** com tooltip "Inferred from context" |
| Label no formulário | Badge verde **AI** | Badge âmbar **AI ?** |

### Quando um campo é marcado como incerto?

A IA marca como incerto quando ela **deduziu** o valor, e não quando o usuário o disse diretamente:

| Situação | Resultado |
|---|---|
| Usuário diz *"Sou da Leroy Merlin Brazil"* | `requesterBU` → verde AI |
| Usuário diz *"Sou do time de e-commerce do Brasil"* — IA infere LM Brazil | `requesterBU` → âmbar AI ? |
| Usuário diz *"É urgente por questão legal"* | `timeSensitive` → verde AI |
| Usuário descreve urgência sem mencionar legal/security — IA infere | `timeSensitive` → âmbar AI ? |

### Como usar para extrair o máximo

- **Foque a revisão nos badges âmbares** — são os campos com maior probabilidade de erro
- **Para eliminar incerteza, seja explícito no chat** — quanto mais direto, menos a IA precisa inferir
  - Em vez de: *"Trabalho no Brasil com produtos de construção"*
  - Prefira: *"Sou da BU Leroy Merlin Brazil"*
- **Na confirmação, verifique os valores âmbares antes de clicar em Confirmar** — é o momento mais fácil para corrigir

---

## Fluxo completo recomendado

```
1. Abra o chat clicando no FAB (botão flutuante verde)

2. Descreva sua demanda livremente — a IA vai fazer perguntas

3. Responda as perguntas da IA até ela ter contexto suficiente
   → Dica: quanto mais detalhes você der logo no início, menos perguntas a IA fará

4. O card de revisão aparece no chat
   → Verifique os campos, especialmente os marcados com ❓
   → Se algo estiver errado: clique em "Corrigir pelo chat" e descreva a mudança
   → Se estiver tudo certo: clique em "Confirmar e preencher"

5. O formulário é preenchido automaticamente
   → Campos preenchidos pela IA têm badge "AI" (verde) ou "AI ?" (âmbar)
   → Você pode editar qualquer campo manualmente — o badge desaparece ao editar

6. Precisa ajustar algo depois? Basta digitar no chat
   → A IA corrige e preenche novamente sem precisar repetir toda a conversa
```

---

## Dicas gerais

| Dica | Por quê |
|---|---|
| Mencione sua BU e business line cedo na conversa | A IA não precisa perguntar e marca como confirmado (verde) |
| Use nomes completos das BUs | Reduz inferências e elimina badges âmbares |
| Descreva o problema de negócio antes de mencionar a solução | Ajuda a IA a preencher `whyDemand` com mais precisão |
| Separe contexto, situação atual e dor em partes distintas | O campo `whyDemand` cobre os três — quanto mais claro, melhor |
| Dê estimativas de impacto de usuários | A IA usa isso para `whoIsImpacted` sem precisar inferir |

---

## Referência rápida de campos

| Campo | O que a IA precisa saber |
|---|---|
| `title` | O que você quer fazer + onde |
| `businessLine` | Área organizacional responsável |
| `requesterBU` | Qual BU está fazendo o pedido |
| `busInterested` | Quais outras BUs já concordaram com a demanda |
| `timeSensitive` | Há urgência legal ou de segurança? |
| `whyDemand` | Situação atual + ferramentas usadas + dores + contexto motivador |
| `whoIsImpacted` | Tipos de usuários afetados + estimativa de quantidade |
| `benefitCategory` | Qual tipo de benefício esperado |
| `benefitHypothesis` | Como você acredita que a demanda vai gerar esse benefício |
| `measureBenefits` | KPIs e prazo para medir o resultado |
