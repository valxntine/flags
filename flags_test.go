package flags

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

const (
	jsonFlagFileName    = "flags.goff.json"
	yamlFlagFileName    = "flags.goff.yaml"
	notExistsFlagName   = "not-exists"
	floatFlagName       = "ff-float"
	numberFlagName      = "ff-number"
	descriptionFlagName = "ff-description"
	isEnabledFlagName   = "is-enabled"
	timeFlagName        = "cr-start"
	jsonFlagName        = "ff-json"
	enabledByIDFlagName = "is-enabled-for-user"
	idListIntFlag       = "ff-json-list"
	idListStringFlag    = "ff-json-list-string"
)

func setupClient(t *testing.T, filename string) *ffclient.GoFeatureFlag {
	t.Helper()
	fileType := "yaml"
	if strings.Contains(filename, "json") {
		fileType = "json"
	}
	c, err := ffclient.New(ffclient.Config{
		PollingInterval: 10 * time.Minute,
		Retrievers: []retriever.Retriever{
			&fileretriever.Retriever{Path: filename},
		},
		FileFormat: fileType,
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}
	return c
}

func TestGetFloat(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected float64
		flag     string
	}{
		{"flag exists, returns value", "1", 3.14159, floatFlagName},
		{"flag doesnt exist, returns default", "1", 1.11, notExistsFlagName},
		{"no user id, still returns flag value", "", 3.14159, floatFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			f, err := GetFloat(tt.flag, tt.userID, 1.11, c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if f != tt.expected {
				t.Errorf("unexpected float: got %f want %f", f, tt.expected)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected int
		flag     string
	}{
		{"flag exists, returns value", "1", 9081, numberFlagName},
		{"flag doesnt exist, returns default", "1", 69, notExistsFlagName},
		{"no user id, still returns flag value", "", 9081, numberFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			i, err := GetInt(tt.flag, tt.userID, 69, c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if i != tt.expected {
				t.Errorf("unexpected float: got %d want %d", i, tt.expected)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected string
		flag     string
	}{
		{"flag exists, returns value", "1", "Something about chocolate eggs", descriptionFlagName},
		{"flag doesnt exist, returns default", "1", "hello", notExistsFlagName},
		{"no user id, still returns flag value", "", "Something about chocolate eggs", descriptionFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			s, err := GetString(tt.flag, tt.userID, "hello", c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if s != tt.expected {
				t.Errorf("unexpected float: got %s want %s", s, tt.expected)
			}
		})
	}
}

func TestGetJSONMap(t *testing.T) {
	defaultResponseTimes := map[string]any{
		"p50":     50,
		"p75":     75,
		"p95":     95,
		"p99":     99,
		"p99_5":   995,
		"p99_999": 99999,
		"default": 0,
	}

	expectedResponseTimes := map[string]any{
		"p50":     40,
		"p75":     50,
		"p95":     70,
		"p99":     150,
		"p99_5":   225,
		"p99_999": 500,
		"default": 1200,
	}
	tests := []struct {
		name     string
		userID   string
		expected map[string]any
		flag     string
	}{
		{"flag exists, returns value", "1", expectedResponseTimes, jsonFlagName},
		{"flag doesnt exist, returns default", "1", defaultResponseTimes, notExistsFlagName},
		{"no user id, still returns flag value", "", expectedResponseTimes, jsonFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			j, err := GetJSONMap(
				tt.flag,
				tt.userID,
				defaultResponseTimes,
				c,
			)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if diff := cmp.Diff(j, tt.expected); diff != "" {
				t.Errorf("unexpected struct (-got +want)\n%s", diff)
			}
		})
	}
}

func TestGetJSONStruct(t *testing.T) {
	type ResponseTimes struct {
		P50     int `json:"p50"`
		P75     int `json:"p75"`
		P95     int `json:"p95"`
		P99     int `json:"p99"`
		P99_5   int `json:"p99_5"`
		P99_999 int `json:"p99_999"`
		Def     int `json:"default"`
	}

	defaultResponseTimes := ResponseTimes{
		P50:     50,
		P75:     75,
		P95:     95,
		P99:     99,
		P99_5:   995,
		P99_999: 99999,
		Def:     0,
	}

	expectedResponseTimes := ResponseTimes{
		P50:     40,
		P75:     50,
		P95:     70,
		P99:     150,
		P99_5:   225,
		P99_999: 500,
		Def:     1200,
	}
	tests := []struct {
		name     string
		userID   string
		expected ResponseTimes
		flag     string
	}{
		{"flag exists, returns value", "1", expectedResponseTimes, jsonFlagName},
		{"flag doesnt exist, returns default", "1", defaultResponseTimes, notExistsFlagName},
		{"no user id, still returns flag value", "", expectedResponseTimes, jsonFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			j, err := GetJSONStruct[ResponseTimes](
				tt.flag,
				tt.userID,
				defaultResponseTimes,
				c,
			)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if diff := cmp.Diff(j, tt.expected); diff != "" {
				t.Errorf("unexpected struct (-got +want)\n%s", diff)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	expectedTime, err := time.Parse(time.RFC3339, "2025-07-18T22:37:22.176Z")
	if err != nil {
		t.Fatalf("failed to parse expected time: %v", err)
	}

	defaultTime, err := time.Parse(time.RFC3339, "2099-12-31T23:59:59.999Z")
	if err != nil {
		t.Fatalf("failed to parse default time: %v", err)
	}
	tests := []struct {
		name     string
		userID   string
		expected time.Time
		flag     string
	}{
		{"flag exists, returns value", "1", expectedTime, timeFlagName},
		{"flag doesnt exist, returns default", "1", defaultTime, notExistsFlagName},
		{"no user id, still returns flag value", "", expectedTime, timeFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			ti, err := GetTime(tt.flag, tt.userID, time.RFC3339, defaultTime, c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if ti != tt.expected {
				t.Errorf("unexpected float: got %s want %s", ti, tt.expected)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected bool
		flag     string
	}{
		{"flag exists, returns value", "1", true, isEnabledFlagName},
		{"flag doesnt exist, returns default", "1", false, notExistsFlagName},
		{"no user id, still returns flag value", "", true, isEnabledFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			b, err := IsEnabled(tt.flag, tt.userID, false, c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if b != tt.expected {
				t.Errorf("unexpected float: got %t want %t", b, tt.expected)
			}
		})
	}
}

func TestIsEnabledByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		ctxUserID string
		expected  bool
		flag      string
	}{
		{"flag exists, user id in list, returns enabled", "1", "1", true, enabledByIDFlagName},
		{"flag exists, user id not in list, returns disabled", "4", "4", false, enabledByIDFlagName},
		{"flag doesnt exist, returns default", "1", "1", false, notExistsFlagName},
		{"no user id in context eval, user id is in list, returns enabled", "1", "", true, enabledByIDFlagName},
		{"no user id in context eval, user id not in list, returns disabled", "4", "", false, enabledByIDFlagName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t, yamlFlagFileName)
			b, err := IsEnabledByID(tt.flag, tt.ctxUserID, tt.userID, "user-id", false, c)
			if tt.flag == notExistsFlagName {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error getting flag value: %v", err)
				}
			}

			if b != tt.expected {
				t.Errorf("unexpected float: got %t want %t", b, tt.expected)
			}
		})
	}
}

func TestIsEnabledByIDList(t *testing.T) {
	t.Run("id is an int", func(t *testing.T) {
		tests := []struct {
			name      string
			lookupID  int
			ctxUserID string
			expected  bool
			flag      string
		}{
			{"flag exists, user id in list, returns enabled", 1, "1", true, idListIntFlag},
			{"flag exists, user id not in list, returns disabled", 4, "4", false, idListIntFlag},
			{"flag doesnt exist, returns default", 1, "1", false, notExistsFlagName},
			{"no user id in context eval, user id is in list, returns enabled", 1, "", true, idListIntFlag},
			{"no user id in context eval, user id not in list, returns disabled", 4, "", false, idListIntFlag},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				c := setupClient(t, yamlFlagFileName)
				b, err := IsEnabledByIDList(tt.flag, tt.ctxUserID, tt.lookupID, false, c)
				if tt.flag == notExistsFlagName {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("unexpected error getting flag value: %v", err)
					}
				}

				if b != tt.expected {
					t.Errorf("unexpected float: got %t want %t", b, tt.expected)
				}
			})
		}
	})
	t.Run("id is a string", func(t *testing.T) {
		tests := []struct {
			name      string
			lookupID  string
			ctxUserID string
			expected  bool
			flag      string
		}{
			{"flag exists, user id in list, returns enabled", "1", "1", true, idListStringFlag},
			{"flag exists, user id not in list, returns disabled", "4", "4", false, idListStringFlag},
			{"flag doesnt exist, returns default", "1", "1", false, notExistsFlagName},
			{"no user id in context eval, user id is in list, returns enabled", "1", "", true, idListStringFlag},
			{"no user id in context eval, user id not in list, returns disabled", "4", "", false, idListStringFlag},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				c := setupClient(t, yamlFlagFileName)
				b, err := IsEnabledByIDList(tt.flag, tt.ctxUserID, tt.lookupID, false, c)
				if tt.flag == notExistsFlagName {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("unexpected error getting flag value: %v", err)
					}
				}

				if b != tt.expected {
					t.Errorf("unexpected float: got %t want %t", b, tt.expected)
				}
			})
		}
	})
}

