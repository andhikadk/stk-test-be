package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/andhikadk/stk-test-be/internal/database"
	"github.com/andhikadk/stk-test-be/internal/dto"
	"github.com/andhikadk/stk-test-be/internal/models"
	"github.com/andhikadk/stk-test-be/internal/routes"
	"github.com/andhikadk/stk-test-be/internal/testutil"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*fiber.App, *gorm.DB, func()) {
	db := testutil.SetupTestDB(t)

	originalDB := database.DB
	database.DB = db

	testutil.InitTestLogger()

	app := fiber.New()
	routes.SetupRoutes(app)

	cleanup := func() {
		database.DB = originalDB
		testutil.TeardownTestDB(db)
	}

	return app, db, cleanup
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func TestGetMenus_EmptyDatabase(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/menus", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menus retrieved successfully", result.Message)
	testutil.AssertEqual(t, fiber.StatusOK, result.Status)

	menus, ok := result.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected Data to be array, got %T", result.Data)
	}
	testutil.AssertLen(t, menus, 0, "Expected empty menu array")
}

func TestGetMenus_WithSingleMenu(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	testutil.CreateMenuFixture(db, "Dashboard", nil, 0)

	req := httptest.NewRequest("GET", "/api/menus", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menus := result.Data.([]interface{})
	testutil.AssertLen(t, menus, 1)

	menu := menus[0].(map[string]interface{})
	testutil.AssertEqual(t, "Dashboard", menu["title"])
}

func TestGetMenus_WithHierarchy(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent, children := testutil.CreateMenuHierarchy(db)

	req := httptest.NewRequest("GET", "/api/menus", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menus := result.Data.([]interface{})
	testutil.AssertLen(t, menus, 1, "Should have 1 root menu")

	rootMenu := menus[0].(map[string]interface{})
	testutil.AssertEqual(t, parent.Title, rootMenu["title"])

	childrenData := rootMenu["children"].([]interface{})
	testutil.AssertLen(t, childrenData, len(children), "Should have 3 children")

	for i, child := range children {
		childData := childrenData[i].(map[string]interface{})
		testutil.AssertEqual(t, child.Title, childData["title"])
	}
}

func TestGetMenus_WithMultiLevelHierarchy(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	hierarchy := testutil.CreateMultiLevelHierarchy(db)

	req := httptest.NewRequest("GET", "/api/menus", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menus := result.Data.([]interface{})
	testutil.AssertLen(t, menus, 2, "Should have 2 root menus")

	root1 := menus[0].(map[string]interface{})
	testutil.AssertEqual(t, hierarchy["root1"].Title, root1["title"])

	root1Children := root1["children"].([]interface{})
	testutil.AssertLen(t, root1Children, 2, "Root 1 should have 2 children")

	child1_1 := root1Children[0].(map[string]interface{})
	grandchildren := child1_1["children"].([]interface{})
	testutil.AssertLen(t, grandchildren, 1, "Child 1.1 should have 1 grandchild")
}

func TestGetMenu_Success(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuWithPath(db, "Dashboard", "/dashboard", "icon-dashboard", nil)

	url := fmt.Sprintf("/api/menus/%s", menu.ID)
	req := httptest.NewRequest("GET", url, nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu retrieved successfully", result.Message)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, menu.Title, menuData["title"])
	testutil.AssertEqual(t, *menu.Path, menuData["path"])
	testutil.AssertEqual(t, *menu.Icon, menuData["icon"])
}

func TestGetMenu_NotFound(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	nonExistentID := uuid.New()
	url := fmt.Sprintf("/api/menus/%s", nonExistentID)
	req := httptest.NewRequest("GET", url, nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusNotFound, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu not found", result.Message)
	testutil.AssertNotEmpty(t, result.Error)
}

func TestGetMenu_InvalidID(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/menus/invalid", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Invalid menu ID", result.Message)
}

func TestGetMenu_WithChildren(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent, _ := testutil.CreateMenuHierarchy(db)

	url := fmt.Sprintf("/api/menus/%s", parent.ID)
	req := httptest.NewRequest("GET", url, nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	children := menuData["children"].([]interface{})
	testutil.AssertLen(t, children, 3, "Parent should have 3 children")
}

func TestCreateMenu_Success(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	reqBody := dto.CreateMenuRequest{
		Title: "New Menu",
		Path:  stringPtr("/new-menu"),
		Icon:  stringPtr("icon-new"),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("Expected status 201, got %d. Message: %s, Error: %s", resp.StatusCode, result.Message, result.Error)
	}

	testutil.AssertEqual(t, "Menu created successfully", result.Message)

	if result.Data == nil {
		t.Fatalf("Expected menu data, got nil. Message: %s, Error: %s", result.Message, result.Error)
	}

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, reqBody.Title, menuData["title"])
	testutil.AssertEqual(t, *reqBody.Path, menuData["path"])
	testutil.AssertEqual(t, *reqBody.Icon, menuData["icon"])
	testutil.AssertNotNil(t, menuData["id"])
}

func TestCreateMenu_WithParent(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent := testutil.CreateMenuFixture(db, "Parent", nil, 0)

	reqBody := dto.CreateMenuRequest{
		Title:    "Child Menu",
		ParentID: &parent.ID,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusCreated, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, parent.ID.String(), menuData["parent_id"])
}

func TestCreateMenu_WithCustomOrderIndex(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	testutil.CreateMenuFixture(db, "Menu 1", nil, 0)
	testutil.CreateMenuFixture(db, "Menu 2", nil, 1)

	reqBody := dto.CreateMenuRequest{
		Title:      "Menu Inserted",
		OrderIndex: intPtr(1),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusCreated, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, float64(1), menuData["order_index"])
}

func TestCreateMenu_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		request dto.CreateMenuRequest
		errMsg  string
	}{
		{
			name:    "empty title",
			request: dto.CreateMenuRequest{Title: ""},
			errMsg:  "title is required",
		},
		{
			name:    "title too long",
			request: dto.CreateMenuRequest{Title: string(make([]byte, 256))},
			errMsg:  "title cannot exceed 255 characters",
		},
		{
			name:    "negative order index",
			request: dto.CreateMenuRequest{Title: "Test", OrderIndex: intPtr(-1)},
			errMsg:  "order_index must be a non-negative integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, cleanup := setupTest(t)
			defer cleanup()

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/menus", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}

			testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

			var result models.APIResponse
			testutil.ParseJSONResponse(t, resp.Body, &result)

			testutil.AssertEqual(t, "Validation failed", result.Message)
			testutil.AssertContains(t, result.Error, tt.errMsg)
		})
	}
}

func TestCreateMenu_InvalidJSON(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/api/menus", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Invalid request body", result.Message)
}

