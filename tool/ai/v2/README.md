# AI v2

`tool/ai/v2` is the newer AI package for provider capability checks, plain chat, streaming chat, and lightweight stateful agents.

The directory name is `v2`, but the Go package name is `ai`:

```go
import ai "github.com/intmian/mian_go_lib/tool/ai/v2"
```

This package does not replace the older `github.com/intmian/mian_go_lib/tool/ai` wrapper yet. Use it when you need provider capabilities, reasoning controls, streaming, or agent settings.

## Core concepts

1. `ProviderID` is a `uint32`.
   - `0` means unset and is rejected.
   - Register providers once with `AddProvider`.
2. `AgentID` identifies an agent setting.
   - `AgentIDBase` is the built-in base chat agent setting ID.
3. `IAgentSetting` is the persisted/configurable agent setting surface.
   - Settings support JSON export/import and JSON doc generation.
4. `BaseAgent` is a stateful one-on-one chat agent.
   - Each `BaseAgent` instance remembers only its own successful chat history.
   - Failed chats do not enter history.
5. The package uses one package-level registry singleton.
   - Callers do not create or pass registry objects.
   - Tests inside this package reset the singleton with package-private helpers.

## Quick start

Register a provider, register the default base agent setting, then create an agent:

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

	if err := ai.AddAgentSetting(ai.AgentIDBase, &ai.BaseAgentSetting{
		ProviderID:      1,
		SysPrompt:       "You are a helpful assistant.",
		Models:          []string{"gpt-5.4", "gpt-5-chat-latest"},
		ReasoningEffort: ai.ReasoningEffortMedium,
	}); err != nil {
		return "", err
	}

	agent, err := ai.NewBaseAgent()
	if err != nil {
		return "", err
	}
	return agent.Chat(ctx, "hello")
}
```

Use `NewBaseAgent()` when the package singleton already has `AgentIDBase` registered.

## Create with explicit settings

Use `NewBaseAgentWithSetting` for a temporary or child agent. It does not read registered agent settings, but it still resolves providers from the package singleton.

```go
agent, err := ai.NewBaseAgentWithSetting(ai.BaseAgentSetting{
	ProviderID:      1,
	SysPrompt:       "Answer briefly.",
	Models:          []string{"gpt-5.4-mini"},
	ReasoningEffort: ai.ReasoningEffortLow,
})
```

This is the preferred path for sub-agents that share the provider registry but need their own prompt/model list.

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
if ok {
	setting.Models = []string{"gpt-5.4-mini"}
}
```

`GetAgentSettingAs` returns the registered pointer. Mutating it changes future `NewBaseAgent()` calls. Existing `BaseAgent` instances keep the setting copy they were created with.

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

## Errors and invariants

Important validation rules:

- `ProviderID == 0` is invalid.
- Duplicate provider IDs are rejected.
- Duplicate agent setting IDs are rejected.
- `BaseAgentSetting.Models` must not be empty.
- `NewBaseAgent()` requires a registered `*BaseAgentSetting` at `AgentIDBase`.
- `NewBaseAgentWithSetting` requires the referenced provider to be registered before chat, not before construction.

Provider lookup happens at chat time. This allows providers to be registered after agent construction, but before the first successful chat.

## Package notes

- The package singleton is process-local.
- There is no public registry object by design.
- If a caller needs isolation, add explicit test seams or a separate package-level lifecycle before broadening the public API.
- Keep business policy outside this package. Pass prompts, models, provider IDs, and reasoning effort through settings.

