package service

// WebSvc
type WebSvc interface {
	Run() error
	Stop() error
}

// SSHSvc
type SSHSvc interface {
	// Запуск сервера трансляции SSH с сервера
	Run() error
}
