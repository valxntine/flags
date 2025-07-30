package flags

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

func TestGetFloatJSON(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		f, err := GetFloat("ff-float", "1", 1.11, c, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if f != 3.14159 {
			t.Errorf("unexpected float: got %f want %f", f, 3.14159)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		f, err := GetFloat("not-there", "1", 1.11, c, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if f != 1.11 {
			t.Errorf("unexpected float: got %f want %f", f, 1.11)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		f, err := GetFloat("ff-float", "", 1.11, c, c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if f != 3.14159 {
			t.Errorf("unexpected float: got %f want %f", f, 3.14159)
		}
	})
}

func TestGetIntJSON(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		i, err := GetInt("ff-number", "1", 69, c, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if i != 9081 {
			t.Errorf("unexpected int: got %d want %d", i, 9081)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		i, err := GetInt("not-there", "1", 69, c, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if i != 69 {
			t.Errorf("unexpected int: got %d want %d", i, 69)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		i, err := GetInt("ff-number", "", 69, c, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if i != 9081 {
			t.Errorf("unexpected int: got %d want %d", i, 9081)
		}
	})
}

func TestGetStringJSON(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetString("ff-description", "1", "hello", c, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != "Something about chocolate eggs" {
			t.Errorf("unexpected string: got %s want %s", s, "Something about chocolate eggs")
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetString("not-there", "1", "hello", c, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if s != "hello" {
			t.Errorf("unexpected string: got %s want %s", s, "hello")
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetString("ff-description", "", "hello", c, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != "Something about chocolate eggs" {
			t.Errorf("unexpected string: got %s want %s", s, "Something about chocolate eggs")
		}
	})
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
	t.Run("flag exists, returns struct", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONMap(
			"ff-json",
			"1",
			defaultResponseTimes,
			c,
		)
		if err != nil {
			t.Fatalf("unexpected error getting json: %v", err)
		}

		if diff := cmp.Diff(s, expectedResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
	t.Run("flag doesnt exist, returns default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONMap(
			"not-exists",
			"1",
			defaultResponseTimes,
			c,
		)
		// error is returned, but so is default value
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		if diff := cmp.Diff(s, defaultResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
	t.Run("no user id provided, returns flag value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONMap(
			"ff-json",
			"",
			defaultResponseTimes,
			c,
		)
		// error is returned, but so is default value
		if err != nil {
			t.Fatalf("unexpected error getting json: %v", err)
		}

		if diff := cmp.Diff(s, expectedResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
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
	t.Run("flag exists, returns struct", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONStruct[ResponseTimes](
			"ff-json",
			"1",
			defaultResponseTimes,
			c,
		)
		if err != nil {
			t.Fatalf("unexpected error getting json: %v", err)
		}

		if diff := cmp.Diff(s, expectedResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
	t.Run("flag doesnt exist, returns default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONStruct[ResponseTimes](
			"not-exists",
			"1",
			defaultResponseTimes,
			c,
		)
		// error is returned, but so is default value
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		if diff := cmp.Diff(s, defaultResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
	t.Run("no user id provided, returns flag value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetJSONStruct[ResponseTimes](
			"ff-json",
			"",
			defaultResponseTimes,
			c,
		)
		// error is returned, but so is default value
		if err != nil {
			t.Fatalf("unexpected error getting json: %v", err)
		}

		if diff := cmp.Diff(s, expectedResponseTimes); diff != "" {
			t.Errorf("unexpected struct (-got +want)\n%s", diff)
		}
	})
}

func TestGetTimeJSON(t *testing.T) {
	expected, err := time.Parse(time.RFC3339, "2025-07-18T22:37:22.176Z")
	if err != nil {
		t.Fatalf("failed to parse expected time: %v", err)
	}

	defaultTime, err := time.Parse(time.RFC3339, "2099-12-31T23:59:59.999Z")
	if err != nil {
		t.Fatalf("failed to parse default time: %v", err)
	}
	t.Run("flag exists, returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetTime(
			"cr-start",
			"1",
			time.RFC3339,
			defaultTime,
			c,
		)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != expected {
			t.Errorf("unexpected time: got %s want %s", s, expected)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetTime(
			"not-there",
			"1",
			time.RFC3339,
			defaultTime,
			c,
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if s != defaultTime {
			t.Errorf("unexpected time: got %s want %s", s, defaultTime)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		s, err := GetTime(
			"cr-start",
			"",
			time.RFC3339,
			defaultTime,
			c,
		)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != expected {
			t.Errorf("unexpected time: got %s want %s", s, expected)
		}
	})
}

func TestIsEnabledJSON(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabled("is-enabled", "1", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if !b {
			t.Errorf("unexpected bool: got %t want %t", b, true)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabled("not-there", "1", false, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if b {
			t.Errorf("unexpected bool: got %t want %t", b, false)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabled("is-enabled", "", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if !b {
			t.Errorf("unexpected bool: got %t want %t", b, true)
		}
	})
}

func TestIsEnabledByIDJSON(t *testing.T) {
	t.Run("flag exists, user id is in list, returns enabled", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabledByID("is-enabled-for-user", "9", "3", "user-id", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if !b {
			t.Errorf("unexpected bool: got %t want %t", b, true)
		}
	})
	t.Run("flag exists, user id not in list, returns disabled", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabledByID("is-enabled-for-user", "9", "4", "user-id", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if b {
			t.Errorf("unexpected bool: got %t want %t", b, false)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabledByID("not-exists", "9", "4", "user-id", false, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if b {
			t.Errorf("unexpected bool: got %t want %t", b, false)
		}
	})
	t.Run("no user id in context eval, user id is in list, returns enabled", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabledByID("is-enabled-for-user", "", "3", "user-id", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if !b {
			t.Errorf("unexpected bool: got %t want %t", b, true)
		}
	})
	t.Run("no user id in context eval, user id not in list, returns disabled", func(t *testing.T) {
		c := setupClient(t, "flags.goff.json")
		b, err := IsEnabledByID("is-enabled-for-user", "", "4", "user-id", false, c)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if b {
			t.Errorf("unexpected bool: got %t want %t", b, true)
		}
	})
}

func TestIsEnabledByIDListJSON(t *testing.T) {
	t.Run("item is an int", func(t *testing.T) {
		t.Run("flag exists, user id is in list, returns enabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list", "9", 3, false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if !b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
		t.Run("flag exists, user id not in list, returns disabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list", "9", 4, false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, false)
			}
		})
		t.Run("flag doesnt exist, return default", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("not-exists", "9", 4, false, c)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, false)
			}
		})
		t.Run("no user id in context eval, user id is in list, returns enabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list", "", 3, false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if !b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
		t.Run("no user id in context eval, user id not in list, returns disabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list", "", 4, false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
	})
	t.Run("item is a string", func(t *testing.T) {
		t.Run("flag exists, user id is in list, returns enabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list-string", "9", "3", false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if !b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
		t.Run("flag exists, user id not in list, returns disabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list-string", "9", "4", false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, false)
			}
		})
		t.Run("flag doesnt exist, return default", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("not-exists", "9", "4", false, c)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, false)
			}
		})
		t.Run("no user id in context eval, user id is in list, returns enabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list-string", "", "3", false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if !b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
		t.Run("no user id in context eval, user id not in list, returns disabled", func(t *testing.T) {
			c := setupClient(t, "flags.goff.json")
			b, err := IsEnabledByIDList("ff-json-list-string", "", "4", false, c)
			if err != nil {
				t.Fatalf("unexpected error getting flag value: %v", err)
			}

			if b {
				t.Errorf("unexpected bool: got %t want %t", b, true)
			}
		})
	})
}

func TestRefreshJSON(t *testing.T) {
	t.Run("refresh is called and refreshed time is updated", func(t *testing.T) {
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

		err = NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: tmp.Name()},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error creating temp file: %v", err)
		}
		defer Close()

		refreshTime := ffclient.GetCacheRefreshDate()

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

		if refreshTime != ffclient.GetCacheRefreshDate() {
			t.Fatal("cache refreshed unexpectedly")
		}

		Refresh()

		if refreshTime == ffclient.GetCacheRefreshDate() {
			t.Error("expected refresh time to have updated")
		}
	})
}

func TestNewClientJSON(t *testing.T) {
	t.Run("flag file isnt available", func(t *testing.T) {
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
