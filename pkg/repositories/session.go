package sparklerepo

type SessionRepository interface {
	Create()
	Update()
	Delete()
	FindOne()
}
