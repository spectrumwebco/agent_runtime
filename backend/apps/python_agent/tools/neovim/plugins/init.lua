
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({
    "git",
    "clone",
    "--filter=blob:none",
    "https://github.com/folke/lazy.nvim.git",
    "--branch=stable",
    lazypath,
  })
end
vim.opt.rtp:prepend(lazypath)

local plugins = {
  {
    "aserowy/tmux.nvim",
    config = function()
      require("tmux").setup({
        copy_sync = {
          enable = true,
        },
        navigation = {
          enable_default_keybindings = true,
        },
        resize = {
          enable_default_keybindings = true,
        }
      })
    end
  },

  {
    "ibhagwan/fzf-lua",
    dependencies = { "nvim-tree/nvim-web-devicons" },
    config = function()
      require("fzf-lua").setup({
        global_resume = true,
        global_resume_query = true,
        winopts = {
          height = 0.85,
          width = 0.80,
          preview = {
            scrollbar = 'float',
          },
        },
        keymap = {
          builtin = {
            ["<c-d>"] = "preview-page-down",
            ["<c-u>"] = "preview-page-up",
          },
        },
      })
      
      vim.keymap.set("n", "<leader>ff", "<cmd>lua require('fzf-lua').files()<CR>", { desc = "Find files" })
      vim.keymap.set("n", "<leader>fg", "<cmd>lua require('fzf-lua').grep_project()<CR>", { desc = "Grep project" })
      vim.keymap.set("n", "<leader>fb", "<cmd>lua require('fzf-lua').buffers()<CR>", { desc = "Find buffers" })
    end
  },

  {
    "benlubas/molten-nvim",
    dependencies = {
      "3rd/image.nvim",
    },
    build = ":UpdateRemotePlugins",
    init = function()
      vim.g.molten_image_provider = "image.nvim"
      vim.g.molten_output_win_max_height = 20
      vim.g.molten_auto_open_output = true
    end,
    config = function()
      vim.keymap.set("n", "<leader>mi", ":MoltenInit<CR>", { desc = "Initialize Molten" })
      vim.keymap.set("n", "<leader>me", ":MoltenEvaluateOperator<CR>", { desc = "Evaluate operator" })
      vim.keymap.set("n", "<leader>ml", ":MoltenEvaluateLine<CR>", { desc = "Evaluate line" })
      vim.keymap.set("n", "<leader>mc", ":MoltenReevaluateCell<CR>", { desc = "Reevaluate cell" })
    end
  },

  {
    "iamcco/markdown-preview.nvim",
    cmd = { "MarkdownPreview", "MarkdownPreviewStop" },
    build = function() vim.fn["mkdp#util#install"]() end,
    init = function()
      vim.g.mkdp_auto_start = 0
      vim.g.mkdp_auto_close = 1
      vim.g.mkdp_refresh_slow = 0
      vim.g.mkdp_command_for_global = 0
      vim.g.mkdp_open_to_the_world = 1  -- Allow access from other containers
      vim.g.mkdp_browser = ""  -- Use default browser
      vim.g.mkdp_echo_preview_url = 1
      vim.g.mkdp_page_title = '${name}'
    end,
    config = function()
      vim.keymap.set("n", "<leader>mp", ":MarkdownPreview<CR>", { desc = "Start Markdown preview" })
      vim.keymap.set("n", "<leader>ms", ":MarkdownPreviewStop<CR>", { desc = "Stop Markdown preview" })
    end,
    ft = { "markdown" }
  },

  {
    "folke/tokyonight.nvim",
    lazy = false,
    priority = 1000,
    config = function()
      require("tokyonight").setup({
        style = "storm",
        transparent = false,
        terminal_colors = true,
        styles = {
          comments = { italic = true },
          keywords = { italic = true },
          functions = {},
          variables = {},
          sidebars = "dark",
          floats = "dark",
        },
        sidebars = { "qf", "help", "terminal", "packer" },
        day_brightness = 0.3,
        hide_inactive_statusline = false,
        dim_inactive = false,
        lualine_bold = false,
      })
      
      vim.cmd[[colorscheme tokyonight]]
    end
  },

  {
    "neoclide/coc.nvim",
    branch = "release",
    config = function()
      vim.opt.backup = false
      vim.opt.writebackup = false
      
      vim.opt.updatetime = 300
      
      vim.opt.signcolumn = "yes"
      
      local keyset = vim.keymap.set
      
      local opts = {silent = true, noremap = true, expr = true, replace_keycodes = false}
      keyset("i", "<TAB>", 'coc#pum#visible() ? coc#pum#next(1) : v:lua.check_back_space() ? "<TAB>" : coc#refresh()', opts)
      keyset("i", "<S-TAB>", [[coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"]], opts)
      
      keyset("i", "<cr>", [[coc#pum#visible() ? coc#pum#confirm() : "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"]], opts)
      
      keyset("i", "<c-space>", "coc#refresh()", {silent = true, expr = true})
      
      keyset("n", "gd", "<Plug>(coc-definition)", {silent = true})
      keyset("n", "gy", "<Plug>(coc-type-definition)", {silent = true})
      keyset("n", "gi", "<Plug>(coc-implementation)", {silent = true})
      keyset("n", "gr", "<Plug>(coc-references)", {silent = true})
      
      function _G.show_docs()
          local cw = vim.fn.expand('<cword>')
          if vim.fn.index({'vim', 'help'}, vim.bo.filetype) >= 0 then
              vim.api.nvim_command('h ' .. cw)
          elseif vim.api.nvim_eval('coc#rpc#ready()') then
              vim.fn.CocActionAsync('doHover')
          else
              vim.api.nvim_command('!' .. vim.o.keywordprg .. ' ' .. cw)
          end
      end
      keyset("n", "K", '<CMD>lua _G.show_docs()<CR>', {silent = true})
      
      vim.api.nvim_create_autocmd("VimEnter", {
        command = "CocInstall coc-tsserver",
        once = true
      })
    end
  }
}

require("lazy").setup(plugins)
