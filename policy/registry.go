package policy

import "sync"

// Registry stores named retry policies that can be looked up by route key.
type Registry struct {
	mu       sync.RWMutex
	policies map[string]*Policy
	fallback *Policy
}

// NewRegistry creates a Registry with the given fallback policy.
// If fallback is nil, DefaultPolicy is used.
func NewRegistry(fallback *Policy) *Registry {
	if fallback == nil {
		fallback = DefaultPolicy()
	}
	return &Registry{
		policies: make(map[string]*Policy),
		fallback: fallback,
	}
}

// Register associates a named policy with a route key.
func (r *Registry) Register(key string, p *Policy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[key] = p
}

// Lookup returns the policy for the given key.
// If no policy is registered for the key, the fallback policy is returned.
func (r *Registry) Lookup(key string) *Policy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.policies[key]; ok {
		return p
	}
	return r.fallback
}

// Remove deletes a registered policy by key.
func (r *Registry) Remove(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.policies, key)
}

// SetFallback updates the fallback policy used when no key matches.
func (r *Registry) SetFallback(p *Policy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = p
}
