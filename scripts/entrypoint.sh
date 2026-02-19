#!/usr/bin/env sh

file_env() {
   var="$1"
   fileVar="${var}_FILE"

   val_var="$(printenv "$var")"
   val_fileVar="$(printenv "$fileVar")"

   if [ -n "$val_var" ] && [ -n "$val_fileVar" ]; then
      echo "error: both $var and $fileVar are set (but are exclusive)" >&2
      exit 1
   fi

   if [ -n "$val_var" ]; then
      val="$val_var"
   elif [ -n "$val_fileVar" ]; then
      if [ ! -r "$val_fileVar" ]; then
         echo "error: file '$val_fileVar' does not exist or is not readable" >&2
         exit 1
      fi
      val="$(cat "$val_fileVar")"
   fi

   export "$var"="$val"
   unset "$fileVar"
}

file_env "DIVINESENSE_DSN"

# Geek Mode: Install Claude Code CLI if enabled
if [ "$(printenv DIVINESENSE_CLAUDE_CODE_ENABLED)" = "true" ]; then
   echo "Geek Mode enabled: Checking Claude Code CLI installation..."
   if ! command -v claude >/dev/null 2>&1; then
      echo "Installing Claude Code CLI..."
      npm install -g @anthropic-ai/claude-code
      if [ $? -eq 0 ]; then
         echo "Claude Code CLI installed successfully"
         claude --version 2>/dev/null || echo "Claude Code CLI installed (version check skipped)"
      else
         echo "Warning: Failed to install Claude Code CLI. Geek Mode may not work properly."
      fi
   else
      echo "Claude Code CLI already installed"
      claude --version 2>/dev/null || true
   fi
fi

exec "$@"
