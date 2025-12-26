# Git Safety Rules

## Never Auto-Commit or Push

**CRITICAL:** Never automatically commit or push code without explicit user approval.

### Required Workflow

1. **Before committing:**
   - Show all changes with `git diff` or `git status`
   - Ask user: "Ready to commit these changes?"
   - Wait for explicit confirmation

2. **Before pushing:**
   - Show what will be pushed with `git log origin/main..HEAD`
   - Ask user: "Ready to push to remote?"
   - Wait for explicit confirmation

3. **Never run these without asking:**
   - `git commit` (unless user explicitly says "commit")
   - `git push` (unless user explicitly says "push")
   - `git push --force` (ALWAYS ask, even if user says push)
   - `git reset --hard`
   - `git rebase`

### Safe Operations (can run without asking)

- `git status`
- `git diff`
- `git log`
- `git branch`
- `git stash` (for temporary work)
