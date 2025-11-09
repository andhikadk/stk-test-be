package testutil

import (
	"github.com/andhikadk/stk-test-be/internal/models"

	"gorm.io/gorm"
)

func CreateMenuFixture(db *gorm.DB, title string, parentID *uint, orderIndex int) *models.Menu {
	menu := &models.Menu{
		Title:      title,
		ParentID:   parentID,
		OrderIndex: orderIndex,
	}
	db.Create(menu)
	return menu
}

func CreateMenuWithPath(db *gorm.DB, title string, path string, icon string, parentID *uint) *models.Menu {
	pathPtr := &path
	iconPtr := &icon
	menu := &models.Menu{
		Title:      title,
		Path:       pathPtr,
		Icon:       iconPtr,
		ParentID:   parentID,
		OrderIndex: 0,
	}
	db.Create(menu)
	return menu
}

func CreateMenuHierarchy(db *gorm.DB) (*models.Menu, []*models.Menu) {
	parent := CreateMenuFixture(db, "Parent Menu", nil, 0)

	children := []*models.Menu{
		CreateMenuFixture(db, "Child 1", &parent.ID, 0),
		CreateMenuFixture(db, "Child 2", &parent.ID, 1),
		CreateMenuFixture(db, "Child 3", &parent.ID, 2),
	}

	return parent, children
}

func CreateMultiLevelHierarchy(db *gorm.DB) map[string]*models.Menu {
	root1 := CreateMenuFixture(db, "Root 1", nil, 0)
	root2 := CreateMenuFixture(db, "Root 2", nil, 1)

	child1_1 := CreateMenuFixture(db, "Child 1.1", &root1.ID, 0)
	child1_2 := CreateMenuFixture(db, "Child 1.2", &root1.ID, 1)

	grandchild1_1_1 := CreateMenuFixture(db, "Grandchild 1.1.1", &child1_1.ID, 0)

	return map[string]*models.Menu{
		"root1":           root1,
		"root2":           root2,
		"child1_1":        child1_1,
		"child1_2":        child1_2,
		"grandchild1_1_1": grandchild1_1_1,
	}
}
