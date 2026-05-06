# AI v2

`tool/ai/v2` is the newer AI package for provider capability checks, plain chat, streaming chat, and lightweight stateful agents.

The directory name is `v2`, but the Go package name is `ai`:

```go
import ai "github.com/intmian/mian_go_lib/tool/ai/v2"
```

This package does not replace the older `github.com/intmian/mian_go_lib/tool/ai` wrapper yet. Use it when you need provider capabilities, reasoning controls, streaming, base agent components, or agent setting JSON helpers.

## Core concepts

1. `ProviderID` and `AgentID` are `uint32`.
   - `0` means unset and is rejected.
   - Register providers once with `AddProvider`.
2. `AgentID` identifies an agent setting.
   - `AgentIDBase` is the built-in base chat agent setting ID.
   - Agent IDs are for setting registration, not agent construction.
3. `IAgentSetting` is the persisted/configurable agent setting surface.
   - Settings support JSON export/import and JSON doc generation.
4. `IAgent[S]` is the minimal typed agent lifecycle surface.
   - It exposes `GetID`, initialization from registered settings, and initialization from an explicit typed setting.
   - It does not include chat, streaming, factory, scheduler, or middleware behavior.
5. `BaseAgent` is a stateful one-on-one chat agent and the built-in minimal implementation.
   - Each `BaseAgent` instance remembers only its own successful chat history.
   - Failed chats do not enter history.
   - It is composed from the same public provider binding, setting state, and message history components that custom agents can reuse.
6. The package uses one package-level registry singleton.
   - Callers register providers and agent settings once.
   - Callers still initialize agents explicitly; the package does not create agents by ID.
   - Tests inside this package reset the singleton with package-private helpers.

## Quick start

Register a provider, then initialize a base agent with explicit settings:

```go
package example

import (
	"context"

	ai "github.com/intmian/mian_go_lib/tool/ai/v2"
)

func Run(ctx context.Context, baseURL string, token string) (string, error) {
	provider, err := ai.NewOpenAIProvider(baseURL, token, ai.ProviderSourceTypeOpenAI)
	if err != nil {
		return "", err
	}
	if err := ai.AddProvider(1, provider); err != nil {
		return "", err
	}

	agent := ai.NewBaseAgent()
	if err := agent.InitWithSetting(&ai.BaseAgentSetting{
		ProviderID:      1,
		SysPrompt:       "You are a helpful assistant.",
		Models:          []string{"gpt-5.4", "gpt-5-chat-latest"},
		ReasoningEffort: ai.ReasoningEffortMedium,
	}); err != nil {
		return "", err
	}
	return agent.Chat(ctx, "hello")
}
```

## Create Agents

Use `NewBaseAgent` to create the thin built-in agent, then call `InitWithSetting` for explicit settings or `Init` to read `AgentIDBase` from the setting registry. The agent still resolves providers from the package singleton.

```go
agent := ai.NewBaseAgent()
err := agent.InitWithSetting(&ai.BaseAgentSetting{
	ProviderID:      1,
	SysPrompt:       "Answer briefly.",
	Models:          []string{"gpt-5.4-mini"},
	ReasoningEffort: ai.ReasoningEffortLow,
})
```

The package does not create agents by ID. Callers own agent selection and pass prompts, models, provider IDs, and reasoning effort through settings.

`BaseAgent` implements:

```go
var _ ai.IAgent[*ai.BaseAgentSetting] = (*ai.BaseAgent)(nil)
```

Reusable public components are available for custom agents:

- `AgentProviderBinding`: binds a `ProviderID` and resolves the registered provider.
- `AgentSettingState[S]`: reads typed settings by `agent.GetID()`, stores typed settings, and returns the stored setting.
- `BaseAgentSettingState`: alias for `AgentSettingState[*BaseAgentSetting]`.
- `AgentMessageHistory`: stores successful chat history and builds message lists with an optional system prompt.

`BaseAgent` is a reference composition and the built-in minimal agent; external agents should compose the public components directly rather than embedding or extending `BaseAgent`.

## Register and read settings

Register settings with `AddAgentSetting`:

```go
err := ai.AddAgentSetting(ai.AgentIDBase, &ai.BaseAgentSetting{
	ProviderID: 1,
	Models:     []string{"gpt-5.4"},
})
```

Read the untyped setting when the caller does not know its concrete type:

```go
setting, ok := ai.GetAgentSetting(ai.AgentIDBase)
```

Read the typed setting when the caller knows the type:

