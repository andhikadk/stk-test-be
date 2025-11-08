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

func (s *MenuService) getSiblingCount(parentID *uint) (int64, error) {
	var count int64
	query := s.db.Model(&models.Menu{})

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *MenuService) ReorderMenu(id uint, newIndex int, oldIndex *int) error {
	var menu models.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("menu not found")
		}
		return err
	}

	siblingCount, err := s.getSiblingCount(menu.ParentID)
	if err != nil {
		return err
	}

	if newIndex < 0 || int64(newIndex) >= siblingCount {
		return errors.New("invalid target position")
	}

	actualOldIndex := menu.OrderIndex
	if oldIndex != nil {
		actualOldIndex = *oldIndex
	}

	if actualOldIndex == newIndex {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		baseQuery := tx.Model(&models.Menu{}).Where("id != ?", id)

		if menu.ParentID == nil {
			baseQuery = baseQuery.Where("parent_id IS NULL")
		} else {
			baseQuery = baseQuery.Where("parent_id = ?", *menu.ParentID)
		}

		if actualOldIndex < newIndex {
			if err := baseQuery.
				Where("order_index > ?", actualOldIndex).
				Where("order_index <= ?", newIndex).
				Update("order_index", gorm.Expr("order_index - 1")).Error; err != nil {
				return err
			}
		} else {
			if err := baseQuery.
				Where("order_index >= ?", newIndex).
				Where("order_index < ?", actualOldIndex).
				Update("order_index", gorm.Expr("order_index + 1")).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&models.Menu{}).Where("id = ?", id).Update("order_index", newIndex).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *MenuService) buildChildren(parentID uint, menuMap map[uint]*models.Menu, allMenus []models.Menu) []models.Menu {
	var children []models.Menu

	for i := range allMenus {
		if allMenus[i].ParentID != nil && *allMenus[i].ParentID == parentID {
			child := allMenus[i]
			child.Children = s.buildChildren(child.ID, menuMap, allMenus)
			children = append(children, child)
		}
	}

	return children
}

func (s *MenuService) GetMenuTree() ([]models.Menu, error) {
	var allMenus []models.Menu
	if err := s.db.Order("order_index ASC").Find(&allMenus).Error; err != nil {
		return nil, err
	}

	menuMap := make(map[uint]*models.Menu)
	for i := range allMenus {
		menuMap[allMenus[i].ID] = &allMenus[i]
	}

	var rootMenus []models.Menu
	for i := range allMenus {
		if allMenus[i].ParentID == nil {
			menu := allMenus[i]
			menu.Children = s.buildChildren(menu.ID, menuMap, allMenus)
			rootMenus = append(rootMenus, menu)
		}
	}

	return rootMenus, nil
}
