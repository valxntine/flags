package flags

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffcontext"
	"github.com/thomaspoignant/go-feature-flag/retriever"
	"golang.org/x/exp/slices"
)

type Config struct {
	PollingInterval time.Duration
	Retrievers      []retriever.Retriever
	FileFormat      string
}

func NewClient(cfg Config) error {
	if cfg.Retrievers == nil {
		return fmt.Errorf("ffclient expects at least 1 retriever")
	}

	format := "yaml"
	if cfg.FileFormat != "" {
		format = cfg.FileFormat
	}
	err := ffclient.Init(ffclient.Config{
		PollingInterval: cfg.PollingInterval,
		Retrievers:      cfg.Retrievers,
		FileFormat:      format,
	})
	if err != nil {
		return fmt.Errorf("failed to init goff: %v", err)
	}
	return nil
}

func Close() {
	ffclient.Close()
}

func IsEnabledByID(flag, userID, id, lookup string, defaultValue bool) (bool, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContextBuilder(userID).AddCustom(lookup, id).Build()
	return ffclient.BoolVariation(flag, c, defaultValue)
}

func IsEnabled(flag, userID string, defaultValue bool) (bool, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	return ffclient.BoolVariation(flag, c, defaultValue)
}

func GetTime(flag, userID, layout string, defaultValue time.Time) (time.Time, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	s, err := ffclient.StringVariation(flag, c, defaultValue.Format(layout))
	if err != nil {
		return defaultValue, fmt.Errorf("failed to get flag %s: %w", flag, err)
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to parse time %s into layout %s: %w", s, layout, err)
	}
	return t, nil
}

func GetInt(flag, userID string, defaultValue int) (int, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	return ffclient.IntVariation(flag, c, defaultValue)
}

func GetFloat(flag, userID string, defaultValue float64) (float64, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	return ffclient.Float64Variation(flag, c, defaultValue)
}

func GetString(flag, userID string, defaultValue string) (string, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	return ffclient.StringVariation(flag, c, defaultValue)
}

func GetJSONStruct[T any](flag, userID string, defaultValue T) (T, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)

	defaultBytes, err := json.Marshal(defaultValue)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to marshal default value: %w", err)
	}

	var defaultMap map[string]any
	if err = json.Unmarshal(defaultBytes, &defaultMap); err != nil {
		return defaultValue, fmt.Errorf("failed to unmarshal default value to map: %w", err)
	}

	j, err := ffclient.JSONVariation(flag, c, defaultMap)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to get flag %s: %w", flag, err)
	}

	result, err := json.Marshal(j)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to marshal result to target: %w", err)
	}

	var v T
	if err = json.Unmarshal(result, &v); err != nil {
		return defaultValue, fmt.Errorf("failed to unmarshal flag to target: %w", err)
	}

	return v, nil
}

func GetJSONMap(flag, userID string, defaultValue map[string]any) (map[string]any, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	return ffclient.JSONVariation(flag, c, defaultValue)
}

func IsEnabledByIDList[T comparable](flag, userID string, lookup T, defaultValue bool) (bool, error) {
	if userID == "" {
		userID = "anonymous"
	}
	c := ffcontext.NewEvaluationContext(userID)
	l, err := ffclient.JSONArrayVariation(flag, c, []any{})
	if err != nil {
		return defaultValue, err
	}
	t := reflect.TypeOf(lookup)
	if slices.ContainsFunc(l, func(i any) bool {
		// assuming ID's are always ints or strings, so convert json numbers to int
		if f, fOk := i.(float64); fOk {
			i = int(f)
		}
		v, ok := i.(T)
		if !ok {
			return false
		}
		if v == lookup {
			return true
		}
		return false
	}) {
		return true, nil
	}
	return false, nil
}

func Refresh() {
	ffclient.ForceRefresh()
}
