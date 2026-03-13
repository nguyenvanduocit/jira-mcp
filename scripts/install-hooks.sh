#!/usr/bin/env bash
set -euo pipefail

# Installs git hooks for the jira-mcp project.
# Run once after cloning: ./scripts/install-hooks.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Installing git hooks..."

cat > "$HOOKS_DIR/pre-commit" << 'HOOK'
#!/usr/bin/env bash
# Pre-commit hook: scan staged changes for secrets
exec ./scripts/security-scan.sh --staged
HOOK

chmod +x "$HOOKS_DIR/pre-commit"
echo "Installed pre-commit hook (secret scanning)."
echo "Done."
