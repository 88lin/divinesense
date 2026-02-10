# 测试规范

> DivineSense 测试最佳实践 — 单元测试、集成测试、E2E

---

## 核心原则

1. **测试金字塔**：大量单元测试 + 少量集成测试 + 最少 E2E
2. **测试驱动**：关键功能先写测试
3. **覆盖率目标**：核心模块 > 80%，整体 > 60%

---

## Go 测试

### 文件组织

```
package/
├── module.go
└── module_test.go
```

### 表格驱动测试

```go
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{"valid input", "valid", true, false},
		{"invalid input", "invalid", false, true},
		{"empty input", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### Mock 外部依赖

```go
// 使用 gomock
type MockStore struct {
	ctrl *gomock.Controller
}

func (m *MockStore) Get(id int64) (*Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	return ret[0].(*Item), ret[1].(error)
}
```

### 并发测试

```go
func TestConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 测试代码
		}()
	}
	wg.Wait()
}
```

---

## React 测试

### 组件测试

```tsx
import { render, screen } from '@testing-library/react';

test('renders button', () => {
  render(<Button>Click me</Button>);
  expect(screen.getByText('Click me')).toBeInTheDocument();
});
```

### Hook 测试

```tsx
import { renderHook, act } from '@testing-library/react';

test('useCounter increments', () => {
  const { result } = renderHook(() => useCounter());
  act(() => {
    result.current.increment();
  });
  expect(result.current.count).toBe(1);
});
```

---

## 测试命令

```bash
# Go 测试
make test              # 所有测试
make test-ai           # AI 相关测试
go test ./path -v      # 详细输出
go test -race ./...    # 竞态检测

# 前端测试
pnpm test              # 所有测试
pnpm test:coverage     # 覆盖率报告
```

---

## 测试检查清单

- [ ] 单元测试覆盖核心逻辑
- [ ] 边界条件测试（空值、极值）
- [ ] 错误路径测试
- [ ] 并发安全测试（如适用）
- [ ] 集成测试覆盖关键流程
- [ ] CI 中自动运行

---

*文档路径：.claude/rules/testing.md*
