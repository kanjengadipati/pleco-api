package role

import "gorm.io/gorm"

type Repository interface {
	FindByID(id uint) (*Role, error)
}

type GormRepository struct {
	db *gorm.DB
}

var _ Repository = (*GormRepository)(nil)

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindByID(id uint) (*Role, error) {
	var role Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
