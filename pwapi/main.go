package pwapi

import (
	"fmt"
	"net"
	"time"
)

//GetOnlineList retorna uma lista de RoleID de usuários online
//
//Parâmetros:
//	Não há parâmetros
//
//Retorno:
//	[]RoleID - Retorna uma lista de RoleID de usuários online
//
//Observações:
//	RoleID é o ID do personagem no jogo
//	Esta função envia uma requisição para o gdeliveryd para obter a lista de usuários online
//	Esta função é utilizada para obter a lista de usuários online para realizar sorteios
//  As informações utilizadas para escrever esta função foram obtidas através de engenharia reversa realizada por desenvolvedores da comunidade
//  Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/GMQueryOnline

func GetOnlineList() []RoleID {

	//Configuração do pacote GMQueryOnline
	GMQueryOnlinePacket := GMQueryOnline{
		QType: 0,
	}
	pack := createPack(GMQueryOnlinePacket)
	opcode := 0x189
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	// Envio do pacote e tratamento de erros
	recvAfterSend := true
	justSend := false
	response, err := SendToDelivery(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o gdeliveryd: %v\n", err)
	}

	//Deleta o cabeçalho do pacote
	//Ao invés de deletar o cabeçalho, o struct GMQueryOnlineRe foi modificado para evidenciar os dados do cabeçalho
	//data := deleteHeader(response)

	//GMQueryOnlineRe é a estrutura do pacote que será recebido do gdeliveryd
	var usersOnline GMQueryOnlineRe

	//Desempacota os dados do pacote recebido
	unpackData(response, &usersOnline)

	//OnlineList é a lista de RoleID de usuários online que será retornada
	var OnlineList []RoleID
	OnlineList = usersOnline.RoleIDS
	return OnlineList
}

//IsServerOnline verifica se o servidor está online
//
//Parâmetros:
//	Não há parâmetros
//
//Retorno:
//	bool - Retorna true se o servidor estiver online, false caso contrário
//
//Observações:
//	Esta função tenta se conectar ao gamedbd para verificar se o servidor está online
//	O gamedbd é um serviço do servidor do Perfect World que gerencia as informações dos personagens logados

func IsServerOnline() bool {

	//Tenta se conectar ao gamedbd
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", AppConfig.IP, AppConfig.Ports["gamedbd"]), 5*time.Second)
	if err != nil {
		fmt.Printf("Erro ao conectar ao gamedbd: %v\n", err)
		return false
	}

	return true
}

//GetRoleStatus retorna o status de um personagem
//
//Parâmetros:
//	roleID: RoleID - ID do personagem
//
//Retorno:
//	RoleStatus - Retorna o status do personagem
//
//Observações:
//	RoleStatus é a estrutura que contém as informações dos status do personagem
//	Status nesse contexto se refere a informações que se alteram durante o jogo, como nível e cultivo que são utilizadas para verificar se o personagem é elegível para o sorteio
//  Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/GetRoleStatusArg

func GetRoleStatus(roleID RoleID) RoleStatus {

	//configuração do pacote GetRoleStatusArg
	GetRoleStatusArgPacket := GetRoleStatusArg{
		Handler: -1,
		RoleID:  roleID,
	}

	// Cria e prepara o pacote para envio
	pack := createPack(GetRoleStatusArgPacket)
	opcode := 0xbc7
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	// Envio do pacote e tratamento de erros
	recvAfterSend := false
	justSend := false
	response, err := SendToGamedBD(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o gdeliveryd: %v\n", err)
	}

	// Desempacota os dados para obter o RoleStatus
	data := deleteHeader(response)
	var roleStatus RoleStatus
	unpackData(data, &roleStatus)

	return roleStatus
}

//GetRoleBase retorna as informações básicas de um personagem
//
//Parâmetros:
//	roleID: RoleID - ID do personagem
//
//Retorno:
//	RoleBase: Retorna as informações básicas do personagem
//
//Observações:
//	RoleBase é a estrutura que contém as informações básicas do personagem
//	Esta função é utilizada para obter as informações básicas do personagem para realizar sorteios
//  As informações utilizadas para escrever esta função foram obtidas através de engenharia reversa realizada por desenvolvedores da comunidade
//  Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/GetRoleBaseArg