func TestUpdateMenu_Success(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuFixture(db, "Original Title", nil, 0)

	reqBody := dto.UpdateMenuRequest{
		Title: stringPtr("Updated Title"),
		Path:  stringPtr("/updated"),
		Icon:  stringPtr("icon-updated"),
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s", menu.ID)
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu updated successfully", result.Message)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, *reqBody.Title, menuData["title"])
	testutil.AssertEqual(t, *reqBody.Path, menuData["path"])
	testutil.AssertEqual(t, *reqBody.Icon, menuData["icon"])
}

func TestUpdateMenu_ChangeParent(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent1 := testutil.CreateMenuFixture(db, "Parent 1", nil, 0)
	parent2 := testutil.CreateMenuFixture(db, "Parent 2", nil, 1)
	child := testutil.CreateMenuFixture(db, "Child", &parent1.ID, 0)

	reqBody := dto.UpdateMenuRequest{
		ParentID: &parent2.ID,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s", child.ID)
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, parent2.ID.String(), menuData["parent_id"])
}

func TestUpdateMenu_MoveToRoot(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent := testutil.CreateMenuFixture(db, "Parent", nil, 0)
	child := testutil.CreateMenuFixture(db, "Child", &parent.ID, 0)

	reqBody := dto.UpdateMenuRequest{
		ParentID: nil,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s", child.ID)
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertNil(t, menuData["parent_id"])
}

func TestUpdateMenu_NotFound(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	reqBody := dto.UpdateMenuRequest{
		Title: stringPtr("Updated"),
	}

	body, _ := json.Marshal(reqBody)
	nonExistentID := uuid.New()
	url := fmt.Sprintf("/api/menus/%s", nonExistentID)
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusInternalServerError, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Failed to update menu", result.Message)
}

func TestUpdateMenu_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		request dto.UpdateMenuRequest
		errMsg  string
	}{
		{
			name:    "empty title",
			request: dto.UpdateMenuRequest{Title: stringPtr("")},
			errMsg:  "title cannot be empty",
		},
		{
			name:    "title too long",
			request: dto.UpdateMenuRequest{Title: stringPtr(string(make([]byte, 256)))},
			errMsg:  "title cannot exceed 255 characters",
		},
		{
			name:    "negative order index",
			request: dto.UpdateMenuRequest{OrderIndex: intPtr(-1)},
			errMsg:  "order_index must be a non-negative integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, db, cleanup := setupTest(t)
			defer cleanup()

			menu := testutil.CreateMenuFixture(db, "Test Menu", nil, 0)

			body, _ := json.Marshal(tt.request)
			url := fmt.Sprintf("/api/menus/%s", menu.ID)
			req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}

			testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

			var result models.APIResponse
			testutil.ParseJSONResponse(t, resp.Body, &result)

			testutil.AssertEqual(t, "Validation failed", result.Message)
			testutil.AssertContains(t, result.Error, tt.errMsg)
		})
	}
}

