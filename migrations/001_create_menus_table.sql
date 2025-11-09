-- Create menus table with UUID primary key
-- Created at: 2025-11-09
-- Purpose: Hierarchical menu structure for navigation with UUID identifiers

-- Enable uuid-ossp extension for PostgreSQL (if not already enabled)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create menus table
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES menus(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    path VARCHAR(255),
    icon VARCHAR(100),
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_order_index ON menus(order_index);
CREATE INDEX IF NOT EXISTS idx_menus_deleted_at ON menus(deleted_at);

-- Create composite index for querying children by parent
CREATE INDEX IF NOT EXISTS idx_menus_parent_order ON menus(parent_id, order_index) WHERE deleted_at IS NULL;

-- Add comment to table
COMMENT ON TABLE menus IS 'Hierarchical menu structure for navigation';
COMMENT ON COLUMN menus.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN menus.parent_id IS 'Reference to parent menu item (NULL for root menus)';
COMMENT ON COLUMN menus.title IS 'Menu item title displayed in UI';
COMMENT ON COLUMN menus.path IS 'URL path for navigation (NULL for parent menus)';
COMMENT ON COLUMN menus.icon IS 'Icon identifier for UI display';
COMMENT ON COLUMN menus.order_index IS 'Order position within same parent level';
COMMENT ON COLUMN menus.deleted_at IS 'Soft delete timestamp (NULL if not deleted)';
