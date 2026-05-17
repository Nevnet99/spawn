local M = {}

-- General system wide keybindings
M.general = {
  n = { -- Normal mode navigation overrides
    -- Set up relative line numbers for lightning-fast cursor jumps
    ["<leader>nu"] = { "<cmd> set rnu! <CR>", "Toggle relative line numbers" },
  },
}

return M
