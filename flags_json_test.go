package flags

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetFloatJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestGetIntJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestGetStringJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestGetJSONMapJSON(t *testing.T) {
	defaultResponseTimes := map[string]any{
		"p50":     50,
		"p75":     75,
		"p95":     95,
		"p99":     99,
		"p99_5":   995,
		"p99_999": 99999,
		"default": 0,
	}

	//TODO: JSON treats numbers as floats by default
	// what can we do to mitigate this - if you're choosing to
	// just get a map[string]any then you might expect ints in the value
	// address with the team
	expectedResponseTimes := map[string]any{
		"p50":     40.,
		"p75":     50.,
		"p95":     70.,
		"p99":     150.,
		"p99_5":   225.,
		"p99_999": 500.,
		"default": 1200.,
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
			c := setupClient(t, jsonFlagFileName)
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

func TestGetJSONStructJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestGetTimeJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestIsEnabledJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestIsEnabledByIDJSON(t *testing.T) {
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
			c := setupClient(t, jsonFlagFileName)
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

func TestIsEnabledByIDListJSON(t *testing.T) {
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
				c := setupClient(t, jsonFlagFileName)
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
				c := setupClient(t, jsonFlagFileName)
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
