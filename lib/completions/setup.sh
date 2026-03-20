#!/usr/bin/env bash
# Setup shell completions for micepad CLI
# Usage: ./lib/completions/setup.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

detect_shell() {
  local shell_name
  shell_name="$(basename "${SHELL:-/bin/bash}")"
  echo "$shell_name"
}

install_bash() {
  local target="${HOME}/.bash_completion.d"
  mkdir -p "$target"
  cp "$SCRIPT_DIR/micepad.bash" "$target/micepad"
  echo "Installed bash completion to $target/micepad"

  local rc="${HOME}/.bashrc"
  local source_line="[ -f ~/.bash_completion.d/micepad ] && source ~/.bash_completion.d/micepad"
  if ! grep -qF "bash_completion.d/micepad" "$rc" 2>/dev/null; then
    echo "" >> "$rc"
    echo "$source_line" >> "$rc"
    echo "Added source line to $rc"
  else
    echo "Source line already in $rc"
  fi
  echo "Run: source $rc"
}

install_zsh() {
  # Prefer oh-my-zsh plugin if oh-my-zsh is installed
  if [[ -d "${ZSH:-$HOME/.oh-my-zsh}/custom/plugins" ]]; then
    install_zsh_omz
  else
    install_zsh_fpath
  fi
}

install_zsh_omz() {
  local target="${ZSH:-$HOME/.oh-my-zsh}/custom/plugins/micepad"
  local source_dir="$SCRIPT_DIR/oh-my-zsh"

  if [[ -L "$target" ]]; then
    echo "Symlink already exists: $target"
  elif [[ -d "$target" ]]; then
    rm -rf "$target"
    ln -s "$source_dir" "$target"
    echo "Replaced directory with symlink: $target -> $source_dir"
  else
    ln -s "$source_dir" "$target"
    echo "Created symlink: $target -> $source_dir"
  fi

  local rc="${HOME}/.zshrc"
  if grep -qE '^\s*plugins=\(' "$rc" 2>/dev/null; then
    if grep -qE 'plugins=.*micepad' "$rc" 2>/dev/null; then
      echo "Plugin 'micepad' already in plugins list"
    else
      echo ""
      echo "Add 'micepad' to your plugins list in $rc:"
      echo "  plugins=(... micepad)"
    fi
  fi
  echo "Run: source $rc"
}

install_zsh_fpath() {
  local target="${HOME}/.zsh/completions"
  mkdir -p "$target"
  cp "$SCRIPT_DIR/micepad.zsh" "$target/_micepad"
  echo "Installed zsh completion to $target/_micepad"

  local rc="${HOME}/.zshrc"
  local fpath_line='fpath=(~/.zsh/completions $fpath)'
  if ! grep -qF ".zsh/completions" "$rc" 2>/dev/null; then
    echo "" >> "$rc"
    echo "$fpath_line" >> "$rc"
    echo "autoload -Uz compinit && compinit" >> "$rc"
    echo "Added fpath and compinit to $rc"
  else
    echo "fpath already configured in $rc"
  fi
  echo "Run: source $rc"
}

install_fish() {
  local target="${HOME}/.config/fish/completions"
  mkdir -p "$target"
  cp "$SCRIPT_DIR/micepad.fish" "$target/micepad.fish"
  echo "Installed fish completion to $target/micepad.fish"
  echo "Fish completions are loaded automatically."
}

shell_name="${1:-$(detect_shell)}"

case "$shell_name" in
  bash)
    install_bash
    ;;
  zsh)
    install_zsh
    ;;
  fish)
    install_fish
    ;;
  all)
    echo "Installing completions for all shells..."
    echo ""
    echo "=== Bash ==="
    install_bash
    echo ""
    echo "=== Zsh ==="
    install_zsh
    echo ""
    echo "=== Fish ==="
    install_fish
    ;;
  *)
    echo "Usage: $0 [bash|zsh|fish|all]"
    echo ""
    echo "Installs shell completions for the micepad CLI."
    echo "If no argument is given, detects your current shell."
    exit 1
    ;;
esac

echo ""
echo "Done! Restart your shell or source your rc file."
