// Capítulo 01 — Sintaxe essencial
//
// Rode com: make run CAP=01
//
// Arquivos deste pacote:
//
//	variaveis.go   → variáveis, tipos, conversões, constantes
//	controle.go    → if, switch, for
//	funcoes.go     → funções, (valor, error), closures, defer
//	formatacao.go  → verbos do fmt
//
// Todos os arquivos de uma pasta formam UM pacote e enxergam tudo um do outro,
// sem import. Para o Go, o pacote é a pasta, não o arquivo.
package main

import "fmt"

func main() {
	secao("Variáveis e tipos", exemploVariaveis)
	secao("Conversões", exemploConversoes)
	secao("Constantes", exemploConstantes)
	secao("if e switch", exemploCondicionais)
	secao("for", exemploLoops)
	secao("Funções", exemploFuncoes)
	secao("Closures", exemploClosures)
	secao("defer", exemploDefer)
	secao("fmt", exemploFmt)
}

func secao(titulo string, exemplo func()) {
	fmt.Printf("\n━━━ %s ━━━\n", titulo)
	exemplo()
}
