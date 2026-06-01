package generate

func EffectIcon(t EffectType) string {
	switch t {
	case EffectWet:
		return "≈"
	case EffectFire:
		return "⚹"
	case EffectSmoke:
		return "~"
	default:
		return "?"
	}
}

func HasEffect(effects []Effect, t EffectType) bool {
	for _, e := range effects {
		if e.Type == t {
			return true
		}
	}
	return false
}

func AddEffect(effects []Effect, e Effect) []Effect {
	for i, ex := range effects {
		if ex.Type == e.Type {
			effects[i].Duration = e.Duration
			return effects
		}
	}
	return append(effects, e)
}

func RemoveEffect(effects []Effect, t EffectType) []Effect {
	out := make([]Effect, 0, len(effects))
	for _, e := range effects {
		if e.Type != t {
			out = append(out, e)
		}
	}
	return out
}

func ResolveEffects(effects []Effect, new Effect) []Effect {
	switch new.Type {
	case EffectFire:
		effects = RemoveEffect(effects, EffectWet)
	case EffectWet:
		effects = RemoveEffect(effects, EffectFire)
	case EffectSmoke:
		effects = RemoveEffect(effects, EffectFire)
		effects = RemoveEffect(effects, EffectWet)
	}
	return AddEffect(effects, new)
}

func TickEffects(effects []Effect) []Effect {
	out := make([]Effect, 0, len(effects))
	for _, e := range effects {
		e.Duration--
		if e.Duration > 0 {
			out = append(out, e)
		}
	}
	return out
}
