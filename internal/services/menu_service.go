package services

import (
	"errors"
	"go-fiber-boilerplate/internal/models"
	"gorm.io/gorm"
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{db: db}
}

func (s *MenuService) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	if err := s.db.Where("parent_id IS NULL").Preload("Children").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (s *MenuService) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	if err := s.db.Preload("Children").First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu not found")
		}
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) CreateMenu(menu *models.Menu) error {
	return s.db.Create(menu).Error
}

func (s *MenuService) UpdateMenu(id uint, menu *models.Menu) error {
	return s.db.Model(&models.Menu{}).Where("id = ?", id).Updates(menu).Error
}

func (s *MenuService) DeleteMenu(id uint) error {
	if err := s.db.Where("parent_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.Menu{}, id).Error
}

func (s *MenuService) MoveMenu(id uint, newParentID *uint) error {
	if newParentID != nil && *newParentID != 0 {
		var parent models.Menu
		if err := s.db.First(&parent, *newParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parent menu not found")
			}
			return err
		}
	}

	return s.db.Model(&models.Menu{}).Where("id = ?", id).Update("parent_id", newParentID).Error
}

func (s *MenuService) ReorderMenu(id uint, newOrder int) error {
	return s.db.Model(&models.Menu{}).Where("id = ?", id).Update("order_index", newOrder).Error
}

func (s *MenuService) GetMenuTree() ([]models.Menu, error) {
	var menus []models.Menu
	if err := s.db.Where("parent_id IS NULL").Order("order_index ASC").Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_index ASC")
	}).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
