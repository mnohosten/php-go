package parallel

import (
	"testing"
)

// TestNewSafetyAnalyzer tests creating a new analyzer
func TestNewSafetyAnalyzer(t *testing.T) {
	analyzer := NewSafetyAnalyzer()
	if analyzer == nil {
		t.Error("NewSafetyAnalyzer should return non-nil analyzer")
	}
}

// TestIsFileIOFunction tests the file I/O function detection
func TestIsFileIOFunction(t *testing.T) {
	analyzer := NewSafetyAnalyzer()

	fileIOFuncs := []string{"fopen", "file_get_contents", "fwrite", "mkdir", "unlink"}
	for _, fn := range fileIOFuncs {
		if !analyzer.isFileIOFunction(fn) {
			t.Errorf("%s should be detected as file I/O function", fn)
		}
	}

	nonFileIOFuncs := []string{"strlen", "array_map", "json_encode"}
	for _, fn := range nonFileIOFuncs {
		if analyzer.isFileIOFunction(fn) {
			t.Errorf("%s should not be detected as file I/O function", fn)
		}
	}
}

// TestIsNetworkIOFunction tests the network I/O function detection
func TestIsNetworkIOFunction(t *testing.T) {
	analyzer := NewSafetyAnalyzer()

	networkIOFuncs := []string{"curl_exec", "fsockopen", "socket_connect", "mail"}
	for _, fn := range networkIOFuncs {
		if !analyzer.isNetworkIOFunction(fn) {
			t.Errorf("%s should be detected as network I/O function", fn)
		}
	}
}

// TestIsDatabaseFunction tests the database function detection
func TestIsDatabaseFunction(t *testing.T) {
	analyzer := NewSafetyAnalyzer()

	dbFuncs := []string{"mysqli_query", "pg_query", "pdo"}
	for _, fn := range dbFuncs {
		if !analyzer.isDatabaseFunction(fn) {
			t.Errorf("%s should be detected as database function", fn)
		}
	}
}

// TestHasSideEffects tests the side effects function detection
func TestHasSideEffects(t *testing.T) {
	analyzer := NewSafetyAnalyzer()

	sideEffectFuncs := []string{"echo", "print", "header", "session_start", "rand", "time"}
	for _, fn := range sideEffectFuncs {
		if !analyzer.hasSideEffects(fn) {
			t.Errorf("%s should be detected as having side effects", fn)
		}
	}

	pureFuncs := []string{"strlen", "abs", "max", "min"}
	for _, fn := range pureFuncs {
		if analyzer.hasSideEffects(fn) {
			t.Errorf("%s should not be detected as having side effects", fn)
		}
	}
}
