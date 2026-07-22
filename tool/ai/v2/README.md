# AI v2

`tool/ai/v2` is the newer AI package for provider/model capability checks, chat, streaming chat, OpenAI web search, and lightweight stateful agents.

The directory name is `v2`, but the Go package name is `ai`:

```go
import ai "github.com/intmian/mian_go_lib/tool/ai/v2"
```

This package does not replace the older `github.com/intmian/mian_go_lib/tool/ai` wrapper yet. Use it when you need provider capabilities, reasoning controls, OpenAI Responses tools, streaming, base agent components, or agent setting JSON helpers.

## Core concepts

1. `ProviderID` and `AgentID` are `uint32`.
   - `0` means unset and is rejected.
   - Register providers once with `AddProvider`.
   - `OpenAIProvider.RegisterModel` registers model IDs, upstream names, model types, optional call-protocol metadata, reasoning efforts, and tools on that provider instance.
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
	if err := provider.RegisterModel(ai.ModelConfig{
		ID:        "main",
		Name:      "gpt-5.4",
		Type:      ai.ModelTypeText,
		Reasoning: []ai.ReasoningEffort{ai.ReasoningEffortMedium},
	}); err != nil {
		return "", err
	}
	if err := ai.AddProvider(1, provider); err != nil {
		return "", err
	}

	agent := ai.NewBaseAgent()
	if err := agent.InitWithSetting(&ai.BaseAgentSetting{
		ProviderID:      1,
		SysPrompt:       "You are a helpful assistant.",
		Models:          []string{"main"},
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

`Models` is an ordered model fallback list, not a mode name list. For providers with registered models, entries are registered model IDs; providers without a model catalog retain the legacy direct upstream-name behavior.

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

`AvailableReasoning` and `AvailableTools` describe the adapter's upper bound. Register each configured model with a `ModelConfig`; registered chat calls resolve the model ID to its upstream name and reject reasoning efforts or tools that the model did not declare. `ProviderSourceTypeOpenAI` can register `web_search`; the current DeepSeek adapter cannot.

```go
err := openaiProvider.RegisterModel(ai.ModelConfig{
	ID:        "research",
	Name:      "gpt-5.4",
	Type:      ai.ModelTypeText,
	Reasoning: []ai.ReasoningEffort{ai.ReasoningEffortLow, ai.ReasoningEffortMedium},
	Tools:     []ai.ChatTool{ai.ChatToolWebSearch},
})
model, ok := openaiProvider.GetModel("research")
registered := openaiProvider.RegisteredModels()
```

`ModelTypeText` and `ModelTypeSTT` are supported configuration types. Reasoning and chat tools only apply to text models. STT registration is capability/config metadata; audio processing still belongs to the owning caller until a matching v2 runtime surface is added. `ModelTypeChat` remains a source-compatible alias of `ModelTypeText`.

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

Enable hosted OpenAI web search by adding `ChatToolWebSearch`. `MaxToolCalls` is optional and must only be set when at least one tool is present:

```go
resp, err := provider.Chat(ctx, ai.ChatRequest{
	Model: "gpt-5.4",
	Messages: []ai.ChatMessage{
		{Role: ai.ChatRoleUser, Content: "Search the web for the latest OpenAI API changes."},
	},
	Tools:        []ai.ChatTool{ai.ChatToolWebSearch},
	MaxToolCalls: 3,
})
if err != nil {
	return err
}
for _, citation := range resp.Citations {
	fmt.Printf("%s: %s\n", citation.Title, citation.URL)
}
```

`ChatResponse.Citations` contains URL citations returned in OpenAI output-text annotations, including the upstream character offsets, title, and URL. Streaming requests expose the same citations on the final `ChatStreamEventDone` response.

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
		// event.Response.Citations contains final URL citations.
	}
}
```

## Manual OpenAI web search test

The live test is opt-in and is excluded from normal test runs. It reads credentials from environment variables and never prints the API key:

- required: `OPENAI_API_KEY`
- optional: `OPENAI_BASE_URL` (empty uses the official OpenAI endpoint)
- optional: `OPENAI_MODEL` (default `gpt-5.4`)
- optional: `OPENAI_REASONING_EFFORT` (default `low`)

From the `mian_go_lib` module root:

```powershell
$secure = Read-Host "OpenAI API Key" -AsSecureString
$env:OPENAI_API_KEY = [System.Net.NetworkCredential]::new("", $secure).Password
$env:OPENAI_MODEL = "gpt-5.4"
$env:OPENAI_REASONING_EFFORT = "low"
go test -tags live_ai -run '^TestLiveOpenAIWebSearch$' -count=1 -v ./tool/ai/v2
Remove-Item Env:OPENAI_API_KEY
```

The test performs a real, billable API request and requires at least one URL citation in the response. It keeps the configured reasoning effort while omitting `max_tool_calls` and `max_output_tokens`, which are not required for the search check and may be unsupported by some OpenAI-compatible gateways.

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
- Registered model IDs and names are required; duplicate model IDs are rejected per provider.
- A registered model may only declare capabilities supported by its provider adapter.
- When a provider has registered models, chat requests must use a registered chat model ID and may only request that model's declared reasoning efforts and tools.
- Duplicate agent setting IDs are rejected.
- `BaseAgentSetting.Models` must contain at least one non-empty model after trimming whitespace.
- `Init` requires a registered `*BaseAgentSetting` at `AgentIDBase`.
- `InitWithSetting` requires `ProviderID` and at least one model.
- The referenced provider must be registered before provider access or chat, not before initialization.
- `MaxToolCalls` cannot be negative and requires at least one `ChatRequest.Tools` entry when non-zero.
- Duplicate or provider-unsupported chat tools are rejected before an upstream request.

Provider lookup happens at chat time. This allows providers to be registered after agent construction, but before the first successful chat.

## Package notes

- The package singleton is process-local.
- The package singleton stores providers and agent settings, but it does not contain agent factories.
- If a caller needs isolation, add explicit test seams or a separate package-level lifecycle before broadening the public API.
- Keep business policy outside this package. Pass prompts, models, provider IDs, and reasoning effort through settings.
