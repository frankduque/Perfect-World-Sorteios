package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"pwapi/pwapi"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

//Carrega as configurações do arquivo config.yaml e popula o struct AppConfig
//
//A função recebe o caminho do arquivo de configuração e retorna um erro caso ocorra algum problema
//
//Parâmetros:
//	filename: string - Caminho relativo do arquivo de configuração
//
//Retorno:
//	error - Retorna um erro caso ocorra algum problema

func loadConfig(filename string) error {

	//Obtém o caminho absoluto do arquivo de configuração
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	//Abre o arquivo de configuração
	file, err := os.Open(absPath)
	if err != nil {
		return err
	}

	//Fecha o arquivo após o término da função
	defer file.Close()

	//Decodifica o arquivo de configuração e popula o struct AppConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&pwapi.AppConfig); err != nil {
		return fmt.Errorf("error decoding YAML: %v", err)
	}

	return nil
}

//removeUser remove um usuário de um slice de RoleID
//
//A função recebe um slice de RoleID e um índice e retorna um novo slice sem o usuário no índice informado
//
//Parâmetros:
//	slice: []pwapi.RoleID - Slice de RoleID
//	index: int - Índice do usuário a ser removido
//
//Retorno:
//	[]pwapi.RoleID - Retorna um novo slice sem o usuário no índice informado

func removeUser(slice []pwapi.RoleID, index int) []pwapi.RoleID {
	return append(slice[:index], slice[index+1:]...)
}

//removerCaracteresIndesejados remove caracteres indesejados de uma string
//
//A função recebe uma string e retorna uma nova string sem os caracteres indesejados
//
//Parâmetros:
//	nome: string - String a ser tratada
//
//Retorno:
//	string - Retorna uma nova string sem os caracteres indesejados
//
//Observação:
//	Neste caso, internamente o servidor do Perfect World utiliza o caractere "&" como marcador de usuário, muito parecido com o "@" utilizado em redes sociais
//	Para evitar problemas com a formatação da mensagem, é necessário remover este caractere do nome do usuário antes de exibir a notificação internamente no jogo.

func removerCaracteresIndesejados(nome string) string {
	// Define a expressão regular para encontrar os caracteres indesejados
	re := regexp.MustCompile("[&]")

	// Substitui os caracteres indesejados por uma string vazia
	novoNome := re.ReplaceAllString(nome, "")

	return novoNome
}

