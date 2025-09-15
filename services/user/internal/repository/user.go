package repository

import (
	"kube/services/user/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepository) List(limit, offset int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// QRLoginSession methods
func (r *UserRepository) CreateQRLoginSession(session *model.QRLoginSession) error {
	return r.db.Create(session).Error
}

func (r *UserRepository) GetQRLoginSessionByToken(token string) (*model.QRLoginSession, error) {
	var session model.QRLoginSession
	err := r.db.Where("token = ?", token).First(&session).Error
	return &session, err
}

func (r *UserRepository) UpdateQRLoginSession(session *model.QRLoginSession) error {
	return r.db.Save(session).Error
}

func (r *UserRepository) DeleteQRLoginSession(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.QRLoginSession{}).Error
}
