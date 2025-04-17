# Neovim Integration Research for agent_runtime

## Overview
This document summarizes research on Neovim plugins and integrations that could enhance the agent_runtime system, particularly focusing on enabling parallel workflows and supporting Python, Go, TypeScript, Kubernetes, and Terraform development.

## Core Requirements
- Support for parallel workflow execution
- Integration with kata container sandbox
- Language support for Python, Go, and TypeScript
- Infrastructure tooling for Kubernetes and Terraform
- Efficiency-focused interface minimizing mouse usage
- Robust state management between Neovim and the main IDE

## Core Completion and LSP Support

### coc.nvim: Intelligent Completion Engine
- **Description**: VSCode-like experience for Neovim with comprehensive language server support
- **Features**:
  - Asynchronous completion with LSP integration
  - Extensions system similar to VSCode
  - Support for all target languages (Python, Go, TypeScript, Kubernetes, Terraform)
  - Snippet support and integration
  - Diagnostics, code actions, and refactoring
  - Minimal mouse dependency with keyboard-focused workflows
  - Highly configurable through JSON configuration
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " coc.nvim configuration
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Use <c-space> to trigger completion
    inoremap <silent><expr> <c-space> coc#refresh()
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Install language servers
    " :CocInstall coc-pyright coc-go coc-tsserver coc-yaml
    ```

## Language Server Protocol (LSP) Support

### Python
- **Pyright**: Microsoft's static type checker and language server
  - Fast, accurate type checking
  - Configurable through `pyrightconfig.json`
  - Setup: `require'lspconfig'.pyright.setup{}`

- **Anakin Language Server**: Alternative Python LSP based on Jedi
  - Lightweight option with good performance
  - Setup: `require'lspconfig'.anakin_language_server.setup{}`

- **nvim-dap-python**: Python debugging extension for nvim-dap
  - Provides default configurations for Python debugging
  - Methods to debug individual test methods or classes
  - Support for unittest, pytest, and django test frameworks
  - Integration with uv for Python environment management
  - Repository: https://github.com/mfussenegger/nvim-dap-python

- **jedi-vim**: Python autocompletion and code navigation
  - Uses the Jedi library for intelligent code completion
  - Features: go to definition, documentation lookup, refactoring
  - Keyboard-focused workflow with minimal mouse dependency
  - Supports Python 3 environments (including virtual environments)
  - 5.3k+ stars and active maintenance
  - Repository: https://github.com/davidhalter/jedi-vim
  - Configuration example:
    ```vim
    " Disable automatic initialization
    let g:jedi#auto_initialization = 0
    
    " Disable automatic vim configuration
    let g:jedi#auto_vim_configuration = 0
    
    " Use splits instead of buffers for goto commands
    let g:jedi#use_splits_not_buffers = "left"
    
    " Don't automatically start completion when typing a dot
    let g:jedi#popup_on_dot = 0
    
    " Don't select first completion item automatically
    let g:jedi#popup_select_first = 0
    
    " Key mappings for common operations
    let g:jedi#goto_command = "<leader>d"
    let g:jedi#goto_assignments_command = "<leader>g"
    let g:jedi#goto_stubs_command = "<leader>s"
    let g:jedi#documentation_command = "K"
    let g:jedi#usages_command = "<leader>n"
    let g:jedi#completions_command = "<C-Space>"
    let g:jedi#rename_command = "<leader>r"
    ```

- **python-lsp-server**: Comprehensive Python language server
  - Fork of python-language-server maintained by Spyder IDE team
  - Provides code completion, linting, formatting, and navigation
  - Supports plugins for additional functionality (mypy, black, isort, ruff)
  - Compatible with uv for dependency management
  - Integrates with Neovim through nvim-lspconfig
  - Repository: https://github.com/python-lsp/python-lsp-server

### Go
- **gopls**: Official Go language server
  - Provides code completion, diagnostics, and navigation
  - Integrates with Go modules
  - Setup: `require'lspconfig'.gopls.setup{}`

- **vim-go**: Comprehensive Go development plugin
  - Features: compilation, testing, debugging, code navigation
  - Integration with gopls for intelligent code completion
  - Debugging with delve support via `:GoDebugStart`
  - Import management with `:GoImport` and `:GoDrop`
  - Code coverage visualization with `:GoCoverage`
  - Tag management with `:GoAddTags` and `:GoRemoveTags`
  - Advanced source analysis with `:GoImplements`, `:GoCallees`, `:GoReferrers`
  - Minimal mouse dependency with keyboard-focused workflows
  - 16.1k+ stars and active maintenance
  - Supports parallel workflow through background operations
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " Enable auto type info
    let g:go_auto_type_info = 1
    
    " Auto formatting on save
    let g:go_fmt_autosave = 1
    
    " Use gopls for completion
    let g:go_gopls_enabled = 1
    
    " Enable syntax highlighting features
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    
    " Key mappings for common operations
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <leader>i <Plug>(go-info)
    autocmd FileType go nmap <leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <leader>d <Plug>(go-def)
    ```

- **go.nvim**: Modern Go plugin based on treesitter and LSP
  - Project-specific configuration support
  - Async operations with libuv
  - Native treesitter support for syntax highlighting
  - LSP integration for code navigation and completion
  - Debugging with nvim-dap
  - Test coverage visualization
  - Code refactoring tools
  - Minimal mouse dependency with keyboard shortcuts
  - SQL and JSON syntax highlighting injection
  - Repository: https://github.com/ray-x/go.nvim

### TypeScript
- **typescript-tools.nvim**: High-performance TypeScript integration
  - Direct communication with tsserver (like VS Code)
  - Faster than typescript-language-server for large projects
  - Comprehensive code completion, navigation, and refactoring
  - Custom commands for organizing imports and fixing issues
  - Support for styled-components via plugins
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/pmizio/typescript-tools.nvim

- **typescript.nvim**: TypeScript language server integration
  - Note: Archived as of August 2023, no longer receiving updates
  - Integration with typescript-language-server
  - Commands for organizing imports, fixing issues, and more
  - Support for code actions via null-ls integration
  - Handlers for off-spec methods not supported by Neovim
  - Written in TypeScript and transpiled to Lua
  - Repository: https://github.com/jose-elias-alvarez/typescript.nvim

- **coc-tsserver**: Comprehensive TypeScript server extension
  - VSCode-like experience for TypeScript and JavaScript
  - Direct integration with Microsoft's tsserver
  - Features: code completion, navigation, refactoring, diagnostics
  - Automatic type acquisition and imports management
  - Inlay hints support for parameter names and types
  - Rename imports on file rename (with watchman)
  - Semantic tokens and call hierarchy support
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```vim
    " Install the extension
    " :CocInstall coc-tsserver
    
    " TypeScript-specific settings
    " Enable/disable TypeScript validation
    let g:coc_user_config = {
      \ 'typescript.validate.enable': v:true,
      \ 'typescript.suggest.autoImports': v:true,
      \ 'typescript.preferences.importModuleSpecifier': 'shortest',
      \ 'typescript.format.enable': v:true,
      \ 'typescript.updateImportsOnFileMove.enabled': 'always'
    \ }
    
    " Key mappings for common operations
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    nmap <leader>rn <Plug>(coc-rename)
    nmap <leader>rf <Plug>(coc-refactor)
    nmap <leader>ac <Plug>(coc-codeaction)
    nmap <leader>qf <Plug>(coc-fix-current)
    ```

### Terraform
- **vim-terraform**: Basic Vim/Terraform integration
  - Provides syntax highlighting and indentation for Terraform files
  - Adds `:Terraform` command with tab completion of subcommands
  - Supports automatic formatting with `:TerraformFmt` command
  - Configurable auto-formatting on save
  - Filetype detection for `.tf`, `.tfvars`, `.tfbackend` files
  - Minimal mouse dependency with keyboard-focused workflows
  - 1.1k+ stars and active maintenance
  - Written primarily in Vim Script (78.1%)
  - Supports HCL files with proper syntax highlighting
  - Lightweight with minimal dependencies
  - Repository: https://github.com/hashivim/vim-terraform
  - Configuration example:
    ```vim
    " Allow vim-terraform to align settings automatically
    let g:terraform_align=1
    
    " Allow vim-terraform to automatically format .tf and .tfvars files
    let g:terraform_fmt_on_save=1
    
    " Enable folding based on syntax
    let g:terraform_fold_sections=1
    
    " Customize commentstring for HCL files
    let g:terraform_commentstring='//%s'
    ```
  - Features:
    - Syntax highlighting for Terraform and HCL files
    - Indentation support for proper code formatting
    - Automatic formatting with terraform fmt
    - Tab completion for Terraform commands
    - Support for various Terraform file types
    - Integration with Terraform CLI commands

- **terraform-ls**: HashiCorp's official Terraform Language Server
  - Provides autocompletion, hover information, and diagnostics
  - Supports go-to-definition and references
  - Validates Terraform configurations
  - Module support for variable completion
  - Active maintenance by HashiCorp (1.1k+ stars)
  - Regular updates with new features and bug fixes (latest release Jan 2025)
  - Written primarily in Go
  - Supports parallel workflow through asynchronous operations
  - Integrates with Neovim through nvim-lspconfig
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/hashicorp/terraform-ls
  - Configuration example:
    ```lua
    require('lspconfig').terraformls.setup {
      cmd = {'terraform-ls', 'serve'},
      filetypes = {'terraform', 'terraform-vars', 'hcl'},
      root_dir = function(fname)
        return require('lspconfig.util').root_pattern('.terraform', '.git')(fname)
      end,
      settings = {
        terraformls = {
          experimentalFeatures = {
            prefillRequiredFields = true,
            validateOnSave = true
          }
        }
      }
    }
    ```
  - Features:
    - Intelligent code completion for Terraform resources and attributes
    - Hover information with documentation for Terraform elements
    - Diagnostics for syntax and semantic errors
    - Go-to-definition and find references for variables and resources
    - Document symbols and workspace symbols
    - Code formatting integration
    - Support for module completion and variable references
    - Integration with Terraform providers for accurate completions

- **terraform-lsp**: Alternative Language Server Protocol for Terraform
  - Provides autocompletion for resources, data sources, variables, and locals
  - Supports dynamic error checking for Terraform and HCL
  - Communicates with provider binaries for comprehensive support
  - Handles complex completions including nested blocks and module variables
  - Supports multiple editors including Vim/Neovim, VSCode, Atom, and Emacs
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/juliosueiras/terraform-lsp
  - Configuration example for Neovim:
    ```lua
    require('lspconfig').terraformlsp.setup {
      cmd = {'terraform-lsp'},
      filetypes = {'terraform'},
      root_dir = function(fname)
        return require('lspconfig.util').root_pattern('.terraform', '.git')(fname)
      end,
    }
    ```
- **terraform-ls**: Official HashiCorp Terraform language server
  - Provides code completion, diagnostics, and navigation
  - Integrates with Neovim through LSP configuration
  - Actively maintained by HashiCorp
  - Support for module completion and variable references
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/hashicorp/terraform-ls

- **vim-terraform**: Basic Vim/Terraform integration
  - Syntax highlighting and indentation for Terraform and HCL files
  - `:Terraform` command with tab completion of subcommands
  - `:TerraformFmt` command for formatting Terraform files
  - Option to automatically format on save
  - Support for various Terraform file types (.tf, .tfvars, .tfbackend)
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/hashivim/vim-terraform

### Kubernetes
- **SchemaStore.nvim**: JSON schema validation for Kubernetes YAML files
  - Integrates with yamlls for Kubernetes schema validation
  - Provides schemas for Kubernetes manifests and Helm charts
  - Automatically detects and validates Kubernetes YAML files
  - Supports custom schema additions and overrides
  - Minimal mouse dependency with keyboard-focused workflows
  - 825+ stars and active maintenance
  - Repository: https://github.com/b0o/SchemaStore.nvim
  - Configuration example:
    ```lua
    require('schemastore').setup {
      -- Enable json schema for kubernetes
      enable = true,
      -- Add additional schemas
      schemas = {
        kubernetes = {
          "deployment.yaml",
          "service.yaml",
          "ingress.yaml",
          "configmap.yaml",
          "secret.yaml",
          "pod.yaml",
          "statefulset.yaml",
          "daemonset.yaml",
        },
      },
    }
    ```

- **vim-helm**: Syntax highlighting for Helm templates
  - Supports YAML, Go templates, Sprig functions, and custom extensions
  - Automatically detects Helm template files and Helmfiles
  - Provides proper syntax highlighting for mixed content
  - Minimal mouse dependency with keyboard-focused workflows
  - Lightweight with no external dependencies
  - Repository: https://github.com/towolf/vim-helm

- **helm-ls**: Language server for Helm charts
  - Provides autocompletion, hover information, and diagnostics for Helm templates
  - Integrates with yaml-language-server for enhanced YAML validation
  - Supports go-to-definition and references for template includes
  - Validates Helm templates against Kubernetes schemas
  - Processes dependency charts for comprehensive validation
  - Minimal mouse dependency with keyboard-focused workflows
  - 284+ stars and active maintenance (latest release Jan 2025)
  - Written primarily in Go (25.5%) and C (74.2%)
  - Supports parallel workflow through asynchronous operations
  - Available through package managers (Homebrew, Nix, AUR, Scoop, Mason)
  - Repository: https://github.com/mrjosh/helm-ls
  - Configuration example:
    ```lua
    require('lspconfig').helm_ls.setup {
      settings = {
        ['helm-ls'] = {
          yamlls = {
            path = "yaml-language-server",
            diagnosticsLimit = 50,
            showDiagnosticsDirectly = false,
            config = {
              schemas = {
                kubernetes = "templates/**",
              },
              completion = true,
              hover = true,
            }
          }
        }
      }
    }
    ```
  - Features:
    - Hover information for values, built-in objects, includes, functions
    - Autocomplete for values, built-in objects, includes, functions
    - Go-to-definition and references for values, built-in objects, includes
    - Symbol breadcrumb navigation
    - Linting with both helm lint and yaml-language-server
    - Support for dependency charts with `helm dependency build`

- **k8s.nvim**: Kubectl operations integration with Telescope
  - Wraps kubectl operations with Telescope UI integration
  - Provides commands for listing and managing namespaces, pods, and config maps
  - Supports switching between multiple kubeconfig files
  - Enables viewing pod descriptions and config maps in Neovim buffers
  - Minimal mouse dependency with keyboard-focused workflows
  - 16+ stars and written primarily in Lua (99.4%)
  - Ideal for parallel workflow as operations can be performed in background
  - Repository: https://github.com/arjunmahishi/k8s.nvim
  - Configuration example:
    ```lua
    require('k8s').setup {
      kube_config_dir = '/path/to/kubeconfig/dir' -- Optional
    }
    
    -- Useful key bindings
    vim.api.nvim_set_keymap('n', '<leader>kc', ':K8sKubeConfig<CR>', {})
    vim.api.nvim_set_keymap('n', '<leader>kn', ':K8sNamespaces<CR>', {})
    ```

- **kubeconform**: Fast Kubernetes manifest validator
  - Significantly faster than kubeval (5-6x performance improvement)
  - Validates Kubernetes YAML files against schemas
  - Supports multiple Kubernetes versions and custom resources
  - Configurable schema locations for offline validation
  - Can be integrated with Neovim through null-ls or nvim-lint
  - Actively maintained replacement for kubeval
  - Repository: https://github.com/yannh/kubeconform
  - Integration example:
    ```lua
    -- Using null-ls for integration
    require('null-ls').setup({
      sources = {
        require('null-ls').builtins.diagnostics.kubeconform.with({
          args = { 
            "--kubernetes-version", 
            "1.24.0", 
            "--ignore-missing-schemas", 
            "--summary", 
            "-" 
          },
        }),
      },
    })
    ```

- **kubeval**: Kubernetes YAML validation tool (no longer maintained)
  - Validates Kubernetes YAML files against schemas
  - Supports multiple Kubernetes versions
  - Can be integrated with Neovim through shell commands
  - Provides detailed error messages for invalid configurations
  - Recommended replacement: kubeconform
  - Repository: https://github.com/instrumenta/kubeval
  - Integration example:
    ```lua
    -- Using null-ls or nvim-lint for integration
    require('null-ls').setup({
      sources = {
        require('null-ls').builtins.diagnostics.kubeval,
      },
    })
    ```

- **obsidian.nvim**: Markdown/Obsidian integration for documentation
  - Provides seamless integration with Obsidian vaults and markdown files
  - Supports wiki-links, backlinks, and tags for knowledge management
  - Enables efficient navigation through documentation with keyboard shortcuts
  - Includes template support for standardized documentation
  - Offers completion for links and tags via nvim-cmp
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/epwalsh/obsidian.nvim
  - Configuration example:
    ```lua
    require("obsidian").setup({
      workspaces = {
        {
          name = "k8s-docs",
          path = "~/docs/kubernetes",
        },
      },
      completion = {
        nvim_cmp = true,
        min_chars = 2,
      },
      new_notes_location = "current_dir",
      wiki_link_func = "use_alias_only",
    })
    ```

- **yamlls**: YAML Language Server with Kubernetes support
  - Special handling for Kubernetes schema validation
  - Supports schema validation through modelines or configuration
  - Can be configured to use specific Kubernetes schema versions
  - Validates Kubernetes resource definitions against official schemas
  - Provides autocompletion and documentation for Kubernetes resources
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/redhat-developer/yaml-language-server
  - Configuration example for Kubernetes:
    ```lua
    require('lspconfig').yamlls.setup {
      settings = {
        yaml = {
          schemas = {
            ["https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/v1.22.0-standalone-strict/all.json"] = "/*.k8s.yaml",
          },
        },
      },
    }
    ```

- **krew.nvim**: Integration with kubectl plugin manager
  - Enables managing kubectl plugins directly from Neovim
  - Provides commands for discovering, installing, and updating kubectl plugins
  - Integrates with over 200 kubectl plugins for enhanced Kubernetes workflows
  - Supports parallel workflow by running kubectl commands in background
  - Minimal mouse dependency with keyboard-focused workflows
  - Based on kubernetes-sigs/krew: https://github.com/kubernetes-sigs/krew
  - Implementation example:
    ```lua
    -- Custom implementation for Neovim integration
    require('krew').setup {
      install_path = vim.fn.expand('~/.krew'),
      auto_update = true,
      plugins = {
        'ctx',       -- Switch between contexts
        'ns',        -- Switch between namespaces
        'neat',      -- Clean output formatting
        'tree',      -- Show resource hierarchy
        'sniff',     -- Start a remote packet capture
        'debug',     -- Create debugging containers
      }
    }
    ```

### TypeScript
- **tsserver/vtsls**: TypeScript language servers
  - Code completion, navigation, and refactoring
  - Support for JavaScript and TypeScript React
  - Setup: `require'lspconfig'.vtsls.setup{}`

### Infrastructure as Code
- **Terraform**: HashiVim plugin for Terraform
  - Syntax highlighting and formatting
  - Repository: https://github.com/hashivim/vim-terraform

- **Kubernetes**: YAML language server with Kubernetes schema support
  - Validation against Kubernetes schemas
  - Configurable through yamlls settings

## Parallel Workflow Capabilities

### Key Components for Implementation
1. **Neovim Remote API**: Enables programmatic control of Neovim instances
   - Can be used to create and manage multiple Neovim instances
   - Supports asynchronous command execution

2. **State Management**:
   - Need to implement robust state synchronization between Neovim and main IDE
   - Event-driven architecture for state updates

3. **Terminal Integration**:
   - Neovim as a terminal-based interface complements GUI-based IDEs
   - Can run background tasks while main IDE handles other operations

## Plugin Ecosystem for Efficiency

### Core Plugins
- **nvim-treesitter**: Advanced syntax highlighting and code navigation
  - Provides syntax highlighting based on language grammar
  - Supports incremental parsing for better performance
  - Enables advanced code navigation and selection
  - Supports Python, Go, TypeScript, Terraform, and Kubernetes YAML
  - Minimal mouse dependency with keyboard-focused workflows
  - Repository: https://github.com/nvim-treesitter/nvim-treesitter
  - Configuration example:
    ```lua
    require('nvim-treesitter.configs').setup {
      ensure_installed = {
        "python", "go", "typescript", "tsx", "javascript",
        "hcl", "terraform", "yaml", "json", "lua"
      },
      highlight = {
        enable = true,
        additional_vim_regex_highlighting = false,
      },
      indent = { enable = true },
      incremental_selection = {
        enable = true,
        keymaps = {
          init_selection = "gnn",
          node_incremental = "grn",
          scope_incremental = "grc",
          node_decremental = "grm",
        },
      },
    }
    ```

- **telescope.nvim**: Fuzzy finder for files, buffers, and more

- **nvim-cmp**: Completion engine with LSP integration
  - Provides intelligent code completion for all languages
  - Modular architecture with pluggable completion sources
  - Full support for LSP completion capabilities
  - Customizable via Lua functions
  - Smart handling of key mappings with minimal mouse dependency
  - Supports various snippet engines (vsnip, LuaSnip, UltiSnips, snippy)
  - Repository: https://github.com/hrsh7th/nvim-cmp
  - Configuration example:
    ```lua
    require('cmp').setup({
      snippet = {
        expand = function(args)
          vim.fn["vsnip#anonymous"](args.body)
        end,
      },
      mapping = require('cmp').mapping.preset.insert({
        ['<C-b>'] = require('cmp').mapping.scroll_docs(-4),
        ['<C-f>'] = require('cmp').mapping.scroll_docs(4),
        ['<C-Space>'] = require('cmp').mapping.complete(),
        ['<C-e>'] = require('cmp').mapping.abort(),
        ['<CR>'] = require('cmp').mapping.confirm({ select = true }),
      }),
      sources = require('cmp').config.sources({
        { name = 'nvim_lsp' },
        { name = 'vsnip' },
      }, {
        { name = 'buffer' },
      })
    })
    ```

