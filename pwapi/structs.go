package pwapi

type Config struct {
	Debug                 bool           `yaml:"Debug"`
	IP                    string         `yaml:"IP"`
	Ports                 map[string]int `yaml:"Ports"`
	MySQL                 MySQLConfig    `yaml:"MySQL"`
	QuantidadeDeSorteados int            `yaml:"QuantidadeDeSorteados"`
	GmReceber             bool           `yaml:"GmReceber"`
	LevelMinimo           int            `yaml:"LevelMinimo"`
	CultivoMinimo         int            `yaml:"CultivoMinimo"`
	CanalMensagem         int            `yaml:"CanalMensagem"`
	Moedas                []int          `yaml:"Moedas"`
	Golds                 []int          `yaml:"Golds"`
	ItensSortear          []ItemNome     `yaml:"ItensSortear"`
}

type MySQLConfig struct {
	Host    string `yaml:"Host"`
	Usuario string `yaml:"Usuario"`
	Senha   string `yaml:"Senha"`
	DB      string `yaml:"DB"`
}

type Sorteio struct {
	Tipo       string
	Quantidade int
	Nome       string
	Item       Item
}

type ItemNome struct {
	ID         int    `yaml:"ID"`
	Nome       string `yaml:"Nome"`
	Pos        int    `yaml:"Pos"`
	Count      int    `yaml:"Count"`
	MaxCount   int    `yaml:"MaxCount"`
	Data       string `yaml:"Data"`
	ProcType   int    `yaml:"ProcType"`
	ExpireDate int    `yaml:"ExpireDate"`
	GUID1      int    `yaml:"GUID1"`
	GUID2      int    `yaml:"GUID2"`
	Mask       int    `yaml:"Mask"`
}

type Item struct {
	ID         int
	Pos        int
	Count      int
	MaxCount   int
	Data       []byte
	ProcType   int
	ExpireDate int
	GUID1      int
	GUID2      int
	Mask       int
}

type UserOnline struct {
	UserID   UserID
	RoleID   RoleID
	LinkID   int
	LocalSID int
	GSID     int
	Status   byte
	Name     string `pw:"string"`
}

type RoleStatus struct {
	Sversion         byte
	Level            int
	Level2           int
	Exp              int
	Sp               int
	Pp               int
	Hp               int
	Mp               int
	Posx             float32
	Posy             float32
	Posz             float32
	Worldtag         int
	InvaderState     int
	InvaderTime      int
	PariahTime       int
	Reputation       int
	CustomStatus     []byte
	FilterData       []byte
	Charactermode    []byte
	Instancekeylist  []byte
	DbltimeExpire    int
	DbltimeMode      int
	DbltimeBegin     int
	DbltimeUsed      int
	DbltimeMax       int
	TimeUsed         int
	DbltimeData      []byte
	Storesize        uint16
	Petcorral        []byte
	Property         []byte
	VarData          []byte
	Skills           []byte
	Storehousepasswd []byte
	Waypointlist     []byte
	Coolingtime      []byte
	Reserved1        uint
	Reserved2        int
	Reserved3        int
	Reserved4        int
}

type GetRoleBaseArg struct {
	Handler int
	RoleID  RoleID
}

type GRoleForbid struct {
	Type       byte
	Time       int
	CreateTime int
	Reason     string `pw:"string"`
}

type RoleBase struct {
	Version       byte
	ID            int
	Name          string `pw:"string"`
	Race          int
	CLS           int
	Gender        byte
	CustomData    []byte
	ConfigData    []byte
	CustomStamp   int
	Status        byte
	DeleteTime    int
	CreateTime    int
	LastLoginTime int
	ForbidSize    Cuint
	Forbid        []GRoleForbid
	HelpStates    []byte
	Spouse        int
	UserID        UserID
	CrossData     []byte
	Reserved2     byte
	Reserved3     byte
	Reserved4     byte
}

type ChatBroadCast struct {
	Channel   byte
	Emotion   byte
	SrcRoleID int
	Msg       string
	Data      []byte
}

type UserID int
type RoleID struct {
	RoleID int
}

type GMQueryOnline struct {
	QType int
}
type GMQueryOnlineRe struct {
	Header     int
	QType      int
	UsersCount Cuint
	RoleIDS    []RoleID
}

type GetRoleStatusArg struct {
	Handler int
	RoleID  RoleID
}

type DebugAddCash struct {
	UserID UserID
	Cash   int
}

type SysSendMail struct {
	TID         int
	SysID       int
	SysType     byte
	Receiver    RoleID
	Title       string
	Content     string
	AttachObj   Item
	AttachMoney int
}

type TestAPI struct {
	Handler int
}