func GetRoleBase(roleID RoleID) RoleBase {

	//Configuração do pacote GetRoleBaseArg
	GetRoleBaseArgPacket := GetRoleBaseArg{
		Handler: -1,
		RoleID:  roleID,
	}

	// Cria e prepara o pacote para envio
	pack := createPack(GetRoleBaseArgPacket)
	opcode := 0xbc5
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	// Envio do pacote e tratamento de erros
	recvAfterSend := false
	justSend := false
	response, err := SendToGamedBD(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o gdeliveryd: %v\n", err)
	}

	//	Deleta o cabeçalho do pacote
	data := deleteHeader(response)

	// Desempacota os dados para obter o RoleBase
	var roleBase RoleBase
	unpackData(data, &roleBase)
	return roleBase
}

// ChatItem envia uma mensagem para o chat do jogo
// o local que a mensagem será enviada é definido pela variável AppConfig.CanalMensagem
//
// Parâmetros:
//
//	text: string - Mensagem a ser enviada
//
// Retorno:
//
//	Não há retorno
//
// Observações:
//
//	Esta função envia uma mensagem para o chat do jogo
//	As informações utilizadas para escrever esta função foram obtidas através de engenharia reversa realizada por desenvolvedores da comunidade
//	Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/ChatBroadCast
func ChatItem(text string) {

	//ChatBroadCastAPI é a estrutura do pacote que será enviado para o gdeliveryd
	ChatBroadCastPacket := ChatBroadCast{
		Channel:   byte(AppConfig.CanalMensagem),
		Emotion:   0,
		SrcRoleID: 0,
		Msg:       text,
		Data:      []byte{},
	}

	// Cria e prepara o pacote para envio
	pack := createPack(ChatBroadCastPacket)
	opcode := 0x78
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	// Envio do pacote e tratamento de erros
	recvAfterSend := true
	justSend := true
	_, err := SendToProvider(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o gdeliveryd: %v\n", err)
	}
}

// AddCash adiciona cash a um personagem
//
// Parâmetros:
// 	userID: UserID - ID do usuário
// 	cash: int - Quantidade de cash a ser adicionada
//
// Retorno:
// 	Não há retorno
//
// Observações:
// 	Cash é um termo mais utilizado em servidores oficiais do Perfect World para se referir a moeda premium
// 	Em servidores privados, o termo mais utilizado é Gold
// 	Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/DebugAddCash

func AddCash(userID UserID, cash int) {

	//DebugAddCash é a estrutura do pacote que será enviado para o gamedbd
	DebugAddCashPacket := DebugAddCash{
		UserID: userID,
		Cash:   cash * 100,
	}

	// Cria e prepara o pacote para envio
	pack := createPack(DebugAddCashPacket)
	opcode := 0x209
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	// Envio do pacote e tratamento de erros
	recvAfterSend := false
	justSend := true
	_, err := SendToGamedBD(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o Gamedbd: %v\n", err)
	}
}

// SendMail envia um e-mail para um personagem
//
// Parâmetros:
//
//	RoleID: RoleID - ID do personagem
//	title: string - Título do e-mail
//	content: string - Conteúdo do e-mail
//	item: Item - Item a ser enviado
//	money: int - Quantidade de dinheiro a ser enviada
//
// Retorno:
//
//	Não há retorno
//
// Observações:
//
//	Esta função envia um e-mail para um personagem dentro do jogo
//	Diferente de mensagens, e-mails podem conter itens e dinheiro
//	Mais informações em sobre o Opcode e detalhes do pacote em: http://pwdev.ru/index.php/SysSendMail
func SendMail(RoleID RoleID, title string, content string, item Item, money int) {

	// Configuração do pacote SysSendMailAPI
	// valores hardcoded definidos pela comunidade
	SysSendMailpacket := SysSendMail{
		TID:         344,
		SysID:       1025,
		SysType:     3,
		Receiver:    RoleID,
		Title:       title,
		Content:     content,
		AttachObj:   item,
		AttachMoney: money,
	}

	//Cria o pacote
	pack := createPack(SysSendMailpacket)
	opcode := 0x1076
	opcodeHex := fmt.Sprintf("%X", opcode)
	pack = createHeader(opcodeHex, pack)

	//envia o pacote e verifica se houve erro
	recvAfterSend := false
	justSend := true
	_, err := SendToDelivery(pack, recvAfterSend, justSend)
	if err != nil {
		fmt.Printf("Erro ao enviar para o gdeliveryd: %v\n", err)
	}
}
