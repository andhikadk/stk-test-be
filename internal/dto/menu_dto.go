package dto

import (
	"errors"
	"strings"
)

type CreateMenuRequest struct {
	ParentID   *uint   `json:"parent_id" example:"1"`
	Title      string  `json:"title" example:"Dashboard"`
	Path       *string `json:"path,omitempty" example:"/dashboard"`
	Icon       *string `json:"icon,omitempty" example:"icon-dashboard"`
	OrderIndex *int    `json:"order_index,omitempty" example:"0"`
}

func (r *CreateMenuRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required and cannot be empty")
	}

	if len(r.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}

	if r.Path != nil && len(*r.Path) > 255 {
		return errors.New("path cannot exceed 255 characters")
	}

	if r.Icon != nil && len(*r.Icon) > 100 {
		return errors.New("icon cannot exceed 100 characters")
	}

	if r.OrderIndex != nil && *r.OrderIndex < 0 {
		return errors.New("order_index must be a non-negative integer")
	}

	return nil
}

type UpdateMenuRequest struct {
	ParentID   *uint   `json:"parent_id,omitempty" example:"1"`
	Title      *string `json:"title,omitempty" example:"Dashboard"`
	Path       *string `json:"path,omitempty" example:"/dashboard"`
	Icon       *string `json:"icon,omitempty" example:"icon-dashboard"`
	OrderIndex *int    `json:"order_index,omitempty" example:"0"`
}

func (r *UpdateMenuRequest) Validate() error {
	if r.Title != nil {
		trimmedTitle := strings.TrimSpace(*r.Title)
		if trimmedTitle == "" {
			return errors.New("title cannot be empty if provided")
		}
		if len(trimmedTitle) > 255 {
			return errors.New("title cannot exceed 255 characters")
		}
	}

	if r.Path != nil && len(*r.Path) > 255 {
		return errors.New("path cannot exceed 255 characters")
	}

	if r.Icon != nil && len(*r.Icon) > 100 {
		return errors.New("icon cannot exceed 100 characters")
	}

	if r.OrderIndex != nil && *r.OrderIndex < 0 {
		return errors.New("order_index must be a non-negative integer")
	}

	return nil
}

type MoveMenuRequest struct {
	ParentID *uint `json:"parent_id" example:"1"`
}

func (r *MoveMenuRequest) Validate() error {
	return nil
}

type ReorderMenuRequest struct {
	NewIndex int  `json:"new_index" example:"2"`
	OldIndex *int `json:"old_index,omitempty" example:"0"`
}

func (r *ReorderMenuRequest) Validate() error {
	if r.NewIndex < 0 {
		return errors.New("new_index must be a non-negative integer")
	}

	if r.OldIndex != nil && *r.OldIndex < 0 {
		return errors.New("old_index must be a non-negative integer if provided")
	}

	return nil
}
