package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Função para processar uma parte dos números
func processarNumeros(id int, nums []int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	soma := 0
	for _, num := range nums {
		soma += num
	}

	// Envia o resultado de volta através do channel
	result <- soma
	fmt.Printf("Go-routine %d processou %d números.\n", id, len(nums))
}

func main() {
	// Configura o número total de elementos e a quantidade de goroutines
	const totalElementos = 1_000_000_000
	const numeroGoroutines = 100
	elementosPorGo := totalElementos / numeroGoroutines

	// Gerar números aleatórios para processamento
	fmt.Println("Gerando números...")
	numeros := make([]int, totalElementos)
	for i := range numeros {
		numeros[i] = rand.Intn(100) // Simula dados de entrada
	}

	// Channel para coletar resultados
	result := make(chan int, numeroGoroutines)
	var wg sync.WaitGroup

	// Divide o trabalho entre várias go-routines
	fmt.Println("Iniciando processamento...")
	start := time.Now()

	for i := 0; i < numeroGoroutines; i++ {
		inicio := i * elementosPorGo
		fim := inicio + elementosPorGo

		// Adiciona uma goroutine ao WaitGroup
		wg.Add(1)

		// Inicia uma goroutine para processar cada pedaço dos dados
		go processarNumeros(i+1, numeros[inicio:fim], result, &wg)
	}

	// Goroutine para esperar todas as outras completarem e fechar o channel
	go func() {
		wg.Wait()
		close(result)
	}()

	// Coleta os resultados do channel
	total := 0
	for soma := range result {
		total += soma
	}

	fmt.Printf("Processamento concluído em %v segundos.\n", time.Since(start).Seconds())
	fmt.Printf("A soma total dos números é: %d\n", total)
}
