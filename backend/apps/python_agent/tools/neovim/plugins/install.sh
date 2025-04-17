
set -e

NEOVIM_CONFIG_DIR="$HOME/.config/nvim"
PLUGINS_DIR="$(dirname "$0")"

echo "Installing Neovim plugins for agent_runtime..."

mkdir -p "$NEOVIM_CONFIG_DIR"

echo "Copying init.lua to $NEOVIM_CONFIG_DIR..."
cp "$PLUGINS_DIR/init.lua" "$NEOVIM_CONFIG_DIR/init.lua"

echo "Installing dependencies..."
apt-get update
apt-get install -y git curl unzip nodejs npm python3-pip

echo "Installing Python dependencies for molten-nvim..."
pip3 install pynvim jupyter_client

echo "Installing Node.js dependencies for coc.nvim..."
npm install -g neovim

echo "Installing fzf..."
if [ ! -d "$HOME/.fzf" ]; then
  git clone --depth 1 https://github.com/junegunn/fzf.git "$HOME/.fzf"
  "$HOME/.fzf/install" --all
fi

echo "Creating minimal init.vim to bootstrap lazy.nvim..."
cat > "$NEOVIM_CONFIG_DIR/init.vim" << EOF
" Set leader key
let mapleader = " "

" Load Lua configuration
lua require('init')
EOF

echo "Creating coc-settings.json for TypeScript support..."
cat > "$NEOVIM_CONFIG_DIR/coc-settings.json" << EOF
{
  "tsserver.enable": true,
  "typescript.suggest.enabled": true,
  "typescript.format.enabled": true,
  "typescript.updateImportsOnFileMove.enabled": "always",
  "typescript.preferences.importModuleSpecifier": "relative",
  "coc.preferences.formatOnSaveFiletypes": ["typescript", "typescriptreact", "javascript", "javascriptreact"],
  "suggest.noselect": false,
  "suggest.enablePreselect": false,
  "suggest.completionItemKindLabels": {
    "class": "📦",
    "color": "🎨",
    "constant": "🔒",
    "default": "📎",
    "enum": "🔢",
    "enumMember": "🔢",
    "event": "🔆",
    "field": "🔑",
    "file": "📄",
    "folder": "📁",
    "function": "⚙️",
    "interface": "🔗",
    "keyword": "🔑",
    "method": "⚙️",
    "module": "📦",
    "operator": "➕",
    "property": "🔑",
    "reference": "🔍",
    "snippet": "✂️",
    "struct": "🔧",
    "text": "📝",
    "typeParameter": "🔠",
    "unit": "📏",
    "value": "💯",
    "variable": "📌"
  }
}
EOF

echo "Creating tmux configuration for tmux.nvim integration..."
cat > "$HOME/.tmux.conf" << EOF
unbind C-b
set -g prefix C-Space
bind C-Space send-prefix

set -g mouse on

set -g default-terminal "screen-256color"
set -ga terminal-overrides ",*256col*:Tc"

set -g base-index 1
setw -g pane-base-index 1

set -g renumber-windows on

set -g history-limit 10000

is_vim="ps -o state= -o comm= -t '#{pane_tty}' \
    | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?(view|n?vim?x?)(diff)?$'"
bind-key -n 'C-h' if-shell "$is_vim" 'send-keys C-h'  'select-pane -L'
bind-key -n 'C-j' if-shell "$is_vim" 'send-keys C-j'  'select-pane -D'
bind-key -n 'C-k' if-shell "$is_vim" 'send-keys C-k'  'select-pane -U'
bind-key -n 'C-l' if-shell "$is_vim" 'send-keys C-l'  'select-pane -R'
tmux_version='$(tmux -V | sed -En "s/^tmux ([0-9]+(.[0-9]+)?).*/\1/p")'
if-shell -b '[ "$(echo "$tmux_version < 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\'  'select-pane -l'"
if-shell -b '[ "$(echo "$tmux_version >= 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\\\'  'select-pane -l'"

bind -r H resize-pane -L 5
bind -r J resize-pane -D 5
bind -r K resize-pane -U 5
bind -r L resize-pane -R 5

bind | split-window -h -c "#{pane_current_path}"
bind - split-window -v -c "#{pane_current_path}"
EOF

echo "Installation complete! Neovim plugins are ready to use."
echo "To verify installation, run 'nvim' and check for plugin loading messages."
