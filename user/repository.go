package user

import "gorm.io/gorm"

// untuk object/struct lain akan mengacu ke repository
// anggep kaya kontrak interface tuh, harus ada func yang udahdi deklarasi disini
type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindById(ID int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// register user (add new user)
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// check registered user
func (r *repository) FindByEmail(email string) (User, error) {
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil

}

// find user by id
func (r *repository) FindById(ID int) (User, error){
	var user User

	err := r.db.Where("id = ?", ID).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// update data tabel users
func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil

}