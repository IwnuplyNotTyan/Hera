package tests

import (
	"testing"

	generate "hera/core"

	"github.com/stretchr/testify/assert"
)

func TestEffectIcon(t *testing.T) {
	assert.Equal(t, "≈", generate.EffectIcon(generate.EffectWet))
	assert.Equal(t, "⚹", generate.EffectIcon(generate.EffectFire))
	assert.Equal(t, "~", generate.EffectIcon(generate.EffectSmoke))
	assert.Equal(t, "?", generate.EffectIcon(generate.EffectType("unknown")))
}

func TestHasEffect_True(t *testing.T) {
	effects := []generate.Effect{
		{Type: generate.EffectFire, Duration: 2},
		{Type: generate.EffectWet, Duration: 1},
	}
	assert.True(t, generate.HasEffect(effects, generate.EffectFire))
	assert.True(t, generate.HasEffect(effects, generate.EffectWet))
}

func TestHasEffect_False(t *testing.T) {
	effects := []generate.Effect{
		{Type: generate.EffectFire, Duration: 2},
	}
	assert.False(t, generate.HasEffect(effects, generate.EffectSmoke))
}

func TestHasEffect_Empty(t *testing.T) {
	assert.False(t, generate.HasEffect(nil, generate.EffectFire))
	assert.False(t, generate.HasEffect([]generate.Effect{}, generate.EffectWet))
}

func TestAddEffect_New(t *testing.T) {
	effects := generate.AddEffect(nil, generate.Effect{Type: generate.EffectFire, Duration: 2})
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectFire, effects[0].Type)
	assert.Equal(t, 2, effects[0].Duration)
}

func TestAddEffect_RefreshDuration(t *testing.T) {
	effects := []generate.Effect{{Type: generate.EffectFire, Duration: 1}}
	effects = generate.AddEffect(effects, generate.Effect{Type: generate.EffectFire, Duration: 3})
	assert.Len(t, effects, 1)
	assert.Equal(t, 3, effects[0].Duration)
}

func TestRemoveEffect_RemovesMatching(t *testing.T) {
	effects := []generate.Effect{
		{Type: generate.EffectFire, Duration: 2},
		{Type: generate.EffectWet, Duration: 1},
	}
	effects = generate.RemoveEffect(effects, generate.EffectFire)
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectWet, effects[0].Type)
}

func TestRemoveEffect_NoMatch(t *testing.T) {
	effects := []generate.Effect{{Type: generate.EffectFire, Duration: 2}}
	effects = generate.RemoveEffect(effects, generate.EffectSmoke)
	assert.Len(t, effects, 1)
}

func TestRemoveEffect_Empty(t *testing.T) {
	effects := generate.RemoveEffect(nil, generate.EffectFire)
	assert.Len(t, effects, 0)
}

func TestResolveEffects_FireRemovesWet(t *testing.T) {
	effects := []generate.Effect{{Type: generate.EffectWet, Duration: 2}}
	effects = generate.ResolveEffects(effects, generate.Effect{Type: generate.EffectFire, Duration: 2})
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectFire, effects[0].Type)
}

func TestResolveEffects_WetRemovesFire(t *testing.T) {
	effects := []generate.Effect{{Type: generate.EffectFire, Duration: 2}}
	effects = generate.ResolveEffects(effects, generate.Effect{Type: generate.EffectWet, Duration: 2})
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectWet, effects[0].Type)
}

func TestResolveEffects_SmokeRemovesAll(t *testing.T) {
	effects := []generate.Effect{
		{Type: generate.EffectFire, Duration: 2},
		{Type: generate.EffectWet, Duration: 1},
	}
	effects = generate.ResolveEffects(effects, generate.Effect{Type: generate.EffectSmoke, Duration: 3})
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectSmoke, effects[0].Type)
	assert.Equal(t, 3, effects[0].Duration)
}

func TestResolveEffects_NewEffect(t *testing.T) {
	effects := generate.ResolveEffects(nil, generate.Effect{Type: generate.EffectFire, Duration: 2})
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectFire, effects[0].Type)
}

func TestTickEffects_Decrements(t *testing.T) {
	effects := []generate.Effect{
		{Type: generate.EffectFire, Duration: 2},
		{Type: generate.EffectWet, Duration: 1},
	}
	effects = generate.TickEffects(effects)
	assert.Len(t, effects, 1)
	assert.Equal(t, generate.EffectFire, effects[0].Type)
	assert.Equal(t, 1, effects[0].Duration)
}

func TestTickEffects_RemovesExpired(t *testing.T) {
	effects := []generate.Effect{{Type: generate.EffectFire, Duration: 1}}
	effects = generate.TickEffects(effects)
	assert.Len(t, effects, 0)
}

func TestTickEffects_Empty(t *testing.T) {
	effects := generate.TickEffects(nil)
	assert.Len(t, effects, 0)
}
