package main

import (
  "fmt"
  "math/rand"
  "time"
  "sync"
)

// Conta do Usuario
type Account struct {
  Nome string
  Saldo int
}

// Funcao de Saque direto da conta
func Saque(conta *Account, valor int, signal chan bool, mutex *sync.Mutex) {
  // Se o saldo for maior que o valor requrido, realizar o saque
  mutex.Lock()
  if conta.Saldo >= valor {
    PrintOp(conta, "Saque", valor)
    conta.Saldo = conta.Saldo - valor
    PrintSaldo(conta)
  }
  mutex.Unlock()
  signal <- true
}

// Funcao de Deposito direto em conta
func Deposito(conta *Account, valor int, signal chan bool, mutex *sync.Mutex) {
  mutex.Lock()
  PrintOp(conta, "Deposito", valor)
  conta.Saldo = conta.Saldo + valor
  PrintSaldo(conta)
  mutex.Unlock()
  signal <- true
}

// Printa a operacao realizada em conta
func PrintOp(conta *Account, op string, valor int) {
  fmt.Println("-------------------------------------------------")
  fmt.Println(op, " --> Conta: ", conta.Nome, " valor: ", valor)
  fmt.Println("-------------------------------------------------")
}

// Printa o Saldo da conta
func PrintSaldo(conta *Account) {
  fmt.Println("-------------------------------------------------")
  fmt.Println("Conta: ", conta.Nome, " --> Saldo: ", conta.Saldo)
  fmt.Println("-------------------------------------------------")
}

var mutex sync.Mutex

func main() {
  rand.Seed(time.Now().UnixNano())
  // sinal de termino
  exitSignal := make(chan bool)
  // Numero de conjunto de operacoes saque/deposito
  transcations := 5
  // Conta teste
  minhaConta := Account{Nome: "Thiago Boeker", Saldo: 500}

  for i := 0; i < transcations; i++ {
    go Deposito(&minhaConta, 500, exitSignal, &mutex)
    go Saque(&minhaConta,500, exitSignal, &mutex)
  }

  for i := 0; i < 2 * transcations; i++ {
    <-exitSignal
  }
}
