package handlers

import (
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
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
// @Param        menu  body      models.Menu  true  "Menu object"
// @Success      201   {object}  models.APIResponse{data=models.Menu}
// @Failure      400   {object}  models.APIResponse
// @Failure      500   {object}  models.APIResponse
// @Router       /api/menus [post]
func CreateMenu(c *fiber.Ctx) error {
	var menu models.Menu

	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if menu.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Title is required",
			Error:   "missing required field: title",
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.CreateMenu(&menu); err != nil {
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
// @Param        id    path      int          true  "Menu ID"
// @Param        menu  body      models.Menu  true  "Menu object"
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

	var menu models.Menu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.UpdateMenu(uint(id), &menu); err != nil {
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
// @Param        id       path      int     true  "Menu ID"
// @Param        request  body      object  true  "Move request"
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

	var req struct {
		ParentID *uint `json:"parent_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.MoveMenu(uint(id), req.ParentID); err != nil {
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
// @Param        id       path      int     true  "Menu ID"
// @Param        request  body      object  true  "Reorder request"
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

	var req struct {
		OrderIndex int `json:"order_index"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	menuService := services.NewMenuService(database.GetDB())
	if err := menuService.ReorderMenu(uint(id), req.OrderIndex); err != nil {
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
