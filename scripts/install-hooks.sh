#!/bin/bash
# install-hooks.sh - Install git hooks from scripts directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_DIR="$(git rev-parse --git-common-dir)/hooks"

echo "📦 Installing git hooks..."
echo ""

# Copy pre-commit hook (lightweight - runs on every commit)
if [ -f "$SCRIPT_DIR/pre-commit" ]; then
    cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
    chmod +x "$HOOKS_DIR/pre-commit"
    echo "  ✓ pre-commit  → 快速检查 (fmt + vet)，~2秒"
else
    echo "  ✗ pre-commit hook not found in $SCRIPT_DIR"
    exit 1
fi

# Copy pre-push hook (full CI checks - runs on git push)
if [ -f "$SCRIPT_DIR/pre-push" ]; then
    cp "$SCRIPT_DIR/pre-push" "$HOOKS_DIR/pre-push"
    chmod +x "$HOOKS_DIR/pre-push"
    echo "  ✓ pre-push   → 完整 CI 检查 (golangci-lint + test + build)，~1分钟"
else
    echo "  ✗ pre-push hook not found in $SCRIPT_DIR"
    exit 1
fi

echo ""
echo "✅ Git hooks installed successfully!"
echo ""
echo "检查时机:"
echo "  • pre-commit  → 每次 commit 时"
echo "  • pre-push     → 每次 push 到远程时"
echo ""
echo "跳过检查:"
echo "  • commit:  git commit --no-verify -m 'msg'"
echo "  • push:   git push --no-verify"
