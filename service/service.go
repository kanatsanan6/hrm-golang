package service

type Service interface {
	Mailer
	Policy
}

type service struct {
	Mailer
	Policy
}

func NewService() Service {
	return &service{
		&mailer{},
		newPolicy(),
	}
}
