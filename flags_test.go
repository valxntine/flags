package flags

import (
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
	"os"
	"testing"
	"time"
)

func TestGetFloat(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		f, err := GetFloat("ff-float", "1", 1.11)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if f != 3.14159 {
			t.Errorf("unexpected float: got %f want %f", f, 3.14159)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		f, err := GetFloat("not-there", "1", 1.11)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if f != 1.11 {
			t.Errorf("unexpected float: got %f want %f", f, 1.11)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		f, err := GetFloat("ff-float", "", 1.11)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if f != 3.14159 {
			t.Errorf("unexpected float: got %f want %f", f, 3.14159)
		}
	})
}

func TestGetInt(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		i, err := GetInt("ff-number", "1", 69)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if i != 9081 {
			t.Errorf("unexpected int: got %d want %d", i, 9081)
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		i, err := GetInt("not-there", "1", 69)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if i != 69 {
			t.Errorf("unexpected int: got %d want %d", i, 69)
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		i, err := GetInt("ff-number", "", 69)
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if i != 9081 {
			t.Errorf("unexpected int: got %d want %d", i, 9081)
		}
	})
}

func TestGetString(t *testing.T) {
	t.Run("flag exists, returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		s, err := GetString("ff-description", "1", "hello")
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != "Something about chocolate eggs" {
			t.Errorf("unexpected string: got %s want %s", s, "Something about chocolate eggs")
		}
	})
	t.Run("flag doesnt exist, return default", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		s, err := GetString("not-there", "1", "hello")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if s != "hello" {
			t.Errorf("unexpected string: got %s want %s", s, "hello")
		}
	})
	t.Run("no user id, still returns value", func(t *testing.T) {
		err := NewClient(Config{
			PollingInterval: 10 * time.Minute,
			Retrievers: []retriever.Retriever{
				&fileretriever.Retriever{Path: "flags.goff.yaml"},
			},
		})
		defer Close()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}

		s, err := GetString("ff-description", "", "hello")
		if err != nil {
			t.Fatalf("unexpected error getting flag value: %v", err)
		}

		if s != "Something about chocolate eggs" {
			t.Errorf("unexpected string: got %s want %s", s, "Something about chocolate eggs")
		}
	})
}

//	func TestGetJSONMap(t *testing.T) {
//		type args struct {
//			flag         string
//			userID       string
//			defaultValue map[string]any
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    map[string]any
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := flags.GetJSONMap(tt.args.flag, tt.args.userID, tt.args.defaultValue)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetJSONMap() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetJSONMap() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
//	func TestGetJSONStruct(t *testing.T) {
//		type args[T any] struct {
//			flag         string
//			userID       string
//			defaultValue flags.T
//		}
//		type testCase[T any] struct {
//			name    string
//			args    args[T]
//			want    flags.T
//			wantErr bool
//		}
//		tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := flags.GetJSONStruct(tt.args.flag, tt.args.userID, tt.args.defaultValue)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetJSONStruct() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetJSONStruct() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
//	func TestGetTime(t *testing.T) {
//		type args struct {
//			flag         string
//			userID       string
//			layout       string
//			defaultValue time.Time
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    time.Time
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := flags.GetTime(tt.args.flag, tt.args.userID, tt.args.layout, tt.args.defaultValue)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("GetTime() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("GetTime() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
//	func TestIsEnabled(t *testing.T) {
//		type args struct {
//			flag         string
//			userID       string
//			defaultValue bool
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    bool
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := flags.IsEnabled(tt.args.flag, tt.args.userID, tt.args.defaultValue)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("IsEnabled() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if got != tt.want {
//					t.Errorf("IsEnabled() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//
//	func TestIsEnabledByID(t *testing.T) {
//		type args struct {
//			flag         string
//			userID       string
//			id           string
//			lookup       string
//			defaultValue bool
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    bool
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := flags.IsEnabledByID(tt.args.flag, tt.args.userID, tt.args.id, tt.args.lookup, tt.args.defaultValue)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("IsEnabledByID() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if got != tt.want {
//					t.Errorf("IsEnabledByID() got = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
func TestRefresh(t *testing.T) {
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

func TestNewClient(t *testing.T) {
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
