// Capítulo 00 — Ambiente
//
// Este programa confirma que o ambiente está funcionando e mostra
// informações sobre o runtime do Go que está executando o código.
//
// Rode com: make run CAP=00
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("🚀 Ambiente Go funcionando!")
	fmt.Println()

	// Em JS (Node) o equivalente seria:
	//   process.version, process.platform, process.arch, os.cpus().length
	fmt.Printf("Versão do Go ........: %s\n", runtime.Version())
	fmt.Printf("Sistema operacional .: %s\n", runtime.GOOS)
	fmt.Printf("Arquitetura .........: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs disponíveis ....: %d\n", runtime.NumCPU())
}
