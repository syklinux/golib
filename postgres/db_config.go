package postgres

type Conf struct {
	WriteServer string   `json:"WriteServer"`
	ReadServer  []string `json:"ReadServer"`
	User        string   `json:"User"`
	Password    string   `json:"Password"`
	DataBase    string   `json:"DataBase"`
	MaxIdleTime int      `json:"MaxIdleTime"`
	MaxLiftTime int      `json:"MaxLiftTime"`
	MaxIdleConn int      `json:"MaxIdleConn"`
	MaxOpenConn int      `json:"maxOpenConn"`
	LogLevel    string   `json:"logLevel"`
}