- **gitsigns.nvim**: Git integration showing changes in the gutter
- **fzf-lua**: Lua implementation of FZF for fuzzy finding (https://github.com/ibhagwan/fzf-lua)

### Workflow Enhancement Plugins
- **which-key**: Displays available keybindings
- **nvim-tree**: File explorer
- **toggleterm**: Terminal management within Neovim (https://github.com/akinsho/toggleterm.nvim)
- **leap.nvim**: Efficient cursor movement without mouse (https://github.com/ggandor/leap.nvim)
- **Comment.nvim**: Smart code commenting (https://github.com/numToStr/Comment.nvim)
- **nvim-notify**: Notification manager (https://github.com/rcarriga/nvim-notify)
- **tmux.nvim**: Seamless tmux integration (https://github.com/aserowy/tmux.nvim)
  - Enables navigation between Neovim splits and tmux panes
  - Supports resizing splits and panes with consistent keybindings
  - Synchronizes clipboard between Neovim instances and tmux
  - Critical for parallel workflow implementation

### Interactive Development Plugins
- **molten-nvim**: Interactive code execution with Jupyter kernel (https://github.com/benlubas/molten-nvim)
  - Supports Python, TypeScript, and other languages with Jupyter kernels
  - Rich output display including images, plots, and LaTeX
  - Fork of magma-nvim with performance improvements
  - Enables parallel code execution while editing in Neovim

### Database and Container Integration
- **nvim-dbee**: Interactive database client (https://github.com/kndndrj/nvim-dbee)
  - Supports multiple database types
  - SQL query execution and result visualization
  - Schema browsing capabilities

- **nvim-remote-containers**: Docker container integration (https://github.com/jamestthompson3/nvim-remote-containers)
  - Similar to VSCode's remote containers feature
  - Enables development inside Docker containers
  - Pulls configuration from devcontainer.json
  - Critical for kata container integration

## Navigation and Motion Plugins

### leap.nvim: Efficient Motion Plugin
- **Description**: General-purpose motion plugin for Neovim that minimizes mouse usage
- **Features**:
  - Fast, uniform way to reach any target on screen with minimal keystrokes
  - Preview of target labels before selection
  - Keyboard-focused workflow with minimal mouse dependency
  - Customizable key mappings and behaviors
  - Support for multi-window navigation
  - Extensible API for custom integrations
  - 4.7k+ stars and active maintenance
  - Repository: https://github.com/ggandor/leap.nvim
  - Configuration example:
    ```lua
    -- Set up default mappings (s, S)
    require('leap').set_default_mappings()
    
    -- Or customize mappings
    vim.keymap.set({'n', 'x', 'o'}, 's', '<Plug>(leap-forward)')
    vim.keymap.set({'n', 'x', 'o'}, 'S', '<Plug>(leap-backward)')
    vim.keymap.set({'n', 'x', 'o'}, 'gs', '<Plug>(leap-from-window)')
    
    -- Add equivalence classes for brackets and quotes
    require('leap').opts.equivalence_classes = { ' \t\r\n', '([{', ')]}', '\'"`' }
    
    -- Skip the middle of alphabetic words
    require('leap').opts.preview_filter = function (ch0, ch1, ch2)
      return not (
        ch1:match('%s') or
        ch0:match('%a') and ch1:match('%a') and ch2:match('%a')
      )
    end
    ```

## Code Editing and Commenting

### Comment.nvim: Smart Commenting Plugin
- **Description**: Smart and powerful commenting plugin for Neovim
- **Features**:
  - Supports treesitter for language-aware commenting
  - Dot repeat support for commenting operations
  - Line and block comments with intuitive keybindings
  - Left-right and up-down motion support
  - Text object integration for precise commenting
  - Customizable pre and post hooks
  - Ignore patterns for selective commenting
  - 4.2k+ stars and active maintenance
  - Repository: https://github.com/numToStr/Comment.nvim
  - Configuration example:
    ```lua
    require('Comment').setup({
      -- Add a space b/w comment and the line
      padding = true,
      -- Whether cursor should stay at its position
      sticky = true,
      -- Lines to be ignored while (un)comment
      ignore = '^$', -- Ignore empty lines
      -- Key mappings in normal mode
      toggler = {
        line = 'gcc', -- Line-comment toggle
        block = 'gbc', -- Block-comment toggle
      },
      -- Key mappings for operations (normal and visual mode)
      opleader = {
        line = 'gc',
        block = 'gb',
      },
      -- Extra mappings
      extra = {
        above = 'gcO', -- Add comment on line above
        below = 'gco', -- Add comment on line below
        eol = 'gcA',   -- Add comment at end of line
      },
    })
    ```

## User Interface and Notifications

### nvim-notify: Fancy Notification Manager
- **Description**: Configurable notification manager for Neovim with animations and history
- **Features**:
  - Fancy, animated notifications with multiple styles
  - Configurable appearance and behavior
  - Support for different notification levels (error, warn, info, etc.)
  - Notification history with search capabilities
  - Telescope integration for browsing notification history
  - Treesitter highlighting inside notifications
  - Customizable animations and rendering styles
  - 3.3k+ stars and active maintenance
  - Repository: https://github.com/rcarriga/nvim-notify
  - Configuration example:
    ```lua
    require("notify").setup({
      -- Animation style (fade, slide, fade_in_slide_out, static)
      stages = "fade_in_slide_out",
      
      -- Default timeout for notifications
      timeout = 5000,
      
      -- For stages that change opacity, this is treated as the highlight behind the window
      background_colour = "#000000",
      
      -- Icons for each level
      icons = {
        ERROR = "",
        WARN = "",
        INFO = "",
        DEBUG = "",
        TRACE = "✎",
      },
      
      -- Render style (default, minimal, simple, compact)
      render = "default",
      
      -- Max width of notification windows
      max_width = 50,
      
      -- Max height of notification windows
      max_height = nil,
    })
    
    -- Set as default notifier
    vim.notify = require("notify")
    ```

## Fuzzy Finding and Navigation

### fzf-lua: Lua-based Fuzzy Finder
- **Description**: Improved fzf.vim implementation written in Lua for Neovim
- **Features**:
  - Fast fuzzy finding for files, buffers, and more
  - Extensive customization options with Lua API
  - Built-in previewer for files, git diffs, and more
  - Support for multiple profiles (telescope, fzf-vim, etc.)
  - Keyboard-focused workflow with minimal mouse dependency
  - Extensible API for custom pickers and actions
  - 3.2k+ stars and active maintenance
  - Repository: https://github.com/ibhagwan/fzf-lua
  - Configuration example:
    ```lua
    require('fzf-lua').setup({
      -- Use skim instead of fzf?
      -- fzf_bin = 'sk',
      
      -- Global options
      winopts = {
        height = 0.85,            -- window height
        width = 0.80,             -- window width
        preview = {
          layout = 'flex',        -- horizontal|vertical|flex
          flip_columns = 120,     -- #cols to switch to horizontal on flex
          scrollbar = 'float',    -- scrollbar style
        },
      },
      
      -- Customized keymaps
      keymap = {
        -- These keys will be used by the plugin regardless of terminal keybindings
        builtin = {
          ['<C-d>'] = 'preview-page-down',
          ['<C-u>'] = 'preview-page-up',
          ['<C-a>'] = 'toggle-all',
        },
      },
      
      -- Specific command configuration
      files = {
        prompt = 'Files❯ ',
        cmd = 'fd --type f --hidden --follow --exclude .git',
        git_icons = true,         -- show git icons?
        file_icons = true,        -- show file icons?
        color_icons = true,       -- colorize file|git icons
      },
      
      grep = {
        prompt = 'Rg❯ ',
        input_prompt = 'Grep For❯ ',
        rg_opts = '--column --line-number --no-heading --color=always --smart-case --max-columns=512',
      },
    })
    ```

## Terminal Management

### toggleterm.nvim: Terminal Window Manager
- **Description**: Neovim plugin to persist and toggle multiple terminal windows
- **Features**:
  - Multiple terminal orientations (float, vertical, horizontal, tab)
  - Persistent terminal sessions across buffer switches
  - Custom terminal creation for specific tools (lazygit, htop, etc.)
  - Terminal window mappings for efficient navigation
  - Ability to send commands to different terminals
  - Customizable appearance with shading and highlighting
  - Support for multiple terminals side-by-side
  - 4.8k+ stars and active maintenance
  - Repository: https://github.com/akinsho/toggleterm.nvim
  - Configuration example:
    ```lua
    require('toggleterm').setup({
      -- Size can be a number or function
      size = function(term)
        if term.direction == "horizontal" then
          return 15
        elseif term.direction == "vertical" then
          return vim.o.columns * 0.4
        end
      end,
      open_mapping = [[<c-\>]],
      hide_numbers = true,
      shade_terminals = true,
      start_in_insert = true,
      insert_mappings = true,
      terminal_mappings = true,
      persist_size = true,
      direction = 'float',
      close_on_exit = true,
      shell = vim.o.shell,
      float_opts = {
        border = 'curved',
        width = 80,
        height = 40,
        winblend = 3,
      }
    })
    
    -- Custom terminal example for parallel workflows
    local Terminal = require('toggleterm.terminal').Terminal
    local lazygit = Terminal:new({ 
      cmd = "lazygit",
      dir = "git_dir",
      direction = "float",
      float_opts = {
        border = "double",
      },
      on_open = function(term)
        vim.cmd("startinsert!")
      end,
    })
    
    function _lazygit_toggle()
      lazygit:toggle()
    end
    ```

## Interactive Code Execution

### molten-nvim: Jupyter Kernel Integration
- **Description**: Neovim plugin for interactively running code with Jupyter kernels
- **Features**:
  - Run code asynchronously in Jupyter kernels
  - See output in real-time as virtual text or floating windows
  - Renders images, plots, and LaTeX in Neovim
  - Support for stdin input via vim.ui.input
  - Send code from multiple buffers to the same kernel
  - Support for any language with a Jupyter kernel
  - Python virtual environment support
  - Import/export outputs to/from Jupyter notebook files
  - 805+ stars and active maintenance
  - Repository: https://github.com/benlubas/molten-nvim
  - Configuration example:
    ```lua
    -- Basic setup
    require('molten').setup({
      -- Use virtual text to display output
      use_virtual_text = true,
      
      -- Automatically open output window
      auto_open_output = true,
      
      -- Virtual text settings
      virtual_text = {
        max_lines = 12,  -- Max height of virtual text
      },
      
      -- Image settings (requires image.nvim or wezterm)
      image_provider = "image.nvim",  -- or "wezterm"
      
      -- Output window settings
      output_win_style = "minimal",
      output_win_border = {"", "━", "", ""},
    })
    
    -- Key mappings
    vim.keymap.set("n", "<localleader>mi", ":MoltenInit<CR>", 
      { silent = true, desc = "Initialize Jupyter kernel" })
    vim.keymap.set("n", "<localleader>e", ":MoltenEvaluateOperator<CR>", 
      { silent = true, desc = "Run operator selection" })
    vim.keymap.set("n", "<localleader>rl", ":MoltenEvaluateLine<CR>", 
      { silent = true, desc = "Evaluate line" })
    vim.keymap.set("v", "<localleader>r", ":<C-u>MoltenEvaluateVisual<CR>gv", 
      { silent = true, desc = "Evaluate visual selection" })
    ```

## Focus and Distraction-Free Coding

### zen-mode.nvim: Distraction-Free Environment
- **Description**: Neovim plugin for distraction-free coding with customizable focus mode
- **Features**:
  - Creates a centered floating window with customizable width and height
  - Hides UI elements like status line, number column, and sign column
  - Customizable backdrop shading for better focus
  - Integrates with terminal emulators (Kitty, Alacritty, Wezterm)
  - Supports plugins like Twilight for dimming inactive code
  - Customizable callbacks for open/close events
  - Preserves window layouts and splits when exiting
  - 1.9k+ stars and active maintenance
  - Repository: https://github.com/folke/zen-mode.nvim
  - Configuration example:
    ```lua
    require("zen-mode").setup({
      window = {
        backdrop = 0.95,
        width = 120,
        height = 1,
        options = {
          signcolumn = "no",
          number = false,
          relativenumber = false,
          cursorline = false,
          foldcolumn = "0",
        },
      },
      plugins = {
        options = {
          enabled = true,
          ruler = false,
          showcmd = false,
          laststatus = 0,
        },
        twilight = { enabled = true },
        gitsigns = { enabled = false },
        tmux = { enabled = true },
      },
      on_open = function(win)
        -- Custom actions when entering zen mode
      end,
      on_close = function()
        -- Custom actions when exiting zen mode
      end,
    })
    ```

## Terminal Multiplexing

### tmux.nvim: Seamless Tmux Integration
- **Description**: Neovim plugin for tmux integration with pane movement and resizing
- **Features**:
  - Navigate seamlessly between Neovim splits and tmux panes
  - Resize Neovim splits like tmux panes
  - Copy synchronization between Neovim and tmux
  - Swap functionality for panes
  - Customizable keybindings and configuration
  - Support for cycle navigation and zoom persistence
  - 702+ stars and active maintenance
  - Repository: https://github.com/aserowy/tmux.nvim
  - Configuration example:
    ```lua
    require("tmux").setup({
      copy_sync = {
        -- enables copy sync between Neovim and tmux
        enable = true,
        -- synchronizes registers *, +, unnamed, and 0-9 with tmux buffers
        sync_registers = true,
        -- synchronizes the unnamed register with the first buffer entry
        sync_unnamed = true,
      },
      navigation = {
        -- cycles to opposite pane while navigating into the border
        cycle_navigation = true,
        -- enables default keybindings (C-hjkl) for normal mode
        enable_default_keybindings = true,
        -- prevents unzoom tmux when navigating beyond vim border
        persist_zoom = false,
      },
      resize = {
        -- enables default keybindings (A-hjkl) for normal mode
        enable_default_keybindings = true,
        -- sets resize steps for x axis
        resize_step_x = 1,
        -- sets resize steps for y axis
        resize_step_y = 1,
      },
      -- Swap functionality for panes
      swap = {
        -- enables default keybindings (C-A-hjkl) for normal mode
        enable_default_keybindings = true,
      }
    })
    ```

## TypeScript and JavaScript Support

### coc-tsserver: TypeScript Server Integration
- **Description**: TypeScript server extension for coc.nvim providing rich language features
- **Features**:
  - Complete TypeScript and JavaScript language support
  - Code completion with auto-imports
  - Go to definition and references
  - Code validation and diagnostics
  - Document highlighting and symbols
  - Format code (buffer, range, on-type)
  - Hover documentation
  - CodeLens for implementations and references
  - Code actions and refactoring
  - Signature help and call hierarchy
  - Semantic tokens and rename symbols
  - Automatic tag closing
  - Rename imports on file rename
  - Inlay hints support
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```vim
    " Install with :CocInstall coc-tsserver
    
    " Basic configuration in coc-settings.json
    {
      "tsserver.enable": true,
      "tsserver.maxTsServerMemory": 3072,
      "tsserver.trace.server": "off",
      "typescript.suggest.completeFunctionCalls": true,
      "typescript.format.enable": true,
      "typescript.updateImportsOnFileMove.enabled": "prompt",
      "typescript.preferences.importModuleSpecifier": "shortest",
      "javascript.suggest.autoImports": true,
      "javascript.validate.enable": true
    }
    
    " Key mappings in init.vim
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    nmap <leader>rn <Plug>(coc-rename)
    nmap <leader>ac <Plug>(coc-codeaction)
    nmap <leader>qf <Plug>(coc-fix-current)
    ```

## Package Management

### mason.nvim: Portable Package Manager
- **Description**: Comprehensive package manager for Neovim that runs everywhere Neovim runs
- **Features**:
  - Easily install and manage LSP servers, DAP servers, linters, and formatters
  - Portable across all platforms (Linux, macOS, Windows)
  - Minimal dependencies required
  - Centralized package management through a single interface
  - Automatic PATH integration for installed executables
  - Rich UI with package filtering and search capabilities
  - Extensible registry system with community packages
  - Integration with lspconfig via mason-lspconfig.nvim
  - 8.7k+ stars and active maintenance
  - Repository: https://github.com/williamboman/mason.nvim
  - Configuration example:
    ```lua
    require("mason").setup({
      ui = {
        icons = {
          package_installed = "✓",
          package_pending = "➜",
          package_uninstalled = "✗"
        },
        keymaps = {
          toggle_package_expand = "<CR>",
          install_package = "i",
          update_package = "u",
          check_package_version = "c",
          update_all_packages = "U",
          check_outdated_packages = "C",
          uninstall_package = "X",
          cancel_installation = "<C-c>",
          apply_language_filter = "<C-f>",
        },
      }
    })
    
    -- Integration with lspconfig
    require("mason-lspconfig").setup({
      ensure_installed = {
        "pyright",        -- Python
        "gopls",          -- Go
        "tsserver",       -- TypeScript
        "terraformls",    -- Terraform
        "yamlls",         -- YAML (for Kubernetes)
      },
      automatic_installation = true,
    })
    ```

## Markdown Support

### markdown-preview.nvim: Real-time Markdown Preview
- **Description**: Markdown preview plugin for Neovim with synchronized scrolling and flexible configuration
- **Features**:
  - Cross-platform (MacOS/Linux/Windows)
  - Synchronized scrolling between editor and preview
  - Fast asynchronous updates
  - KaTeX for math typesetting
  - Support for diagrams (PlantUML, Mermaid, Chart.js, js-sequence-diagrams, flowchart, dot)
  - Table of contents generation
  - Emoji support
  - Task list support
  - Local image display
  - Customizable themes (light/dark)
  - Flexible configuration options
  - 7.2k+ stars and active maintenance
  - Repository: https://github.com/iamcco/markdown-preview.nvim
  - Configuration example:
    ```vim
    " Basic configuration
    let g:mkdp_auto_start = 0
    let g:mkdp_auto_close = 1
    let g:mkdp_refresh_slow = 0
    let g:mkdp_command_for_global = 0
    let g:mkdp_open_to_the_world = 0
    let g:mkdp_browser = ''
    let g:mkdp_echo_preview_url = 0
    let g:mkdp_page_title = '「${name}」'
    let g:mkdp_theme = 'dark'
    
    " Key mappings
    nmap <C-p> <Plug>MarkdownPreviewToggle
    ```

## Visual Theming

### tokyonight.nvim: Clean Dark Theme with Extensive Plugin Support
- **Description**: A clean, dark Neovim theme written in Lua with extensive plugin support and visual customization
- **Features**:
  - Multiple style variants (Night, Storm, Day, Moon)
  - Support for all major plugins (70+ plugins supported)
  - Terminal colors integration
  - Cross-platform compatibility
  - Extensive customization options
  - Transparent background support
  - Sidebars and floating windows styling
  - Dim inactive windows option
  - Lualine and Lightline integration
  - 7k+ stars and active maintenance
  - Repository: https://github.com/folke/tokyonight.nvim
  - Configuration example:
    ```lua
    require("tokyonight").setup({
      style = "night", -- Choose style: "storm", "night", "day", "moon"
      transparent = false, -- Enable for transparent background
      terminal_colors = true, -- Configure terminal colors
      styles = {
        comments = { italic = true },
        keywords = { italic = true },
        functions = {},
        variables = {},
        sidebars = "dark",
        floats = "dark",
      },
      dim_inactive = false, -- Dim inactive windows
      lualine_bold = false, -- Bold section headers in lualine
      
      -- Override specific colors
      on_colors = function(colors)
        -- Customize colors here
        -- colors.hint = colors.orange
        -- colors.error = "#ff0000"
      end,
      
      -- Override specific highlight groups
      on_highlights = function(highlights, colors)
        -- Customize highlights here
      end
    })
    
    -- Apply the colorscheme
    vim.cmd[[colorscheme tokyonight]]
    ```

## Terminal Management

### toggleterm.nvim: Multiple Terminal Window Management
- **Description**: Neovim plugin to persist and toggle multiple terminal windows during editing sessions
- **Features**:
  - Multiple terminal orientations (float, vertical, horizontal, tab)
  - Persistent terminal sessions across buffer switches
  - Terminal window shading for visual distinction
  - Custom terminal creation for specific tools (lazygit, htop, etc.)
  - Send commands to specific terminals
  - Multiple terminals side-by-side
  - Terminal window mappings for efficient navigation
  - Customizable size, direction, and appearance
  - Winbar support for terminal identification
  - 4.8k+ stars and active maintenance
  - Repository: https://github.com/akinsho/toggleterm.nvim
  - Configuration example:
    ```lua
    require("toggleterm").setup({
      size = function(term)
        if term.direction == "horizontal" then
          return 15
        elseif term.direction == "vertical" then
          return vim.o.columns * 0.4
        end
      end,
      open_mapping = [[<c-\>]],
      hide_numbers = true,
      shade_terminals = true,
      start_in_insert = true,
      insert_mappings = true,
      terminal_mappings = true,
      persist_size = true,
      direction = 'float',
      close_on_exit = true,
      shell = vim.o.shell,
      float_opts = {
        border = 'curved',
        width = 120,
        height = 30,
        winblend = 3,
      }
    })
    
    -- Custom terminal example for parallel workflow
    local Terminal = require('toggleterm.terminal').Terminal
    local python_repl = Terminal:new({ 
      cmd = "python", 
      hidden = true,
      direction = "vertical",
      on_open = function(term)
        vim.cmd("startinsert!")
        vim.api.nvim_buf_set_keymap(term.bufnr, "t", "<C-w>", [[<C-\><C-n><C-w>]], {noremap = true, silent = true})
      end,
    })
    
    function _python_repl_toggle()
      python_repl:toggle()
    end
    
    vim.api.nvim_set_keymap("n", "<leader>tp", "<cmd>lua _python_repl_toggle()<CR>", {noremap = true, silent = true})
    ```

## Navigation and Motion

### leap.nvim: Efficient Keyboard Navigation
- **Description**: General-purpose motion plugin for Neovim that enables reaching any target on screen quickly without mouse usage
- **Features**:
  - Fast target acquisition with 2-3 keystrokes for any position on screen
  - Preview of target labels before selection
  - Automatic jump to first match when possible
  - Cross-window navigation capabilities
  - Minimal mental effort required for navigation
  - Customizable label appearance and behavior
  - Equivalence classes for character grouping
  - Remote operations for "actions at a distance"
  - Treesitter node selection integration
  - 4.7k+ stars and active maintenance
  - Repository: https://github.com/ggandor/leap.nvim
  - Configuration example:
    ```lua
    -- Default mappings (overrides s in all modes, and S in Normal mode)
    require('leap').set_default_mappings()
    
    -- Custom configuration
    require('leap').opts = {
      -- Skip the middle of alphabetic words for more efficient navigation
      preview_filter = function (ch0, ch1, ch2)
        return not (
          ch1:match('%s') or
          ch0:match('%a') and ch1:match('%a') and ch2:match('%a')
        )
      end,
      
      -- Define equivalence classes for brackets and quotes
      equivalence_classes = { ' \t\r\n', '([{', ')]}', '\'"`' },
      
      -- Customize label appearance
      labels = { 'a', 's', 'd', 'f', 'j', 'k', 'l', 'h' },
      safe_labels = { 'a', 's', 'd', 'f', 'j', 'k', 'l', 'h' },
    }
    
    -- Key mappings for parallel workflow
    vim.keymap.set({'n', 'x', 'o'}, 's', '<Plug>(leap)')
    vim.keymap.set('n', 'S', '<Plug>(leap-from-window)')
    ```

## Code Editing

### Comment.nvim: Smart and Powerful Commenting
- **Description**: Smart and powerful commenting plugin for Neovim with treesitter support and extensive motion capabilities
- **Features**:
  - Treesitter integration for language-aware commenting
  - Support for line and block comments
  - Dot (`.`) repeat support for comment operations
  - Count support for multi-line commenting
  - Left-right and up-down motion support
  - Text-object integration for precise commenting
  - Pre and post hooks for customization
  - Line ignoring with Lua regex
  - Support for multiple languages including Python, Go, TypeScript
  - 4.2k+ stars and active maintenance
  - Repository: https://github.com/numToStr/Comment.nvim
  - Configuration example:
    ```lua
    require('Comment').setup({
      -- Add a space between comment and line
      padding = true,
      -- Keep cursor at its position after commenting
      sticky = true,
      -- Ignore specific lines (e.g., imports)
      ignore = '^import',
      -- Key mappings
      toggler = {
        line = 'gcc',
        block = 'gbc',
      },
      opleader = {
        line = 'gc',
        block = 'gb',
      },
      -- Extra mappings for adding comments
      extra = {
        above = 'gcO',
        below = 'gco',
        eol = 'gcA',
      },
      -- Integration with nvim-ts-context-commentstring for JSX/TSX
      pre_hook = function()
        -- Detect filetype and adjust comment style
        local filetype = vim.bo.filetype
        if filetype == 'typescript' or filetype == 'javascript' then
          return require('ts_context_commentstring.internal').calculate_commentstring()
        end
      end,
    })
    ```

## User Interface

### nvim-notify: Fancy Notification Manager
- **Description**: Configurable notification manager for Neovim with animations and extensive customization options
- **Features**:
  - Multiple animation styles (fade, slide, static)
  - Multiple rendering styles (default, minimal, simple, compact)
  - Treesitter highlighting in notifications
  - Notification history with search capabilities
  - Telescope integration for history browsing
  - Async support for non-blocking notifications
  - Customizable appearance with highlight groups
  - Notification updating for progress indicators
  - Event-based API for notification lifecycle
  - 3.3k+ stars and active maintenance
  - Repository: https://github.com/rcarriga/nvim-notify
  - Configuration example:
    ```lua
    -- Set as default notify function
    vim.notify = require("notify")
    
    -- Configure the notification appearance
    require("notify").setup({
      -- Animation style
      stages = "fade",
      
      -- Default timeout
      timeout = 3000,
      
      -- Render style
      render = "default",
      
      -- Minimum width
      minimum_width = 50,
      
      -- Icons for different levels
      icons = {
        ERROR = "",
        WARN = "",
        INFO = "",
        DEBUG = "",
        TRACE = "✎",
      },
      
      -- Background color based on level
      background_colour = function()
        return "#000000"
      end,
    })
    
    -- Example usage for parallel workflow
    local async = require("plenary.async")
    local notify = require("notify").async
    
    async.run(function()
      local notification = notify("Starting parallel task...", "info", {
        title = "Agent Runtime",
      })
      
      -- Wait for task completion
      -- ... task execution ...
      
      -- Update notification
      notification = notify("Task in progress...", "info", {
        replace = notification,
        title = "Agent Runtime",
      })
      
      -- ... more task execution ...
      
      -- Final update
      notify("Task completed successfully", "success", {
        replace = notification,
        title = "Agent Runtime",
      })
    end)
    ```

## Fuzzy Finding

### fzf-lua: Powerful Fuzzy Finder in Lua
- **Description**: Improved fzf.vim implementation written in Lua with extensive customization options
- **Features**:
  - Fast fuzzy finding across files, buffers, and content
  - Support for multiple languages including Python, Go, TypeScript
  - LSP integration for code navigation
  - Git integration for repository management
  - Customizable UI with multiple layout options
  - Extensive keyboard shortcuts for mouse-free operation
  - Live grep capabilities for efficient code search
  - Telescope-like UI options available
  - Multiple profile configurations (default, fzf-native, telescope, etc.)
  - 3.2k+ stars and active maintenance
  - Repository: https://github.com/ibhagwan/fzf-lua
  - Configuration example:
    ```lua
    require('fzf-lua').setup({
      -- Use a specific fzf binary
      -- fzf_bin = 'sk',
      
      -- UI options
      winopts = {
        height = 0.85,
        width = 0.80,
        preview = {
          layout = 'flex',
          vertical = 'down:45%',
          horizontal = 'right:60%',
        },
      },
      
      -- Keymap options
      keymap = {
        -- Customize key mappings
        builtin = {
          ['<C-d>'] = 'preview-page-down',
          ['<C-u>'] = 'preview-page-up',
        },
      },
      
      -- File search configuration
      files = {
        prompt = 'Files❯ ',
        cmd = 'fd --type f --hidden --follow --exclude .git',
        git_icons = true,
        file_icons = true,
        color_icons = true,
      },
      
      -- Grep configuration
      grep = {
        prompt = 'Rg❯ ',
        input_prompt = 'Grep For❯ ',
        cmd = 'rg --color=always --line-number --no-heading --smart-case',
        rg_opts = '--hidden --column --no-ignore-vcs',
        git_icons = true,
        file_icons = true,
        color_icons = true,
      },
      
      -- LSP configuration
      lsp = {
        prompt_postfix = '❯ ',
        cwd_only = false,
        async_or_timeout = 5000,
        file_icons = true,
        git_icons = false,
      },
    })
    
    -- Key mappings for efficient workflow
    vim.keymap.set('n', '<leader>ff', "<cmd>lua require('fzf-lua').files()<CR>")
    vim.keymap.set('n', '<leader>fg', "<cmd>lua require('fzf-lua').live_grep()<CR>")
    vim.keymap.set('n', '<leader>fb', "<cmd>lua require('fzf-lua').buffers()<CR>")
    vim.keymap.set('n', '<leader>fh', "<cmd>lua require('fzf-lua').help_tags()<CR>")
    ```

## Knowledge Management

### obsidian.nvim: Obsidian Integration for Neovim
- **Description**: Neovim plugin for writing and navigating Obsidian vaults with extensive markdown support
- **Features**:
  - Markdown-based note management with wiki-style links
  - Completion for note references and tags via nvim-cmp
  - Navigation throughout vaults with keyboard shortcuts
  - Image pasting and management capabilities
  - Enhanced markdown syntax highlighting and concealing
  - Template system for consistent note creation
  - Search functionality powered by ripgrep
  - Customizable frontmatter management
  - Extensive command set for note manipulation
  - 5.1k+ stars and active maintenance
  - Repository: https://github.com/epwalsh/obsidian.nvim
  - Configuration example:
    ```lua
    require("obsidian").setup({
      workspaces = {
        {
          name = "work",
          path = "~/vaults/work",
        },
      },
      
      -- Completion for wiki links and tags
      completion = {
        nvim_cmp = true,
        min_chars = 2,
      },
      
      -- Key mappings
      mappings = {
        ["gf"] = {
          action = function()
            return require("obsidian").util.gf_passthrough()
          end,
          opts = { noremap = false, expr = true, buffer = true },
        },
        ["<leader>ch"] = {
          action = function()
            return require("obsidian").util.toggle_checkbox()
          end,
          opts = { buffer = true },
        },
      },
      
      -- UI customization
      ui = {
        enable = true,
        update_debounce = 200,
        checkboxes = {
          [" "] = { char = "󰄱", hl_group = "ObsidianTodo" },
          ["x"] = { char = "", hl_group = "ObsidianDone" },
        },
        bullets = { char = "•", hl_group = "ObsidianBullet" },
        external_link_icon = { char = "", hl_group = "ObsidianExtLinkIcon" },
      },
    })
    ```

## Language Support

### marksman: Markdown Language Server
- **Description**: Language server for Markdown that provides intelligent code assistance and navigation
- **Features**:
  - Completion for links, references, and headings
  - Goto definition for links and references
  - Find references for headings and link targets
  - Rename refactoring for headings and references
  - Diagnostics for broken links and duplicate headings
  - Support for inline links, reference links, and wiki-links
  - Hover information for links and references
  - Document symbols for headings and references
  - Multi-file workspace support
  - Customizable configuration via .marksman.toml
  - 2.4k+ stars and active maintenance
  - Repository: https://github.com/artempyanykh/marksman
  - Integration example:
    ```lua
    -- Neovim LSP configuration
    require('lspconfig').marksman.setup({
      -- Default settings
      cmd = { "marksman", "server" },
      filetypes = { "markdown" },
      root_dir = function(fname)
        return require('lspconfig.util').find_git_ancestor(fname) or
               require('lspconfig.util').path.dirname(fname)
      end,
      single_file_support = true,
      
      -- Additional settings for our workflow
      init_options = {
        -- Enable diagnostics for broken links
        diagnostics = true,
      },
    })
    
    -- Key mappings for Markdown files
    vim.api.nvim_create_autocmd("FileType", {
      pattern = "markdown",
      callback = function()
        -- Local mappings for Markdown files
        vim.keymap.set("n", "gd", vim.lsp.buf.definition, { buffer = true })
        vim.keymap.set("n", "gr", vim.lsp.buf.references, { buffer = true })
        vim.keymap.set("n", "K", vim.lsp.buf.hover, { buffer = true })
        vim.keymap.set("n", "<leader>rn", vim.lsp.buf.rename, { buffer = true })
      end,
    })
    ```

## Focus and Productivity

### zen-mode.nvim: Distraction-free Coding
- **Description**: Distraction-free coding environment for Neovim that creates a focused workspace
- **Features**:
  - Opens current buffer in a full-screen floating window
  - Preserves existing window layouts and splits
  - Works with other floating windows (LSP hover, WhichKey)
  - Dynamically adjustable window size
  - Backdrop shading for enhanced focus
  - Hides status line and optionally other UI elements
  - Customizable with Lua callbacks for open/close events
  - Plugin integrations for git, tmux, terminal emulators
  - Automatic closing when new non-floating windows open
  - Compatible with Telescope for buffer switching
  - 1.9k+ stars and active maintenance
  - Repository: https://github.com/folke/zen-mode.nvim
  - Configuration example:
    ```lua
    require("zen-mode").setup({
      window = {
        backdrop = 0.95,
        width = 120,
        height = 1,
        options = {
          signcolumn = "no",
          number = false,
          relativenumber = false,
          cursorline = false,
          cursorcolumn = false,
          foldcolumn = "0",
          list = false,
        },
      },
      plugins = {
        options = {
          enabled = true,
          ruler = false,
          showcmd = false,
          laststatus = 0,
        },
        twilight = { enabled = true },
        gitsigns = { enabled = false },
        tmux = { enabled = true },
        kitty = { enabled = false, font = "+4" },
        alacritty = { enabled = false, font = "14" },
        wezterm = { enabled = false, font = "+4" },
      },
      on_open = function(win)
        -- Custom actions when Zen Mode opens
        vim.cmd("echo 'Entering focus mode'")
      end,
      on_close = function()
        -- Custom actions when Zen Mode closes
        vim.cmd("echo 'Exiting focus mode'")
      end,
    })
    
    -- Key mapping for toggling Zen Mode
    vim.keymap.set("n", "<leader>z", "<cmd>ZenMode<CR>", { silent = true })
    
    -- Toggle with custom settings
    vim.keymap.set("n", "<leader>Z", function()
      require("zen-mode").toggle({
        window = { width = 0.85 },
        plugins = { twilight = { enabled = false } }
      })
    end, { silent = true })
    ```

## Coding Practice

### leetcode.nvim: LeetCode Integration
- **Description**: Neovim plugin for solving LeetCode problems directly within the editor
- **Features**:
  - Intuitive dashboard for problem navigation
  - Formatted question descriptions with enhanced readability
  - LeetCode profile statistics within Neovim
  - Support for daily and random questions
  - Caching for optimized performance
  - Multiple language support including Python, Go, TypeScript
  - Code execution and submission directly from Neovim
  - Customizable UI with multiple themes
  - Telescope and fzf-lua integration for problem selection
  - Session management for multiple LeetCode accounts
  - 1.5k+ stars and active maintenance
  - Repository: https://github.com/kawre/leetcode.nvim
  - Configuration example:
    ```lua
    require("leetcode").setup({
      -- Language to use for LeetCode problems
      lang = "typescript", -- For TypeScript support
      
      -- Storage directories
      storage = {
        home = vim.fn.stdpath("data") .. "/leetcode",
        cache = vim.fn.stdpath("cache") .. "/leetcode",
      },
      
      -- Inject code before or after solution
      injector = {
        ["python3"] = {
          before = true, -- Use default imports
        },
        ["typescript"] = {
          before = "import { ListNode, TreeNode } from './definitions';",
        },
        ["go"] = {
          before = "package main\n\nimport (\n\t\"fmt\"\n)",
        },
      },
      
      -- Use fzf-lua for picker
      picker = { provider = "fzf-lua" },
      
      -- Console configuration
      console = {
        open_on_runcode = true,
        dir = "row",
        size = {
          width = "90%",
          height = "75%",
        },
      },
      
      -- Description configuration
      description = {
        position = "left",
        width = "40%",
        show_stats = true,
      },
    })
    
    -- Key mappings
    vim.keymap.set("n", "<leader>ll", "<cmd>Leet<CR>", { desc = "Open LeetCode" })
    vim.keymap.set("n", "<leader>lt", "<cmd>Leet tabs<CR>", { desc = "LeetCode tabs" })
    vim.keymap.set("n", "<leader>lr", "<cmd>Leet run<CR>", { desc = "Run LeetCode solution" })
    vim.keymap.set("n", "<leader>ls", "<cmd>Leet submit<CR>", { desc = "Submit LeetCode solution" })
    ```

## Command Menu Enhancement

### wilder.nvim: Enhanced Command Menu
- **Description**: Advanced wildmenu replacement for Neovim that enhances command-line and search functionality
- **Features**:
  - Automatic command suggestions as you type
  - Support for cmdline, search, and substitute modes
  - Fuzzy matching with multiple highlighter options
  - Customizable UI with popupmenu and wildmenu renderers
  - Gradient highlighting for better visual feedback
  - Devicons integration for file types
  - Border themes and palette mode for command menu
  - Pipeline architecture for flexible customization
  - Support for both Vim and Neovim
  - Performance optimization options
  - 1.4k+ stars and active maintenance
  - Repository: https://github.com/gelguy/wilder.nvim
  - Configuration example:
    ```lua
    local wilder = require('wilder')
    wilder.setup({
      modes = {':', '/', '?'},
      next_key = '<Tab>',
      previous_key = '<S-Tab>',
      accept_key = '<Down>',
      reject_key = '<Up>'
    })
    
    -- Fuzzy search pipeline
    wilder.set_option('pipeline', {
      wilder.branch(
        wilder.cmdline_pipeline({
          fuzzy = 1,
          set_pcre2_pattern = 1,
        }),
        wilder.python_search_pipeline({
          pattern = 'fuzzy',
        })
      ),
    })
    
    -- Popupmenu renderer with devicons
    wilder.set_option('renderer', wilder.renderer_mux({
      [':'] = wilder.popupmenu_renderer({
        highlighter = wilder.basic_highlighter(),
        left = {' ', wilder.popupmenu_devicons()},
        right = {' ', wilder.popupmenu_scrollbar()},
      }),
      ['/'] = wilder.wildmenu_renderer({
        highlighter = wilder.basic_highlighter(),
      }),
    }))
    ```

## TypeScript Development

### coc-tsserver: TypeScript Language Server
- **Description**: TypeScript server extension for coc.nvim that provides rich language features for JavaScript and TypeScript
- **Features**:
  - Code completion with auto imports
  - Go to definition and references
  - Code validation and diagnostics
  - Document highlighting and symbols
  - Folding and folding range support
  - Formatting (buffer, range, and on-type)
  - Hover documentation
  - CodeLens for implementations and references
  - Organize imports command
  - Quickfix using code actions
  - Code refactoring
  - Source code actions (fix issues, remove unused code, etc.)
  - Signature help
  - Call hierarchy
  - Semantic tokens
  - Rename symbols
  - Automatic tag closing
  - Rename imports on file rename
  - Workspace symbol search
  - Inlay hints support
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```lua
    -- Install coc.nvim first
    -- Then install coc-tsserver with :CocInstall coc-tsserver
    
    -- coc-settings.json configuration
    {
      "tsserver.enable": true,
      "tsserver.log": "verbose",
      "tsserver.trace.server": "off",
      "tsserver.maxTsServerMemory": 4096,
      "tsserver.useLocalTsdk": true,
      "tsserver.tsdk": "./node_modules/typescript/lib",
      
      "typescript.suggest.enabled": true,
      "typescript.suggest.paths": true,
      "typescript.suggest.autoImports": true,
      "typescript.format.enable": true,
      "typescript.updateImportsOnFileMove.enabled": "always",
      "typescript.preferences.importModuleSpecifier": "relative",
      "typescript.preferences.quoteStyle": "single",
      
      "javascript.suggest.enabled": true,
      "javascript.suggest.paths": true,
      "javascript.suggest.autoImports": true,
      "javascript.format.enable": true
    }
    
    -- Key mappings in init.vim/init.lua
    -- Use CocConfig to edit coc-settings.json
    vim.api.nvim_set_keymap('n', 'gd', '<Plug>(coc-definition)', {silent = true})
    vim.api.nvim_set_keymap('n', 'gy', '<Plug>(coc-type-definition)', {silent = true})
    vim.api.nvim_set_keymap('n', 'gi', '<Plug>(coc-implementation)', {silent = true})
    vim.api.nvim_set_keymap('n', 'gr', '<Plug>(coc-references)', {silent = true})
    vim.api.nvim_set_keymap('n', '<leader>rn', '<Plug>(coc-rename)', {silent = true})
    vim.api.nvim_set_keymap('n', '<leader>f', '<Plug>(coc-format)', {silent = true})
    vim.api.nvim_set_keymap('n', '<leader>a', '<Plug>(coc-codeaction)', {silent = true})
    vim.api.nvim_set_keymap('n', '<leader>qf', '<Plug>(coc-fix-current)', {silent = true})
    ```

## Database Management

### nvim-dbee: Interactive Database Client
- **Description**: Interactive database client for Neovim that provides a rich UI for database operations
- **Features**:
  - Support for multiple database types (PostgreSQL, MySQL, SQLite, etc.)
  - Interactive query execution with results pagination
  - Connection management with secure credential handling
  - Scratchpad system for saving and organizing queries
  - Result export in multiple formats (JSON, CSV, Table)
  - Syntax highlighting for SQL
  - Keyboard-driven interface (minimal mouse usage)
  - Split window UI with drawer, editor, and results panes
  - Template-based secret management
  - Go-powered backend for performance
  - Lua-based frontend for flexibility
  - 984+ stars and active maintenance
  - Repository: https://github.com/kndndrj/nvim-dbee
  - Configuration example:
    ```lua
    require("dbee").setup({
      -- Sources for loading connections
      sources = {
        -- Load connections from a file with editing support
        require("dbee.sources").FileSource:new(vim.fn.stdpath("data") .. "/dbee/connections.json"),
        -- Load connections from environment variable
        require("dbee.sources").EnvSource:new("DBEE_CONNECTIONS"),
        -- Hardcoded connections
        require("dbee.sources").MemorySource:new({
          {
            name = "Local PostgreSQL",
            type = "postgres",
            url = "postgres://{{ env `DB_USER` }}:{{ env `DB_PASS` }}@localhost:5432/mydb?sslmode=disable",
          },
          {
            name = "Project SQLite DB",
            type = "sqlite",
            url = "~/repos/agent_runtime/data/agent.db",
          },
        }),
      },
      
      -- UI configuration
      ui = {
        -- Border style for floating windows
        border = "rounded",
        -- Width of the drawer (connections/scratchpads)
        drawer_width = 40,
        -- Default width of result window
        result_width = 40,
        -- Default height of result window
        result_height = 10,
      },
      
      -- Result view configuration
      results = {
        -- Number of rows to fetch at once
        fetch_count = 100,
        -- Maximum number of results to show
        max_results = 1000,
      },
      
      -- Key mappings
      mappings = {
        -- Execute query under cursor or visual selection
        execute = "<leader>be",
        -- Close the UI
        close = "<leader>bq",
        -- Navigate result pages
        results = {
          next_page = "L",
          prev_page = "H",
          first_page = "F",
          last_page = "E",
        },
      },
    })
    
    -- Key mappings for global access
    vim.keymap.set("n", "<leader>bd", "<cmd>lua require('dbee').toggle()<CR>", { desc = "Toggle DB UI" })
    vim.keymap.set("n", "<leader>bs", "<cmd>lua require('dbee').store('csv', 'buffer')<CR>", { desc = "Store results as CSV" })
    ```

## Container Development

### nvim-remote-containers: Docker Container Integration
- **Description**: Neovim plugin that provides VSCode-like remote container development capabilities
- **Features**:
  - Development inside Docker containers directly from Neovim
  - Support for devcontainer.json configuration
  - Container attachment and management commands
  - Docker image building with progress visualization
  - Docker Compose integration (up, down, destroy)
  - Container status display in statusline
  - Dockerfile parsing and configuration
  - Container lifecycle management
  - Seamless workspace mounting
  - Lua API for programmatic container control
  - 907+ stars and active maintenance
  - Repository: https://github.com/jamestthompson3/nvim-remote-containers
  - Configuration example:
    ```lua
    -- Load the plugin
    use 'jamestthompson3/nvim-remote-containers'
    
    -- Optional statusline integration
    vim.cmd([[
      hi Container guifg=#BADA55 guibg=Black
      set statusline+=%#Container#%{g:currentContainer}
    ]])
    
    -- Key mappings for container operations
    vim.api.nvim_set_keymap('n', '<leader>ca', ':AttachToContainer<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cb', ':BuildImage true<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cs', ':StartImage<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cu', ':ComposeUp<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cd', ':ComposeDown<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cx', ':ComposeDestroy<CR>', { noremap = true, silent = true })
    
    -- Example of programmatic usage with Lua API
    local function attach_to_dev_container()
      -- Parse the devcontainer.json configuration
      local config = vim.fn.luaeval('require("nvim-remote-containers").parseConfig("dockerFile")')
      
      -- Build the container image if needed
      vim.fn.luaeval('require("nvim-remote-containers").buildImage(true)')
      
      -- Attach to the container
      vim.fn.luaeval('require("nvim-remote-containers").attachToContainer()')
    end
    ```

## Interactive Coding

### molten-nvim: Jupyter Kernel Integration
- **Description**: Neovim plugin for interactively running code with the Jupyter kernel, providing a notebook-like experience
- **Features**:
  - Asynchronous code execution with real-time output
  - Support for multiple languages with Jupyter kernels (Python, Go, TypeScript, etc.)
  - Image, plot, and LaTeX rendering directly in Neovim
  - Virtual text output display for persistent results
  - Multiple buffer support with the same kernel
  - Multiple kernel support from the same buffer
  - Python virtual environment support
  - Import/export to Jupyter notebook files
  - Keyboard-driven interface with minimal mouse usage
  - Cell-based code execution with highlighting
  - Floating output windows with execution state indicators
  - Stdin input support via vim.ui.input
  - 805+ stars and active maintenance
  - Repository: https://github.com/benlubas/molten-nvim
  - Configuration example:
    ```lua
    -- Basic setup
    require('molten').setup({
      -- Use virtual text for output
      virt_text_output = true,
      -- Limit virtual text to 10 lines
      virt_text_max_lines = 10,
      -- Use a custom image provider
      image_provider = "image.nvim", -- or "wezterm"
      -- Automatically open output window when entering a cell
      auto_open_output = true,
      -- Show execution time
      output_show_exec_time = true,
      -- Wrap output text
      wrap_output = true,
      -- Border style for output window
      output_win_border = { "", "━", "", "" },
      -- Key mappings
      mappings = {
        -- Execute the current line
        evaluate_line = "<localleader>rl",
        -- Execute the current cell
        evaluate_cell = "<localleader>rc",
        -- Execute visual selection
        evaluate_selection = "<localleader>rs",
        -- Delete the current cell
        delete_cell = "<localleader>rd",
        -- Show output for the current cell
        show_output = "<localleader>ro",
        -- Hide output for the current cell
        hide_output = "<localleader>rh",
      },
    })
    
    -- Key mappings
    vim.keymap.set("n", "<localleader>mi", ":MoltenInit<CR>", { silent = true, desc = "Initialize Jupyter kernel" })
    vim.keymap.set("n", "<localleader>e", ":MoltenEvaluateOperator<CR>", { silent = true, desc = "Run operator selection" })
    vim.keymap.set("n", "<localleader>rl", ":MoltenEvaluateLine<CR>", { silent = true, desc = "Evaluate line" })
    vim.keymap.set("n", "<localleader>rr", ":MoltenReevaluateCell<CR>", { silent = true, desc = "Re-evaluate cell" })
    vim.keymap.set("v", "<localleader>r", ":<C-u>MoltenEvaluateVisual<CR>gv", { silent = true, desc = "Evaluate visual selection" })
    vim.keymap.set("n", "<localleader>rd", ":MoltenDelete<CR>", { silent = true, desc = "Delete cell" })
    vim.keymap.set("n", "<localleader>oh", ":MoltenHideOutput<CR>", { silent = true, desc = "Hide output" })
    vim.keymap.set("n", "<localleader>os", ":noautocmd MoltenEnterOutput<CR>", { silent = true, desc = "Show/enter output" })
    ```

## Terminal Integration

### tmux.nvim: Seamless Tmux Integration
- **Description**: Neovim plugin that provides seamless integration with tmux for navigation, resizing, and clipboard synchronization
- **Features**:
  - Navigation between Neovim splits and tmux panes with the same keybindings
  - Resizing Neovim splits like tmux panes
  - Copy sync between Neovim and tmux buffers
  - Register synchronization across Neovim instances
  - Swapping panes with keyboard shortcuts
  - Clipboard harmonization across environments
  - Window navigation with cycle support
  - Configurable keybindings and behavior
  - Support for tmux plugin manager (tpm)
  - Lua API for programmatic control
  - 702+ stars and active maintenance
  - Repository: https://github.com/aserowy/tmux.nvim
  - Configuration example:
    ```lua
    require("tmux").setup({
      copy_sync = {
        -- enables copy sync between Neovim and tmux
        enable = true,
        -- ignore specific tmux buffers
        ignore_buffers = { empty = false },
        -- redirect to system clipboard
        redirect_to_clipboard = false,
        -- register offset for sync
        register_offset = 0,
        -- sync clipboard with tmux
        sync_clipboard = true,
        -- sync registers
        sync_registers = true,
        -- sync registers on put operations
        sync_registers_keymap_put = true,
        -- sync registers on register access
        sync_registers_keymap_reg = true,
        -- sync deletes
        sync_deletes = true,
        -- sync unnamed register
        sync_unnamed = true,
      },
      navigation = {
        -- cycle to opposite pane while navigating into the border
        cycle_navigation = true,
        -- enable default keybindings (C-hjkl)
        enable_default_keybindings = true,
        -- prevent unzoom tmux when navigating beyond vim border
        persist_zoom = false,
      },
      resize = {
        -- enable default keybindings (A-hjkl)
        enable_default_keybindings = true,
        -- resize steps for x axis
        resize_step_x = 1,
        -- resize steps for y axis
        resize_step_y = 1,
      },
      swap = {
        -- cycle to opposite pane while navigating into the border
        cycle_navigation = false,
        -- enable default keybindings (C-A-hjkl)
        enable_default_keybindings = true,
      },
    })
    
    -- Example tmux.conf configuration
    -- is_vim="ps -o state= -o comm= -t '#{pane_tty}' | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?\.?(view|n?vim?x?)(-wrapped)?(diff)?$'"
    -- bind-key -n 'C-h' if-shell "$is_vim" 'send-keys C-h' 'select-pane -L'
    -- bind-key -n 'C-j' if-shell "$is_vim" 'send-keys C-j' 'select-pane -D'
    -- bind-key -n 'C-k' if-shell "$is_vim" 'send-keys C-k' 'select-pane -U'
    -- bind-key -n 'C-l' if-shell "$is_vim" 'send-keys C-l' 'select-pane -R'
    ```

## Focus and Productivity

### zen-mode.nvim: Distraction-free Coding
- **Description**: Neovim plugin that provides a distraction-free coding environment by creating a centered window and hiding UI elements
- **Features**:
  - Opens current buffer in a full-screen floating window
  - Preserves existing window layouts and splits
  - Compatible with other floating windows (LSP hover, WhichKey, etc.)
  - Dynamic window size adjustment
  - Automatic realignment when editor or window is resized
  - Backdrop shading for better focus
  - Hides status line, number column, sign column, fold column
  - Customizable with Lua callbacks (on_open, on_close)
  - Integration with other plugins (gitsigns, tmux, kitty, alacritty, wezterm, neovide)
  - Automatic closing when a new non-floating window is opened
  - Works with Telescope for opening new buffers inside the Zen window
  - 1.9k+ stars and active maintenance
  - Repository: https://github.com/folke/zen-mode.nvim
  - Configuration example:
    ```lua
    require("zen-mode").setup({
      window = {
        backdrop = 0.95, -- shade the backdrop (0-1)
        width = 120,     -- width of the Zen window
        height = 1,      -- height of the Zen window (as a percentage)
        options = {
          signcolumn = "no",       -- disable signcolumn
          number = false,          -- disable number column
          relativenumber = false,  -- disable relative numbers
          cursorline = false,      -- disable cursorline
          cursorcolumn = false,    -- disable cursor column
          foldcolumn = "0",        -- disable fold column
          list = false,            -- disable whitespace characters
        },
      },
      plugins = {
        options = {
          enabled = true,
          ruler = false,           -- disable ruler
          showcmd = false,         -- disable command display
          laststatus = 0,          -- disable status line
        },
        twilight = { enabled = true },  -- dim inactive code
        gitsigns = { enabled = false }, -- disable git signs
        tmux = { enabled = false },     -- disable tmux statusline
        -- Additional terminal integrations
        kitty = { enabled = false, font = "+4" },
        alacritty = { enabled = false, font = "14" },
        wezterm = { enabled = false, font = "+4" },
        neovide = { 
          enabled = false, 
          scale = 1.2,
          disable_animations = {
            neovide_animation_length = 0,
            neovide_cursor_animate_command_line = false,
            neovide_scroll_animation_length = 0,
            neovide_cursor_vfx_mode = "",
          }
        },
      },
      -- Callbacks for custom actions
      on_open = function(win)
        -- Actions when Zen window opens
      end,
      on_close = function()
        -- Actions when Zen window closes
      end,
    })
    
    -- Key mapping to toggle Zen Mode
    vim.keymap.set("n", "<leader>z", ":ZenMode<CR>", { silent = true, desc = "Toggle Zen Mode" })
    ```

## TypeScript Development

### coc-tsserver: TypeScript Language Server
- **Description**: TypeScript server extension for coc.nvim that provides rich language features similar to VSCode
- **Features**:
  - Support for JavaScript, TypeScript, JSX, and TSX
  - Automatic typings installation
  - Code completion with intelligent suggestions
  - Go to definition and references
  - Code validation and error reporting
  - Document highlighting and symbols
  - Folding and folding range support
  - Format current buffer, range format, and format on type
  - Hover documentation
  - CodeLens for implementations and references
  - Organize imports command
  - Quickfix using code actions
  - Code refactoring
  - Source code actions (fix issues, remove unused code, etc.)
  - Signature help
  - Call hierarchy
  - Semantic tokens
  - Rename symbols
  - Automatic tag closing
  - Rename imports on file rename
  - Workspace symbol search
  - Inlay hints support
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```json
    {
      "tsserver.enable": true,
      "tsserver.maxTsServerMemory": 4096,
      "tsserver.log": "verbose",
      "tsserver.trace.server": "messages",
      "typescript.suggest.completeFunctionCalls": true,
      "typescript.suggest.autoImports": true,
      "typescript.preferences.importModuleSpecifier": "relative",
      "typescript.preferences.quoteStyle": "single",
      "typescript.format.enable": true,
      "typescript.inlayHints.parameterNames.enabled": "all",
      "typescript.inlayHints.variableTypes.enabled": true,
      "typescript.inlayHints.functionLikeReturnTypes.enabled": true,
      "javascript.suggest.completeFunctionCalls": true,
      "javascript.suggest.autoImports": true,
      "javascript.preferences.importModuleSpecifier": "relative",
      "javascript.preferences.quoteStyle": "single",
      "javascript.format.enable": true
    }
    ```

## Go Development

### vim-go: Go Development Plugin
- **Description**: Comprehensive Go development plugin for Vim that provides rich features for Go programming
- **Features**:
  - Compile, install, and test Go packages with `:GoBuild`, `:GoInstall`, `:GoTest`
  - Run current file(s) with `:GoRun`
  - Improved syntax highlighting and folding
  - Integrated debugging with Delve support via `:GoDebugStart`
  - Code completion and language features via `gopls`
  - Automatic formatting on save with cursor position preservation
  - Go to symbol/declaration with `:GoDef`
  - Documentation lookup with `:GoDoc` or `:GoDocBrowser`
  - Package management with `:GoImport` and `:GoDrop`
  - Type-safe renaming with `:GoRename`
  - Code coverage visualization with `:GoCoverage`
  - Struct tag management with `:GoAddTags` and `:GoRemoveTags`
  - Linting with `:GoLint`, `:GoMetaLinter`, and `:GoVet`
  - Advanced source analysis with `:GoImplements`, `:GoCallees`, and `:GoReferrers`
  - Integration with Tagbar via gotags
  - Snippet engine integration
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " Basic settings
    let g:go_fmt_command = "goimports"
    let g:go_auto_type_info = 1
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    
    " Auto formatting and importing
    let g:go_fmt_autosave = 1
    let g:go_fmt_command = "goimports"
    
    " Status line types/signatures
    let g:go_auto_type_info = 1
    
    " Run :GoBuild or :GoTestCompile based on the go file
    function! s:build_go_files()
      let l:file = expand('%')
      if l:file =~# '^\f\+_test\.go$'
        call go#test#Test(0, 1)
      elseif l:file =~# '^\f\+\.go$'
        call go#cmd#Build(0)
      endif
    endfunction
    
    " Key mappings
    autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <Leader>i <Plug>(go-info)
    autocmd FileType go nmap <Leader>d <Plug>(go-def)
    autocmd FileType go nmap <Leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <Leader>v <Plug>(go-doc-vertical)
    ```

## Package Management

### mason.nvim: Portable Package Manager
- **Description**: Comprehensive package manager for Neovim that runs everywhere Neovim runs
- **Features**:
  - Easily install and manage LSP servers, DAP servers, linters, and formatters
  - Portable across all platforms (Linux, macOS, Windows)
  - Minimal dependencies required
  - Automatic PATH integration for installed tools
  - Unified interface for package management
  - Extensible registry system
  - Automatic version management
  - Configurable installation directory
  - Concurrent package installation
  - Detailed installation logs
  - Comprehensive UI with keymaps for common actions
  - 8.7k+ stars and active maintenance
  - Repository: https://github.com/williamboman/mason.nvim
  - Configuration example:
    ```lua
    require("mason").setup({
      ui = {
        icons = {
          package_installed = "✓",
          package_pending = "➜",
          package_uninstalled = "✗"
        },
        keymaps = {
          toggle_package_expand = "<CR>",
          install_package = "i",
          update_package = "u",
          check_package_version = "c",
          update_all_packages = "U",
          check_outdated_packages = "C",
          uninstall_package = "X",
          cancel_installation = "<C-c>",
          apply_language_filter = "<C-f>",
        },
        border = "rounded",
        width = 0.8,
        height = 0.9
      },
      max_concurrent_installers = 4,
      github = {
        download_url_template = "https://github.com/%s/releases/download/%s/%s"
      },
      log_level = vim.log.levels.INFO
    })
    
    -- Recommended extension for LSP integration
    require("mason-lspconfig").setup({
      ensure_installed = {
        -- Python
        "pyright",
        "ruff_lsp",
        -- Go
        "gopls",
        -- TypeScript
        "tsserver",
        -- Kubernetes
        "yamlls",
        -- Terraform
        "terraformls",
        "tflint"
      },
      automatic_installation = true
    })
    ```

## Navigation and Search

### telescope.nvim: Fuzzy Finder
- **Description**: Highly extensible fuzzy finder over lists with powerful search capabilities
- **Features**:
  - Find files, buffers, git files, and more with fuzzy matching
  - Live grep for searching text across your project
  - LSP integration for code navigation (definitions, references, etc.)
  - Git integration (commits, branches, status)
  - Treesitter integration for code symbols
  - Extensible architecture for custom pickers
  - Multiple layout strategies (horizontal, vertical, cursor, etc.)
  - Customizable themes (dropdown, ivy, cursor)
  - Comprehensive keybindings for mouse-free operation
  - Preview window for files, grep results, and more
  - Multiple selection with actions (quickfix list, etc.)
  - 17.3k+ stars and active maintenance
  - Repository: https://github.com/nvim-telescope/telescope.nvim
  - Configuration example:
    ```lua
    require('telescope').setup({
      defaults = {
        mappings = {
          i = {
            ["<C-j>"] = "move_selection_next",
            ["<C-k>"] = "move_selection_previous",
            ["<C-n>"] = "move_selection_next",
            ["<C-p>"] = "move_selection_previous",
            ["<C-c>"] = "close",
            ["<C-u>"] = "preview_scrolling_up",
            ["<C-d>"] = "preview_scrolling_down",
          }
        },
        file_ignore_patterns = {
          "node_modules",
          ".git",
          "dist",
          "build"
        },
        vimgrep_arguments = {
          "rg",
          "--color=never",
          "--no-heading",
          "--with-filename",
          "--line-number",
          "--column",
          "--smart-case",
          "--hidden"
        }
      },
      pickers = {
        find_files = {
          hidden = true,
          find_command = { "fd", "--type", "f", "--strip-cwd-prefix" }
        },
        live_grep = {
          additional_args = function()
            return { "--hidden" }
          end
        }
      },
      extensions = {
        fzf = {
          fuzzy = true,
          override_generic_sorter = true,
          override_file_sorter = true,
          case_mode = "smart_case"
        }
      }
    })
    
    -- Load extensions
    require('telescope').load_extension('fzf')
    
    -- Keymaps
    vim.keymap.set('n', '<leader>ff', require('telescope.builtin').find_files)
    vim.keymap.set('n', '<leader>fg', require('telescope.builtin').live_grep)
    vim.keymap.set('n', '<leader>fb', require('telescope.builtin').buffers)
    vim.keymap.set('n', '<leader>fh', require('telescope.builtin').help_tags)
    vim.keymap.set('n', '<leader>fs', require('telescope.builtin').lsp_document_symbols)
    vim.keymap.set('n', '<leader>fr', require('telescope.builtin').lsp_references)
    ```

## Infrastructure as Code

### terraform-ls: Terraform Language Server
- **Description**: Official HashiCorp Terraform Language Server that provides IDE features to any LSP-compatible editor
- **Features**:
  - Code completion for Terraform resources, data sources, and attributes
  - Hover information for resource properties and types
  - Go to definition for resources and variables
  - Find references across Terraform files
  - Document symbols for quick navigation
  - Diagnostics for syntax and validation errors
  - Format on save with terraform fmt
  - Workspace symbol search
  - Semantic token highlighting
  - Auto-completion for provider schemas
  - Module navigation and completion
  - Variable completion across files
  - 1.1k+ stars and active maintenance by HashiCorp
  - Repository: https://github.com/hashicorp/terraform-ls
  - Configuration example:
    ```lua
    require('lspconfig').terraformls.setup {
      cmd = { "terraform-ls", "serve" },
      filetypes = { "terraform", "terraform-vars", "tf" },
      root_dir = function(fname)
        return util.root_pattern(".terraform", ".git")(fname)
      end,
      settings = {
        terraform = {
          path = "terraform",
          logLevel = "info",
          excludeRootModules = false,
          ignoreDirectoryNames = [],
          telemetry = {
            enable = false
          },
          experimentalFeatures = {
            validateOnSave = true,
            prefillRequiredFields = true
          }
        }
      }
    }
    
    -- Format on save
    vim.api.nvim_create_autocmd({"BufWritePre"}, {
      pattern = {"*.tf", "*.tfvars"},
      callback = function()
        vim.lsp.buf.format()
      end,
    })
    ```

## Kubernetes and YAML Support

### yaml-language-server: YAML Language Server
- **Description**: Language Server for YAML files that provides comprehensive IDE features to any LSP-compatible editor
- **Features**:
  - YAML validation for syntax and schema compliance
  - Auto-completion for YAML properties and values
  - Hover information for documentation and type details
  - Document outlining for navigation
  - Schema association via configuration
  - Support for JSON Schema 7 and below
  - Custom tags support for extended YAML syntax
  - Kubernetes schema integration
  - Configurable formatting options
  - Support for multiple YAML versions (1.1, 1.2)
  - Telemetry-free operation
  - 1.2k+ stars and active maintenance by Red Hat
  - Repository: https://github.com/redhat-developer/yaml-language-server
  - Configuration example:
    ```lua
    require('lspconfig').yamlls.setup {
      settings = {
        yaml = {
          schemas = {
            kubernetes = "/*.yaml",
            ["https://json.schemastore.org/github-workflow.json"] = "/.github/workflows/*",
            ["https://raw.githubusercontent.com/instrumenta/kubernetes-json-schema/master/v1.18.0-standalone-strict/all.json"] = "/*.k8s.yaml",
            ["https://raw.githubusercontent.com/hashicorp/terraform-json-schema/main/v1.1.0/schema.json"] = "/*.tf.yaml"
          },
          format = {
            enable = true,
            singleQuote = true,
            bracketSpacing = true,
            proseWrap = "preserve",
            printWidth = 120
          },
          validate = true,
          completion = true,
          hover = true,
          schemaStore = {
            enable = true,
            url = "https://www.schemastore.org/api/json/catalog.json"
          },
          customTags = [
            "!include scalar",
            "!reference sequence"
          ],
          keyOrdering = false,
          maxItemsComputed = 5000
        }
      }
    }
    ```

## Core Extensions

### coc.nvim: Intelligent Coding Extension Host
- **Description**: Nodejs extension host for Vim & Neovim that provides VSCode-like features and language server protocol support
- **Features**:
  - Language server protocol (LSP) integration for all major languages
  - Intelligent code completion with popup menu
  - Diagnostics and error checking in real-time
  - Code navigation (go to definition, references, implementations)
  - Symbol search and workspace symbol support
  - Code actions and quick fixes
  - Hover documentation
  - Signature help for function parameters
  - Rename symbol across files
  - Document formatting and range formatting
  - Document highlighting
  - Extensions ecosystem similar to VSCode
  - Snippet support
  - Multi-cursor support
  - Git integration
  - Extensible through JavaScript/TypeScript
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " Install with vim-plug
    Plug 'neoclide/coc.nvim', {'branch': 'release'}
    
    " Basic configuration
    set hidden
    set nobackup
    set nowritebackup
    set updatetime=300
    set shortmess+=c
    set signcolumn=yes
    
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Make <CR> auto-select the first completion item
    inoremap <silent><expr> <cr> coc#pum#visible() ? coc#pum#confirm() : "\<C-g>u\<CR>"
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Use K to show documentation in preview window
    nnoremap <silent> K :call ShowDocumentation()<CR>
    
    " Symbol renaming
    nmap <leader>rn <Plug>(coc-rename)
    
    " Formatting selected code
    xmap <leader>f  <Plug>(coc-format-selected)
    nmap <leader>f  <Plug>(coc-format-selected)
    
    " Extensions for our languages
    let g:coc_global_extensions = [
      \ 'coc-pyright',           " Python
      \ 'coc-go',                " Go
      \ 'coc-tsserver',          " TypeScript/JavaScript
      \ 'coc-json',              " JSON
      \ 'coc-yaml',              " YAML (for Kubernetes)
      \ 'coc-docker',            " Docker
      \ 'coc-sh',                " Shell
      \ 'coc-snippets',          " Snippets
      \ 'coc-pairs',             " Auto pairs
      \ 'coc-explorer',          " File explorer
      \ 'coc-git',               " Git integration
      \ 'coc-highlight',         " Document highlighting
      \ ]
    ```

## Language-Specific Support

### vim-go: Go Development Plugin
- **Description**: Comprehensive Go development plugin for Vim & Neovim with extensive features for Go programming
- **Features**:
  - Compile, install, and test Go packages directly from Vim
  - Run Go files with `:GoRun`
  - Improved syntax highlighting and folding for Go
  - Integrated debugging with Delve support
  - Gopls integration for intelligent code features
  - Automatic formatting on save with cursor position preservation
  - Go to definition with `:GoDef`
  - Documentation lookup with `:GoDoc` and `:GoDocBrowser`
  - Package management with `:GoImport` and `:GoDrop`
  - Precise renaming with `:GoRename`
  - Code coverage visualization with `:GoCoverage`
  - Struct tag management with `:GoAddTags` and `:GoRemoveTags`
  - Linting and static analysis with `:GoLint`, `:GoMetaLinter`, and `:GoVet`
  - Advanced source analysis with `:GoImplements`, `:GoCallees`, and `:GoReferrers`
  - Integration with Tagbar via gotags
  - Snippet support
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " vim-go configuration
    let g:go_fmt_command = "goimports"
    let g:go_auto_type_info = 1
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    let g:go_metalinter_autosave = 1
    let g:go_metalinter_autosave_enabled = ['vet', 'golint']
    let g:go_metalinter_deadline = "5s"
    let g:go_def_mode = 'gopls'
    let g:go_info_mode = 'gopls'
    let g:go_rename_command = 'gopls'
    let g:go_gopls_enabled = 1
    let g:go_gopls_options = ['-remote=auto']
    let g:go_doc_popup_window = 1
    let g:go_test_timeout = '10s'
    
    " Key mappings
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <leader>i <Plug>(go-info)
    autocmd FileType go nmap <leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <leader>d <Plug>(go-def)
    autocmd FileType go nmap <leader>v <Plug>(go-def-vertical)
    autocmd FileType go nmap <leader>s <Plug>(go-def-split)
    autocmd FileType go nmap <leader>dt <Plug>(go-def-tab)
    autocmd FileType go nmap <leader>e <Plug>(go-rename)
    autocmd FileType go nmap <leader>a <Plug>(go-alternate-edit)
    autocmd FileType go nmap <leader>av <Plug>(go-alternate-vertical)
    ```

## TypeScript Support

### coc-tsserver: TypeScript Language Server
- **Description**: TypeScript language server extension for coc.nvim providing VSCode-like TypeScript support
- **Features**:
  - Intelligent code completion for TypeScript and JavaScript
  - Type information on hover
  - Go to definition, type definition, and implementation
  - Find references and document symbols
  - Rename symbols across files
  - Code actions and quick fixes
  - Automatic imports and import organization
  - Signature help for function parameters
  - Diagnostics and error checking
  - Code formatting with prettier integration
  - Support for JSX/TSX
  - Project-wide symbol search
  - Snippet completion for function calls
  - Supports TypeScript project references
  - Repository: https://github.com/neoclide/coc-tsserver
  - Installation and configuration:
    ```vim
    " Install with coc.nvim
    :CocInstall coc-tsserver
    
    " Configuration in coc-settings.json
    {
      "tsserver.enable": true,
      "tsserver.enableJavascript": true,
      "typescript.format.enabled": true,
      "typescript.suggest.enabled": true,
      "typescript.suggest.paths": true,
      "typescript.suggest.autoImports": true,
      "typescript.preferences.importModuleSpecifier": "relative",
      "typescript.preferences.quoteStyle": "single",
      "typescript.updateImportsOnFileMove.enabled": "always",
      "typescript.implementationsCodeLens.enabled": true,
      "typescript.referencesCodeLens.enabled": true,
      "typescript.suggest.completeFunctionCalls": true,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces": true,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces": false,
      "typescript.format.insertSpaceAfterTypeAssertion": true,
      "typescript.format.insertSpaceBeforeFunctionParenthesis": false,
      "typescript.format.placeOpenBraceOnNewLineForFunctions": false,
      "typescript.format.placeOpenBraceOnNewLineForControlBlocks": false
    }
    ```

## Python Support

### jedi-vim: Python Autocompletion
- **Description**: Using the Jedi autocompletion library for Vim & Neovim to provide intelligent Python code features
- **Features**:
  - Intelligent code completion for Python with <C-Space>
  - Go to assignment with `<leader>g`
  - Go to definition with `<leader>d`
  - Go to typing stub with `<leader>s`
  - Show documentation/Pydoc with `K`
  - Renaming with `<leader>r` and `<leader>R`
  - Find usages with `<leader>n`
  - Open module with `:Pyimport module_name`
  - Popup documentation windows
  - Signature help for function parameters
  - Support for virtual environments
  - Case insensitive completion
  - 5.3k+ stars and active maintenance
  - Repository: https://github.com/davidhalter/jedi-vim
  - Configuration example:
    ```vim
    " Basic configuration
    let g:jedi#completions_enabled = 1
    let g:jedi#auto_initialization = 1
    let g:jedi#auto_vim_configuration = 1
    let g:jedi#use_tabs_not_buffers = 0
    let g:jedi#use_splits_not_buffers = "right"
    let g:jedi#popup_on_dot = 1
    let g:jedi#popup_select_first = 1
    let g:jedi#show_call_signatures = "1"
    let g:jedi#case_insensitive_completion = 1
    
    " Key mappings
    let g:jedi#goto_command = "<leader>d"
    let g:jedi#goto_assignments_command = "<leader>g"
    let g:jedi#goto_stubs_command = "<leader>s"
    let g:jedi#documentation_command = "K"
    let g:jedi#usages_command = "<leader>n"
    let g:jedi#completions_command = "<C-Space>"
    let g:jedi#rename_command = "<leader>r"
    let g:jedi#rename_command_keep_name = "<leader>R"
    
    " Virtual environment support
    let g:jedi#environment_path = "venv"  " Point to your virtual environment
    ```

## Python Plugins and LSP Support

### Python Plugins

#### nvim-dap-python: Python Debugging Extension
- **Description**: Extension for nvim-dap providing default configurations for Python debugging
- **Features**:
  - Provides default configurations for Python debugging
  - Supports debugging individual test methods or classes
  - Integrates with debugpy for Python debugging
  - Supports unittest, pytest, and django test frameworks
  - Detects virtual environments automatically
  - Customizable configurations for different Python projects
  - 641+ stars and active maintenance
  - Written primarily in Lua
  - Minimal mouse dependency with keyboard-focused workflows
- **Repository**: https://github.com/mfussenegger/nvim-dap-python
- **Configuration example**:
  ```lua
  require('dap-python').setup('~/.virtualenvs/debugpy/bin/python')
  require('dap-python').test_runner = 'pytest'
  ```

#### deoplete-jedi: Python Code Completion
- **Description**: Python code completion for Neovim through Jedi integration
- **Features**:
  - Provides intelligent autocompletion for Python code
  - Integrates with Jedi for accurate code analysis
  - Supports virtual environments
  - Shows documentation for functions and classes
  - Handles imports and module completion
  - 589+ stars on GitHub
  - Written primarily in Python and Vim Script
  - Minimal mouse dependency with keyboard-focused workflows
- **Repository**: https://github.com/deoplete-plugins/deoplete-jedi
- **Configuration example**:
  ```vim
  let g:deoplete#enable_at_startup = 1
  let g:deoplete#sources#jedi#show_docstring = 1
  let g:deoplete#sources#jedi#python_path = '/path/to/python'
  ```

#### semshi: Semantic Highlighting for Python
- **Description**: Semantic highlighting for Python in Neovim
- **Features**:
  - Provides semantic understanding of Python code
  - Highlights names based on scope and context
  - Differentiates between locals, globals, imports, parameters, etc.
  - Marks and renames related nodes in the same scope
  - Indicates syntax errors as you type
  - Supports jumping between classes, functions, and related names
  - 1k+ stars on GitHub
  - Written primarily in Python (98.6%)
  - Minimal mouse dependency with keyboard-focused workflows
- **Repository**: https://github.com/numirias/semshi
- **Configuration example**:
  ```vim
  " Enable Semshi for Python files
  let g:semshi#filetypes = ['python']
  
  " Customize which highlight groups to use
  let g:semshi#excluded_hl_groups = ['local']
  
  " Enable marking of selected nodes
  let g:semshi#mark_selected_nodes = 2
  
  " Tolerate syntax errors while typing
  let g:semshi#tolerate_syntax_errors = v:true
  
  " Error sign in the sign column
  let g:semshi#error_sign = v:true
  ```

### python-lsp-server: Python Language Server
- **Description**: A Python implementation of the Language Server Protocol providing IDE features to editors and other tools
- **Features**:
  - Auto-completion for Python code
  - Code linting with multiple providers (Pyflakes, pycodestyle, McCabe, etc.)
  - Code formatting with autopep8 or YAPF
  - Code actions and quick fixes
  - Signature help for function parameters
  - Go to definition and find references
  - Document symbols and workspace symbols
  - Hover documentation
  - Rename refactoring
  - Multiple workspace support
  - Jedi-based intelligence
  - Support for third-party plugins
  - Configurable through various sources (pycodestyle, flake8)
  - WebSockets support for remote development
  - 2.2k+ stars and active maintenance by Spyder IDE team
  - Repository: https://github.com/python-lsp/python-lsp-server
  - Configuration example:
    ```lua
    require('lspconfig').pylsp.setup {
      settings = {
        pylsp = {
          plugins = {
            pycodestyle = {
              enabled = true,
              ignore = {"E501"},  -- Ignore line length
              maxLineLength = 100
            },
            pyflakes = { enabled = true },
            pylint = { enabled = false },  -- Can enable if needed
            yapf = { enabled = true },
            autopep8 = { enabled = false },  -- Using YAPF instead
            flake8 = { enabled = false },
            mccabe = { enabled = true, threshold = 15 },
            preload = { enabled = true },
            rope_completion = { enabled = true },
            jedi_completion = {
              enabled = true,
              include_params = true,
              include_class_objects = true,
              include_function_objects = true,
            },
            jedi_definition = {
              enabled = true,
              follow_imports = true,
              follow_builtin_imports = true
            },
            jedi_hover = { enabled = true },
            jedi_references = { enabled = true },
            jedi_signature_help = { enabled = true },
            jedi_symbols = { enabled = true, all_scopes = true }
          }
        }
      }
    }
    ```

## Core Navigation and Search

### telescope.nvim: Fuzzy Finder
- **Description**: Highly extensible fuzzy finder over lists with powerful search capabilities
- **Features**:
  - Find files with fuzzy matching
  - Live grep for searching text across files
  - Buffer management and navigation
  - Project-wide symbol search
  - Git integration (commits, branches, status)
  - LSP integration (references, definitions, symbols)
  - Customizable layouts and themes
  - Extensible through plugins
  - Efficient keyboard-driven workflow
  - Multiple sorting algorithms
  - Preview window for files, git diffs, etc.
  - Customizable key mappings
  - 17.3k+ stars and active maintenance
  - Repository: https://github.com/nvim-telescope/telescope.nvim
  - Configuration example:
    ```lua
    -- Basic setup
    require('telescope').setup{
      defaults = {
        mappings = {
          i = {
            ["<C-j>"] = "move_selection_next",
            ["<C-k>"] = "move_selection_previous",
            ["<C-n>"] = "cycle_history_next",
            ["<C-p>"] = "cycle_history_prev",
            ["<C-c>"] = "close",
            ["<CR>"] = "select_default",
            ["<C-x>"] = "select_horizontal",
            ["<C-v>"] = "select_vertical",
            ["<C-t>"] = "select_tab",
            ["<C-u>"] = "preview_scrolling_up",
            ["<C-d>"] = "preview_scrolling_down"
          }
        },
        vimgrep_arguments = {
          "rg",
          "--color=never",
          "--no-heading",
          "--with-filename",
          "--line-number",
          "--column",
          "--smart-case"
        },
        prompt_prefix = "🔍 ",
        selection_caret = "❯ ",
        entry_prefix = "  ",
        initial_mode = "insert",
        selection_strategy = "reset",
        sorting_strategy = "ascending",
        layout_strategy = "horizontal",
        layout_config = {
          horizontal = {
            prompt_position = "top",
            preview_width = 0.55,
            results_width = 0.8,
          },
          vertical = {
            mirror = false,
          },
          width = 0.87,
          height = 0.80,
          preview_cutoff = 120,
        },
        file_sorter = require("telescope.sorters").get_fuzzy_file,
        file_ignore_patterns = { "node_modules", ".git/", "dist/", "build/" },
        generic_sorter = require("telescope.sorters").get_generic_fuzzy_sorter,
        path_display = { "truncate" },
        winblend = 0,
        border = {},
        borderchars = { "─", "│", "─", "│", "╭", "╮", "╯", "╰" },
        color_devicons = true,
        use_less = true,
        set_env = { ["COLORTERM"] = "truecolor" },
      }
    }
    
    -- Key mappings
    local builtin = require('telescope.builtin')
    vim.keymap.set('n', '<leader>ff', builtin.find_files, {})
    vim.keymap.set('n', '<leader>fg', builtin.live_grep, {})
    vim.keymap.set('n', '<leader>fb', builtin.buffers, {})
    vim.keymap.set('n', '<leader>fh', builtin.help_tags, {})
    vim.keymap.set('n', '<leader>fs', builtin.lsp_document_symbols, {})
    vim.keymap.set('n', '<leader>fr', builtin.lsp_references, {})
    vim.keymap.set('n', '<leader>fd', builtin.lsp_definitions, {})
    vim.keymap.set('n', '<leader>ft', builtin.treesitter, {})
    vim.keymap.set('n', '<leader>fc', builtin.git_commits, {})
    vim.keymap.set('n', '<leader>gb', builtin.git_branches, {})
    ```

## Infrastructure as Code Support

### terraform-ls: Terraform Language Server
- **Description**: Official HashiCorp Terraform Language Server providing IDE features for Terraform files
- **Features**:
  - Code completion for Terraform resources, data sources, and attributes
  - Hover information for resource properties and types
  - Go to definition for resources and variables
  - Find references across Terraform files
  - Document symbols and workspace symbols
  - Diagnostics for syntax and validation errors
  - Format on save with terraform fmt
  - Rename refactoring for resources and variables
  - Document outline for quick navigation
  - Signature help for function parameters
  - Workspace support for multi-module projects
  - Provider schema integration
  - Module completion and navigation
  - Variable completion across modules
  - 1.1k+ stars and active maintenance by HashiCorp
  - Repository: https://github.com/hashicorp/terraform-ls
  - Configuration example:
    ```lua
    require('lspconfig').terraformls.setup {
      cmd = {'terraform-ls', 'serve'},
      filetypes = {'terraform', 'terraform-vars', 'hcl'},
      root_dir = function(fname)
        return util.root_pattern('.terraform', '.git')(fname) or util.path.dirname(fname)
      end,
      settings = {
        terraformls = {
          excludeRootModules = false,
          rootModules = {},
          terraformExecPath = "terraform",
          terraformLogFilePath = "",
          experimentalFeatures = {
            prefillRequiredFields = true,
            validateOnSave = true
          }
        }
      }
    }
    
    -- Key mappings for Terraform files
    vim.api.nvim_create_autocmd('FileType', {
      pattern = {'terraform', 'terraform-vars', 'hcl'},
      callback = function()
        vim.keymap.set('n', 'K', vim.lsp.buf.hover, {buffer=0})
        vim.keymap.set('n', 'gd', vim.lsp.buf.definition, {buffer=0})
        vim.keymap.set('n', 'gi', vim.lsp.buf.implementation, {buffer=0})
        vim.keymap.set('n', 'gr', vim.lsp.buf.references, {buffer=0})
        vim.keymap.set('n', '<leader>rn', vim.lsp.buf.rename, {buffer=0})
        vim.keymap.set('n', '<leader>ca', vim.lsp.buf.code_action, {buffer=0})
        vim.keymap.set('n', '<leader>f', function() vim.lsp.buf.format({async = true}) end, {buffer=0})
      end
    })
    ```

## Kubernetes and YAML Support

### yaml-language-server: YAML Language Server
- **Description**: Language Server for YAML files providing IDE features through LSP integration
- **Features**:
  - YAML validation for syntax and schema compliance
  - Schema-based validation for Kubernetes manifests
  - Auto-completion for YAML properties and values
  - Hover documentation for YAML nodes
  - Document outlining for navigation
  - Support for custom tags
  - Schema association via configuration
  - Integration with JSON Schema Store
  - Support for Kubernetes schema out of the box
  - Format on type and document formatting
  - Customizable validation settings
  - WebSocket support for remote development
  - 1.2k+ stars and active maintenance by Red Hat
  - Repository: https://github.com/redhat-developer/yaml-language-server
  - Neovim integration via coc.nvim or native LSP
  - Configuration example:
    ```lua
    require('lspconfig').yamlls.setup {
      settings = {
        yaml = {
          schemas = {
            kubernetes = "/*.k8s.yaml",
            ["https://raw.githubusercontent.com/instrumenta/kubernetes-json-schema/master/v1.18.0-standalone-strict/all.json"] = "/*.yaml",
            ["https://json.schemastore.org/github-workflow.json"] = "/.github/workflows/*",
            ["https://json.schemastore.org/kustomization.json"] = "kustomization.yaml",
            ["https://json.schemastore.org/helmfile.json"] = "helmfile.yaml",
            ["https://json.schemastore.org/helm-chart.json"] = "Chart.yaml",
            ["https://json.schemastore.org/terraform.json"] = ["*.tf", "*.tfvars"]
          },
          format = {
            enable: true,
            singleQuote: true,
            bracketSpacing: true,
            proseWrap: "always",
            printWidth: 100
          },
          validate: true,
          completion: true,
          hover: true,
          customTags: [
            "!include scalar",
            "!reference sequence"
          ],
          schemaStore: {
            enable: true,
            url: "https://www.schemastore.org/api/json/catalog.json"
          },
          editor: {
            tabSize: 2,
            formatOnType: true
          },
          keyOrdering: false,
          maxItemsComputed: 5000
        }
      }
    }
    ```

## Go Development Support

### vim-go: Go Development Plugin
- **Description**: Comprehensive Go development plugin for Vim/Neovim with extensive IDE-like features
- **Features**:
  - Code compilation with `:GoBuild` and testing with `:GoTest`
  - Quick execution with `:GoRun`
  - Improved syntax highlighting and folding
  - Integrated debugging with Delve support
  - Intelligent code completion via gopls
  - Automatic formatting on save
  - Go to definition with `:GoDef`
  - Documentation lookup with `:GoDoc`
  - Package management with `:GoImport` and `:GoDrop`
  - Type-safe renaming with `:GoRename`
  - Code coverage visualization with `:GoCoverage`
  - Struct tag management with `:GoAddTags` and `:GoRemoveTags`
  - Linting with `:GoLint` and `:GoMetaLinter`
  - Advanced source analysis with `:GoImplements`, `:GoCallees`, and `:GoReferrers`
  - Integration with Tagbar via gotags
  - Snippet engine integration
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " Basic vim-go configuration
    let g:go_fmt_command = "goimports"
    let g:go_auto_type_info = 1
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    let g:go_metalinter_autosave = 1
    let g:go_metalinter_autosave_enabled = ['vet', 'golint']
    let g:go_metalinter_deadline = "5s"
    let g:go_list_type = "quickfix"
    let g:go_auto_sameids = 1
    let g:go_addtags_transform = "camelcase"
    let g:go_def_mode = 'gopls'
    let g:go_info_mode = 'gopls'
    let g:go_rename_command = 'gopls'
    
    " Key mappings
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <leader>i <Plug>(go-info)
    autocmd FileType go nmap <leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <leader>a <Plug>(go-alternate-edit)
    autocmd FileType go nmap <leader>d <Plug>(go-def)
    autocmd FileType go nmap <leader>ds <Plug>(go-def-split)
    autocmd FileType go nmap <leader>dv <Plug>(go-def-vertical)
    autocmd FileType go nmap <leader>dt <Plug>(go-def-tab)
    autocmd FileType go nmap <leader>gd <Plug>(go-doc)
    autocmd FileType go nmap <leader>gv <Plug>(go-doc-vertical)
    autocmd FileType go nmap <leader>gb <Plug>(go-doc-browser)
    autocmd FileType go nmap <leader>s <Plug>(go-implements)
    autocmd FileType go nmap <leader>e <Plug>(go-rename)
    ```

## TypeScript Development Support

### coc.nvim: Intelligent Coding Extension
- **Description**: Nodejs extension host for Vim/Neovim that provides VSCode-like features and language server protocol integration
- **Features**:
  - Intelligent code completion with popup menu
  - Language server protocol (LSP) client implementation
  - Extensions ecosystem similar to VSCode
  - Diagnostics display in real-time
  - Code actions and quick fixes
  - Go to definition, references, and implementation
  - Symbol search in workspace
  - Rename symbol across files
  - Document formatting
  - Document highlighting
  - Document symbols outline
  - Snippet support
  - Multi-cursor support
  - Code lens support
  - Git integration
  - Floating windows for documentation
  - Extensible through JavaScript/TypeScript
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " Basic coc.nvim configuration
    " Install extensions
    let g:coc_global_extensions = [
      \ 'coc-tsserver',
      \ 'coc-json',
      \ 'coc-html',
      \ 'coc-css',
      \ 'coc-yaml',
      \ 'coc-pyright',
      \ 'coc-go',
      \ 'coc-snippets',
      \ 'coc-pairs',
      \ 'coc-prettier',
      \ 'coc-eslint',
      \ 'coc-highlight',
      \ 'coc-explorer',
      \ 'coc-git'
      \ ]
    
    " TextEdit might fail if hidden is not set
    set hidden
    
    " Some servers have issues with backup files
    set nobackup
    set nowritebackup
    
    " Give more space for displaying messages
    set cmdheight=2
    
    " Having longer updatetime (default is 4000 ms) leads to noticeable delays
    set updatetime=300
    
    " Don't pass messages to |ins-completion-menu|
    set shortmess+=c
    
    " Always show the signcolumn
    set signcolumn=yes
    
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Make <CR> auto-select the first completion item
    inoremap <silent><expr> <cr> coc#pum#visible() ? coc#pum#confirm() : "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"
    
    " Use `[g` and `]g` to navigate diagnostics
    nmap <silent> [g <Plug>(coc-diagnostic-prev)
    nmap <silent> ]g <Plug>(coc-diagnostic-next)
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Use K to show documentation in preview window
    nnoremap <silent> K :call ShowDocumentation()<CR>
    
    " Symbol renaming
    nmap <leader>rn <Plug>(coc-rename)
    
    " Formatting selected code
    xmap <leader>f  <Plug>(coc-format-selected)
    nmap <leader>f  <Plug>(coc-format-selected)
    
    " Apply code action to selected region
    xmap <leader>a  <Plug>(coc-codeaction-selected)
    nmap <leader>a  <Plug>(coc-codeaction-selected)
    
    " Apply code action to current buffer
    nmap <leader>ac  <Plug>(coc-codeaction)
    
    " Apply AutoFix to problem on the current line
    nmap <leader>qf  <Plug>(coc-fix-current)
    ```

## Core Neovim Plugins

### nvim-treesitter: Syntax Highlighting and Code Analysis
- **Description**: Treesitter configurations and abstraction layer for Neovim providing advanced syntax highlighting and code analysis
- **Features**:
  - Language-aware syntax highlighting for 100+ languages including Python, Go, TypeScript, YAML, and HCL
  - Incremental selection based on syntax tree
  - Code folding based on syntax structure
  - Indentation based on language semantics
  - Text objects based on syntax nodes
  - Consistent highlighting across languages
  - Extensible query system for custom features
  - Module system for additional functionality
  - 11.7k+ stars and active maintenance
  - Repository: https://github.com/nvim-treesitter/nvim-treesitter
  - Configuration example:
    ```lua
    require('nvim-treesitter.configs').setup {
      ensure_installed = {
        "python", "go", "typescript", "tsx", "javascript", 
        "yaml", "json", "terraform", "hcl", "bash", "lua"
      },
      sync_install = false,
      auto_install = true,
      highlight = {
        enable = true,
        additional_vim_regex_highlighting = false,
      },
      indent = {
        enable = true
      },
      incremental_selection = {
        enable = true,
        keymaps = {
          init_selection = "gnn",
          node_incremental = "grn",
          scope_incremental = "grc",
          node_decremental = "grm",
        },
      },
      textobjects = {
        select = {
          enable = true,
          lookahead = true,
          keymaps = {
            ["af"] = "@function.outer",
            ["if"] = "@function.inner",
            ["ac"] = "@class.outer",
            ["ic"] = "@class.inner",
            ["aa"] = "@parameter.outer",
            ["ia"] = "@parameter.inner",
            ["al"] = "@loop.outer",
            ["il"] = "@loop.inner",
            ["ai"] = "@conditional.outer",
            ["ii"] = "@conditional.inner",
            ["ab"] = "@block.outer",
            ["ib"] = "@block.inner",
          },
        },
        move = {
          enable = true,
          set_jumps = true,
          goto_next_start = {
            ["]m"] = "@function.outer",
            ["]]"] = "@class.outer",
          },
          goto_next_end = {
            ["]M"] = "@function.outer",
            ["]["] = "@class.outer",
          },
          goto_previous_start = {
            ["[m"] = "@function.outer",
            ["[["] = "@class.outer",
          },
          goto_previous_end = {
            ["[M"] = "@function.outer",
            ["[]"] = "@class.outer",
          },
        },
      },
    }
    ```

## Fuzzy Finding and Navigation

### telescope.nvim: Highly Extensible Fuzzy Finder
- **Description**: Powerful fuzzy finder for Neovim that enables efficient navigation and search across files, buffers, and more
- **Features**:
  - Fuzzy finding for files, buffers, git files, and more
  - Live grep for searching text across your codebase
  - Extensible architecture with plugin ecosystem
  - Customizable layouts and themes
  - Highly configurable keybindings
  - Preview window for files, grep results, and more
  - Integration with LSP for code navigation
  - Git integration for branches, commits, and status
  - Treesitter integration for symbol navigation
  - 17.3k+ stars and active maintenance
  - Repository: https://github.com/nvim-telescope/telescope.nvim
  - Configuration example:
    ```lua
    require('telescope').setup {
      defaults = {
        mappings = {
          i = {
            ["<C-j>"] = "move_selection_next",
            ["<C-k>"] = "move_selection_previous",
            ["<C-c>"] = "close",
            ["<C-u>"] = "preview_scrolling_up",
            ["<C-d>"] = "preview_scrolling_down",
          }
        },
        file_sorter = require('telescope.sorters').get_fzy_sorter,
        prompt_prefix = ' 🔍 ',
        color_devicons = true,
        file_previewer = require('telescope.previewers').vim_buffer_cat.new,
        grep_previewer = require('telescope.previewers').vim_buffer_vimgrep.new,
        qflist_previewer = require('telescope.previewers').vim_buffer_qflist.new,
      },
      extensions = {
        fzy_native = {
          override_generic_sorter = false,
          override_file_sorter = true,
        }
      }
    }
    
    -- Key mappings
    local builtin = require('telescope.builtin')
    vim.keymap.set('n', '<leader>ff', builtin.find_files, { desc = 'Find files' })
    vim.keymap.set('n', '<leader>fg', builtin.live_grep, { desc = 'Live grep' })
    vim.keymap.set('n', '<leader>fb', builtin.buffers, { desc = 'Find buffers' })
    vim.keymap.set('n', '<leader>fh', builtin.help_tags, { desc = 'Help tags' })
    vim.keymap.set('n', '<leader>fs', builtin.lsp_document_symbols, { desc = 'Document symbols' })
    vim.keymap.set('n', '<leader>fr', builtin.lsp_references, { desc = 'Find references' })
    vim.keymap.set('n', '<leader>fd', builtin.lsp_definitions, { desc = 'Find definitions' })
    vim.keymap.set('n', '<leader>ft', builtin.treesitter, { desc = 'Treesitter symbols' })
    vim.keymap.set('n', '<leader>fc', builtin.git_commits, { desc = 'Git commits' })
    vim.keymap.set('n', '<leader>fk', builtin.keymaps, { desc = 'Find keymaps' })
    ```

## Terminal Management

### toggleterm.nvim: Terminal Window Manager
- **Description**: Neovim plugin for persisting and toggling multiple terminal windows during an editing session
- **Features**:
  - Multiple terminal orientations (float, vertical, horizontal, tab)
  - Persistent terminal sessions across buffer switches
  - Terminal window shading for visual distinction
  - Custom terminal creation for specific tools (lazygit, htop, etc.)
  - Send commands to specific terminals
  - Terminal window mappings for easy navigation
  - Customizable size, direction, and appearance
  - Floating terminal with border options
  - Terminal selection UI for managing multiple instances
  - Ability to send lines/selections from buffers to terminals
  - Support for parallel terminal workflows
  - 4.8k+ stars and active maintenance
  - Repository: https://github.com/akinsho/toggleterm.nvim
  - Configuration example:
    ```lua
    require('toggleterm').setup {
      size = function(term)
        if term.direction == "horizontal" then
          return 15
        elseif term.direction == "vertical" then
          return vim.o.columns * 0.4
        end
      end,
      open_mapping = [[<c-\>]],
      hide_numbers = true,
      shade_terminals = true,
      start_in_insert = true,
      insert_mappings = true,
      terminal_mappings = true,
      persist_size = true,
      persist_mode = true,
      direction = 'float',
      close_on_exit = true,
      shell = vim.o.shell,
      auto_scroll = true,
      float_opts = {
        border = 'curved',
        width = 120,
        height = 30,
        winblend = 3,
      }
    }
    
    -- Custom terminal example for parallel workflows
    local Terminal = require('toggleterm.terminal').Terminal
    
    -- Create specialized terminals for different tasks
    local python_repl = Terminal:new({ 
      cmd = "python", 
      direction = "horizontal",
      hidden = true,
      count = 1
    })
    
    local go_test = Terminal:new({ 
      cmd = "go test ./...", 
      direction = "vertical",
      hidden = true,
      count = 2
    })
    
    local typescript_watch = Terminal:new({ 
      cmd = "npm run watch", 
      direction = "float",
      hidden = true,
      count = 3
    })
    
    -- Functions to toggle specific terminals
    function _python_toggle()
      python_repl:toggle()
    end
    
    function _go_test_toggle()
      go_test:toggle()
    end
    
    function _typescript_watch_toggle()
      typescript_watch:toggle()
    end
    
    -- Key mappings for toggling terminals
    vim.api.nvim_set_keymap("n", "<leader>tp", "<cmd>lua _python_toggle()<CR>", {noremap = true, silent = true})
    vim.api.nvim_set_keymap("n", "<leader>tg", "<cmd>lua _go_test_toggle()<CR>", {noremap = true, silent = true})
    vim.api.nvim_set_keymap("n", "<leader>tt", "<cmd>lua _typescript_watch_toggle()<CR>", {noremap = true, silent = true})
    ```

## Infrastructure as Code Support

### terraform-ls: Terraform Language Server
- **Description**: Official HashiCorp Terraform Language Server providing IDE features to any LSP-compatible editor
- **Features**:
  - Code completion for Terraform resources, data sources, and attributes
  - Hover information for resource properties and types
  - Go to definition for resources and variables
  - Find references across Terraform files
  - Document symbols for quick navigation
  - Diagnostics for syntax and semantic errors
  - Format on save using `terraform fmt`
  - Workspace symbol search
  - Module navigation and exploration
  - Provider schema support
  - Variable completion across modules
  - Terraform version detection
  - 1.1k+ stars and active maintenance by HashiCorp
  - Repository: https://github.com/hashicorp/terraform-ls
  - Configuration example:
    ```lua
    require('lspconfig').terraformls.setup {
      cmd = { "terraform-ls", "serve" },
      filetypes = { "terraform", "terraform-vars", "hcl" },
      root_dir = function(fname)
        return util.root_pattern(".terraform", ".git")(fname) or util.path.dirname(fname)
      end,
      settings = {
        terraform = {
          path = "terraform",
          logLevel = "info",
          excludeRootModules = false,
          ignoreDirectoryNames = [],
          experimentalFeatures = {
            validateOnSave = true,
            prefillRequiredFields = true
          }
        }
      },
      on_attach = function(client, bufnr)
        -- Enable completion triggered by <c-x><c-o>
        vim.api.nvim_buf_set_option(bufnr, 'omnifunc', 'v:lua.vim.lsp.omnifunc')
        
        -- Key mappings
        local bufopts = { noremap=true, silent=true, buffer=bufnr }
        vim.keymap.set('n', 'gD', vim.lsp.buf.declaration, bufopts)
        vim.keymap.set('n', 'gd', vim.lsp.buf.definition, bufopts)
        vim.keymap.set('n', 'K', vim.lsp.buf.hover, bufopts)
        vim.keymap.set('n', 'gi', vim.lsp.buf.implementation, bufopts)
        vim.keymap.set('n', '<C-k>', vim.lsp.buf.signature_help, bufopts)
        vim.keymap.set('n', '<leader>wa', vim.lsp.buf.add_workspace_folder, bufopts)
        vim.keymap.set('n', '<leader>wr', vim.lsp.buf.remove_workspace_folder, bufopts)
        vim.keymap.set('n', '<leader>wl', function()
          print(vim.inspect(vim.lsp.buf.list_workspace_folders()))
        end, bufopts)
        vim.keymap.set('n', '<leader>D', vim.lsp.buf.type_definition, bufopts)
        vim.keymap.set('n', '<leader>rn', vim.lsp.buf.rename, bufopts)
        vim.keymap.set('n', '<leader>ca', vim.lsp.buf.code_action, bufopts)
        vim.keymap.set('n', 'gr', vim.lsp.buf.references, bufopts)
        vim.keymap.set('n', '<leader>f', function() vim.lsp.buf.format { async = true } end, bufopts)
      end
    }
    ```

## Kubernetes and YAML Support

### yaml-language-server: Language Server for YAML Files
- **Description**: Red Hat-developed Language Server Protocol implementation for YAML files with comprehensive schema support
- **Features**:
  - YAML validation against JSON Schema 7 and below
  - Detailed error and warning detection
  - Auto-completion for all commands and schema defaults
  - Hover support showing descriptions
  - Document outlining for all nodes
  - Schema association via glob patterns
  - Support for custom YAML tags
  - Integration with SchemaStore.org
  - Kubernetes schema support built-in
  - Support for multiple YAML versions (1.1, 1.2)
  - Customizable formatting options
  - 1.2k+ stars and active maintenance by Red Hat
  - Repository: https://github.com/redhat-developer/yaml-language-server
  - Configuration example:
    ```lua
    require('lspconfig').yamlls.setup {
      settings = {
        yaml = {
          schemas = {
            kubernetes = "/*.k8s.yaml",
            ["https://raw.githubusercontent.com/instrumenta/kubernetes-json-schema/master/v1.18.0-standalone-strict/all.json"] = "/*.yaml",
            ["https://json.schemastore.org/github-workflow.json"] = "/.github/workflows/*",
            ["https://json.schemastore.org/terraform.json"] = "/terraform/*.tf",
            ["https://json.schemastore.org/kustomization.json"] = "/kustomization.yaml",
            ["https://json.schemastore.org/helmfile.json"] = "/helmfile.yaml",
            ["https://json.schemastore.org/chart.json"] = "/Chart.yaml"
          },
          format = {
            enable = true,
            singleQuote = true,
            bracketSpacing = true,
            proseWrap = "preserve",
            printWidth = 120
          },
          validate = true,
          completion = true,
          hover = true,
          schemaStore = {
            enable = true,
            url = "https://www.schemastore.org/api/json/catalog.json"
          },
          customTags = [
            "!include scalar",
            "!reference sequence",
            "!Ref scalar"
          ],
          keyOrdering = false,
          maxItemsComputed = 5000,
          yamlVersion = "1.2"
        }
      }
    }
    ```

## Go Development Support

### vim-go: Go Development Plugin for Vim
- **Description**: Comprehensive Go development plugin for Vim with extensive features and tooling integration
- **Features**:
  - Compile, install, and test Go packages with `:GoBuild`, `:GoInstall`, `:GoTest`
  - Run current file(s) with `:GoRun`
  - Improved syntax highlighting and folding for Go files
  - Integrated debugging with Delve support via `:GoDebugStart`
  - LSP integration with gopls for intelligent code completion
  - Format on save with cursor position and undo history preservation
  - Go to symbol/declaration with `:GoDef`
  - Documentation lookup with `:GoDoc` and `:GoDocBrowser`
  - Package management with `:GoImport` and `:GoDrop`
  - Type-safe renaming with `:GoRename`
  - Code coverage visualization with `:GoCoverage`
  - Struct tag management with `:GoAddTags` and `:GoRemoveTags`
  - Linting with `:GoLint`, `:GoMetaLinter`, and `:GoVet`
  - Advanced source analysis with `:GoImplements`, `:GoCallees`, and `:GoReferrers`
  - Integration with Tagbar via gotags
  - Snippet engine integration
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " vim-go configuration
    let g:go_fmt_command = "goimports"
    let g:go_auto_type_info = 1
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    let g:go_metalinter_autosave = 1
    let g:go_metalinter_autosave_enabled = ['vet', 'golint']
    let g:go_metalinter_deadline = "5s"
    let g:go_def_mode = 'gopls'
    let g:go_info_mode = 'gopls'
    let g:go_rename_command = 'gopls'
    let g:go_gopls_enabled = 1
    let g:go_gopls_options = ['-remote=auto']
    let g:go_list_type = "quickfix"
    let g:go_addtags_transform = "camelcase"
    let g:go_term_enabled = 1
    let g:go_term_mode = "split"
    let g:go_term_height = 15
    let g:go_term_reuse = 1
    
    " Key mappings
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <leader>i <Plug>(go-info)
    autocmd FileType go nmap <leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <leader>d <Plug>(go-def)
    autocmd FileType go nmap <leader>v <Plug>(go-def-vertical)
    autocmd FileType go nmap <leader>s <Plug>(go-def-split)
    autocmd FileType go nmap <leader>dt <Plug>(go-def-tab)
    autocmd FileType go nmap <leader>e <Plug>(go-rename)
    autocmd FileType go nmap <leader>a <Plug>(go-alternate-edit)
    autocmd FileType go nmap <leader>av <Plug>(go-alternate-vertical)
    ```

## TypeScript Development Support

### coc.nvim: Intelligent Coding Experience for Vim/Neovim
- **Description**: VSCode-like experience for Vim/Neovim with full language server protocol support
- **Features**:
  - Intelligent code completion with automatic imports
  - Diagnostics and error checking as you type
  - Code actions and quick fixes with auto-imports
  - Go to definition, references, implementation, and type definition
  - Symbol search in workspace and current document
  - Rename symbol across files
  - Document formatting with multiple formatter support
  - Document highlighting and symbol highlighting
  - List support for various operations (symbols, diagnostics, etc.)
  - Snippet support with integration for popular snippet engines
  - Extensions ecosystem similar to VSCode
  - Multi-language support through language servers
  - Floating windows for documentation and signature help
  - Status line integration
  - Extensible through Node.js
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " coc.nvim configuration
    " Install extensions
    let g:coc_global_extensions = [
      \ 'coc-tsserver',
      \ 'coc-json',
      \ 'coc-html',
      \ 'coc-css',
      \ 'coc-yaml',
      \ 'coc-pyright',
      \ 'coc-go',
      \ 'coc-snippets',
      \ 'coc-pairs',
      \ 'coc-prettier',
      \ 'coc-eslint',
      \ 'coc-markdownlint',
      \ 'coc-highlight',
      \ 'coc-explorer',
      \ 'coc-git'
    \ ]
    
    " TextEdit might fail if hidden is not set
    set hidden
    
    " Some servers have issues with backup files
    set nobackup
    set nowritebackup
    
    " Give more space for displaying messages
    set cmdheight=2
    
    " Having longer updatetime (default is 4000 ms = 4 s) leads to noticeable
    " delays and poor user experience
    set updatetime=300
    
    " Don't pass messages to |ins-completion-menu|
    set shortmess+=c
    
    " Always show the signcolumn, otherwise it would shift the text each time
    " diagnostics appear/become resolved
    set signcolumn=yes
    
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Make <CR> auto-select the first completion item and notify coc.nvim to
    " format on enter, <cr> could be remapped by other vim plugin
    inoremap <silent><expr> <cr> coc#pum#visible() ? coc#pum#confirm()
                                  \: "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"
    
    " Use `[g` and `]g` to navigate diagnostics
    nmap <silent> [g <Plug>(coc-diagnostic-prev)
    nmap <silent> ]g <Plug>(coc-diagnostic-next)
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Use K to show documentation in preview window
    nnoremap <silent> K :call ShowDocumentation()<CR>
    
    function! ShowDocumentation()
      if CocAction('hasProvider', 'hover')
        call CocActionAsync('doHover')
      else
        call feedkeys('K', 'in')
      endif
    endfunction
    
    " Highlight the symbol and its references when holding the cursor
    autocmd CursorHold * silent call CocActionAsync('highlight')
    
    " Symbol renaming
    nmap <leader>rn <Plug>(coc-rename)
    
    " Formatting selected code
    xmap <leader>f  <Plug>(coc-format-selected)
    nmap <leader>f  <Plug>(coc-format-selected)
    
    " Apply AutoFix to problem on the current line
    nmap <leader>qf  <Plug>(coc-fix-current)
    
    " Run the Code Lens action on the current line
    nmap <leader>cl  <Plug>(coc-codelens-action)
    ```

### coc-tsserver: TypeScript Language Server Extension for coc.nvim
- **Description**: Official TypeScript language server extension for coc.nvim
- **Features**:
  - TypeScript and JavaScript language support
  - Intelligent code completion with auto-imports
  - Type checking and error reporting
  - Code navigation and refactoring
  - Automatic import organization
  - JSX/TSX support
  - Project-wide symbol search
  - Code formatting
  - Snippet support
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```json
    // coc-settings.json
    {
      "tsserver.enable": true,
      "tsserver.enableJavascript": true,
      "tsserver.formatOnType": true,
      "tsserver.implicitProjectConfig.experimentalDecorators": true,
      "tsserver.implicitProjectConfig.checkJs": true,
      "tsserver.implicitProjectConfig.strictNullChecks": true,
      "tsserver.trace.server": "off",
      "tsserver.maxTsServerMemory": 4096,
      "typescript.format.enabled": true,
      "typescript.format.insertSpaceAfterCommaDelimiter": true,
      "typescript.format.insertSpaceAfterConstructor": false,
      "typescript.format.insertSpaceAfterSemicolonInForStatements": true,
      "typescript.format.insertSpaceBeforeAndAfterBinaryOperators": true,
      "typescript.format.insertSpaceAfterKeywordsInControlFlowStatements": true,
      "typescript.format.insertSpaceAfterFunctionKeywordForAnonymousFunctions": true,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces": true,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets": false,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis": false,
      "typescript.format.insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces": false,
      "typescript.format.insertSpaceAfterTypeAssertion": false,
      "typescript.format.placeOpenBraceOnNewLineForFunctions": false,
      "typescript.format.placeOpenBraceOnNewLineForControlBlocks": false,
      "typescript.preferences.quoteStyle": "single",
      "typescript.preferences.importModuleSpecifier": "relative",
      "typescript.updateImportsOnFileMove.enable": true,
      "typescript.suggest.autoImports": true,
      "typescript.suggest.completeFunctionCalls": true,
      "typescript.suggest.includeCompletionsForImportStatements": true,
      "typescript.suggest.paths": true
    }
    ```

## Navigation and Search

### telescope.nvim: Fuzzy Finder for Neovim
- **Description**: Highly extensible fuzzy finder over lists with powerful search capabilities
- **Features**:
  - Fast fuzzy finding for files, buffers, and various other sources
  - Extensible architecture with plugin support
  - Multiple layout strategies (horizontal, vertical, cursor, center, dropdown, ivy)
  - Customizable key mappings and themes
  - Preview window for files, git diffs, and more
  - Integration with git, LSP, treesitter, and other tools
  - Builtin pickers for files, buffers, grep, git operations, LSP symbols, and more
  - Sorter plugins for improved performance (fzf, fzy)
  - Highly customizable through Lua API
  - 17.3k+ stars and active maintenance
  - Repository: https://github.com/nvim-telescope/telescope.nvim
  - Configuration example:
    ```lua
    -- Telescope setup with defaults
    require('telescope').setup {
      defaults = {
        mappings = {
          i = {
            ["<C-j>"] = "move_selection_next",
            ["<C-k>"] = "move_selection_previous",
            ["<C-c>"] = "close",
            ["<C-u>"] = "preview_scrolling_up",
            ["<C-d>"] = "preview_scrolling_down",
          }
        },
        layout_strategy = 'horizontal',
        layout_config = {
          horizontal = {
            width = 0.9,
            height = 0.85,
            preview_width = 0.5,
          },
          vertical = {
            width = 0.9,
            height = 0.85,
            preview_height = 0.5,
          },
        },
        path_display = { "truncate" },
        winblend = 0,
        border = {},
        borderchars = { '─', '│', '─', '│', '┌', '┐', '┘', '└' },
        color_devicons = true,
        set_env = { ['COLORTERM'] = 'truecolor' },
      },
      pickers = {
        find_files = {
          hidden = true,
          find_command = { "rg", "--files", "--hidden", "--glob", "!**/.git/*" }
        },
        live_grep = {
          additional_args = function() return { "--hidden" } end
        },
        buffers = {
          show_all_buffers = true,
          sort_lastused = true,
          mappings = {
            i = {
              ["<c-d>"] = "delete_buffer",
            }
          }
        }
      },
      extensions = {
        fzf = {
          fuzzy = true,
          override_generic_sorter = true,
          override_file_sorter = true,
          case_mode = "smart_case",
        }
      }
    }
    
    -- Load extensions
    require('telescope').load_extension('fzf')
    
    -- Key mappings
    local builtin = require('telescope.builtin')
    vim.keymap.set('n', '<leader>ff', builtin.find_files, { desc = 'Find files' })
    vim.keymap.set('n', '<leader>fg', builtin.live_grep, { desc = 'Live grep' })
    vim.keymap.set('n', '<leader>fb', builtin.buffers, { desc = 'Buffers' })
    vim.keymap.set('n', '<leader>fh', builtin.help_tags, { desc = 'Help tags' })
    vim.keymap.set('n', '<leader>fs', builtin.lsp_document_symbols, { desc = 'Document symbols' })
    vim.keymap.set('n', '<leader>fr', builtin.lsp_references, { desc = 'References' })
    vim.keymap.set('n', '<leader>fd', builtin.lsp_definitions, { desc = 'Definitions' })
    vim.keymap.set('n', '<leader>fi', builtin.lsp_implementations, { desc = 'Implementations' })
    vim.keymap.set('n', '<leader>ft', builtin.lsp_type_definitions, { desc = 'Type definitions' })
    vim.keymap.set('n', '<leader>gc', builtin.git_commits, { desc = 'Git commits' })
    vim.keymap.set('n', '<leader>gb', builtin.git_branches, { desc = 'Git branches' })
    vim.keymap.set('n', '<leader>gs', builtin.git_status, { desc = 'Git status' })
    ```

## Terminal Management

### toggleterm.nvim: Terminal Window Management for Neovim
- **Description**: Neovim plugin to persist and toggle multiple terminal windows during an editing session
- **Features**:
  - Multiple terminal orientations (horizontal, vertical, float, tab)
  - Persistent terminals that maintain state between toggles
  - Custom terminal creation for specific tools (lazygit, htop, etc.)
  - Terminal window shading for visual distinction
  - Customizable keybindings for terminal toggling
  - Send commands to specific terminals
  - Terminal window mappings for easy navigation
  - Floating terminal support with customizable borders
  - Terminal window sizing and positioning
  - Terminal window naming for better organization
  - Support for multiple terminals side-by-side
  - Automatic scrolling to bottom on terminal output
  - Callbacks for terminal events (open, close, stdout, stderr, exit)
  - 4.8k+ stars and active maintenance
  - Repository: https://github.com/akinsho/toggleterm.nvim
  - Configuration example:
    ```lua
    require('toggleterm').setup {
      size = function(term)
        if term.direction == "horizontal" then
          return 15
        elseif term.direction == "vertical" then
          return vim.o.columns * 0.4
        end
      end,
      open_mapping = [[<c-\>]],
      hide_numbers = true,
      shade_terminals = true,
      start_in_insert = true,
      insert_mappings = true,
      terminal_mappings = true,
      persist_size = true,
      persist_mode = true,
      direction = 'float',
      close_on_exit = true,
      shell = vim.o.shell,
      auto_scroll = true,
      float_opts = {
        border = 'curved',
        width = 120,
        height = 30,
        winblend = 3,
      },
      winbar = {
        enabled = true,
        name_formatter = function(term)
          return term.name or term.id
        end
      }
    }
    
    -- Custom terminal example for lazygit
    local Terminal = require('toggleterm.terminal').Terminal
    local lazygit = Terminal:new({
      cmd = "lazygit",
      dir = "git_dir",
      direction = "float",
      float_opts = {
        border = "double",
      },
      on_open = function(term)
        vim.cmd("startinsert!")
        vim.api.nvim_buf_set_keymap(term.bufnr, "n", "q", "<cmd>close<CR>", {noremap = true, silent = true})
      end,
    })
    
    function _lazygit_toggle()
      lazygit:toggle()
    end
    
    vim.api.nvim_set_keymap("n", "<leader>g", "<cmd>lua _lazygit_toggle()<CR>", {noremap = true, silent = true})
    
    -- Terminal window mappings
    function _G.set_terminal_keymaps()
      local opts = {buffer = 0}
      vim.keymap.set('t', '<esc>', [[<C-\><C-n>]], opts)
      vim.keymap.set('t', 'jk', [[<C-\><C-n>]], opts)
      vim.keymap.set('t', '<C-h>', [[<Cmd>wincmd h<CR>]], opts)
      vim.keymap.set('t', '<C-j>', [[<Cmd>wincmd j<CR>]], opts)
      vim.keymap.set('t', '<C-k>', [[<Cmd>wincmd k<CR>]], opts)
      vim.keymap.set('t', '<C-l>', [[<Cmd>wincmd l<CR>]], opts)
    end
    
    vim.cmd('autocmd! TermOpen term://* lua set_terminal_keymaps()')
    ```

## Syntax Highlighting and Code Analysis

### nvim-treesitter: Syntax Highlighting and Code Analysis for Neovim
- **Description**: Treesitter configurations and abstraction layer for Neovim providing consistent syntax highlighting
- **Features**:
  - Consistent syntax highlighting across multiple languages
  - Incremental selection based on syntax tree
  - Tree-sitter based folding
  - Indentation based on syntax tree
  - Support for over 100 languages including Python, Go, TypeScript, YAML, and HCL (Terraform)
  - Language-specific queries for highlighting, indentation, and folding
  - Extensible architecture with module system
  - Integration with other plugins via API
  - Automatic parser installation and management
  - Customizable highlighting groups
  - 11.7k+ stars and active maintenance
  - Repository: https://github.com/nvim-treesitter/nvim-treesitter
  - Configuration example:
    ```lua
    require('nvim-treesitter.configs').setup {
      -- A list of parser names, or "all" (the listed parsers MUST always be installed)
      ensure_installed = {
        "python", "go", "typescript", "tsx", "javascript", "yaml", "json", "hcl", "terraform",
        "bash", "lua", "vim", "vimdoc", "query", "markdown", "markdown_inline"
      },
    
      -- Install parsers synchronously (only applied to `ensure_installed`)
      sync_install = false,
    
      -- Automatically install missing parsers when entering buffer
      auto_install = true,
    
      -- List of parsers to ignore installing (or "all")
      ignore_install = {},
    
      highlight = {
        enable = true,
        
        -- NOTE: these are the names of the parsers and not the filetype. (for example if you want to
        -- disable highlighting for the `tex` filetype, you need to include `latex` in this list as this is
        -- the name of the parser)
        disable = {},
    
        -- Setting this to true will run `:h syntax` and tree-sitter at the same time.
        -- Set this to `true` if you depend on 'syntax' being enabled (like for indentation).
        -- Using this option may slow down your editor, and you may see some duplicate highlights.
        -- Instead of true it can also be a list of languages
        additional_vim_regex_highlighting = false,
      },
      
      incremental_selection = {
        enable = true,
        keymaps = {
          init_selection = "gnn",
          node_incremental = "grn",
          scope_incremental = "grc",
          node_decremental = "grm",
        },
      },
      
      indent = {
        enable = true
      },
      
      -- Additional modules for specific languages
      playground = {
        enable = true,
        disable = {},
        updatetime = 25, -- Debounced time for highlighting nodes in the playground from source code
        persist_queries = false, -- Whether the query persists across vim sessions
        keybindings = {
          toggle_query_editor = 'o',
          toggle_hl_groups = 'i',
          toggle_injected_languages = 't',
          toggle_anonymous_nodes = 'a',
          toggle_language_display = 'I',
          focus_language = 'f',
          unfocus_language = 'F',
          update = 'R',
          goto_node = '<cr>',
          show_help = '?',
        },
      }
    }
    
    -- Set foldmethod to expr and use treesitter for folding
    vim.wo.foldmethod = 'expr'
    vim.wo.foldexpr = 'nvim_treesitter#foldexpr()'
    ```

## Terraform Development Support

### vim-terraform: Terraform and HCL Integration for Vim/Neovim
- **Description**: Basic vim/terraform integration with syntax highlighting, indentation, and commands
- **Features**:
  - Syntax highlighting for Terraform (`.tf`, `.tfvars`) and HCL files
  - Proper indentation for Terraform and HCL files
  - `:Terraform` command with tab completion of subcommands
  - `:TerraformFmt` command for formatting Terraform files
  - Auto-formatting on save (optional)
  - Support for `.tfbackend` files
  - Support for HCL files
  - Filetype detection for Terraform and HCL files
  - Integration with Terraform CLI
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/hashivim/vim-terraform
  - Configuration example:
    ```vim
    " vim-terraform configuration
    let g:terraform_align = 1
    let g:terraform_fold_sections = 0
    let g:terraform_fmt_on_save = 1
    let g:terraform_binary_path = '/usr/local/bin/terraform'
    
    " Key mappings
    autocmd FileType terraform nnoremap <leader>ti :!terraform init<CR>
    autocmd FileType terraform nnoremap <leader>tv :!terraform validate<CR>
    autocmd FileType terraform nnoremap <leader>tp :!terraform plan<CR>
    autocmd FileType terraform nnoremap <leader>ta :!terraform apply<CR>
    autocmd FileType terraform nnoremap <leader>td :!terraform destroy<CR>
    
    " Comment string for commentary.vim
    autocmd FileType terraform setlocal commentstring=#%s
    ```

## Kubernetes Development Support

### kube-utils-nvim: Kubernetes and Helm Integration for Neovim
- **Description**: Neovim plugin providing seamless integration with Kubernetes and Helm
- **Features**:
  - Kubernetes Context and Namespace Management
  - CRD Viewer for Custom Resource Definitions
  - Helm Integration for chart management and deployment
  - Log Viewer and Formatting with JSON formatting
  - Telescope Integration for context/namespace selection
  - LSP Integration for YAML and Helm
  - K9s Integration directly from Neovim
  - Multiple terminal orientations for K9s (split view)
  - Kubectl command integration
  - 41+ stars and active maintenance (updated 5 days ago)
  - Repository: https://github.com/h4ckm1n-dev/kube-utils-nvim
  - Configuration example:
    ```lua
    -- ~/.config/nvim/lua/plugins/kube-utils.lua
    return {
      {
        "h4ckm1n-dev/kube-utils-nvim",
        dependencies = { "nvim-telescope/telescope.nvim" },
        lazy = true,
        event = "VeryLazy"
      },
    }
    
    -- ~/.config/nvim/lua/config/keymaps.lua
    local kube_utils_mappings = {
      { "<leader>k", group = "Kubernetes" }, -- Main title for all Kubernetes related commands
      -- Helm Commands
      { "<leader>kh", group = "Helm" },
      { "<leader>khT", "<cmd>HelmDryRun<CR>", desc = "Helm DryRun Buffer" },
      { "<leader>khb", "<cmd>HelmDependencyBuildFromBuffer<CR>", desc = "Helm Dependency Build" },
      { "<leader>khd", "<cmd>HelmDeployFromBuffer<CR>", desc = "Helm Deploy Buffer to Context" },
      { "<leader>khr", "<cmd>RemoveDeployment<CR>", desc = "Helm Remove Deployment From Buffer" },
      { "<leader>kht", "<cmd>HelmTemplateFromBuffer<CR>", desc = "Helm Template From Buffer" },
      { "<leader>khu", "<cmd>HelmDependencyUpdateFromBuffer<CR>", desc = "Helm Dependency Update" },
      -- Kubectl Commands
      { "<leader>kk", group = "Kubectl" },
      { "<leader>kkC", "<cmd>SelectSplitCRD<CR>", desc = "Download CRD Split" },
      { "<leader>kkD", "<cmd>DeleteNamespace<CR>", desc = "Kubectl Delete Namespace" },
      { "<leader>kkK", "<cmd>OpenK9s<CR>", desc = "Open K9s" },
      { "<leader>kka", "<cmd>KubectlApplyFromBuffer<CR>", desc = "Kubectl Apply From Buffer" },
      { "<leader>kkc", "<cmd>SelectCRD<CR>", desc = "Download CRD" },
      { "<leader>kkk", "<cmd>OpenK9sSplit<CR>", desc = "Split View K9s" },
      { "<leader>kkl", "<cmd>ToggleYamlHelm<CR>", desc = "Toggle YAML/Helm" },
      -- Logs Commands
      { "<leader>kl", group = "Logs" },
      { "<leader>klf", "<cmd>JsonFormatLogs<CR>", desc = "Format JSON" },
      { "<leader>klv", "<cmd>ViewPodLogs<CR>", desc = "View Pod Logs" },
    }
    -- Register the Kube Utils keybindings
    require('which-key').add(kube_utils_mappings)
    ```

## Go Development Support

### vim-go: Go Development Plugin for Vim/Neovim
- **Description**: Comprehensive Go development plugin for Vim/Neovim
- **Features**:
  - Compile, install, test, and run Go code directly from Vim
  - Improved syntax highlighting and folding
  - Integrated debugging with Delve
  - Completion and other features via gopls
  - Auto-formatting on save with cursor position preservation
  - Go to symbol/declaration with :GoDef
  - Documentation lookup with :GoDoc and :GoDocBrowser
  - Package management with :GoImport and :GoDrop
  - Type-safe renaming with :GoRename
  - Code coverage visualization with :GoCoverage
  - Struct tag management with :GoAddTags and :GoRemoveTags
  - Linting with :GoLint, :GoMetaLinter, and :GoVet
  - Advanced source analysis with gopls
  - Integration with Tagbar via gotags
  - Integration with Ultisnips and other snippet engines
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " vim-go configuration
    let g:go_fmt_command = "goimports"
    let g:go_autodetect_gopath = 0
    let g:go_list_type = "quickfix"
    
    " Highlighting
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_generate_tags = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_build_constraints = 1
    
    " Auto formatting and importing
    let g:go_fmt_autosave = 1
    let g:go_fmt_command = "goimports"
    
    " Status line types/signatures
    let g:go_auto_type_info = 1
    
    " Run :GoBuild or :GoTestCompile based on the go file
    function! s:build_go_files()
      let l:file = expand('%')
      if l:file =~# '^\f\+_test\.go$'
        call go#test#Test(0, 1)
      elseif l:file =~# '^\f\+\.go$'
        call go#cmd#Build(0)
      endif
    endfunction
    
    " Key mappings
    autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <Leader>i <Plug>(go-info)
    autocmd FileType go nmap <Leader>l <Plug>(go-metalinter)
    autocmd FileType go nmap <Leader>d <Plug>(go-def)
    autocmd FileType go nmap <Leader>v <Plug>(go-def-vertical)
    autocmd FileType go nmap <Leader>s <Plug>(go-def-split)
    autocmd FileType go nmap <Leader>dt <Plug>(go-def-tab)
    autocmd FileType go nmap <Leader>e <Plug>(go-rename)
    autocmd FileType go nmap <Leader>a <Plug>(go-alternate-edit)
    autocmd FileType go nmap <Leader>av <Plug>(go-alternate-vertical)
    ```

## TypeScript Development Support

### coc.nvim: Intelligent Coding Experience for Vim/Neovim
- **Description**: VSCode-like experience for Vim/Neovim with language server protocol support
- **Features**:
  - Intellisense code completion with automatic imports
  - Diagnostic errors and warnings with inline text
  - Code actions and quick fixes with lightbulb
  - Hover for documentation and type information
  - Find references, definition, implementation, and type definition
  - Rename symbol across files
  - Format document or selection
  - Code outline and symbol search
  - Snippet support with placeholder jumping
  - List support for fuzzy search
  - Extensions like VSCode for additional functionality
  - Floating windows for documentation and completion details
  - Multiple language server support
  - Workspace folder support
  - Extensive configuration options
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " coc.nvim configuration
    " Install extensions
    let g:coc_global_extensions = [
      \ 'coc-tsserver',
      \ 'coc-json',
      \ 'coc-pyright',
      \ 'coc-go',
      \ 'coc-yaml',
      \ 'coc-snippets',
      \ 'coc-prettier',
      \ 'coc-eslint',
      \ 'coc-markdownlint',
      \ 'coc-highlight',
      \ 'coc-explorer',
      \ ]
    
    " TextEdit might fail if hidden is not set
    set hidden
    
    " Some servers have issues with backup files
    set nobackup
    set nowritebackup
    
    " Give more space for displaying messages
    set cmdheight=2
    
    " Having longer updatetime (default is 4000 ms) leads to noticeable delays
    set updatetime=300
    
    " Don't pass messages to |ins-completion-menu|
    set shortmess+=c
    
    " Always show the signcolumn, otherwise it would shift the text
    set signcolumn=yes
    
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Make <CR> auto-select the first completion item
    inoremap <silent><expr> <cr> coc#pum#visible() ? coc#pum#confirm() : "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"
    
    " Use `[g` and `]g` to navigate diagnostics
    nmap <silent> [g <Plug>(coc-diagnostic-prev)
    nmap <silent> ]g <Plug>(coc-diagnostic-next)
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Use K to show documentation in preview window
    nnoremap <silent> K :call ShowDocumentation()<CR>
    
    " Symbol renaming
    nmap <leader>rn <Plug>(coc-rename)
    
    " Formatting selected code
    xmap <leader>f  <Plug>(coc-format-selected)
    nmap <leader>f  <Plug>(coc-format-selected)
    
    " Apply codeAction to the selected region
    xmap <leader>a  <Plug>(coc-codeaction-selected)
    nmap <leader>a  <Plug>(coc-codeaction-selected)
    
    " Apply codeAction to the current buffer
    nmap <leader>ac  <Plug>(coc-codeaction)
    
    " Apply AutoFix to problem on the current line
    nmap <leader>qf  <Plug>(coc-fix-current)
    ```

## Python Development Support

### jedi-vim: Python Autocompletion for Vim/Neovim
- **Description**: Comprehensive Python autocompletion and code navigation plugin
- **Features**:
  - Intelligent code completion for Python
  - Go to definition, assignment, and stub functionality
  - Documentation/Pydoc integration with popup display
  - Symbol renaming across files
  - Usage finding for variables and functions
  - Module importing with :Pyimport command
  - Integration with the Jedi Python library
  - Support for virtual environments
  - Case-insensitive completion options
  - Popup or preview window documentation display
  - 5.3k+ stars and active maintenance
  - Repository: https://github.com/davidhalter/jedi-vim
  - Configuration example:
    ```vim
    " jedi-vim configuration
    " Disable autocompletion (using deoplete instead)
    let g:jedi#completions_enabled = 0
    
    " Key bindings
    let g:jedi#goto_command = "<leader>d"
    let g:jedi#goto_assignments_command = "<leader>g"
    let g:jedi#goto_stubs_command = "<leader>s"
    let g:jedi#documentation_command = "K"
    let g:jedi#usages_command = "<leader>n"
    let g:jedi#rename_command = "<leader>r"
    
    " Open new split (not buffer) when using goto
    let g:jedi#use_splits_not_buffers = "right"
    
    " Show call signatures in command line instead of popup
    let g:jedi#show_call_signatures = "2"
    
    " Set Python environment path for virtual env support
    let g:jedi#environment_path = "venv"
    
    " Disable automatic initialization
    let g:jedi#auto_initialization = 0
    
    " Disable automatic vim configuration
    let g:jedi#auto_vim_configuration = 0
    ```

## Terraform Development Support

### vim-terraform: Terraform Integration for Vim/Neovim
- **Description**: Basic Vim/Terraform integration for infrastructure-as-code development
- **Features**:
  - Syntax highlighting for Terraform files (.tf, .tfvars, .tfbackend)
  - Automatic indentation for Terraform and HCL files
  - Terraform command integration with `:Terraform` command
  - Tab completion for Terraform subcommands
  - Automatic formatting with `:TerraformFmt` command
  - Option to format on save with `terraform_fmt_on_save`
  - Support for HCL files
  - Filetype detection for Terraform-related files
  - Integration with Terraform CLI tools
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/hashivim/vim-terraform
  - Configuration example:
    ```vim
    " vim-terraform configuration
    " Allow vim-terraform to align settings automatically with Tabularize
    let g:terraform_align=1
    
    " Allow vim-terraform to automatically fold (hide until unfolded) sections
    let g:terraform_fold_sections=1
    
    " Allow vim-terraform to automatically format .tf and .tfvars files with terraform fmt
    let g:terraform_fmt_on_save=1
    
    " Run terraform commands in a horizontal split window
    let g:terraform_command_window_split='horizontal'
    
    " Set the window height for the terraform command window
    let g:terraform_command_window_height=20
    
    " Key mappings
    autocmd FileType terraform nmap <leader>ti :!terraform init<CR>
    autocmd FileType terraform nmap <leader>tp :!terraform plan<CR>
    autocmd FileType terraform nmap <leader>ta :!terraform apply<CR>
    autocmd FileType terraform nmap <leader>td :!terraform destroy<CR>
    autocmd FileType terraform nmap <leader>tv :!terraform validate<CR>
    autocmd FileType terraform nmap <leader>tg :!terraform graph<CR>
    ```

## Core Workflow Plugins

### telescope.nvim: Fuzzy Finder for Neovim
- **Description**: Highly extensible fuzzy finder over lists with powerful search capabilities
- **Features**:
  - Find files with fuzzy matching and preview
  - Live grep across project files
  - Buffer navigation and search
  - Git integration (commits, branches, status)
  - LSP integration (references, symbols, diagnostics)
  - Treesitter integration for code navigation
  - Customizable layouts and themes (dropdown, cursor, ivy)
  - Extensible architecture with many community extensions
  - Fully keyboard-driven interface (no mouse required)
  - Async operations for performance
  - Customizable keymaps and actions
  - 17.3k+ stars and active maintenance
  - Repository: https://github.com/nvim-telescope/telescope.nvim
  - Configuration example:
    ```lua
    -- telescope.nvim configuration
    require('telescope').setup {
      defaults = {
        mappings = {
          i = {
            ["<C-j>"] = "move_selection_next",
            ["<C-k>"] = "move_selection_previous",
            ["<C-q>"] = "send_to_qflist",
            ["<C-s>"] = "select_horizontal",
            ["<C-v>"] = "select_vertical",
          },
        },
        file_ignore_patterns = { "node_modules", ".git", "dist", "build" },
        path_display = { "truncate" },
        sorting_strategy = "ascending",
        layout_config = {
          horizontal = {
            prompt_position = "top",
            preview_width = 0.55,
          },
        },
      },
      pickers = {
        find_files = {
          hidden = true,
          find_command = { "rg", "--files", "--hidden", "--glob", "!.git/*" },
        },
        live_grep = {
          additional_args = function() return { "--hidden" } end,
        },
      },
      extensions = {
        fzf = {
          fuzzy = true,
          override_generic_sorter = true,
          override_file_sorter = true,
          case_mode = "smart_case",
        },
      },
    }
    
    -- Key mappings
    local builtin = require('telescope.builtin')
    vim.keymap.set('n', '<leader>ff', builtin.find_files, { desc = 'Find files' })
    vim.keymap.set('n', '<leader>fg', builtin.live_grep, { desc = 'Live grep' })
    vim.keymap.set('n', '<leader>fb', builtin.buffers, { desc = 'Buffers' })
    vim.keymap.set('n', '<leader>fh', builtin.help_tags, { desc = 'Help tags' })
    vim.keymap.set('n', '<leader>fr', builtin.lsp_references, { desc = 'LSP references' })
    vim.keymap.set('n', '<leader>fs', builtin.lsp_document_symbols, { desc = 'Document symbols' })
    vim.keymap.set('n', '<leader>fd', builtin.diagnostics, { desc = 'Diagnostics' })
    vim.keymap.set('n', '<leader>fc', builtin.git_commits, { desc = 'Git commits' })
    vim.keymap.set('n', '<leader>gb', builtin.git_branches, { desc = 'Git branches' })
    ```

## Core Workflow Plugins

### nvim-tree.lua: File Explorer for Neovim
- **Description**: A file explorer tree for Neovim written in Lua
- **Features**:
  - File system navigation with keyboard shortcuts
  - File operations (create, delete, rename, copy, paste)
  - Git integration with status indicators
  - Diagnostics integration with LSP and COC
  - Live filtering of files
  - Customizable icons and highlighting
  - Directory and file operations
  - Automatic updates when filesystem changes
  - Configurable keymaps and behavior
  - Fully keyboard-driven interface (no mouse required)
  - 7.7k+ stars and active maintenance
  - Repository: https://github.com/nvim-tree/nvim-tree.lua
  - Configuration example:
    ```lua
    -- nvim-tree.lua configuration
    require("nvim-tree").setup({
      sort = {
        sorter = "case_sensitive",
      },
      view = {
        width = 30,
        side = "left",
        adaptive_size = false,
      },
      renderer = {
        group_empty = true,
        highlight_git = true,
        highlight_opened_files = "name",
        icons = {
          show = {
            file = true,
            folder = true,
            folder_arrow = true,
            git = true,
          },
        },
      },
      filters = {
        dotfiles = false,
        custom = { "^.git$" },
      },
      git = {
        enable = true,
        ignore = false,
        timeout = 500,
      },
      actions = {
        open_file = {
          quit_on_open = false,
          resize_window = true,
          window_picker = {
            enable = true,
          },
        },
      },
      on_attach = function(bufnr)
        local api = require("nvim-tree.api")
        
        local function opts(desc)
          return { desc = "nvim-tree: " .. desc, buffer = bufnr, noremap = true, silent = true, nowait = true }
        end
        
        -- Default mappings
        api.config.mappings.default_on_attach(bufnr)
        
        -- Custom mappings
        vim.keymap.set('n', '<C-t>', api.tree.change_root_to_parent, opts('Up'))
        vim.keymap.set('n', '?', api.tree.toggle_help, opts('Help'))
        vim.keymap.set('n', 'l', api.node.open.edit, opts('Open'))
        vim.keymap.set('n', 'h', api.node.navigate.parent_close, opts('Close Directory'))
        vim.keymap.set('n', 'v', api.node.open.vertical, opts('Open: Vertical Split'))
      end,
    })
    
    -- Key mappings to toggle tree
    vim.keymap.set('n', '<C-n>', ':NvimTreeToggle<CR>', { noremap = true, silent = true })
    vim.keymap.set('n', '<leader>e', ':NvimTreeFocus<CR>', { noremap = true, silent = true })
    ```

## Kubernetes Development Support

### kubectl.nvim: Kubernetes Management for Neovim
- **Description**: Streamline Kubernetes management within Neovim
- **Features**:
  - Control and monitor Kubernetes clusters without leaving Neovim
  - View and manage Kubernetes resources (pods, deployments, services)
  - Execute kubectl commands directly from Neovim
  - Real-time resource monitoring
  - Context and namespace switching
  - YAML validation and formatting
  - Resource logs viewing
  - Port forwarding capabilities
  - Interactive resource editing
  - 413+ stars and active maintenance
  - Repository: https://github.com/Ramilito/kubectl.nvim
  - Configuration example:
    ```lua
    -- kubectl.nvim configuration
    require("kubectl").setup({
      -- Load kubectl binary from PATH
      kubectl_path = "kubectl",
      
      -- Default namespace to use
      default_namespace = "default",
      
      -- Timeouts for requests
      timeout = 5000,
      
      -- Logs configuration
      logs = {
        -- Follow logs by default
        follow = true,
        -- Number of lines to show
        lines = 100,
      },
      
      -- Mappings for kubectl commands
      mappings = {
        -- Resource navigation
        next_item = "j",
        prev_item = "k",
        -- Actions
        describe = "d",
        logs = "l",
        delete = "dd",
        edit = "e",
        refresh = "r",
      },
      
      -- Appearance settings
      ui = {
        -- Border style for floating windows
        border = "rounded",
        -- Width of resource list
        width = 40,
        -- Height of logs window
        logs_height = 15,
      },
    })
    
    -- Key mappings
    vim.keymap.set('n', '<leader>kc', '<cmd>KubectlContexts<CR>', { desc = 'Kubectl contexts' })
    vim.keymap.set('n', '<leader>kn', '<cmd>KubectlNamespaces<CR>', { desc = 'Kubectl namespaces' })
    vim.keymap.set('n', '<leader>kp', '<cmd>KubectlPods<CR>', { desc = 'Kubectl pods' })
    vim.keymap.set('n', '<leader>kd', '<cmd>KubectlDeployments<CR>', { desc = 'Kubectl deployments' })
    vim.keymap.set('n', '<leader>ks', '<cmd>KubectlServices<CR>', { desc = 'Kubectl services' })
    ```

### vimkubectl: Vim/Neovim Kubernetes Integration
- **Description**: Manage Kubernetes resources from Vim/Neovim
- **Features**:
  - View and edit Kubernetes resources
  - Support for multiple contexts and namespaces
  - Resource filtering and searching
  - YAML validation and formatting
  - Integration with kubectl command-line tool
  - Resource logs viewing
  - Port forwarding capabilities
  - 80+ stars and active maintenance
  - Repository: https://github.com/rottencandy/vimkubectl
  - Configuration example:
    ```vim
    " vimkubectl configuration
    " Set default namespace
    let g:vimkubectl_default_ns = 'default'
    
    " Set kubectl binary path
    let g:vimkubectl_kubectl_path = 'kubectl'
    
    " Set command output buffer name
    let g:vimkubectl_outbuf = 'KubectlOutput'
    
    " Key mappings
    nnoremap <leader>kk :Kubectl<space>
    nnoremap <leader>kg :KubectlGet<space>
    nnoremap <leader>kd :KubectlDescribe<space>
    nnoremap <leader>kl :KubectlLogs<space>
    nnoremap <leader>ke :KubectlEdit<space>
    nnoremap <leader>kc :KubectlConfig<space>
    ```

### kube-utils-nvim: Kubernetes and Helm Integration for Neovim
- **Description**: Neovim plugin providing seamless integration with Kubernetes and Helm
- **Features**:
  - Deploy and manage Kubernetes resources
  - Helm chart management
  - Context and namespace switching
  - Resource logs viewing
  - Port forwarding capabilities
  - YAML validation and formatting
  - Interactive resource editing
  - 41+ stars and active maintenance
  - Repository: https://github.com/h4ckm1n-dev/kube-utils-nvim
  - Configuration example:
    ```lua
    -- kube-utils-nvim configuration
    require('kube-utils').setup({
      -- Kubectl configuration
      kubectl = {
        -- Path to kubectl binary
        path = "kubectl",
        -- Default namespace
        namespace = "default",
        -- Default context
        context = nil,
      },
      
      -- Helm configuration
      helm = {
        -- Path to helm binary
        path = "helm",
        -- Default namespace
        namespace = "default",
      },
      
      -- UI configuration
      ui = {
        -- Border style for floating windows
        border = "rounded",
        -- Width of resource list
        width = 40,
        -- Height of logs window
        height = 15,
      },
    })
    
    -- Key mappings
    vim.keymap.set('n', '<leader>kc', '<cmd>KubeContext<CR>', { desc = 'Kubernetes contexts' })
    vim.keymap.set('n', '<leader>kn', '<cmd>KubeNamespace<CR>', { desc = 'Kubernetes namespaces' })
    vim.keymap.set('n', '<leader>kr', '<cmd>KubeResources<CR>', { desc = 'Kubernetes resources' })
    vim.keymap.set('n', '<leader>hc', '<cmd>HelmCharts<CR>', { desc = 'Helm charts' })
    vim.keymap.set('n', '<leader>hr', '<cmd>HelmReleases<CR>', { desc = 'Helm releases' })
    ```

## Package Management

### mason.nvim: Portable Package Manager for Neovim
- **Description**: Comprehensive package manager for Neovim that runs everywhere Neovim runs
- **Features**:
  - Install and manage LSP servers for Python, Go, TypeScript, and more
  - Install and manage DAP servers for debugging
  - Install and manage linters for code quality
  - Install and manage formatters for consistent code style
  - Portable across all platforms where Neovim runs
  - User-friendly UI for package management
  - Extensible architecture with registry support
  - Integration with lspconfig via mason-lspconfig.nvim
  - Automatic PATH management for installed tools
  - 8.7k+ stars and active maintenance
  - Repository: https://github.com/williamboman/mason.nvim
  - Configuration example:
    ```lua
    -- mason.nvim configuration
    require("mason").setup({
      ui = {
        icons = {
          package_installed = "✓",
          package_pending = "➜",
          package_uninstalled = "✗"
        },
        border = "rounded",
        width = 0.8,
        height = 0.8
      },
      log_level = vim.log.levels.INFO,
      max_concurrent_installers = 4,
      registries = {
        "github:mason-org/mason-registry"
      },
      providers = {
        "mason.providers.registry-api",
        "mason.providers.client"
      },
      github = {
        download_url_template = "https://github.com/%s/releases/download/%s/%s"
      },
      pip = {
        upgrade_pip = false,
        install_args = {}
      }
    })
    
    -- Integration with lspconfig (requires mason-lspconfig.nvim)
    require("mason-lspconfig").setup({
      -- Automatically install LSP servers when they're needed
      automatic_installation = true,
      -- Ensure these servers are always installed
      ensure_installed = {
        "pyright",           -- Python
        "gopls",             -- Go
        "tsserver",          -- TypeScript
        "terraformls",       -- Terraform
        "yamlls",            -- YAML (for Kubernetes)
      }
    })
    
    -- Key mappings
    vim.keymap.set('n', '<leader>pm', '<cmd>Mason<CR>', { desc = 'Open Mason package manager' })
    vim.keymap.set('n', '<leader>pi', '<cmd>MasonInstall<CR>', { desc = 'Install package with Mason' })
    vim.keymap.set('n', '<leader>pu', '<cmd>MasonUpdate<CR>', { desc = 'Update Mason packages' })
    vim.keymap.set('n', '<leader>pl', '<cmd>MasonLog<CR>', { desc = 'View Mason logs' })
    ```

## Language-Specific Development Support

### vim-go: Go Development Plugin for Vim/Neovim
- **Description**: Comprehensive Go development plugin for Vim/Neovim
- **Features**:
  - Compile, install, and test Go packages directly from Neovim
  - Debug programs with integrated Delve support
  - Code completion via gopls
  - Go to symbol/declaration functionality
  - Documentation lookup with GoDoc
  - Easy package importing and removal
  - Type-safe renaming of identifiers
  - Code coverage visualization
  - Struct tag management
  - Linting and static error checking
  - Advanced source analysis (implements, callees, referrers)
  - Syntax highlighting and folding
  - Integration with Tagbar via gotags
  - Snippet engine integration
  - 16.1k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " vim-go configuration
    " Enable syntax highlighting
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    
    " Auto formatting and importing
    let g:go_fmt_autosave = 1
    let g:go_fmt_command = "goimports"
    
    " Status line types/signatures
    let g:go_auto_type_info = 1
    
    " Run :GoBuild or :GoTestCompile based on the go file
    function! s:build_go_files()
      let l:file = expand('%')
      if l:file =~# '^\f\+_test\.go$'
        call go#test#Test(0, 1)
      elseif l:file =~# '^\f\+\.go$'
        call go#cmd#Build(0)
      endif
    endfunction
    
    " Key mappings
    autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <Leader>i <Plug>(go-info)
    autocmd FileType go nmap <Leader>d <Plug>(go-def)
    autocmd FileType go nmap <Leader>v <Plug>(go-def-vertical)
    autocmd FileType go nmap <Leader>s <Plug>(go-implements)
    autocmd FileType go nmap <Leader>e <Plug>(go-rename)
    autocmd FileType go nmap <Leader>a <Plug>(go-alternate-edit)
    
    " Use gopls for completion
    let g:go_gopls_enabled = 1
    let g:go_gopls_options = ['-remote=auto']
    
    " Debugging settings
    let g:go_debug_windows = {
          \ 'vars':  'rightbelow 60vnew',
          \ 'stack': 'rightbelow 10new',
          \ 'goroutines': 'botright 10new',
          \ 'out':   'botright 5new',
      \ }
    let g:go_debug_address = '127.0.0.1:8181'
    let g:go_debug_log_output = 'debugger,rpc'
    ```

### jedi-vim: Python Autocompletion for Vim/Neovim
- **Description**: Python autocompletion library for Vim/Neovim
- **Features**:
  - Intelligent code completion for Python
  - Go to definition and assignment functionality
  - Documentation lookup with popup display
  - Usages search to find all references
  - Renaming of variables, functions, and classes
  - Support for virtual environments
  - Type inference and signature help
  - Module importing with :Pyimport command
  - Configurable key mappings
  - Integration with popular completion engines
  - 5.3k+ stars and active maintenance
  - Repository: https://github.com/davidhalter/jedi-vim
  - Configuration example:
    ```vim
    " jedi-vim configuration
    " Enable completions
    let g:jedi#completions_enabled = 1
    
    " Auto-completion behavior
    let g:jedi#popup_on_dot = 1
    let g:jedi#popup_select_first = 0
    let g:jedi#show_call_signatures = "1"
    let g:jedi#case_insensitive_completion = 1
    
    " Key mappings
    let g:jedi#goto_command = "<leader>d"
    let g:jedi#goto_assignments_command = "<leader>g"
    let g:jedi#goto_stubs_command = "<leader>s"
    let g:jedi#documentation_command = "K"
    let g:jedi#usages_command = "<leader>n"
    let g:jedi#completions_command = "<C-Space>"
    let g:jedi#rename_command = "<leader>r"
    let g:jedi#rename_command_keep_name = "<leader>R"
    
    " Display options
    let g:jedi#use_splits_not_buffers = "right"
    
    " Environment settings
    let g:jedi#environment_path = "venv"  " Point to virtual environment
    
    " Integration with other plugins
    let g:jedi#smart_auto_mappings = 0  " Disable auto-import mappings
    
    " Use with deoplete (if installed)
    if has('nvim') && exists('*deoplete#custom#source')
      call deoplete#custom#source('jedi', 'rank', 500)
    endif
    ```

## Intelligent Code Completion

### coc.nvim: Intelligent Code Completion for Vim/Neovim
- **Description**: Nodejs extension host for Vim/Neovim that loads extensions like VSCode and hosts language servers
- **Features**:
  - Intelligent code completion for multiple languages (Python, Go, TypeScript, etc.)
  - Language server protocol (LSP) support
  - Diagnostics (errors, warnings)
  - Go to definition, references, implementation
  - Symbol search and workspace symbol search
  - Document formatting and range formatting
  - Code actions and quick fixes
  - Code lens support
  - Snippet support
  - Extensions ecosystem similar to VSCode
  - Hover documentation
  - Signature help
  - Rename symbol
  - Highlight references
  - 24.8k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc.nvim
  - Configuration example:
    ```vim
    " coc.nvim configuration
    " Install extensions
    let g:coc_global_extensions = [
      \ 'coc-tsserver',           " TypeScript support
      \ 'coc-pyright',            " Python support
      \ 'coc-go',                 " Go support
      \ 'coc-yaml',               " YAML support (for Kubernetes)
      \ 'coc-json',               " JSON support
      \ 'coc-terraform',          " Terraform support
      \ 'coc-snippets',           " Snippets support
      \ 'coc-pairs',              " Auto pairs
      \ 'coc-explorer',           " File explorer
      \ 'coc-prettier',           " Code formatting
      \ 'coc-diagnostic'          " Enhanced diagnostics
      \ ]
    
    " TextEdit might fail if hidden is not set
    set hidden
    
    " Some servers have issues with backup files
    set nobackup
    set nowritebackup
    
    " Give more space for displaying messages
    set cmdheight=2
    
    " Having longer updatetime (default is 4000 ms) leads to noticeable delays
    set updatetime=300
    
    " Don't pass messages to |ins-completion-menu|
    set shortmess+=c
    
    " Always show the signcolumn
    set signcolumn=yes
    
    " Use tab for trigger completion with characters ahead and navigate
    inoremap <silent><expr> <TAB>
          \ coc#pum#visible() ? coc#pum#next(1) :
          \ CheckBackspace() ? "\<Tab>" :
          \ coc#refresh()
    inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"
    
    " Make <CR> auto-select the first completion item
    inoremap <silent><expr> <CR> coc#pum#visible() ? coc#pum#confirm()
                                  \: "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"
    
    function! CheckBackspace() abort
      let col = col('.') - 1
      return !col || getline('.')[col - 1]  =~# '\s'
    endfunction
    
    " Use <c-space> to trigger completion
    inoremap <silent><expr> <c-space> coc#refresh()
    
    " Use `[g` and `]g` to navigate diagnostics
    nmap <silent> [g <Plug>(coc-diagnostic-prev)
    nmap <silent> ]g <Plug>(coc-diagnostic-next)
    
    " GoTo code navigation
    nmap <silent> gd <Plug>(coc-definition)
    nmap <silent> gy <Plug>(coc-type-definition)
    nmap <silent> gi <Plug>(coc-implementation)
    nmap <silent> gr <Plug>(coc-references)
    
    " Use K to show documentation in preview window
    nnoremap <silent> K :call ShowDocumentation()<CR>
    
    function! ShowDocumentation()
      if CocAction('hasProvider', 'hover')
        call CocActionAsync('doHover')
      else
        call feedkeys('K', 'in')
      endif
    endfunction
    
    " Highlight the symbol and its references when holding the cursor
    autocmd CursorHold * silent call CocActionAsync('highlight')
    
    " Symbol renaming
    nmap <leader>rn <Plug>(coc-rename)
    
    " Formatting selected code
    xmap <leader>f  <Plug>(coc-format-selected)
    nmap <leader>f  <Plug>(coc-format-selected)
    
    " Apply code actions
    xmap <leader>a  <Plug>(coc-codeaction-selected)
    nmap <leader>a  <Plug>(coc-codeaction-selected)
    
    " Apply code action to the current buffer
    nmap <leader>ac  <Plug>(coc-codeaction)
    
    " Apply AutoFix to problem on the current line
    nmap <leader>qf  <Plug>(coc-fix-current)
    
    " Run the Code Lens action on the current line
    nmap <leader>cl  <Plug>(coc-codelens-action)
    ```

## Visual Themes and UI

### tokyonight.nvim: Clean Dark Theme for Neovim
- **Description**: A clean, dark Neovim theme written in Lua with support for LSP, Treesitter, and many plugins
- **Features**:
  - Multiple style variants (Night, Storm, Day, Moon)
  - Terminal colors for consistent experience
  - Support for 50+ popular plugins
  - Extras for terminal emulators (Kitty, Alacritty, iTerm, etc.)
  - Fish shell integration
  - Tmux theme integration
  - Customizable colors and highlight groups
  - Transparent background option
  - Dim inactive windows option
  - Configurable styles for syntax elements
  - 7k+ stars and active maintenance
  - Repository: https://github.com/folke/tokyonight.nvim
  - Configuration example:
    ```lua
    -- tokyonight.nvim configuration
    require("tokyonight").setup({
      -- Choose style: "storm", "moon", "night", "day"
      style = "moon",
      
      -- Enable transparent background
      transparent = false,
      
      -- Configure terminal colors
      terminal_colors = true,
      
      -- Styles for syntax elements
      styles = {
        comments = { italic = true },
        keywords = { italic = true },
        functions = {},
        variables = {},
        
        -- Background styles for sidebars and floating windows
        sidebars = "dark",
        floats = "dark",
      },
      
      -- Adjust brightness of day style (0-1)
      day_brightness = 0.3,
      
      -- Dim inactive windows
      dim_inactive = false,
      
      -- Bold section headers in lualine
      lualine_bold = false,
      
      -- Override specific colors
      on_colors = function(colors)
        -- Example: make errors more visible
        colors.error = "#ff0000"
      end,
      
      -- Override specific highlight groups
      on_highlights = function(highlights, colors)
        -- Example: customize telescope appearance
        highlights.TelescopeNormal = {
          bg = colors.bg_dark,
          fg = colors.fg_dark,
        }
      end,
      
      -- Enable specific plugin integrations
      plugins = {
        telescope = true,
        nvim_tree = true,
        which_key = true,
        gitsigns = true,
        nvim_cmp = true,
        treesitter = true,
        indent_blankline = true,
      },
    })
    
    -- Apply the colorscheme
    vim.cmd[[colorscheme tokyonight]]
    ```

## Document Editing and Preview

### markdown-preview.nvim: Markdown Preview for Vim/Neovim
- **Description**: Real-time Markdown preview plugin for Vim/Neovim
- **Features**:
  - Synchronized scrolling between editor and preview
  - Fast asynchronous updates
  - Cross-platform support (MacOS/Linux/Windows)
  - KaTeX for math typesetting
  - PlantUML, Mermaid, and js-sequence-diagrams support
  - Chart.js integration
  - Flowchart and dot diagram support
  - Table of contents generation
  - Emoji support
  - Task list support
  - Local image display
  - Customizable themes (light/dark)
  - Flexible configuration options
  - 7.2k+ stars and active maintenance
  - Repository: https://github.com/iamcco/markdown-preview.nvim
  - Configuration example:
    ```vim
    " markdown-preview.nvim configuration
    " Auto start preview when entering markdown buffer
    let g:mkdp_auto_start = 0
    
    " Auto close preview when leaving markdown buffer
    let g:mkdp_auto_close = 1
    
    " Refresh mode (0 = real-time, 1 = save/leave insert mode only)
    let g:mkdp_refresh_slow = 0
    
    " Command availability (0 = markdown only, 1 = all files)
    let g:mkdp_command_for_global = 0
    
    " Open to network (0 = localhost only, 1 = available to network)
    let g:mkdp_open_to_the_world = 0
    
    " Custom browser setting
    let g:mkdp_browser = ''
    
    " Echo preview URL in command line
    let g:mkdp_echo_preview_url = 0
    
    " Custom function to open preview
    let g:mkdp_browserfunc = ''
    
    " Preview options
    let g:mkdp_preview_options = {
        \ 'mkit': {},
        \ 'katex': {},
        \ 'uml': {},
        \ 'maid': {},
        \ 'disable_sync_scroll': 0,
        \ 'sync_scroll_type': 'middle',
        \ 'hide_yaml_meta': 1,
        \ 'sequence_diagrams': {},
        \ 'flowchart_diagrams': {},
        \ 'content_editable': v:false,
        \ 'disable_filename': 0,
        \ 'toc': {}
        \ }
    
    " Custom CSS for markdown content
    let g:mkdp_markdown_css = ''
    
    " Custom CSS for syntax highlighting
    let g:mkdp_highlight_css = ''
    
    " Custom port for server
    let g:mkdp_port = ''
    
    " Preview page title format
    let g:mkdp_page_title = '「${name}」'
    
    " Theme setting (dark or light)
    let g:mkdp_theme = 'dark'
    
    " Key mappings
    nmap <C-p> <Plug>MarkdownPreviewToggle
    ```

## Terminal Management

### toggleterm.nvim: Terminal Management for Neovim
- **Description**: A Neovim Lua plugin to help easily manage multiple terminal windows
- **Features**:
  - Toggle multiple terminals with different orientations (float, vertical, horizontal, tab)
  - Persist terminal state between toggles
  - Send commands to specific terminals
  - Custom terminal creation for specific tools (lazygit, htop, etc.)
  - Terminal window mappings for easy navigation
  - Configurable terminal appearance with shading
  - Winbar support for terminal identification
  - Terminal selection UI
  - Send code/lines from buffer to terminal
  - Responsive layout based on window size
  - 4.8k+ stars and active maintenance
  - Repository: https://github.com/akinsho/toggleterm.nvim
  - Configuration example:
    ```lua
    -- toggleterm.nvim configuration
    require("toggleterm").setup({
      -- Size can be a number or function
      size = function(term)
        if term.direction == "horizontal" then
          return 15
        elseif term.direction == "vertical" then
          return vim.o.columns * 0.4
        end
      end,
      
      -- Key mapping to toggle terminal
      open_mapping = [[<c-\>]],
      
      -- Terminal appearance
      hide_numbers = true,
      shade_terminals = true,
      shading_factor = 2,
      start_in_insert = true,
      insert_mappings = true,
      terminal_mappings = true,
      
      -- Terminal behavior
      persist_size = true,
      persist_mode = true,
      direction = 'float',
      close_on_exit = true,
      shell = vim.o.shell,
      auto_scroll = true,
      
      -- Float window options
      float_opts = {
        border = 'curved',
        winblend = 0,
        highlights = {
          border = "Normal",
          background = "Normal",
        },
      },
      
      -- Winbar configuration
      winbar = {
        enabled = true,
        name_formatter = function(term)
          return term.name or term.id
        end,
      },
      
      -- Responsive layout
      responsiveness = {
        horizontal_breakpoint = 135,
      },
    })
    
    -- Terminal window mappings
    function _G.set_terminal_keymaps()
      local opts = {buffer = 0}
      vim.keymap.set('t', '<esc>', [[<C-\><C-n>]], opts)
      vim.keymap.set('t', '<C-h>', [[<Cmd>wincmd h<CR>]], opts)
      vim.keymap.set('t', '<C-j>', [[<Cmd>wincmd j<CR>]], opts)
      vim.keymap.set('t', '<C-k>', [[<Cmd>wincmd k<CR>]], opts)
      vim.keymap.set('t', '<C-l>', [[<Cmd>wincmd l<CR>]], opts)
    end
    
    -- Auto-apply terminal keymaps when terminal opens
    vim.cmd('autocmd! TermOpen term://* lua set_terminal_keymaps()')
    
    -- Custom terminal for lazygit
    local Terminal = require('toggleterm.terminal').Terminal
    local lazygit = Terminal:new({
      cmd = "lazygit",
      dir = "git_dir",
      direction = "float",
      float_opts = {
        border = "double",
      },
      on_open = function(term)
        vim.cmd("startinsert!")
      end,
    })
    
    -- Function to toggle lazygit terminal
    function _lazygit_toggle()
      lazygit:toggle()
    end
    
    -- Key mapping for lazygit
    vim.keymap.set("n", "<leader>g", "<cmd>lua _lazygit_toggle()<CR>", {noremap = true, silent = true})
    ```

## Navigation and Motion

### leap.nvim: Keyboard-Based Navigation for Neovim
- **Description**: A general-purpose motion plugin for Neovim that provides efficient keyboard-based navigation
- **Features**:
  - Fast navigation to any position in the visible editor area
  - 2-character search pattern with label preview
  - Minimizes required focus level while executing jumps
  - Eliminates the need for mouse usage
  - Uniform way to reach any target on screen
  - Cross-window navigation support
  - Customizable label characters
  - Equivalence classes for character groups
  - Preview filtering for reduced visual noise
  - Treesitter node selection integration
  - Remote operations for "actions at a distance"
  - Dot-repeat support via vim-repeat
  - 4.7k+ stars and active maintenance
  - Repository: https://github.com/ggandor/leap.nvim
  - Configuration example:
    ```lua
    -- leap.nvim configuration
    -- Set up default mappings (s, S)
    require('leap').set_default_mappings()
    
    -- Or set custom mappings
    vim.keymap.set({'n', 'x', 'o'}, 's', '<Plug>(leap-forward)')
    vim.keymap.set({'n', 'x', 'o'}, 'S', '<Plug>(leap-backward)')
    vim.keymap.set({'n', 'x', 'o'}, 'gs', '<Plug>(leap-from-window)')
    
    -- Skip the middle of alphabetic words to reduce visual noise
    require('leap').opts.preview_filter = function (ch0, ch1, ch2)
      return not (
        ch1:match('%s') or
        ch0:match('%a') and ch1:match('%a') and ch2:match('%a')
      )
    end
    
    -- Define equivalence classes for brackets and quotes
    require('leap').opts.equivalence_classes = { 
      ' \t\r\n',  -- whitespace
      '([{',      -- opening brackets
      ')]}',      -- closing brackets
      '\'"`'      -- quotes
    }
    
    -- Use Enter/Backspace to repeat the previous motion
    require('leap.user').set_repeat_keys('<enter>', '<backspace>')
    
    -- Hide the cursor when leaping (for Neovim < 0.10)
    vim.api.nvim_create_autocmd('User', { pattern = 'LeapEnter',
      callback = function()
        vim.cmd.hi('Cursor', 'blend=100')
        vim.opt.guicursor:append { 'a:Cursor/lCursor' }
      end,
    })
    vim.api.nvim_create_autocmd('User', { pattern = 'LeapLeave',
      callback = function()
        vim.cmd.hi('Cursor', 'blend=0')
        vim.opt.guicursor:remove { 'a:Cursor/lCursor' }
      end,
    })
    ```

## Code Editing and Manipulation

### Comment.nvim: Smart Code Commenting for Neovim
- **Description**: Smart and powerful commenting plugin for Neovim with treesitter support
- **Features**:
  - Supports line (`//`) and block (`/* */`) comments
  - Treesitter integration for accurate commenting
  - Dot (`.`) repeat support for comment operations
  - Count support for batch commenting
  - Left-right and up-down motion support
  - Text-object integration for precise commenting
  - Pre and post hooks for custom behavior
  - Line ignoring with Lua regex patterns
  - Multiple language support
  - Operator-pending mappings
  - Extra mappings for comment insertion
  - 4.2k+ stars and active maintenance
  - Repository: https://github.com/numToStr/Comment.nvim
  - Configuration example:
    ```lua
    -- Comment.nvim configuration
    require('Comment').setup({
      -- Add a space between comment and line
      padding = true,
      
      -- Whether cursor should stay at its position
      sticky = true,
      
      -- Lines to be ignored while (un)comment (lua regex)
      ignore = '^$', -- ignore empty lines
      
      -- LHS of toggle mappings in NORMAL mode
      toggler = {
        -- Line-comment toggle keymap
        line = 'gcc',
        -- Block-comment toggle keymap
        block = 'gbc',
      },
      
      -- LHS of operator-pending mappings in NORMAL and VISUAL mode
      opleader = {
        -- Line-comment keymap
        line = 'gc',
        -- Block-comment keymap
        block = 'gb',
      },
      
      -- LHS of extra mappings
      extra = {
        -- Add comment on the line above
        above = 'gcO',
        -- Add comment on the line below
        below = 'gco',
        -- Add comment at the end of line
        eol = 'gcA',
      },
      
      -- Enable keybindings
      mappings = {
        -- Operator-pending mapping
        basic = true,
        -- Extra mapping
        extra = true,
      },
      
      -- Function to call before (un)comment
      pre_hook = function(ctx)
        -- For tsx/jsx files, use context-aware commenting
        local U = require('Comment.utils')
        local type = ctx.ctype == U.ctype.linewise and '__default' or '__multiline'
        
        local location = nil
        if ctx.ctype == U.ctype.blockwise then
          location = require('ts_context_commentstring.utils').get_cursor_location()
        elseif ctx.cmotion == U.cmotion.v or ctx.cmotion == U.cmotion.V then
          location = require('ts_context_commentstring.utils').get_visual_start_location()
        end
        
        return require('ts_context_commentstring.internal').calculate_commentstring({
          key = type,
          location = location,
        })
      end,
      
      -- Function to call after (un)comment
      post_hook = nil,
    })
    ```

## User Interface and Feedback

### nvim-notify: Fancy Notification Manager for Neovim
- **Description**: A fancy, configurable, notification manager for Neovim
- **Features**:
  - Animated notifications with multiple animation styles
  - Configurable rendering styles (default, minimal, simple, compact)
  - Treesitter syntax highlighting in notifications
  - Notification history with search capabilities
  - Telescope integration for browsing notification history
  - Customizable notification appearance and behavior
  - Async notification support
  - Notification updating and chaining
  - Custom filetype support for notifications
  - Notification timeout configuration
  - Event hooks for notification lifecycle
  - 3.3k+ stars and active maintenance
  - Repository: https://github.com/rcarriga/nvim-notify
  - Configuration example:
    ```lua
    -- nvim-notify configuration
    require("notify").setup({
      -- Animation style (fade, slide, fade_in_slide_out, static)
      stages = "fade",
      
      -- Default timeout for notifications
      timeout = 5000,
      
      -- For stages that change opacity, this is treated as the highlight behind the window
      background_colour = "Normal",
      
      -- Icons for the different levels
      icons = {
        ERROR = "",
        WARN = "",
        INFO = "",
        DEBUG = "",
        TRACE = "✎",
      },
      
      -- Render style (default, minimal, simple, compact)
      render = "default",
      
      -- Minimum width for notification windows
      minimum_width = 50,
      
      -- Max number of notifications to show at once
      max_width = nil,
      max_height = nil,
      
      -- Function to get dimensions
      on_open = nil,
      
      -- Function called when a notification is closed
      on_close = nil,
      
      -- Level configuration
      level = 2,
      
      -- Timeout for progress bars (ms)
      fps = 30,
      top_down = true,
    })
    
    -- Set as default notification function
    vim.notify = require("notify")
    
    -- Example usage
    vim.notify("This is an error message", "error", {
      title = "My Plugin",
      timeout = 3000,
      on_open = function()
        -- Do something when notification opens
      end,
      on_close = function()
        -- Do something when notification closes
      end,
    })
    
    -- Load telescope extension
    require("telescope").load_extension("notify")
    
    -- Key mapping for notification history
    vim.keymap.set("n", "<leader>fn", function()
      require("telescope").extensions.notify.notify()
    end, { desc = "View notification history" })
    ```

## Fuzzy Finding and Navigation

### fzf-lua: Fuzzy Finder for Neovim
- **Description**: Improved fzf.vim implementation written in Lua for Neovim
- **Features**:
  - Extensive fuzzy finding capabilities across files, buffers, and git
  - Lua-based configuration for deep customization
  - Multiple UI layouts and themes (floating windows, splits, fullscreen)
  - LSP integration for references, definitions, and code actions
  - Git integration for files, commits, branches, and blame
  - Telescope-like interface option with familiar keybindings
  - Builtin and native previewer options
  - Treesitter integration for syntax highlighting
  - Customizable keymaps and actions
  - Profile system for quick configuration switching
  - Insert-mode completion for paths and code
  - 3.2k+ stars and active maintenance
  - Repository: https://github.com/ibhagwan/fzf-lua
  - Configuration example:
    ```lua
    -- fzf-lua configuration
    require('fzf-lua').setup({
      -- Use a specific fzf binary
      -- fzf_bin = 'sk',
      
      -- Window appearance
      winopts = {
        height = 0.85,
        width = 0.80,
        row = 0.35,
        col = 0.50,
        border = "rounded",
        preview = {
          border = "rounded",
          wrap = false,
          hidden = false,
          vertical = "down:45%",
          horizontal = "right:60%",
          layout = "flex",
          title = true,
          title_pos = "center",
        },
      },
      
      -- Keymap configuration
      keymap = {
        -- Mappings applied to fzf window
        fzf = {
          ["ctrl-z"] = "abort",
          ["ctrl-f"] = "half-page-down",
          ["ctrl-b"] = "half-page-up",
          ["ctrl-a"] = "beginning-of-line",
          ["ctrl-e"] = "end-of-line",
          ["alt-a"] = "toggle-all",
        },
      },
      
      -- Actions for different file types
      actions = {
        files = {
          ["default"] = require('fzf-lua.actions').file_edit,
          ["ctrl-s"] = require('fzf-lua.actions').file_split,
          ["ctrl-v"] = require('fzf-lua.actions').file_vsplit,
          ["ctrl-t"] = require('fzf-lua.actions').file_tabedit,
        },
      },
      
      -- File finding configuration
      files = {
        prompt = 'Files❯ ',
        cmd = "fd --type f --hidden --follow --exclude .git",
        git_icons = true,
        file_icons = true,
        color_icons = true,
      },
      
      -- Git integration
      git = {
        files = {
          prompt = 'GitFiles❯ ',
          cmd = 'git ls-files --exclude-standard',
          git_icons = true,
          file_icons = true,
          color_icons = true,
        },
      },
      
      -- LSP integration
      lsp = {
        prompt_postfix = '❯ ',
        cwd_only = false,
        async_or_timeout = 5000,
        file_icons = true,
        git_icons = false,
      },
    })
    
    -- Key mappings for common operations
    vim.keymap.set('n', '<leader>ff', "<cmd>lua require('fzf-lua').files()<CR>", { silent = true })
    vim.keymap.set('n', '<leader>fg', "<cmd>lua require('fzf-lua').live_grep()<CR>", { silent = true })
    vim.keymap.set('n', '<leader>fb', "<cmd>lua require('fzf-lua').buffers()<CR>", { silent = true })
    vim.keymap.set('n', '<leader>fh', "<cmd>lua require('fzf-lua').help_tags()<CR>", { silent = true })
    ```

## Focus and Productivity

### zen-mode.nvim: Distraction-Free Coding for Neovim
- **Description**: A plugin that provides a distraction-free coding environment in Neovim
- **Features**:
  - Opens current buffer in a full-screen floating window
  - Preserves existing window layouts and splits
  - Compatible with other floating windows (LSP hover, WhichKey)
  - Dynamically adjustable window size
  - Automatic realignment when editor or window is resized
  - Backdrop shading for enhanced focus
  - Hides status line and optionally other UI elements
  - Customizable with Lua callbacks for open/close events
  - Integration with external tools (tmux, Kitty, Alacritty, wezterm, Neovide)
  - Automatic closing when new non-floating window is opened
  - Compatible with Telescope for opening new buffers
  - 1.9k+ stars and active maintenance
  - Repository: https://github.com/folke/zen-mode.nvim
  - Configuration example:
    ```lua
    -- zen-mode.nvim configuration
    require("zen-mode").setup({
      window = {
        backdrop = 0.95, -- shade the backdrop (0-1)
        width = 120,     -- width of the Zen window
        height = 1,      -- height of the Zen window (1 = full height)
        options = {
          signcolumn = "no",       -- disable signcolumn
          number = false,          -- disable number column
          relativenumber = false,  -- disable relative numbers
          cursorline = false,      -- disable cursorline
          cursorcolumn = false,    -- disable cursor column
          foldcolumn = "0",        -- disable fold column
          list = false,            -- disable whitespace characters
        },
      },
      plugins = {
        -- disable global vim options
        options = {
          enabled = true,
          ruler = false,           -- disable ruler
          showcmd = false,         -- disable command display
          laststatus = 0,          -- disable status line
        },
        twilight = { enabled = true },  -- dim inactive code
        gitsigns = { enabled = false }, -- disable git signs
        tmux = { enabled = false },     -- disable tmux statusline
        -- kitty, alacritty, wezterm integration for font size
        kitty = {
          enabled = false,
          font = "+4", -- font size increment
        },
        wezterm = {
          enabled = false,
          font = "+4", -- font size increment
        },
        -- neovide integration
        neovide = {
          enabled = false,
          scale = 1.2, -- scale factor
          -- disable animations
          disable_animations = {
            neovide_animation_length = 0,
            neovide_cursor_animate_command_line = false,
            neovide_scroll_animation_length = 0,
            neovide_cursor_animation_length = 0,
            neovide_cursor_vfx_mode = "",
          }
        },
      },
      -- callbacks for custom behavior
      on_open = function(win)
        -- code to run when Zen window opens
      end,
      on_close = function()
        -- code to run when Zen window closes
      end,
    })
    
    -- Toggle Zen Mode with a keymap
    vim.keymap.set("n", "<leader>z", "<cmd>ZenMode<CR>", { silent = true })
    
    -- Toggle with custom options
    vim.keymap.set("n", "<leader>Z", function()
      require("zen-mode").toggle({
        window = {
          width = 0.85, -- 85% of editor width
        },
        plugins = {
          twilight = { enabled = true },
        }
      })
    end, { silent = true })
    ```

## Development Tools and Practice

### leetcode.nvim: LeetCode Problem Solving in Neovim
- **Description**: A Neovim plugin enabling you to solve LeetCode problems directly within the editor
- **Features**:
  - Intuitive dashboard for problem navigation
  - Question description formatting for better readability
  - LeetCode profile statistics within Neovim
  - Support for daily and random questions
  - Caching for optimized performance
  - Multiple language support (C++, Java, Python, TypeScript, Go, etc.)
  - Telescope or fzf-lua integration for problem selection
  - Console with test case management
  - Customizable layout and appearance
  - Session management for multiple accounts
  - Hooks for custom behavior
  - 1.5k+ stars and active maintenance
  - Repository: https://github.com/kawre/leetcode.nvim
  - Configuration example:
    ```lua
    -- leetcode.nvim configuration
    require('leetcode').setup({
      -- Language to use for LeetCode problems
      lang = "python3",
      
      -- Storage directories
      storage = {
        home = vim.fn.stdpath("data") .. "/leetcode",
        cache = vim.fn.stdpath("cache") .. "/leetcode",
      },
      
      -- Inject code before or after your solution
      injector = {
        ["python3"] = {
          before = true, -- Use default imports
        },
        ["cpp"] = {
          before = { "#include <bits/stdc++.h>", "using namespace std;" },
          after = "int main() {}",
        },
      },
      
      -- Picker provider (telescope or fzf-lua)
      picker = { provider = "telescope" },
      
      -- Console configuration
      console = {
        open_on_runcode = true,
        dir = "row",
        size = {
          width = "90%",
          height = "75%",
        },
        result = {
          size = "60%",
        },
        testcase = {
          virt_text = true,
          size = "40%",
        },
      },
      
      -- Description panel configuration
      description = {
        position = "left",
        width = "40%",
        show_stats = true,
      },
      
      -- Hooks for custom behavior
      hooks = {
        ["enter"] = {},
        ["question_enter"] = {},
        ["leave"] = {},
      },
    })
    
    -- Key mappings for common operations
    vim.keymap.set("n", "<leader>ll", "<cmd>Leet<CR>", { silent = true })
    vim.keymap.set("n", "<leader>lt", "<cmd>Leet test<CR>", { silent = true })
    vim.keymap.set("n", "<leader>ls", "<cmd>Leet submit<CR>", { silent = true })
    vim.keymap.set("n", "<leader>ld", "<cmd>Leet desc<CR>", { silent = true })
    ```

## Command Line and Menu Enhancement

### wilder.nvim: Enhanced Command-Line Interface for Neovim
- **Description**: A more adventurous wildmenu for Neovim with enhanced command-line capabilities
- **Features**:
  - Automatic command suggestions as you type
  - Support for command-line (`:`) and search (`/`) modes
  - Fuzzy matching and filtering capabilities
  - Customizable appearance with multiple renderer options
  - Popupmenu and wildmenu rendering styles
  - Gradient highlighting for matched characters
  - File finder pipeline for project-wide file navigation
  - Devicons integration for file type visualization
  - Command palette mode for VS Code-like experience
  - Extensive pipeline customization options
  - Performance optimization features (debouncing, lazy loading)
  - Compatible with both Vim and Neovim
  - 1.4k+ stars and active maintenance
  - Repository: https://github.com/gelguy/wilder.nvim
  - Configuration example:
    ```lua
    -- wilder.nvim configuration
    local wilder = require('wilder')
    wilder.setup({modes = {':', '/', '?'}})
    
    -- Enable fuzzy matching and highlighting
    wilder.set_option('pipeline', {
      wilder.branch(
        wilder.cmdline_pipeline({
          fuzzy = 1,
          set_pcre2_pattern = 1,
        }),
        wilder.python_search_pipeline({
          pattern = 'fuzzy',
        })
      ),
    })
    
    -- Use popupmenu for command mode and wildmenu for search
    wilder.set_option('renderer', wilder.renderer_mux({
      [':'] = wilder.popupmenu_renderer({
        highlighter = {
          wilder.pcre2_highlighter(),
          wilder.basic_highlighter(),
        },
        left = {' ', wilder.popupmenu_devicons()},
        right = {' ', wilder.popupmenu_scrollbar()},
      }),
      ['/'] = wilder.wildmenu_renderer({
        highlighter = {
          wilder.pcre2_highlighter(),
          wilder.basic_highlighter(),
        },
      }),
    }))
    
    -- For a more fancy look with borders
    wilder.set_option('renderer', wilder.popupmenu_renderer(
      wilder.popupmenu_border_theme({
        highlights = {
          border = 'Normal',
        },
        border = 'rounded',
        highlighter = wilder.basic_highlighter(),
        left = {' ', wilder.popupmenu_devicons()},
        right = {' ', wilder.popupmenu_scrollbar()},
      })
    ))
    ```

## Language Support: TypeScript and JavaScript

### coc-tsserver: TypeScript Server Extension for Neovim
- **Description**: Tsserver language server extension for coc.nvim that provides rich features like VSCode for JavaScript and TypeScript
- **Features**:
  - Complete TypeScript/JavaScript language server integration
  - Code completion with automatic imports
  - Code validation and diagnostics
  - Go to definition and references
  - Document symbols and outline
  - Hover documentation
  - Signature help
  - Code actions and quick fixes
  - Refactoring capabilities
  - Organize imports
  - Format code
  - Rename symbols
  - Automatic tag closing
  - Rename imports on file rename
  - Inlay hints support
  - Semantic tokens
  - Call hierarchy
  - Selection range
  - Workspace symbol search
  - TypeScript version selection
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```json
    // coc-settings.json
    {
      "tsserver.enable": true,
      "tsserver.log": "off",
      "tsserver.trace.server": "off",
      "tsserver.maxTsServerMemory": 3072,
      "tsserver.disableAutomaticTypeAcquisition": false,
      "tsserver.useLocalTsdk": false,
      "tsserver.tsdk": "",
      "typescript.suggest.enabled": true,
      "typescript.suggest.paths": true,
      "typescript.suggest.autoImports": true,
      "typescript.format.enable": true,
      "typescript.validate.enable": true,
      "typescript.showUnused": true,
      "typescript.updateImportsOnFileMove.enabled": "prompt",
      "typescript.preferences.importModuleSpecifier": "shortest",
      "typescript.preferences.quoteStyle": "double",
      "typescript.preferences.useAliasesForRenames": true,
      "typescript.inlayHints.parameterNames.enabled": "none",
      "typescript.inlayHints.variableTypes.enabled": false,
      "typescript.inlayHints.functionLikeReturnTypes.enabled": false,
      "javascript.suggest.enabled": true,
      "javascript.suggest.paths": true,
      "javascript.suggest.autoImports": true,
      "javascript.format.enable": true,
      "javascript.validate.enable": true,
      "javascript.showUnused": true
    }
    ```
  - Installation:
    ```vim
    :CocInstall coc-tsserver
    ```

## Container and Remote Development

### nvim-remote-containers: Docker Container Development in Neovim
- **Description**: A plugin that provides VSCode-like remote container development functionality for Neovim
- **Features**:
  - Develop inside Docker containers directly from Neovim
  - Parse and use devcontainer.json configuration files
  - Build Docker images from Dockerfiles
  - Attach to existing containers
  - Start containers from pre-built images
  - Docker Compose integration (up, down, destroy)
  - Container status indicator in statusline
  - Treesitter integration for parsing configuration files
  - Lua API for container management
  - Vim command wrappers for common operations
  - 900+ stars and active maintenance
  - Repository: https://github.com/jamestthompson3/nvim-remote-containers
  - Configuration example:
    ```lua
    -- nvim-remote-containers configuration
    -- No explicit setup function, but you can create custom mappings
    
    -- Key mappings for common operations
    vim.api.nvim_set_keymap('n', '<leader>cb', ':BuildImage<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>ca', ':AttachToContainer<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cs', ':StartImage<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cu', ':ComposeUp<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cd', ':ComposeDown<CR>', { noremap = true, silent = true })
    vim.api.nvim_set_keymap('n', '<leader>cx', ':ComposeDestroy<CR>', { noremap = true, silent = true })
    
    -- Add container status to statusline
    vim.cmd([[
      hi Container guifg=#BADA55 guibg=Black
      set statusline+=%#Container#%{g:currentContainer}
    ]])
    
    -- Available Lua functions:
    -- parseConfig: Parses devcontainer.json file
    -- attachToContainer: Attaches to a container or builds from image
    -- buildImage: Builds container from Dockerfile
    -- composeUp: Brings docker-compose up
    -- composeDown: Brings docker-compose down
    -- composeDestroy: Destroys docker-compose containers
    
    -- Example usage in Lua:
    -- require('nvim-remote-containers').buildImage(true) -- Show build progress in floating window
    -- require('nvim-remote-containers').attachToContainer()
    ```
  - Requirements:
    - Docker installed and running
    - Treesitter with jsonc module installed (`:TSInstall jsonc`)

## Database Management

### nvim-dbee: Interactive Database Client for Neovim
- **Description**: A full-featured database client for Neovim that allows executing queries and viewing results directly in the editor
- **Features**:
  - Support for multiple database types (PostgreSQL, MySQL, SQLite, etc.)
  - Interactive query execution and result visualization
  - Connection management with support for multiple sources
  - Scratchpad system for storing and organizing queries
  - Pagination for large result sets
  - Result export in various formats (JSON, CSV, table)
  - Syntax highlighting for SQL
  - Secure credential management with template support
  - Go-powered backend for performance
  - Lua frontend for seamless Neovim integration
  - Extensible API for custom integrations
  - 980+ stars and active maintenance
  - Repository: https://github.com/kndndrj/nvim-dbee
  - Configuration example:
    ```lua
    -- nvim-dbee configuration
    require("dbee").setup({
      -- Sources for loading connections
      sources = {
        -- Load connections from memory
        require("dbee.sources").MemorySource:new({
          {
            name = "Local SQLite DB",
            type = "sqlite",
            url = "~/path/to/database.db",
          },
          {
            name = "Development PostgreSQL",
            type = "postgres",
            url = "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
          },
        }),
        -- Load connections from environment variable
        require("dbee.sources").EnvSource:new("DBEE_CONNECTIONS"),
        -- Load connections from file with persistence
        require("dbee.sources").FileSource:new(vim.fn.stdpath("cache") .. "/dbee/persistence.json"),
      },
      
      -- UI configuration
      ui = {
        -- Component layout
        layout = {
          -- Component positions
          position = {
            -- Tree/drawer position
            tree = "left", -- "left" or "right"
            -- Editor position
            editor = "top", -- "top" or "bottom"
            -- Result position
            result = "bottom", -- "top" or "bottom"
          },
          -- Component sizes
          size = {
            -- Width of tree component
            tree = "20%",
            -- Height of editor component
            editor = "30%",
            -- Height of result component
            result = "30%",
          },
        },
      },
      
      -- Result configuration
      result = {
        -- Number of rows to fetch at once
        page_size = 100,
        -- Format options for exporting results
        format = {
          -- CSV format options
          csv = {
            delimiter = ",",
            header = true,
          },
          -- JSON format options
          json = {
            array = true,
            indentation = 2,
          },
        },
      },
      
      -- Key mappings
      mappings = {
        -- Tree/drawer mappings
        tree = {
          -- Toggle node expansion
          toggle = "o",
          -- Refresh tree
          refresh = "r",
          -- Edit connection
          edit = "cw",
          -- Delete connection
          delete = "dd",
          -- Execute action
          execute = "<CR>",
        },
        -- Editor mappings
        editor = {
          -- Execute query
          execute = "BB",
        },
        -- Result mappings
        result = {
          -- Next page
          next_page = "L",
          -- Previous page
          prev_page = "H",
          -- First page
          first_page = "F",
          -- Last page
          last_page = "E",
          -- Yank row as JSON
          yank_row_json = "yaj",
          -- Yank row as CSV
          yank_row_csv = "yac",
          -- Yank all rows as JSON
          yank_all_json = "yaJ",
          -- Yank all rows as CSV
          yank_all_csv = "yaC",
        },
      },
    })
    
    -- Key mappings for common operations
    vim.keymap.set("n", "<leader>db", "<cmd>Dbee toggle<CR>", { silent = true })
    vim.keymap.set("n", "<leader>de", "<cmd>Dbee execute<CR>", { silent = true })
    ```
  - Requirements:
    - Neovim >= 0.10
    - nui.nvim dependency
    - Go (for building the backend)

## Interactive Code Execution

### molten-nvim: Interactive Jupyter Kernel Integration for Neovim
- **Description**: A Neovim plugin for interactively running code with the Jupyter kernel, providing a REPL-like and notebook-like experience
- **Features**:
  - Asynchronous code execution with Jupyter kernels
  - Real-time output display as virtual text or in floating windows
  - Image, plot, and LaTeX rendering directly in Neovim
  - Support for stdin input via vim.ui.input
  - Multi-buffer to single kernel and single buffer to multi-kernel support
  - Compatible with any language that has a Jupyter kernel
  - Python virtual environment support
  - Import/export capabilities with Jupyter notebook files
  - Cell-based code execution with highlighted regions
  - Customizable output display options
  - Automatic output window management
  - Image rendering through multiple providers (image.nvim, wezterm)
  - Status line integration
  - Extensive configuration options
  - 800+ stars and active maintenance
  - Repository: https://github.com/benlubas/molten-nvim
  - Configuration example:
    ```lua
    -- molten-nvim configuration
    vim.g.molten_output_win_max_height = 12
    vim.g.molten_image_provider = "image.nvim" -- or "wezterm"
    vim.g.molten_virt_text_output = true
    vim.g.molten_wrap_output = true
    vim.g.molten_auto_open_output = false
    
    -- Key mappings
    vim.keymap.set("n", "<localleader>mi", ":MoltenInit<CR>", 
      { silent = true, desc = "Initialize Jupyter kernel" })
    vim.keymap.set("n", "<localleader>e", ":MoltenEvaluateOperator<CR>", 
      { silent = true, desc = "Run operator selection" })
    vim.keymap.set("n", "<localleader>rl", ":MoltenEvaluateLine<CR>", 
      { silent = true, desc = "Evaluate line" })
    vim.keymap.set("n", "<localleader>rr", ":MoltenReevaluateCell<CR>", 
      { silent = true, desc = "Re-evaluate cell" })
    vim.keymap.set("v", "<localleader>r", ":<C-u>MoltenEvaluateVisual<CR>gv", 
      { silent = true, desc = "Evaluate visual selection" })
    vim.keymap.set("n", "<localleader>rd", ":MoltenDelete<CR>", 
      { silent = true, desc = "Delete cell" })
    vim.keymap.set("n", "<localleader>os", ":noautocmd MoltenEnterOutput<CR>", 
      { silent = true, desc = "Show/enter output" })
    
    -- Auto-initialization for Python files
    vim.api.nvim_create_autocmd("FileType", {
      pattern = "python",
      callback = function()
        -- Initialize Python kernel when opening Python files
        vim.schedule(function()
          vim.cmd("MoltenInit python3")
        end)
      end,
    })
    ```
  - Requirements:
    - Neovim 9.4+
    - Python 3.10+
    - pynvim and jupyter_client Python packages
    - Optional: image.nvim or wezterm.nvim for image rendering

## Terminal Multiplexing

### tmux.nvim: Seamless Tmux Integration for Neovim
- **Description**: A plugin that provides seamless integration between Neovim and tmux, enabling navigation, resizing, and clipboard synchronization
- **Features**:
  - Navigate between Neovim splits and tmux panes with the same keybindings
  - Resize Neovim splits like tmux panes
  - Synchronize clipboard between Neovim and tmux
  - Swap panes with adjacent tmux panes
  - Cycle navigation between panes
  - Support for window navigation with Ctrl-n and Ctrl-p
  - Configurable keybindings for all operations
  - Automatic detection of tmux environment
  - TPM (Tmux Plugin Manager) support
  - Lua API for custom integrations
  - 700+ stars and active maintenance
  - Repository: https://github.com/aserowy/tmux.nvim
  - Configuration example:
    ```lua
    -- tmux.nvim configuration
    require("tmux").setup({
      copy_sync = {
        -- enables copy sync. by default, all registers are synchronized.
        -- to control which registers are synced, see the `sync_*` options.
        enable = true,
        
        -- ignore specific tmux buffers e.g. buffer0 = true to ignore the
        -- first buffer or named_buffer_name = true to ignore a named tmux
        -- buffer with name named_buffer_name :)
        ignore_buffers = { empty = false },
        
        -- TMUX >= 3.2: all yanks (and deletes) will get redirected to system
        -- clipboard by tmux
        redirect_to_clipboard = false,
        
        -- offset controls where register sync starts
        -- e.g. offset 2 lets registers 0 and 1 untouched
        register_offset = 0,
        
        -- overwrites vim.g.clipboard to redirect * and + to the system
        -- clipboard using tmux. If you sync your system clipboard without tmux,
        -- disable this option!
        sync_clipboard = true,
        
        -- synchronizes registers *, +, unnamed, and 0 till 9 with tmux buffers.
        sync_registers = true,
        
        -- synchronizes registers when pressing p and P.
        sync_registers_keymap_put = true,
        
        -- synchronizes registers when pressing (C-r) and ".
        sync_registers_keymap_reg = true,
        
        -- syncs deletes with tmux clipboard as well, it is adviced to
        -- do so. Nvim does not allow syncing registers 0 and 1 without
        -- overwriting the unnamed register. Thus, ddp would not be possible.
        sync_deletes = true,
        
        -- syncs the unnamed register with the first buffer entry from tmux.
        sync_unnamed = true,
      },
      navigation = {
        -- cycles to opposite pane while navigating into the border
        cycle_navigation = true,
        
        -- enables default keybindings (C-hjkl) for normal mode
        enable_default_keybindings = true,
        
        -- prevents unzoom tmux when navigating beyond vim border
        persist_zoom = false,
      },
      resize = {
        -- enables default keybindings (A-hjkl) for normal mode
        enable_default_keybindings = true,
        
        -- sets resize steps for x axis
        resize_step_x = 1,
        
        -- sets resize steps for y axis
        resize_step_y = 1,
      },
      swap = {
        -- cycles to opposite pane while navigating into the border
        cycle_navigation = false,
        
        -- enables default keybindings (C-A-hjkl) for normal mode
        enable_default_keybindings = true,
      }
    })
    ```
  - Requirements:
    - Neovim >= 0.5
    - tmux >= 3.2a (recommended)
    - Proper tmux configuration in ~/.tmux.conf

## Language Support: Go

### vim-go: Comprehensive Go Development Plugin for Vim/Neovim
- **Description**: A complete Go development environment for Vim/Neovim with syntax highlighting, code completion, formatting, and more
- **Features**:
  - Syntax highlighting and improved folding for Go files
  - Auto-completion support via gopls
  - Code navigation (go to definition, references, implementation)
  - Code formatting on save
  - Build, install, test, and run commands
  - Integrated debugging with Delve
  - Documentation lookup
  - Import management
  - Type-safe renaming
  - Code coverage visualization
  - Struct tag management
  - Linting and static analysis
  - Snippet support
  - Integration with other tools (Tagbar, Ultisnips)
  - 16k+ stars and active maintenance
  - Repository: https://github.com/fatih/vim-go
  - Configuration example:
    ```vim
    " vim-go configuration
    " Enable syntax highlighting
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_operators = 1
    let g:go_highlight_extra_types = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_generate_tags = 1
    
    " Auto formatting and importing
    let g:go_fmt_autosave = 1
    let g:go_fmt_command = "goimports"
    
    " Status line types/signatures
    let g:go_auto_type_info = 1
    
    " Run :GoBuild or :GoTestCompile based on the go file
    function! s:build_go_files()
      let l:file = expand('%')
      if l:file =~# '^\f\+_test\.go$'
        call go#test#Test(0, 1)
      elseif l:file =~# '^\f\+\.go$'
        call go#cmd#Build(0)
      endif
    endfunction
    
    " Key mappings
    autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
    autocmd FileType go nmap <leader>r <Plug>(go-run)
    autocmd FileType go nmap <leader>t <Plug>(go-test)
    autocmd FileType go nmap <leader>c <Plug>(go-coverage-toggle)
    autocmd FileType go nmap <Leader>i <Plug>(go-info)
    autocmd FileType go nmap <Leader>d <Plug>(go-def)
    autocmd FileType go nmap <Leader>v <Plug>(go-def-vertical)
    autocmd FileType go nmap <Leader>s <Plug>(go-def-split)
    autocmd FileType go nmap <Leader>dt <Plug>(go-def-tab)
    autocmd FileType go nmap <Leader>gd <Plug>(go-doc)
    autocmd FileType go nmap <Leader>gv <Plug>(go-doc-vertical)
    autocmd FileType go nmap <Leader>gb <Plug>(go-doc-browser)
    autocmd FileType go nmap <Leader>a <Plug>(go-alternate-edit)
    autocmd FileType go nmap <Leader>ae <Plug>(go-alternate-edit)
    autocmd FileType go nmap <Leader>av <Plug>(go-alternate-vertical)
    autocmd FileType go nmap <Leader>as <Plug>(go-alternate-split)
    
    " Use gopls for completion
    let g:go_gopls_enabled = 1
    let g:go_gopls_options = ['-remote=auto']
    
    " Debugging settings
    let g:go_debug_windows = {
          \ 'vars':  'rightbelow 60vnew',
          \ 'stack': 'rightbelow 10new',
          \ 'goroutines': 'botright 10new',
          \ 'out':   'botright 5new',
      \ }
    let g:go_debug_address = '127.0.0.1:8181'
    let g:go_debug_breakpoint_sign_text = '>'
    ```
  - Requirements:
    - Neovim >= 0.4.0 or Vim >= 8.1.2269
    - Go installed and configured
    - Various Go tools (automatically installed with `:GoInstallBinaries`)

### go.nvim: Modern Go Development Plugin for Neovim
- **Description**: A feature-rich Go plugin for Neovim that leverages treesitter, LSP, and DAP
- **Features**:
  - Modern implementation using Lua
  - Treesitter integration for better syntax highlighting
  - LSP integration with gopls
  - Debug Adapter Protocol (DAP) support
  - Test running and debugging
  - Code generation and refactoring
  - Go module management
  - Struct tag management
  - Code navigation and information
  - Customizable UI elements
  - 2.4k+ stars and active maintenance
  - Repository: https://github.com/ray-x/go.nvim
  - Configuration example:
    ```lua
    -- go.nvim configuration
    require('go').setup({
      -- gopls configuration
      lsp_cfg = {
        settings = {
          gopls = {
            analyses = {
              unusedparams = true,
              shadow = true,
            },
            staticcheck = true,
          },
        },
      },
      -- Run gofmt + goimport on save
      lsp_gofumpt = true,
      -- Function signature settings
      lsp_on_attach = true,
      -- Test settings
      test_runner = 'go',
      run_in_floaterm = true,
      -- Debugging settings
      dap_debug = true,
      dap_debug_gui = true,
      -- Linting settings
      linter = 'golangci-lint',
      -- Treesitter settings
      treesitter = true,
    })
    
    -- Key mappings
    local function map(mode, lhs, rhs, desc)
      vim.keymap.set(mode, lhs, rhs, { silent = true, desc = desc })
    end
    
    map('n', '<leader>gc', '<cmd>GoCoverage<CR>', 'Go Coverage')
    map('n', '<leader>gt', '<cmd>GoTest<CR>', 'Go Test')
    map('n', '<leader>gtf', '<cmd>GoTestFunc<CR>', 'Go Test Function')
    map('n', '<leader>gr', '<cmd>GoRun<CR>', 'Go Run')
    map('n', '<leader>gi', '<cmd>GoInstall<CR>', 'Go Install')
    map('n', '<leader>gta', '<cmd>GoAddTag<CR>', 'Go Add Tags')
    map('n', '<leader>gtr', '<cmd>GoRmTag<CR>', 'Go Remove Tags')
    map('n', '<leader>gfs', '<cmd>GoFillStruct<CR>', 'Go Fill Struct')
    map('n', '<leader>gif', '<cmd>GoIfErr<CR>', 'Go If Err')
    ```
  - Requirements:
    - Neovim >= 0.8.0
    - Treesitter
    - Go installed and configured
    - Optional: Delve debugger for debugging

### nvim-dap-go: Go Debugging Extension for Neovim
- **Description**: An extension for nvim-dap providing configurations for launching the Go debugger (Delve) and debugging individual tests
- **Features**:
  - Auto-launch Delve with zero configuration
  - Debug individual tests directly from the cursor position
  - Debug main functions and running processes
  - Support for command-line arguments and build flags
  - Integration with treesitter for test detection
  - Remote debugging capabilities
  - VSCode launch.json configuration support
  - 544+ stars and active maintenance
  - Repository: https://github.com/leoluz/nvim-dap-go
  - Configuration example:
    ```lua
    -- nvim-dap-go configuration
    require('dap-go').setup({
      -- Additional dap configurations
      dap_configurations = {
        {
          type = "go",
          name = "Attach remote",
          mode = "remote",
          request = "attach",
        },
      },
      -- delve configurations
      delve = {
        -- Path to the delve executable
        path = "dlv",
        -- Timeout for delve initialization
        initialize_timeout_sec = 20,
        -- Port for delve debugger
        port = "${port}",
        -- Build flags for delve
        build_flags = {},
      },
      -- Test options
      tests = {
        -- Enable verbose test output
        verbose = false,
      },
    })
    
    -- Key mappings
    vim.keymap.set('n', '<leader>dt', function() require('dap-go').debug_test() end, { desc = 'Debug Go Test' })
    vim.keymap.set('n', '<leader>dl', function() require('dap-go').debug_last_test() end, { desc = 'Debug Last Go Test' })
    ```
  - Requirements:
    - Neovim >= 0.9.0
    - nvim-dap plugin
    - Delve >= 1.7.0
    - Go treesitter parser

### gopls: Official Go Language Server
- **Description**: The official Go language server developed by the Go team, providing IDE features to any LSP-compatible editor
- **Features**:
  - Code completion and intelligent suggestions
  - Go to definition, references, and implementation
  - Hover information and documentation
  - Diagnostics and error reporting
  - Code formatting and organization
  - Refactoring tools
  - Symbol search
  - Workspace symbol search
  - Signature help
  - Document symbols
  - Rename
  - Find references
  - Code actions
  - Semantic highlighting
  - LSP-compatible with any editor
  - Repository: https://github.com/golang/tools/tree/master/gopls
  - Configuration example for Neovim:
    ```lua
    -- gopls configuration in Neovim with nvim-lspconfig
    require('lspconfig').gopls.setup({
      cmd = {'gopls'},
      settings = {
        gopls = {
          analyses = {
            unusedparams = true,
            shadow = true,
          },
          staticcheck = true,
          gofumpt = true,
          usePlaceholders = true,
          completeUnimported = true,
          matcher = 'fuzzy',
          directoryFilters = {
            "-node_modules",
            "-vendor",
          },
          semanticTokens = true,
        },
      },
      on_attach = function(client, bufnr)
        -- Attach keybindings for gopls commands
        vim.api.nvim_buf_set_keymap(bufnr, 'n', 'gd', '<cmd>lua vim.lsp.buf.definition()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', 'gr', '<cmd>lua vim.lsp.buf.references()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', 'gi', '<cmd>lua vim.lsp.buf.implementation()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', 'K', '<cmd>lua vim.lsp.buf.hover()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', '<C-k>', '<cmd>lua vim.lsp.buf.signature_help()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', '<leader>rn', '<cmd>lua vim.lsp.buf.rename()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', '<leader>ca', '<cmd>lua vim.lsp.buf.code_action()<CR>', {noremap = true})
        vim.api.nvim_buf_set_keymap(bufnr, 'n', '<leader>f', '<cmd>lua vim.lsp.buf.formatting()<CR>', {noremap = true})
      end,
      capabilities = require('cmp_nvim_lsp').default_capabilities(),
    })
    ```
  - Requirements:
    - Go installed
    - Neovim >= 0.5.0
    - nvim-lspconfig plugin

### go.nvim: Modern Go Development Plugin for Neovim
- **Description**: A modern Go Neovim plugin based on treesitter, nvim-lsp and dap debugger, written in Lua and async as much as possible
- **Features**:
  - Per-project setup with project files (.gonvim)
  - Async jobs with libuv
  - Native treesitter support for syntax highlighting and textobjects
  - LSP integration for code navigation and intelligence
  - Gopls commands (fillstruct, organize imports, list modules, etc.)
  - Runtime lint/vet/compile with LSP and golangci-lint
  - Async build/make/test support
  - Test coverage visualization
  - Debugging with nvim-dap and Dap UI
  - VSCode launch configuration support
  - Unit test generation with gotests
  - Struct tag management
  - Code formatting with LSP and gofmt/golines
  - CodeLens support
  - Auto-documentation for packages/functions/structs
  - Go to alternative file (between test and source)
  - Code refactoring tools
  - Go cheatsheet integration
  - Smart build tag detection
  - Mock generation
  - Inlay hints
  - LuaSnip integration
  - Treesitter highlight injection for SQL and JSON
  - 2.4k+ stars and active maintenance
  - Repository: https://github.com/ray-x/go.nvim
  - Configuration example:
    ```lua
    require('go').setup({
      -- Basic setup
      go = 'go',
      goimports = 'gopls',
      gofmt = 'gopls',
      max_line_len = 120,
      tag_transform = false,
      test_dir = '',
      comment_placeholder = '   ',
      
      -- LSP configuration
      lsp_cfg = true,
      lsp_gofumpt = true,
      lsp_on_attach = true,
      lsp_codelens = true,
      
      -- Gopls settings
      gopls_cmd = nil,
      gopls_remote_auto = true,
      
      -- Test settings
      test_runner = 'go',
      verbose_tests = true,
      
      -- Debugging
      dap_debug = true,
      dap_debug_gui = true,
      
      -- Terminal settings
      run_in_floaterm = false,
      
      -- Textobjects
      textobjects = true,
    })
    
    -- Format on save
    local format_sync_grp = vim.api.nvim_create_augroup("GoFormat", {})
    vim.api.nvim_create_autocmd("BufWritePre", {
      pattern = "*.go",
      callback = function()
        require('go.format').goimports()
      end,
      group = format_sync_grp,
    })
    ```
  - Requirements:
    - Neovim >= 0.5.0
    - Go installed
    - Treesitter
    - nvim-lspconfig
    - Optional: nvim-dap for debugging

## Language Support: TypeScript

### coc-tsserver: TypeScript Server Extension for coc.nvim
- **Description**: A comprehensive TypeScript/JavaScript language server extension for coc.nvim that provides rich IDE-like features
- **Features**:
  - Code completion with auto-imports
  - Go to definition, references, and implementation
  - Code validation and diagnostics
  - Document highlighting and symbols
  - Folding and folding range support
  - Formatting (document, range, and on-type)
  - Hover documentation
  - CodeLens for implementations and references
  - Organize imports command
  - Quickfix using code actions
  - Code refactoring
  - Source code actions (fix issues, remove unused code, etc.)
  - Find references
  - Signature help
  - Call hierarchy
  - Selection range
  - Semantic tokens
  - Rename symbols
  - Automatic tag closing
  - Rename imports on file rename
  - Workspace symbol search
  - Inlay hints support
  - 1.1k+ stars and active maintenance
  - Repository: https://github.com/neoclide/coc-tsserver
  - Configuration example:
    ```json
    // coc-settings.json
    {
      "tsserver.enable": true,
      "tsserver.maxTsServerMemory": 3072,
      "tsserver.log": "off",
      "tsserver.trace.server": "off",
      "tsserver.useLocalTsdk": false,
      "tsserver.disableAutomaticTypeAcquisition": false,
      "typescript.suggest.enabled": true,
      "typescript.suggest.paths": true,
      "typescript.suggest.autoImports": true,
      "typescript.suggest.completeFunctionCalls": true,
      "typescript.format.enable": true,
      "typescript.format.insertSpaceAfterCommaDelimiter": true,
      "typescript.format.insertSpaceAfterFunctionKeywordForAnonymousFunctions": true,
      "typescript.implementationsCodeLens.enabled": true,
      "typescript.referencesCodeLens.enabled": true,
      "typescript.preferences.importModuleSpecifier": "shortest",
      "typescript.preferences.quoteStyle": "single",
      "typescript.updateImportsOnFileMove.enabled": "prompt",
      "typescript.inlayHints.parameterNames.enabled": "literals",
      "typescript.inlayHints.variableTypes.enabled": true,
      "typescript.inlayHints.functionLikeReturnTypes.enabled": true
    }
    ```
  - Requirements:
    - Neovim with coc.nvim installed
    - Node.js installed

### typescript-tools.nvim: Modern TypeScript Integration for Neovim
- **Description**: A high-performance TypeScript integration for Neovim that communicates directly with tsserver using its native protocol
- **Features**:
  - Direct communication with tsserver (like VSCode) without typescript-language-server
  - Significantly better performance for large TypeScript/JavaScript projects
  - Support for TypeScript versions 4.0 and above
  - Multiple tsserver instances support
  - Local and global TypeScript installations support
  - Styled-components support
  - Improved code refactoring capabilities
  - Code actions (organize imports, add missing imports, remove unused, etc.)
  - Semantic tokens support
  - Inlay hints
  - Code lens support
  - JSX close tag support
  - Call hierarchy
  - Workspace symbol search
  - File references
  - Rename file with reference updates
  - 1.7k+ stars and active maintenance
  - Repository: https://github.com/pmizio/typescript-tools.nvim
  - Configuration example:
    ```lua
    require("typescript-tools").setup({
      -- Basic setup
      on_attach = function(client, bufnr)
        -- Add your keybindings here
      end,
      
      settings = {
        -- Spawn additional tsserver instance for diagnostics
        separate_diagnostic_server = true,
        -- When to publish diagnostics
        publish_diagnostic_on = "insert_leave",
        -- Code actions to expose
        expose_as_code_action = {},
        -- Custom tsserver path
        tsserver_path = nil,
        -- Plugins for tsserver
        tsserver_plugins = {},
        -- Memory limit for tsserver
        tsserver_max_memory = "auto",
        -- Format options
        tsserver_format_options = {},
        -- File preferences
        tsserver_file_preferences = {},
        -- Locale for tsserver messages
        tsserver_locale = "en",
        -- Complete function calls
        complete_function_calls = false,
        -- Include completions with insert text
        include_completions_with_insert_text = true,
        -- Code lens settings
        code_lens = "off",
        -- Disable member code lens
        disable_member_code_lens = true,
        -- JSX close tag settings
        jsx_close_tag = {
          enable = false,
          filetypes = { "javascriptreact", "typescriptreact" },
        },
      },
    })
    ```
  - Requirements:
    - Neovim >= 0.8.0
    - nvim-lspconfig
    - plenary.nvim
    - Node.js

### typescript.nvim: TypeScript Language Server Integration for Neovim
- **Description**: A Lua plugin written in TypeScript that provides TypeScript language server integration and convenience commands
- **Features**:
  - Easy setup of typescript-language-server via nvim-lspconfig
  - Commands for common TypeScript operations:
    - Add missing imports
    - Organize imports
    - Remove unused variables
    - Fix all (fixes specific issues like non-async functions using await)
    - Rename file with import updates
    - Go to source definition
  - Handlers for off-spec LSP methods
  - Integration with null-ls for code actions
  - Inlay hints support
  - 491 stars (archived as of August 2023)
  - Repository: https://github.com/jose-elias-alvarez/typescript.nvim
  - Configuration example:
    ```lua
    require("typescript").setup({
      disable_commands = false, -- prevent the plugin from creating Vim commands
      debug = false, -- enable debug logging for commands
      go_to_source_definition = {
        fallback = true, -- fall back to standard LSP definition on failure
      },
      server = { -- pass options to lspconfig's setup method
        on_attach = function(client, bufnr)
          -- your keybindings here
        end,
        settings = {
          -- typescript server settings
        }
      },
    })
    ```
  - Requirements:
    - Neovim >= 0.7.0
    - typescript-language-server installed

## Language Support: Kubernetes

### vikube.vim: Kubernetes Cluster Management from Vim
- **Description**: A Vim plugin for operating Kubernetes clusters directly from Vim, providing a terminal-based interface for common kubectl operations
- **Features**:
  - Context management (list, switch, rename, delete)
  - Resource navigation and exploration
  - Pod logs viewing and following
  - Container command execution
  - Resource YAML viewing and editing
  - Resource labeling and deletion
  - Namespace switching
  - Support for multiple resource types (pods, services, nodes, etc.)
  - Top command integration for resource monitoring
  - Wide output option for detailed information
  - Auto-update capability
  - Buffer-scoped context switching
  - 198 stars and active maintenance
  - Repository: https://github.com/c9s/vikube.vim
  - Configuration example:
    ```vim
    " Enable automatic list update
    let g:vikube_autoupdate = 1
    
    " Change default tail lines for logs
    let g:vikube_default_logs_tail = 100
    
    " Use current namespace instead of "default"
    let g:vikube_use_current_namespace = 1
    
    " Disable custom highlight for CursorLine
    let g:vikube_disable_custom_highlight = 1
    ```
  - Key mappings:
    ```vim
    " Default mappings
    " <leader>kc - Open the kubernetes context list buffer
    " <leader>kno - Open the kubernetes node list buffer
    " <leader>kpo - Open the kubernetes pod list buffer
    " <leader>ksv - Open the kubernetes service list buffer
    " <leader>kt - Open the kubernetes top buffer
    " <leader>kpvc - Open the kubernetes persistent volume claim buffer
    ```
  - Requirements:
    - Vim or Neovim
    - kubectl installed and configured
    - Dependencies: helper.vim, treemenu.vim

### kubectl.nvim: Kubernetes Integration for Neovim
- **Description**: A Neovim plugin for interacting with Kubernetes clusters using kubectl, providing a terminal-based interface for common operations
- **Features**:
  - Kubectl command integration
  - Resource management (pods, services, deployments, etc.)
  - Context and namespace switching
  - Log viewing and following
  - Port forwarding
  - Resource YAML editing
  - Exec into containers
  - Resource deletion and scaling
  - Supports multiple clusters
  - Integration with telescope.nvim for fuzzy finding
  - Repository: Various implementations available
  - Configuration example:
    ```lua
    require('kubectl').setup({
      -- Basic setup
      kubectl_path = 'kubectl',
      -- Default namespace
      namespace = 'default',
      -- Default context
      context = nil,
      -- Telescope integration
      use_telescope = true,
      -- Floating window settings
      float = {
        width = 0.8,
        height = 0.8,
        border = 'rounded',
      },
    })
    ```
  - Requirements:
    - Neovim >= 0.5.0
    - kubectl installed and configured
    - Optional: telescope.nvim for enhanced navigation

## Language Support: Terraform

### vim-terraform: Basic Terraform Integration for Vim/Neovim
- **Description**: A Vim plugin for basic Terraform integration providing syntax highlighting, indentation, and formatting
- **Features**:
  - Syntax highlighting for Terraform files (.tf, .tfvars)
  - Automatic indentation and formatting
  - Commands for terraform fmt
  - Integration with terraform validate
  - Support for HCL files
  - 1.1k stars and active maintenance
  - Repository: https://github.com/hashivim/vim-terraform
  - Configuration example:
    ```vim
    " Allow vim-terraform to align settings automatically
    let g:terraform_align=1
    
    " Allow vim-terraform to automatically format .tf and .tfvars files
    let g:terraform_fmt_on_save=1
    ```
  - Requirements:
    - Vim or Neovim
    - Terraform installed (for fmt and validate features)

### terraform-lsp: Language Server Protocol for Terraform
- **Description**: A Language Server Protocol implementation for Terraform that provides intelligent features like autocompletion, error checking, and more
- **Features**:
  - Autocompletion for resources, data sources, variables, and locals
  - Dynamic error checking (Terraform and HCL validation)
  - Provider config completion
  - Resource and data source attribute completion
  - Module nesting variable completion
  - Communication using provider binary (supports any provider)
  - Integration with multiple editors including Vim/Neovim
  - 587 stars and active maintenance
  - Repository: https://github.com/juliosueiras/terraform-lsp
  - Configuration example for Neovim with coc.nvim:
    ```vim
    " Add terraform-lsp to coc.nvim
    let g:coc_global_extensions = [
      \ 'coc-json',
      \ 'coc-terraform',
    \ ]
    ```
  - Configuration example for Neovim with built-in LSP:
    ```lua
    require'lspconfig'.terraformls.setup{
      cmd = {"terraform-lsp"},
      filetypes = {"terraform", "tf"},
      root_dir = function(fname)
        return util.root_pattern(".terraform", ".git")(fname)
      end,
    }
    ```
  - Requirements:
    - Neovim with LSP support or coc.nvim
    - Terraform installed
    - Go 1.14+ (if building from source)

## Implementation Considerations

### Integration with agent_runtime
- Neovim should be treated as a secondary Agent Computer Interface (ACI)
- Need to implement proper state management between Neovim and main IDE
- Consider event stream architecture for state synchronization

### Kata Container Integration
- Neovim needs to be available within the kata container sandbox
- Configuration should be initialized during container startup
- Ensure proper permissions and dependencies

### Performance Optimization
- Minimize plugin overhead to maintain responsiveness
- Consider lazy-loading plugins when possible
- Optimize LSP configurations for performance

## Terminal Multiplexing and Session Management

### Tmux Integration
- **tmux.nvim**: Seamless tmux integration (https://github.com/aserowy/tmux.nvim)
  - Enables navigation between Neovim splits and tmux panes with consistent keybindings
  - Supports resizing splits and panes with the same commands
  - Synchronizes clipboard between Neovim instances and tmux buffers
  - Enables parallel workflow across multiple terminal sessions
  - Configuration options for navigation, resizing, and copy synchronization

### Markdown and Documentation
- **markdown-preview.nvim**: Live Markdown preview (https://github.com/iamcco/markdown-preview.nvim)
  - Real-time preview of Markdown files in browser
  - Support for GitHub-flavored Markdown
  - Automatic refresh on file changes
  - Essential for documentation-heavy workflows

### Visual Enhancements
- **tokyonight.nvim**: Modern colorscheme (https://github.com/folke/tokyonight.nvim)
  - Multiple style variants (storm, night, day, moon)
  - Terminal colors integration
  - Customizable highlights
  - Improves visual appeal and reduces eye strain

### Next Steps
1. Create detailed configuration for each language server
2. Design state management architecture between Neovim and main IDE
3. Implement initialization script for Neovim in kata container
4. Develop testing framework for parallel workflow capabilities
5. Create documentation for Neovim integration in agent_runtime
6. Evaluate additional language-specific plugins for Python, Go, and TypeScript
