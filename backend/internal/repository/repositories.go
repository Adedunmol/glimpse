package repository

import "github.com/Adedunmol/glimpse/internal/server"

type Repositories struct {
	Upload *UploadRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Upload: NewUploadRepository(s),
	}
}
