package xconfig

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"sync"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// Parser defines the universal contract for modules capable of reading,
// parsing, and converting external configuration files or key-value stores
// into a standardized Go map.
type Parser interface {
	// SetOptions injects the necessary configuration properties into the parser instance.
	SetOptions(opts Options) error
	// Read executes the reading strategy and returns a generic key-value map.
	Read() (map[string]any, error)
}

// Config manages the internal global state of loaded configuration variables.
// It maintains a thread-safe registry of parsers and a consolidated map of data.
type Config struct {
	mu      sync.RWMutex
	data    map[string]any
	parsers []Parser
}

// NewConfig initializes an empty Config instance with prepared maps and slices,
// ensuring it is fully ready to safely register parsers or look up keys.
func NewConfig() Config {
	return Config{
		data:    make(map[string]any),
		parsers: make([]Parser, 0),
	}
}

// InitAppConfig acts as a bootstrap orchestrator. It receives a collection of
// 'Parser' and 'Options' objects, pairs them together, registers them to a new
// Config instance, and executes the initial data load in a single sequential flow.
//
// Constraints:
//   - The length of both slices must be identical. If not, it fails immediately.
//   - Registration and loading order follows the exact slice position (index 0 to N).
//
// Returns a fully initialized, populated pointer to Config, or nil and an error if it fails.
func InitAppConfig(parsers []Parser, options []Options) (*Config, error) {
	if len(parsers) != len(options) {
		return nil, xerrors.NewErr(
			xerrors.XERR_ASYMMETRIC_SIZES,
			xerror_CTX,
			"",
			xerrors.MsgArraySize("parsers", len(parsers), "options", len(options)),
		)
	}

	cfg := NewConfig()

	if len(parsers) == 0 {
		return &cfg, nil
	}

	for i := range parsers {
		err := cfg.Register(
			parsers[i],
			options[i],
		)
		if err != nil {
			return nil, err
		}
	}

	err := cfg.Load()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Register adds a custom provider/parser and its respective Options to the
// internal execution queue. This function does not trigger the data read operation.
// Returns an error if the underlying parser fails to ingest the provided options.
func (c *Config) Register(p Parser, opts Options) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := p.SetOptions(opts)
	if err != nil {
		return err
	}

	c.parsers = append(c.parsers, p)

	return nil
}

// Load clears the current state and triggers all registered parsers sequentially.
// Keys extracted from subsequent parsers will overwrite existing values with the
// same name, establishing a strict linear override priority matching the registration order.
// All keys are lowercased and stripped of leading/trailing spaces during ingestion to avoid ambiguity.
func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]any)

	for _, rp := range c.parsers {
		extracted, err := rp.Read()
		if err != nil {
			return err
		}

		// Merge data respecting linear priority and standardizing layout
		for k, v := range extracted {
			cleanKey := strings.ToLower(strings.TrimSpace(k))
			c.data[cleanKey] = v
		}
	}
	return nil
}

// Reload completely clears the internal data map and re-executes the ordered
// queue of registered parsers. This is particularly useful for live-reload strategies.
func (c *Config) Reload() error {
	return c.Load()
}

// Keys returns a thread-safe snapshot slice containing all keys currently loaded.
func (c *Config) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}

	return keys
}

// Has verifies the presence of a specific configuration key in a thread-safe manner.
// The search key parameter is normalized (lowercased and trimmed) to guarantee matching consistency.
func (c *Config) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cleanKey := strings.ToLower(strings.TrimSpace(key))
	_, exists := c.data[cleanKey]
	return exists
}

// Get performs a concurrent-safe lookup for a specific key.
// The target key parameter is automatically normalized (lowercased and trimmed).
// Returns the raw value and a boolean indicating whether the key was found.
func (c *Config) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cleanKey := strings.ToLower(strings.TrimSpace(key))
	val, exists := c.data[cleanKey]
	return val, exists
}

// Populate maps the consolidated internal configuration data into a custom struct pointer.
// It leverages standard JSON marshaling/unmarshaling workflows to map generic string/interface
// structures into typed application structures.
//
// Constraints:
//   - 'target' must be a non-nil pointer.
//   - 'target' must point strictly to a struct category object.
//
// Returns an error if validation fails or if the structural decoding cannot be achieved.
func (c *Config) Populate(target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer {
		return xerrors.NewErr(
			xerrors.XERR_INVALID_TYPE,
			xerror_CTX,
			"",
			"target",
			val.Kind().String(),
			reflect.Pointer.String(),
		)
	}
	if val.IsNil() {
		return xerrors.NewErr(
			xerrors.XERR_NIL_NOT_ALLOWED,
			xerror_CTX,
			"",
			"target",
		)
	}
	if val.Elem().Kind() != reflect.Struct {
		return xerrors.NewErr(
			xerrors.XERR_INVALID_TYPE,
			xerror_CTX,
			"",
			"target",
			val.Elem().Kind().String(),
			reflect.Struct.String(),
		)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	jsonData, err := json.Marshal(c.data)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))

	return decoder.Decode(target)
}