func main() {
	// Carrega as configurações do arquivo config.yaml
	configerr := loadConfig("config.yaml")
	if configerr != nil {
		fmt.Printf("Erro ao carregar config.yaml: %v\n", configerr)
		return
	}

	// Inicializa a conexão com o banco de dados
	pwapi.InitializeDB()
	defer pwapi.CloseDB()

	// Verifica se o servidor está online
	serverOnline := pwapi.IsServerOnline()
	if !serverOnline {
		fmt.Println("Servidor offline")
		return
	}

	// Busca a lista de usuários online
	onlineList := pwapi.GetOnlineList()

	// Verifica se existem usuários online
	if len(onlineList) == 0 {
		fmt.Println("nenhum usuário online")
		return
	}

	// Inicio do sorteio

	//define a variável para verificar se o usuário é válido
	usuarioValido := false

	//define a variável para armazenar o ID do personagem
	var roleID pwapi.RoleID

	// Exibe a quantidade de usuários online caso o modo debug esteja ativado
	if pwapi.AppConfig.Debug {
		fmt.Printf("Total de usuários online: %d\n", len(onlineList))
		//fmt.Printf("onlineList: %v\n\n", onlineList)
	}

	// Exibe a quantidade de usuários a serem sorteados caso o modo debug esteja ativado
	if pwapi.AppConfig.Debug {
		fmt.Printf("Quantidade de usuários a sortear: %d\n", pwapi.AppConfig.QuantidadeDeSorteados)
	}

	//define a variável para armazenar os dados do personagem
	var roleBase pwapi.RoleBase

	// Realiza o sorteio da quantidade de usuários definida no arquivo de configuração
	for i := 0; i < pwapi.AppConfig.QuantidadeDeSorteados; i++ {

		// Sorteia um usuário aleatório e verifica se ele atende aos critérios de level, cultivo e se é um GM
		// de acordo com as configurações do arquivo de configuração
		// Caso o usuário não atenda aos critérios, ele é removido da lista de usuários online e um novo usuário é sorteado
		// até que um usuário válido seja encontrado

		for !usuarioValido {

			// Exibe o número do sorteio caso o modo debug esteja ativado
			if pwapi.AppConfig.Debug {
				fmt.Printf("\nSorteando usuário %d\n", i+1)
			}

			// Verifica se ainda existem usuários online
			if len(onlineList) == 0 {
				fmt.Println("Nenhum usuário restante")
				return
			}

			// Seleciona um usuário aleatório
			key := rand.Intn(len(onlineList))
			roleID = onlineList[key]

			// Busca os dados do personagem (role)
			role := pwapi.GetRoleStatus(roleID)
			if pwapi.AppConfig.Debug {
				fmt.Printf("Usuário sorteado: %v\n", roleID)
			}

			// Verifica se o personagem possui o level mínimo
			if role.Level < pwapi.AppConfig.LevelMinimo {
				if pwapi.AppConfig.Debug {
					fmt.Print("Personagem não possui o level mínimo\n\n")
				}
				// Remove o usuário da lista de usuários online
				onlineList = removeUser(onlineList, key)
				continue
			}

			// Verifica se o personagem possui o cultivo mínimo
			if role.Level2 < pwapi.AppConfig.CultivoMinimo {
				if pwapi.AppConfig.Debug {
					fmt.Print("Personagem não possui o cultivo mínimo\n\n")
				}
				// Remove o usuário da lista de usuários online
				onlineList = removeUser(onlineList, key)
				continue
			}

			//Busca o nome do personagem
			roleBase = pwapi.GetRoleBase(roleID)

			// Verifica se o usuário é um gm
			ehGm := pwapi.UsuarioEGM(roleBase.UserID)
			if !pwapi.AppConfig.GmReceber && ehGm {
				if pwapi.AppConfig.Debug {
					fmt.Print("Usuário é um GM\n\n")
				}
				// Remove o usuário da lista de usuários online
				onlineList = removeUser(onlineList, key)
				continue
			}

			// Se o usuário atender a todos os critérios, a variável usuarioValido é setada como true
			// e o loop é encerrado
			usuarioValido = true

			//remove o usuário sorteado do proximo sorteio
			onlineList = removeUser(onlineList, key)
		}

		// ajuste para o próximo sorteio
		usuarioValido = false

		// Cria um slice de Sorteio com as moedas e golds a serem sorteados
		var Sorteio []pwapi.Sorteio

		// Adiciona as moedas ao sorteio
		for _, moeda := range pwapi.AppConfig.Moedas {
			Sorteio = append(Sorteio, pwapi.Sorteio{
				Tipo:       "moedas",
				Nome:       "Moedas",
				Quantidade: moeda,
			})
		}

		// Adiciona os golds ao sorteio
		for _, gold := range pwapi.AppConfig.Golds {
			Sorteio = append(Sorteio, pwapi.Sorteio{
				Tipo:       "gold",
				Nome:       "Gold",
				Quantidade: gold,
			})
		}

		// Adiciona os itens ao sorteio
		for _, item := range pwapi.AppConfig.ItensSortear {
			//convert item.Data from string to []byte
			stringData := item.Data
			Octets, err := hex.DecodeString(stringData)
			if err != nil {
				fmt.Println("Erro:", err)
				return
			}

			Sorteio = append(Sorteio, pwapi.Sorteio{
				Tipo:       "item",
				Nome:       item.Nome,
				Quantidade: item.Count,
				Item: pwapi.Item{
					ID:         item.ID,
					Pos:        item.Pos,
					Count:      item.Count,
					MaxCount:   item.MaxCount,
					Data:       Octets,
					ProcType:   item.ProcType,
					ExpireDate: item.ExpireDate,
					GUID1:      item.GUID1,
					GUID2:      item.GUID2,
					Mask:       item.Mask,
				},
			})
		}

		// Sorteia um item do slice Sorteio
		key := rand.Intn(len(Sorteio))
		Sorteado := Sorteio[key]

		// Exibe o item sorteado caso o modo debug esteja ativado
		if pwapi.AppConfig.Debug {
			fmt.Println("Item sorteado:", Sorteado)
		}

		// Abrir ou criar o arquivo de log
		arquivoLog, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Erro ao abrir o arquivo de log:", err)
		}
		defer arquivoLog.Close()

		// Configurar o logger para escrever no arquivo
		log.SetOutput(arquivoLog)

		// instancia a variável mensagem para exibir no chat do jogo e no log
		var mensagem string

		// Remove os caracteres indesejados do nome do personagem
		roleName := removerCaracteresIndesejados(roleBase.Name)

		if Sorteado.Tipo == "moedas" {
			// prepara a mensagem para exibir no chat do jogo e no log
			mensagem = fmt.Sprintf("^ffffffO jogador &%s& acabou de ganhar ^33cc33 %d Moedas", roleName, Sorteado.Quantidade)

			// Adiciona as moedas ao personagem
			pwapi.SendMail(roleID, "Logue e ganhe", "Parabens, você ganhou moedas no logue e ganhe", pwapi.Item{}, Sorteado.Quantidade)
		}

		if Sorteado.Tipo == "gold" {
			// prepara a mensagem para exibir no chat do jogo e no log
			mensagem = fmt.Sprintf("^ffffffO jogador &%s& acabou de ganhar ^33cc33 %d Golds", roleName, Sorteado.Quantidade)

			// Adiciona os golds ao personagem
			pwapi.AddCash(roleBase.UserID, Sorteado.Quantidade)

		}

		if Sorteado.Tipo == "item" {
			// prepara a mensagem para exibir no chat do jogo e no log
			mensagem = fmt.Sprintf("^ffffffO jogador &%s& acabou de ganhar ^33cc33 %d %s", roleName, Sorteado.Quantidade, Sorteado.Nome)

			// Adiciona o item ao personagem
			pwapi.SendMail(roleID, "Logue e ganhe", "Parabens, você ganhou um item no logue e ganhe", Sorteado.Item, 0)

		}

		// Exibe a mensagem no chat do jogo
		pwapi.ChatItem(mensagem)
		log.Println(mensagem)

	}

}
