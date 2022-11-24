package conn

type DbConfig struct {
	Dns            string
	MaxConnNumber  int
	IdleConnNumber int
}
