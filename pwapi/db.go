package pwapi

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// UsuarioEGM verifica se o usuário está na tabela auth
//
// Parâmetros:
// 	userID: UserID - ID do usuário
//
// Retorno:
// 	bool - Retorna true se o usuário estiver na tabela auth, false caso contrário
//
//Observações:
//	Auth é o nome escolhido pelos desenvolvedores do Perfect World para a tabela que armazena as permissões especiais de players no jogo.
//  Neste caso a tabela Auth não é utilizada no processo de autenticação, mas sim para armazenar informações de usuários especiais, como administradores.
//	O termo Gm é uma abreviação de Game Master, que é o nome dado aos administradores do jogo.

func UsuarioEGM(userID UserID) bool {
	// Inicializa a conexão com o banco de dados, se ainda não estiver inicializada
	rows, err := db.Query("SELECT DISTINCT userid FROM auth WHERE userid = ?", userID)
	if err != nil {
		fmt.Printf("Erro ao consultar o banco de dados: %v\n", err)
		os.Exit(1)
	}

	//Contador para verificar se o usuário está na tabela auth
	var numRows int

	// Verifica se o usuário está na tabela auth
	for rows.Next() {
		var userID UserID
		err = rows.Scan(&userID)
		if err != nil {
			fmt.Printf("Erro ao consultar o banco de dados: %v\n", err)
			os.Exit(1)
		}
		numRows++
	}

	// Retorna true se o usuário estiver na tabela auth, false caso contrário
	if numRows > 0 {
		return true
	} else {
		return false
	}
}

// InitializeDB inicializa a conexão com o banco de dados
func InitializeDB() {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", AppConfig.MySQL.Usuario, AppConfig.MySQL.Senha, AppConfig.MySQL.Host, AppConfig.MySQL.DB)
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

// CloseDB fecha a conexão com o banco de dados
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
