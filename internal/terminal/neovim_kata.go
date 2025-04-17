package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	
	"github.com/spectrumwebco/agent_runtime/internal/containers"
	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type NeovimKataTerminal struct {
	id            string
	apiURL        string
	running       bool
	mu            sync.RWMutex
	client        *http.Client
	options       map[string]interface{}
	kataManager   *containers.KataManager
	containerName string
}

func NewNeovimKataTerminal(id string, apiURL string, options map[string]interface{}, 
                          kataManager *containers.KataManager) *NeovimKataTerminal {
	return &NeovimKataTerminal{
		id:            id,
		apiURL:        apiURL,
		running:       false,
		client:        &http.Client{Timeout: 10 * time.Second},
		options:       options,
		kataManager:   kataManager,
		containerName: fmt.Sprintf("neovim-%s", id),
	}
}

func (t *NeovimKataTerminal) ID() string {
	return t.id
}

func (t *NeovimKataTerminal) GetType() string {
	return "neovim-kata"
}

func (t *NeovimKataTerminal) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

func (t *NeovimKataTerminal) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if t.running {
		return fmt.Errorf("terminal %s is already running", t.id)
	}
	
	containerOptions := map[string]interface{}{
		"image": "neovim/neovim:latest", // Default image
	}
	
	if t.options != nil {
		for k, v := range t.options {
			containerOptions[k] = v
		}
	}
	
	pluginsVolume := "/tmp/neovim-plugins"
	containerOptions["volumes"] = map[string]string{
		pluginsVolume: "/root/.config/nvim",
	}
	
	if err := t.setupNeovimPlugins(pluginsVolume); err != nil {
		return fmt.Errorf("setup neovim plugins: %w", err)
	}
	
	if err := t.kataManager.CreateContainer(ctx, t.containerName, containerOptions); err != nil {
		return fmt.Errorf("create kata container: %w", err)
	}
	
	startCmd := []string{"nvim", "--headless", "--listen", "0.0.0.0:9999"}
	if _, err := t.kataManager.ExecInContainer(ctx, t.containerName, startCmd); err != nil {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("start neovim in container: %w", err)
	}
	
	payload := map[string]interface{}{
		"id":            t.id,
		"container_id":  t.containerName,
		"container_type": "kata",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("marshal start payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", 
				fmt.Sprintf("%s/terminals/register", t.apiURL), strings.NewReader(string(jsonPayload)))
	if err != nil {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("create register request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := t.client.Do(req)
	if err != nil {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("register terminal: %w", err)
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("decode register response: %w", err)
	}
	
	if !result.Success {
		_ = t.kataManager.RemoveContainer(ctx, t.containerName)
		return fmt.Errorf("register terminal failed: %s", result.Error)
	}
	
	t.running = true
	return nil
}

func (t *NeovimKataTerminal) Stop(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if !t.running {
		return nil
	}
	
	payload := map[string]interface{}{
		"id": t.id,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal stop payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", 
				fmt.Sprintf("%s/terminals/unregister", t.apiURL), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("create unregister request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("unregister terminal: %w", err)
	}
	defer resp.Body.Close()
	
	if err := t.kataManager.StopContainer(ctx, t.containerName); err != nil {
		return fmt.Errorf("stop kata container: %w", err)
	}
	
	if err := t.kataManager.RemoveContainer(ctx, t.containerName); err != nil {
		return fmt.Errorf("remove kata container: %w", err)
	}
	
	t.running = false
	return nil
}

func (t *NeovimKataTerminal) Execute(ctx context.Context, command string) (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	if !t.running {
		return "", fmt.Errorf("terminal %s is not running", t.id)
	}
	
	nvimCmd := []string{"nvim", "--server", "0.0.0.0:9999", "--remote-send", command}
	output, err := t.kataManager.ExecInContainer(ctx, t.containerName, nvimCmd)
	if err != nil {
		return "", fmt.Errorf("execute command in container: %w", err)
	}
	
	return output, nil
}

func (t *NeovimKataTerminal) setupNeovimPlugins(pluginsDir string) error {
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return fmt.Errorf("create plugins directory: %w", err)
	}
	
	initLua := `
-- Neovim configuration for agent_runtime
-- This configuration includes plugins for enhanced functionality

-- Install packer if not installed
local install_path = vim.fn.stdpath('data')..'/site/pack/packer/start/packer.nvim'
if vim.fn.empty(vim.fn.glob(install_path)) > 0 then
  vim.fn.system({'git', 'clone', '--depth', '1', 'https://github.com/wbthomason/packer.nvim', install_path})
  vim.cmd [[packadd packer.nvim]]
end

-- Plugin configuration
require('packer').startup(function(use)
  -- Packer can manage itself
  use 'wbthomason/packer.nvim'
  
  -- Color scheme
  use 'folke/tokyonight.nvim'
  
  -- Tmux integration
  use 'aserowy/tmux.nvim'
  
  -- Fuzzy finder (Lua implementation)
  use 'ibhagwan/fzf-lua'
  
  -- Interactive notebook functionality
  use 'benlubas/molten-nvim'
  
  -- Markdown preview
  use {
    'iamcco/markdown-preview.nvim',
    run = function() vim.fn['mkdp#util#install']() end,
  }
  
  -- Terraform and HCL support
  use 'hashivim/vim-terraform'
  use 'jvirtanen/vim-hcl'
end)

-- Configure color scheme
vim.cmd[[colorscheme tokyonight]]

-- Configure tmux.nvim
require('tmux').setup({
  copy_sync = {
    enable = true,
  },
  navigation = {
    enable_default_keybindings = true,
  },
  resize = {
    enable_default_keybindings = true,
  },
})

-- Configure fzf-lua
require('fzf-lua').setup({
  winopts = {
    height = 0.85,
    width = 0.80,
    row = 0.35,
    col = 0.50,
  }
})

-- Configure molten-nvim
vim.g.molten_output_win_max_height = 20
vim.g.molten_image_provider = 'image.nvim'

-- Set up key mappings
vim.api.nvim_set_keymap('n', '<leader>ff', "<cmd>lua require('fzf-lua').files()<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>fg', "<cmd>lua require('fzf-lua').git_files()<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>fb', "<cmd>lua require('fzf-lua').buffers()<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>fr', "<cmd>lua require('fzf-lua').live_grep()<CR>", {noremap = true, silent = true})

-- Molten key mappings
vim.api.nvim_set_keymap('n', '<leader>mi', "<cmd>MoltenInit<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>me', "<cmd>MoltenEvaluateOperator<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>ml', "<cmd>MoltenEvaluateLine<CR>", {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<leader>mr', "<cmd>MoltenReevaluateCell<CR>", {noremap = true, silent = true})

-- General settings
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.expandtab = true
vim.opt.shiftwidth = 2
vim.opt.tabstop = 2
vim.opt.autoindent = true
vim.opt.smartindent = true
vim.opt.wrap = false
vim.opt.ignorecase = true
vim.opt.smartcase = true
vim.opt.termguicolors = true
vim.opt.mouse = 'a'
vim.opt.clipboard = 'unnamedplus'
vim.opt.backup = false
vim.opt.writebackup = false
vim.opt.swapfile = false
vim.opt.updatetime = 300
vim.opt.timeoutlen = 500
vim.opt.completeopt = {'menuone', 'noselect'}
vim.opt.signcolumn = 'yes'
vim.opt.scrolloff = 8
vim.opt.sidescrolloff = 8
vim.opt.cursorline = true
`
	
	if err := os.WriteFile(filepath.Join(pluginsDir, "init.lua"), []byte(initLua), 0644); err != nil {
		return fmt.Errorf("write init.lua: %w", err)
	}
	
	return nil
}
