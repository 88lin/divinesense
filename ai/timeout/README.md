# AI Timeout Constants (`ai/timeout`)

The `timeout` package centrally manages all timeout constant configurations for AI modules, ensuring system robustness when facing external calls (such as LLM APIs).

## Core Constants

### Agent Execution
| Constant | Default | Description |
| :-------- | :------ | :---------- |
| `AgentTimeout` | 2 min | Maximum duration for entire Agent processing (thinking + tool calls) |
| `AgentExecutionTimeout` | 2 min | Alias for AgentTimeout (backward compatibility) |
| `MaxIterations` | 5 | Maximum ReAct loop iterations, prevents infinite loops |
| `MaxToolIterations` | 5 | Alias for MaxIterations |

### LLM & Tools
| Constant | Default | Description |
| :-------- | :------ | :---------- |
| `StreamTimeout` | 5 min | Maximum hold time for streaming responses |
| `ToolExecutionTimeout` | 30 s | Timeout for single tool execution |
| `EmbeddingTimeout` | 30 s | Timeout for vector generation |

### Fault Tolerance
| Constant | Default | Description |
| :-------- | :------ | :---------- |
| `MaxToolFailures` | 3 | Maximum consecutive failures before aborting |
| `MaxRecentToolCalls` | 10 | Number of recent tool calls to track for loop detection |
| `MaxTruncateLength` | 200 | Maximum length for truncating strings in logs |

## Usage

Import this package directly to use constants instead of hardcoding numbers, facilitating unified system strategy adjustment.

```go
ctx, cancel := context.WithTimeout(parentCtx, timeout.AgentTimeout)
defer cancel()
```
