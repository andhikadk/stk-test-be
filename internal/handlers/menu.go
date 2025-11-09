package handlers

import (
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetMenus godoc
// @Summary      Get all menu items
// @Description  Get all menu items in hierarchical tree structure
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.APIResponse{data=[]models.Menu}
// @Failure      500  {object}  models.APIResponse
// @Router       /api/menus [get]
func GetMenus(c *fiber.Ctx) error {
	menuService := services.NewMenuService(database.GetDB())
	menus, err := menuService.GetMenuTree()
	if err != nil {
		utils.ErrorLogger.Printf("[GetMenus] Failed to fetch menu tree: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to fetch menus",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menus retrieved successfully",
		Data:    menus,
	})
}

// GetMenu godoc
// @Summary      Get single menu item
// @Description  Get a single menu item by ID
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Menu ID"
// @Success      200  {object}  models.APIResponse{data=models.Menu}
// @Failure      400  {object}  models.APIResponse
// @Failure      404  {object}  models.APIResponse
// @Router       /api/menus/{id} [get]
func GetMenu(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu ID",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	menu, err := menuService.GetMenuByID(uint(id))
	if err != nil {
		utils.ErrorLogger.Printf("[GetMenu] menuID=%d error: %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  fiber.StatusNotFound,
			Message: "Menu not found",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menu retrieved successfully",
		Data:    menu,
	})
}

// CreateMenu godoc
// @Summary      Create new menu item
// @Description  Create a new menu item
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        menu  body      dto.CreateMenuRequest  true  "Menu creation data"
// @Success      201   {object}  models.APIResponse{data=models.Menu}
// @Failure      400   {object}  models.APIResponse
// @Failure      500   {object}  models.APIResponse
// @Router       /api/menus [post]
func CreateMenu(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[CreateMenu] Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	menu := models.Menu{
		ParentID:   req.ParentID,
		Title:      req.Title,
		Path:       req.Path,
		Icon:       req.Icon,
		OrderIndex: 0,
	}

	if req.OrderIndex != nil {
		menu.OrderIndex = *req.OrderIndex
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.CreateMenu(&menu); err != nil {
		utils.ErrorLogger.Printf("[CreateMenu] Failed to create menu '%s': %v", req.Title, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create menu",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  fiber.StatusCreated,
		Message: "Menu created successfully",
		Data:    menu,
	})
}

// UpdateMenu godoc
// @Summary      Update menu item
// @Description  Update a menu item
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Menu ID"
// @Param        menu  body      dto.UpdateMenuRequest  true  "Menu update data"
// @Success      200   {object}  models.APIResponse{data=models.Menu}
// @Failure      400   {object}  models.APIResponse
// @Failure      500   {object}  models.APIResponse
// @Router       /api/menus/{id} [put]
func UpdateMenu(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu ID",
			Error:   err.Error(),
		})
	}

	var req dto.UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[UpdateMenu] menuID=%d validation failed: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	menu := models.Menu{}
	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}
	if req.Title != nil {
		menu.Title = *req.Title
	}
	if req.Path != nil {
		menu.Path = req.Path
	}
	if req.Icon != nil {
		menu.Icon = req.Icon
	}
	if req.OrderIndex != nil {
		menu.OrderIndex = *req.OrderIndex
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.UpdateMenu(uint(id), &menu); err != nil {
		utils.ErrorLogger.Printf("[UpdateMenu] menuID=%d error: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update menu",
			Error:   err.Error(),
		})
	}

	updated, _ := menuService.GetMenuByID(uint(id))
	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menu updated successfully",
		Data:    updated,
	})
}

// DeleteMenu godoc
// @Summary      Delete menu item
// @Description  Delete a menu item and its children
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Menu ID"
// @Success      200  {object}  models.APIResponse
// @Failure      400  {object}  models.APIResponse
// @Failure      500  {object}  models.APIResponse
// @Router       /api/menus/{id} [delete]
func DeleteMenu(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu ID",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.DeleteMenu(uint(id)); err != nil {
		utils.ErrorLogger.Printf("[DeleteMenu] menuID=%d error: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete menu",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menu deleted successfully",
	})
}

// MoveMenu godoc
// @Summary      Move menu item to different parent
// @Description  Move a menu item to a different parent
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Menu ID"
// @Param        request  body      dto.MoveMenuRequest  true  "Move request"
// @Success      200      {object}  models.APIResponse{data=models.Menu}
// @Failure      400      {object}  models.APIResponse
// @Router       /api/menus/{id}/move [patch]
func MoveMenu(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu ID",
			Error:   err.Error(),
		})
	}

	var req dto.MoveMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[MoveMenu] menuID=%d validation failed: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.MoveMenu(uint(id), req.ParentID); err != nil {
		utils.ErrorLogger.Printf("[MoveMenu] menuID=%d error: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Failed to move menu",
			Error:   err.Error(),
		})
	}

	updated, _ := menuService.GetMenuByID(uint(id))
	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menu moved successfully",
		Data:    updated,
	})
}

// ReorderMenu godoc
// @Summary      Reorder menu item within same level
// @Description  Change the order index of a menu item
// @Tags         Menus
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "Menu ID"
// @Param        request  body      dto.ReorderMenuRequest  true  "Reorder request"
// @Success      200      {object}  models.APIResponse{data=models.Menu}
// @Failure      400      {object}  models.APIResponse
// @Failure      500      {object}  models.APIResponse
// @Router       /api/menus/{id}/reorder [patch]
func ReorderMenu(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu ID",
			Error:   err.Error(),
		})
	}

	var req dto.ReorderMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[ReorderMenu] menuID=%d validation failed: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.ReorderMenu(uint(id), req.NewIndex, req.OldIndex); err != nil {
		utils.ErrorLogger.Printf("[ReorderMenu] menuID=%d newIndex=%d error: %v", id, req.NewIndex, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to reorder menu",
			Error:   err.Error(),
		})
	}

	updated, _ := menuService.GetMenuByID(uint(id))
	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Menu reordered successfully",
		Data:    updated,
	})
}