func TestRefresh(t *testing.T) {
	t.Run("refresh is called and refreshed time is updated", func(t *testing.T) {
		c := setupClient(t, "flags.goff.yaml")
		defer c.Close()

		tmp, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatalf("unexpected error creating temp file: %v", err)
		}
		defer func() {
			_ = os.Remove(tmp.Name())
		}()

		flags, err := os.ReadFile("flags.goff.yaml")
		if err != nil {
			t.Fatalf("unexpected error reading flags: %v", err)
		}

		err = os.WriteFile(tmp.Name(), flags, os.ModePerm)
		if err != nil {
			t.Fatalf("unexpected error writing to temp file: %v", err)
		}

		refreshTime := c.GetCacheRefreshDate()

		newFlag := `ff-test:
  metadata:
    description: test
  variations:
    e: 2.71828
  defaultRule:
    variation: e
`

		f, err := os.OpenFile(tmp.Name(), os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			t.Fatalf("failed to open file for appending: %v", err)
		}
		defer f.Close()

		_, err = f.WriteString(newFlag)
		if err != nil {
			t.Fatalf("failed to write new flag to file: %v", err)
		}

		if refreshTime != c.GetCacheRefreshDate() {
			t.Fatal("cache refreshed unexpectedly")
		}

		Refresh(c)

		if refreshTime == c.GetCacheRefreshDate() {
			t.Error("expected refresh time to have updated")
		}
	})
}

func TestNewClient(t *testing.T) {
	t.Run("flag file isnt available", func(t *testing.T) {
		// because we're mixing singleton and `New` in this file,
		// we need to ensure the instance is closed before this test
		ffclient.Close()
		err := NewClient(Config{
			PollingInterval: 10 * time.Second,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{
					Path: "non-existent.goff.yaml",
				},
			},
		})

		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})
	t.Run("no retrievers passed in", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Second,
		})

		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})
}
