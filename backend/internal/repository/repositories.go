package repository

import "github.com/Adedunmol/glimpse/internal/server"

type Repositories struct {
	Upload         *UploadRepository
	Photo          *PhotoRepository
	UserRepository UserRepository
}

func NewRepositories(s *server.Server) *Repositories {
	userRepo := NewPostgresRepository(s)
	return &Repositories{
		UserRepository: userRepo,
		Upload:         NewUploadRepository(s),
		Photo:          NewPhotoRepository(s),
	}
}
