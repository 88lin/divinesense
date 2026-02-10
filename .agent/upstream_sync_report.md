# Upstream Synchronization Report (2026-02-10)

## 1. Overview
The upstream repository (`usememos/memos`) has been focused on **stability, bug fixes, and minor refactoring** in the last 50-100 commits. No significant new features (`feat`) were detected in the recent history.

## 2. Recent Highlights

### 🛠️ Important Fixes
- **Auth/Security**: 
  - `cf0a285e` fix(auth): make PKCE optional for OAuth sign-in (#5570) - *Critical for OAuth compatibility*
  - `81ef53b3` fix: prevent 401 errors on window focus when token expires
  - `d9e8387d` fix(postgres): handle missing PAT data gracefully
- **UI/UX**:
  - `b5108b4f` fix(web): calendar navigation should use current page path (#5605)
  - `74b63b27` perf: disable tooltips in year calendar to fix lag
  - `e7605d90` fix: shortcut edit button opens create dialog instead of edit dialog
  - `c4176b4e` fix: videos attachment handling

### 🧹 Refactoring & Chores
- `d9dc5be2` fix: replace `echo.NewHTTPError` with `status.Errorf`
- `cf65f086` refactor: remove hide-scrollbar utility
- `b623162d` chore: fix static check linter warnings

## 3. Analysis & Recommendation
Since there are no major features to port, the synchronization strategy should focus on **cherry-picking critical bug fixes** and performance improvements.

**Recommended Actions:**
1.  **Cherry-pick Auth Fixes**: The OAuth PKCE fix (`cf0a285e`) seems relevant if we use OAuth.
2.  **Apply UI Fixes**: The calendar and shortcut fixes improve user experience.
3.  **Update Dependencies**: Check `go.mod` and `frontend/package.json` for updates.

## 4. Next Steps
Please select which category you would like to proceed with:
- [ ] **Technical Debt**: Apply linter fixes and refactors.
- [ ] **UI Polish**: Apply calendar and dialog fixes.
- [ ] **Security/Auth**: Apply OAuth and PAT fixes.
