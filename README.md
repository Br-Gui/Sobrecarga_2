# Go HTTP Load Tester

Um sistema simples e altamente configurável em Go para realizar testes de carga em uma URL, com suporte a milhares de requisições simultâneas por ciclo e geração de relatórios detalhados.

## Funcionalidades

- Envia múltiplas requisições HTTP simultâneas por ciclo.
- Controla o número de goroutines ativas simultaneamente (`maxThreads`).
- Mede tempo de resposta individual de cada requisição.
- Gera relatório por ciclo com:
    - Número de acertos e erros
    - Códigos HTTP retornados
    - Duração média, mínima e máxima (em milissegundos, segundos e minutos)
- Gera um relatório final consolidado em JSON.

## Como usar

### 1. Clone o projeto

```bash
git clone https://github.com/Br-Gui/Sobrecarga_2.git
cd Sobrecarga_2
```

### 2. Execute o projeto

```bash
go run main.go
```

> **Requisitos**: Go 1.18 ou superior instalado

## Configurações principais

Você pode ajustar os parâmetros principais diretamente no `main.go`:

```go
url := "teste.com" // URL a ser testada
numGoroutines := 2000 // Requisições por ciclo
maxThreads := 500 // Limite de goroutines simultâneas
maxCycles := 3000 // Quantidade de ciclos de teste
```

## Exemplo de relatório gerado

Ao final da execução, um arquivo chamado `api_test_detailed_report.json` será salvo com o seguinte conteúdo estruturado:

```json
{
    "total_cycles": 3,
    "total_requests": 6000,
    "success_count": 5972,
    "error_count": 28,
    "avg_duration_ms": "100ms",
    "response_codes": {
        "200": 5972,
        "500": 28
    },
    "cycle_details": [...]
}
```

## Tecnologias usadas

- **Go** — linguagem principal
- `net/http` — para requisições HTTP
- `sync` — controle de concorrência com goroutines
- `encoding/json` — geração de relatórios

## Observações

- Aumentar `numGoroutines` e `maxThreads` pode exigir mais da sua CPU e rede.
- Cuidado ao fazer testes em domínios que não são seus: isso pode ser interpretado como ataque.


---