func TestDeleteMenu_Success(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuFixture(db, "To Delete", nil, 0)

	url := fmt.Sprintf("/api/menus/%s", menu.ID)
	req := httptest.NewRequest("DELETE", url, nil)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu deleted successfully", result.Message)

	var count int64
	db.Model(&models.Menu{}).Where("id = ?", menu.ID).Count(&count)
	testutil.AssertEqual(t, int64(0), count, "Menu should be deleted")
}

func TestDeleteMenu_WithChildren(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent, children := testutil.CreateMenuHierarchy(db)

	url := fmt.Sprintf("/api/menus/%s", parent.ID)
	req := httptest.NewRequest("DELETE", url, nil)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var parentCount int64
	db.Model(&models.Menu{}).Where("id = ?", parent.ID).Count(&parentCount)
	testutil.AssertEqual(t, int64(0), parentCount, "Parent should be deleted")

	var childCount int64
	db.Model(&models.Menu{}).Where("parent_id = ?", parent.ID).Count(&childCount)
	testutil.AssertEqual(t, int64(0), childCount, "Children should be deleted")

	var totalCount int64
	db.Model(&models.Menu{}).Count(&totalCount)
	testutil.AssertEqual(t, int64(0), totalCount, fmt.Sprintf("All menus should be deleted (parent + %d children)", len(children)))
}

func TestDeleteMenu_NotFound(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	nonExistentID := uuid.New()
	url := fmt.Sprintf("/api/menus/%s", nonExistentID)
	req := httptest.NewRequest("DELETE", url, nil)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)
}

func TestDeleteMenu_InvalidID(t *testing.T) {
	app, _, cleanup := setupTest(t)
	defer cleanup()

	req := httptest.NewRequest("DELETE", "/api/menus/invalid", nil)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Invalid menu ID", result.Message)
}

func TestMoveMenu_Success(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent1 := testutil.CreateMenuFixture(db, "Parent 1", nil, 0)
	parent2 := testutil.CreateMenuFixture(db, "Parent 2", nil, 1)
	child := testutil.CreateMenuFixture(db, "Child", &parent1.ID, 0)

	reqBody := dto.MoveMenuRequest{
		ParentID: &parent2.ID,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/move", child.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu moved successfully", result.Message)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, parent2.ID.String(), menuData["parent_id"])
}

