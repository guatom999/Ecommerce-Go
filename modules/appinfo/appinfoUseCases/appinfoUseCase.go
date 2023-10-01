package appinfousecases

import (
	"github.com/guatom999/Ecommerce-Go/modules/appinfo"
	appinforepositories "github.com/guatom999/Ecommerce-Go/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(category int) error
}

type appinfoUsecase struct {
	appinfoRepo appinforepositories.IAppinfoRepository
}

func AppinfoUseCase(appinfoRepo appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepo: appinfoRepo,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepo.FindCategory(req)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {

	if err := u.appinfoRepo.InsertCategory(req); err != nil {
		return err
	}

	return nil
}

func (u *appinfoUsecase) DeleteCategory(category int) error {
	if err := u.appinfoRepo.DeleteCetegory(category); err != nil {
		return err
	}

	return nil
}