```go
setting, ok := ai.GetAgentSettingAs[*ai.BaseAgentSetting](ai.AgentIDBase)
```

`GetAgentSettingAs` returns the registered pointer. `AgentSettingState[S]` stores the setting it receives and does not deep-copy it, so callers that need immutable snapshots should pass their own copy.

## BaseAgent behavior

`BaseAgent.Chat(ctx, content)`:

1. Trims and validates `content`.
2. Resolves `setting.ProviderID` from the package provider registry.
3. Builds messages from `SysPrompt`, successful chat history, and the current user message.
4. Calls models in `Models` order.
5. Returns the first non-empty successful response.
6. Records user and assistant messages only after success.
7. Returns joined model errors when every model fails.

`Models` is a concrete model-name fallback list, not a mode name list.

`BaseAgent.Chat` serializes calls on the same agent instance, including the provider network call. This preserves strict history ordering for the stateful one-on-one chat use case; callers that need parallel requests should use separate agent instances or a custom agent composition.

## Provider usage

Create OpenAI-compatible providers with `NewOpenAIProvider`:

```go
openaiProvider, err := ai.NewOpenAIProvider("", token, ai.ProviderSourceTypeOpenAI)
deepSeekProvider, err := ai.NewOpenAIProvider("https://api.deepseek.com", token, ai.ProviderSourceTypeDeepSeek)
```

Supported source types:

- `ProviderSourceTypeOpenAI`
- `ProviderSourceTypeDeepSeek`

Provider capabilities:

```go
models, err := provider.AvailableModels(ctx)
reasoning := provider.AvailableReasoning()
tools := provider.AvailableTools()
```

`AvailableTools()` currently returns no callable tools. The first v2 chat surface intentionally supports plain messages plus reasoning/thinking only.

## Direct provider chat

Use a provider directly when you do not need agent history:

```go
resp, err := provider.Chat(ctx, ai.ChatRequest{
	Model: "gpt-5.4",
	Messages: []ai.ChatMessage{
		{Role: ai.ChatRoleSystem, Content: "You are concise."},
		{Role: ai.ChatRoleUser, Content: "Summarize this."},
	},
	ReasoningEffort: ai.ReasoningEffortMedium,
})
```

For streaming:

```go
stream, err := provider.ChatStream(ctx, req)
if err != nil {
	return err
}
defer stream.Close()

for {
	event, err := stream.Recv()
	if err != nil {
		break
	}
	switch event.Type {
	case ai.ChatStreamEventTextDelta:
		// event.TextDelta
	case ai.ChatStreamEventReasoningDelta:
		// event.ReasoningDelta
	case ai.ChatStreamEventDone:
		// event.Response
	}
}
```

## Setting JSON helpers

Agent settings support:

```go
data, err := setting.ExportJSON()
err = setting.ImportJSON(data)
docs, err := setting.ExportJSONDoc()
```

Generic helpers are also available:

```go
data, err := ai.ExportSettingJSON(setting)
next, err := ai.ImportSettingJSON(setting, data)
docs, err := ai.ExportSettingDoc(setting)
```

Supported setting fields:

- Exported scalar fields: string, integer, float, bool, and aliases of those kinds.
- Exported scalar slices/arrays.

Unsupported setting fields:

- maps
- structs
- pointers
- interfaces
- nested slices/arrays

Zero values are treated as not configured:

- JSON export omits zero values.
- JSON import does not overwrite existing values with zero values.
- JSON import cannot clear a previously configured scalar or slice by passing a zero value such as `""`, `0`, `false`, or `[]`; omit/zero means "leave the existing value unchanged".

## Errors and invariants

Important validation rules:

- `ProviderID == 0` is invalid.
- Duplicate provider IDs are rejected.
- Duplicate agent setting IDs are rejected.
- `BaseAgentSetting.Models` must contain at least one non-empty model after trimming whitespace.
- `Init` requires a registered `*BaseAgentSetting` at `AgentIDBase`.
- `InitWithSetting` requires `ProviderID` and at least one model.
- The referenced provider must be registered before provider access or chat, not before initialization.

Provider lookup happens at chat time. This allows providers to be registered after agent construction, but before the first successful chat.

## Package notes

- The package singleton is process-local.
- The package singleton stores providers and agent settings, but it does not contain agent factories.
- If a caller needs isolation, add explicit test seams or a separate package-level lifecycle before broadening the public API.
- Keep business policy outside this package. Pass prompts, models, provider IDs, and reasoning effort through settings.
