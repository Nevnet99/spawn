local M = {}

-- Highlights, UI, and core statusline configuration
M.ui = {
  theme = "onedark",       -- Premium dark aesthetic baseline
  transparency = false,    -- Keep a solid background for maximum code legibility
  
  -- Custom statusline configuration
  statusline = {
    theme = "default",
    separator_style = "round", -- Smooth, pill-shaped terminal status bars
  },
}

-- Point NvChad to your custom keymaps file
M.mappings = require "custom.mappings"

return M
