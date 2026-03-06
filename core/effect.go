package generate

func effectIcon(t EffectType) string {
	switch t {
	case EffectWet:
		return "≈"
	default:
		return "?"
	}
}

func hasEffect(effects []Effect, t EffectType) bool {
	for _, e := range effects {
		if e.Type == t {
			return true
		}
	}
	return false
}

func addEffect(effects []Effect, e Effect) []Effect {
	for i, ex := range effects {
		if ex.Type == e.Type {
			effects[i].Duration = e.Duration
			return effects
		}
	}
	return append(effects, e)
}

func tickEffects(effects []Effect) []Effect {
	result := effects[:0]
	for _, e := range effects {
		e.Duration--
		if e.Duration > 0 {
			result = append(result, e)
		}
	}
	return result
}
