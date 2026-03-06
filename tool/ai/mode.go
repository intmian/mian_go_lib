package ai

type ModelMode string

const (
	ModelModeCheap  ModelMode = "cheap"
	ModelModeFast   ModelMode = "fast"
	ModelModeNormal ModelMode = "normal"
)

func NormalizeModelMode(raw string, fallback ModelMode) ModelMode {
	mode := ModelMode(raw)
	switch mode {
	case ModelModeCheap, ModelModeFast, ModelModeNormal:
		return mode
	default:
		return fallback
	}
}

func SelectModels(mode ModelMode, modelPools map[ModelMode][]string) []string {
	if modelPools == nil {
		return nil
	}
	models := modelPools[mode]
	if len(models) == 0 {
		return nil
	}
	return models
}
