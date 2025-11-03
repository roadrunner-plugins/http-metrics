package prometheus

import (
	"regexp"
	"sync"
)

// EndpointMatcher matches request paths to endpoint patterns
// Uses LRU cache to minimize regex matching overhead
type EndpointMatcher struct {
	rules    []compiledRule
	cache    *endpointCache
	fallback string
	mu       sync.RWMutex
}

type compiledRule struct {
	regex *regexp.Regexp
	name  string
}

// NewEndpointMatcher creates a new endpoint matcher from configuration
func NewEndpointMatcher(config EndpointPatternsConfig) (*EndpointMatcher, error) {
	if !config.Enabled {
		return &EndpointMatcher{
			fallback: "other",
			cache:    newEndpointCache(1),
		}, nil
	}

	rules := make([]compiledRule, 0, len(config.Rules))
	for _, rule := range config.Rules {
		regex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, err
		}
		rules = append(rules, compiledRule{
			regex: regex,
			name:  rule.Name,
		})
	}

	cacheSize := config.CacheSize
	if cacheSize <= 0 {
		cacheSize = 10000
	}

	return &EndpointMatcher{
		rules:    rules,
		cache:    newEndpointCache(cacheSize),
		fallback: "other",
	}, nil
}

// Match returns the endpoint pattern for a given path
// Uses cached results when available, falls back to regex matching
func (em *EndpointMatcher) Match(path string) string {
	// Check cache first (read lock)
	if cached, ok := em.cache.Get(path); ok {
		return cached
	}

	// Try rules
	em.mu.RLock()
	defer em.mu.RUnlock()

	for _, rule := range em.rules {
		if rule.regex.MatchString(path) {
			em.cache.Set(path, rule.name)
			return rule.name
		}
	}

	// No match - use fallback
	em.cache.Set(path, em.fallback)
	return em.fallback
}

// endpointCache is a simple LRU cache for endpoint matching results
type endpointCache struct {
	capacity int
	items    map[string]*cacheNode
	head     *cacheNode
	tail     *cacheNode
	mu       sync.RWMutex
}

type cacheNode struct {
	key   string
	value string
	prev  *cacheNode
	next  *cacheNode
}

func newEndpointCache(capacity int) *endpointCache {
	return &endpointCache{
		capacity: capacity,
		items:    make(map[string]*cacheNode, capacity),
	}
}

func (c *endpointCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if node, ok := c.items[key]; ok {
		return node.value, true
	}
	return "", false
}

func (c *endpointCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key exists
	if node, ok := c.items[key]; ok {
		node.value = value
		c.moveToFront(node)
		return
	}

	// Create new node
	node := &cacheNode{
		key:   key,
		value: value,
	}

	// Add to cache
	if c.head == nil {
		c.head = node
		c.tail = node
	} else {
		node.next = c.head
		c.head.prev = node
		c.head = node
	}

	c.items[key] = node

	// Evict if over capacity
	if len(c.items) > c.capacity {
		c.removeTail()
	}
}

func (c *endpointCache) moveToFront(node *cacheNode) {
	if node == c.head {
		return
	}

	// Remove from current position
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == c.tail {
		c.tail = node.prev
	}

	// Move to front
	node.prev = nil
	node.next = c.head
	c.head.prev = node
	c.head = node
}

func (c *endpointCache) removeTail() {
	if c.tail == nil {
		return
	}

	delete(c.items, c.tail.key)

	if c.tail.prev != nil {
		c.tail.prev.next = nil
		c.tail = c.tail.prev
	} else {
		c.head = nil
		c.tail = nil
	}
}