func TestMoveMenu_ToRoot(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent := testutil.CreateMenuFixture(db, "Parent", nil, 0)
	child := testutil.CreateMenuFixture(db, "Child", &parent.ID, 0)

	reqBody := dto.MoveMenuRequest{
		ParentID: nil,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/move", child.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertNil(t, menuData["parent_id"])
}

func TestMoveMenu_InvalidParent(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuFixture(db, "Menu", nil, 0)

	invalidParentID := uuid.New()
	reqBody := dto.MoveMenuRequest{
		ParentID: &invalidParentID,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/move", menu.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Failed to move menu", result.Message)
	testutil.AssertContains(t, result.Error, "parent menu not found")
}

func TestReorderMenu_Success(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu0 := testutil.CreateMenuFixture(db, "Menu 0", nil, 0)
	testutil.CreateMenuFixture(db, "Menu 1", nil, 1)
	testutil.CreateMenuFixture(db, "Menu 2", nil, 2)

	reqBody := dto.ReorderMenuRequest{
		NewIndex: 2,
		OldIndex: intPtr(0),
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/reorder", menu0.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Menu reordered successfully", result.Message)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, float64(2), menuData["order_index"])
}

func TestReorderMenu_ToFirst(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	testutil.CreateMenuFixture(db, "Menu 0", nil, 0)
	testutil.CreateMenuFixture(db, "Menu 1", nil, 1)
	menu2 := testutil.CreateMenuFixture(db, "Menu 2", nil, 2)

	reqBody := dto.ReorderMenuRequest{
		NewIndex: 0,
		OldIndex: intPtr(2),
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/reorder", menu2.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, float64(0), menuData["order_index"])
}

func TestReorderMenu_AutoClamp(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuFixture(db, "Menu 0", nil, 0)
	testutil.CreateMenuFixture(db, "Menu 1", nil, 1)
	testutil.CreateMenuFixture(db, "Menu 2", nil, 2)

	reqBody := dto.ReorderMenuRequest{
		NewIndex: 100,
		OldIndex: intPtr(0),
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/reorder", menu.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, float64(2), menuData["order_index"], "Should auto-clamp to max index (2)")
}

func TestReorderMenu_NegativeIndex(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	menu := testutil.CreateMenuFixture(db, "Menu 0", nil, 0)

	reqBody := dto.ReorderMenuRequest{
		NewIndex: -1,
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/reorder", menu.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusBadRequest, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	testutil.AssertEqual(t, "Validation failed", result.Message)
	testutil.AssertContains(t, result.Error, "new_index must be a non-negative integer")
}

func TestReorderMenu_WithinSiblings(t *testing.T) {
	app, db, cleanup := setupTest(t)
	defer cleanup()

	parent := testutil.CreateMenuFixture(db, "Parent", nil, 0)
	child0 := testutil.CreateMenuFixture(db, "Child 0", &parent.ID, 0)
	testutil.CreateMenuFixture(db, "Child 1", &parent.ID, 1)
	testutil.CreateMenuFixture(db, "Child 2", &parent.ID, 2)

	testutil.CreateMenuFixture(db, "Other Root", nil, 1)

	reqBody := dto.ReorderMenuRequest{
		NewIndex: 2,
		OldIndex: intPtr(0),
	}

	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("/api/menus/%s/reorder", child0.ID)
	req := httptest.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	testutil.AssertStatusCode(t, fiber.StatusOK, resp)

	var result models.APIResponse
	testutil.ParseJSONResponse(t, resp.Body, &result)

	menuData := result.Data.(map[string]interface{})
	testutil.AssertEqual(t, float64(2), menuData["order_index"])
	testutil.AssertEqual(t, parent.ID.String(), menuData["parent_id"])
}
