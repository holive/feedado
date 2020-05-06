package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config is where all the config logic is.
type Config struct {
	viper   *viper.Viper
	order   []string
	profile map[string]string
}

// GetDuration returns the duration value of a key.
func (c *Config) GetDuration(key string) (time.Duration, error) {
	s := c.GetString(key)
	return time.ParseDuration(s)
}

// GetString return the value of a key.
func (c *Config) GetString(key string) string {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}
		value := c.viper.GetString(pk)
		return c.replaceProfileVars(value)
	}
	return ""
}

// GetInt return the value of a key.
func (c *Config) GetInt(key string) int {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}
		return c.viper.GetInt(pk)
	}
	return 0
}

// GetIntSlice return the value of a key.
func (c *Config) GetIntSlice(key string) []int {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}

		rawValues, ok := c.viper.Get(key).([]interface{})
		if !ok {
			return []int{}
		}

		values := make([]int, len(rawValues))
		for i, rawValue := range rawValues {
			val, ok := rawValue.(int64)
			if !ok {
				return []int{}
			}

			values[i] = int(val)
		}

		return values
	}
	return []int{}
}

// GetBool return the value of a key.
func (c *Config) GetBool(key string) bool {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}
		return c.viper.GetBool(pk)
	}
	return false
}

// GetStringSlice return the value of a key.
func (c *Config) GetStringSlice(key string) []string {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}
		values := c.viper.GetStringSlice(pk)
		for i := range values {
			values[i] = c.replaceProfileVars(values[i])
		}
		return values
	}
	return []string{}
}

// GetRawKey return the value of a key.
func (c *Config) GetRawKey(key string) interface{} {
	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}

		return c.viper.Get(pk)
	}

	return nil
}

// GetMapStringString return the value of a map of keys.
func (c *Config) GetMapStringString(key string) map[string]string {
	result := make(map[string]string)

	for _, pk := range c.genCombinations(key) {
		if !c.viper.IsSet(pk) {
			continue
		}

		var rawValues []map[string]interface{}
		err := c.viper.UnmarshalKey(pk, &rawValues)
		if err != nil {
			return nil
		}

		for _, rawValue := range rawValues {
			key, ok := rawValue["key"].(string)
			if !ok {
				return nil
			}
			value, ok := rawValue["value"].(string)
			if !ok {
				return nil
			}

			if _, ok = result[key]; ok {
				continue
			}

			result[key] = c.replaceProfileVars(value)
		}
	}

	return result
}

// GetSliceMapStringInterface return the value of a slice of map[string]interface.
func (c *Config) GetSliceMapStringInterface(key string) ([]map[string]interface{}, error) {
	content := make([]map[string]interface{}, 0)
	err := c.viper.UnmarshalKey(key, &content)
	return content, errors.Wrap(err, fmt.Sprintf("error during key '%s' unmarshal", key))
}

// GetStringMapString return a map[string]string of a key.
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

// With this method we generate all the combinations that a given key can have. And we follow the
// order of less to most significant. Take this order and key as example:
//
//	order: []string{"aws-va", "production"}
//	key: log.format
//
// The method generate this keys in the following order:
//
//	aws-va.production.log.format
//	production.log.format
//	aws-va.log.format
//	log.format
//
// The first key is a full combination, then we follow the order and then we use the original key.
func (c *Config) genCombinations(key string) []string {
	fragments := make([]string, len(c.order))
	for i, order := range c.order {
		fragments[i] = c.profile[order]
	}

	var keys []string
	keys = append(keys, strings.Join(append(fragments, key), "."))
	for i := len(c.order) - 1; i >= 0; i-- {
		preKey := []string{fragments[i], key}
		keys = append(keys, strings.Join(preKey, "."))
	}

	return append(keys, key)
}

func (c *Config) replaceProfileVars(str string) string {
	for key, value := range c.profile {
		r := regexp.MustCompile(fmt.Sprintf(`\${\s*?(%s)\s*}`, key))
		str = r.ReplaceAllString(str, value)
	}
	return str
}

func (c *Config) loadProfile(name string) error {
	var rawProfiles []map[string]interface{}
	err := c.viper.UnmarshalKey("profiles", &rawProfiles)
	if err != nil {
		return err
	}

	for _, rawProfile := range rawProfiles {
		profile, err := c.convertMap(rawProfile)
		if err != nil {
			return err
		}

		if profile["name"] == name {
			c.profile = profile
			break
		}
	}

	return nil
}

func (c *Config) convertMap(source map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)

	for key, rawValue := range source {
		value, ok := rawValue.(string)
		if !ok {
			return result, errors.New("value expected to be a string")
		}
		result[key] = value
	}

	return result, nil
}

func (c *Config) loadMerge() {
	c.order = c.viper.GetStringSlice("merge.order")
}

// New returns a configured client to do config processing.
func New(profile string, content string) (Config, error) {
	v := viper.New()
	v.SetConfigType("toml")

	if err := v.ReadConfig(bytes.NewBufferString(content)); err != nil {
		return Config{}, errors.Wrap(err, "error during viper setup")
	}

	c := Config{viper: v}
	c.loadMerge()

	if err := c.loadProfile(profile); err != nil {
		return c, errors.Wrap(err, "error during config client initialization")
	}

	return c, nil
}
