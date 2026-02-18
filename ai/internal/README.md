# AI Internal Utilities (`ai/internal`)

The `internal` package contains common utility functions and auxiliary code used internally by the AI module. These codes are not exposed as public APIs to external modules.

## Subpackages

### `strutil`

Provides string processing utilities, especially for safe handling of multi-byte characters (such as Chinese).

- **`Truncate(s string, maxLen int)`**: Safely truncate string by character count (Rune), preventing cutting multi-byte characters causing garbled text or panic.

## Usage Example

```go
import "github.com/hrygo/divinesense/ai/internal/strutil"

// Safe truncation for CJK text
truncated := strutil.Truncate("你好世界", 3) // Returns "你好..."
```
