package generate

func effectIcon(t EffectType) string {
	switch t {
	case EffectWet:
		return "≈"
	case EffectFire:
		return "⽕"
	case EffectSteam:
		return "~"
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

func removeEffect(effects []Effect, t EffectType) []Effect {
	out := make([]Effect, 0, len(effects))
	for _, e := range effects {
		if e.Type != t {
			out = append(out, e)
		}
	}
	return out
}

func resolveEffects(effects []Effect, new Effect) []Effect {
	switch new.Type {
	case EffectFire:
		effects = removeEffect(effects, EffectWet)
	case EffectWet:
		effects = removeEffect(effects, EffectFire)
	case EffectSteam:
		effects = removeEffect(effects, EffectFire)
		effects = removeEffect(effects, EffectWet)
	}
	return addEffect(effects, new)
}

func tickEffects(effects []Effect) []Effect {
	out := make([]Effect, 0, len(effects))
	for _, e := range effects {
		e.Duration--
		if e.Duration > 0 {
			out = append(out, e)
		}
	}
	return out
}
