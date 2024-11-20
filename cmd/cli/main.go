package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Definindo os parâmetros de entrada via CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests a serem realizados")
	concurrency := flag.Int("concurrency", 0, "Número de requisições simultâneas")

	flag.Parse()

	// Validação de parâmetros obrigatórios
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		log.Fatal("Por favor, forneça os parâmetros necessários: --url, --requests e --concurrency.")
	}

	// Inicializando as variáveis para o relatório
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalRequests, successfulRequests int
	statusCodes := make(map[int]int)

	// Criação de um canal para controlar as requisições simultâneas
	concurrencyChan := make(chan struct{}, *concurrency)

	// Definindo o tempo de início do teste
	startTime := time.Now()

	// Realizando as requisições em paralelo
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		concurrencyChan <- struct{}{} // Limita a concorrência
		go func() {
			defer wg.Done()
			defer func() { <-concurrencyChan }() // Libera o canal

			// Realizando a requisição HTTP
			resp, err := http.Get(*url)
			if err != nil {
				log.Println("Erro na requisição:", err)
				return
			}
			defer resp.Body.Close()

			// Atualizando o contador de requisições
			mu.Lock()
			totalRequests++
			statusCodes[resp.StatusCode]++
			if resp.StatusCode == 200 {
				successfulRequests++
			}
			mu.Unlock()
		}()
	}

	// Espera todas as goroutines terminarem
	wg.Wait()

	// Calculando o tempo total de execução
	duration := time.Since(startTime)

	// Gerando o relatório
	fmt.Printf("\nRelatório de Teste de Carga\n")
	fmt.Printf("====================================\n")
	fmt.Printf("URL Testada: %s\n", *url)
	fmt.Printf("Total de Requests: %d\n", totalRequests)
	fmt.Printf("Requests com Status 200: %d\n", successfulRequests)
	fmt.Printf("Tempo Total de Execução: %s\n", duration)

	// Exibindo a distribuição dos status codes
	fmt.Println("Distribuição de Status Codes:")
	for code, count := range statusCodes {
		fmt.Printf("%d: %d requisições\n", code, count)
	}
}
