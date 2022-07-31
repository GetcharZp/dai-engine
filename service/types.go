package service

type EtcdConfig struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}
