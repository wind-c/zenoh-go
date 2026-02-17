package zenoh

import (
	"testing"
)

func TestNewKeyExpr(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"valid simple", "demo/example", false},
		{"valid with wildcard single", "demo/*", false},
		{"valid with wildcard multi", "demo/**", false},
		{"valid with multiple wildcards", "demo/*/test/*", false},
		{"empty", "", true},
		{"invalid chars", "demo\ntest", true},
		{"invalid triple star", "demo/***", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKeyExpr(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKeyExprFromStr(t *testing.T) {
	ke, err := KeyExprFromStr("test/expr")
	if err != nil {
		t.Fatalf("KeyExprFromStr() error = %v", err)
	}
	if ke.String() != "test/expr" {
		t.Errorf("KeyExprFromStr() = %v, want test/expr", ke.String())
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		want    string
		wantErr bool
	}{
		{"normal join", "demo", "example", "demo/example", false},
		{"trailing slash", "demo/", "example", "demo/example", false},
		{"leading slash", "demo", "/example", "demo/example", false},
		{"both slashes", "demo/", "/example", "demo/example", false},
		{"empty a", "", "example", "example", false},
		{"empty b", "demo", "", "demo", false},
		{"both empty", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Join(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanonize(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    string
		wantErr bool
	}{
		{"no change", "demo/example", "demo/example", false},
		{"trailing slash", "demo/example/", "demo/example", false},
		{"multiple slashes", "demo//example", "demo/example", false},
		{"empty segments", "demo//example//test", "demo/example/test", false},
		{"just slash", "/", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ke, err := NewKeyExpr(tt.expr)
			if err != nil {
				t.Fatalf("NewKeyExpr() error = %v", err)
			}
			got, err := ke.Canonize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Canonize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Canonize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanonizeString(t *testing.T) {
	got, err := CanonizeString("demo/example/")
	if err != nil {
		t.Fatalf("CanonizeString() error = %v", err)
	}
	want := "demo/example"
	if got != want {
		t.Errorf("CanonizeString() = %v, want %v", got, want)
	}
}

func TestKeyExprIncludes(t *testing.T) {
	tests := []struct {
		name    string
		expr1   string
		expr2   string
		want    bool
		wantErr bool
	}{
		{"exact match", "demo/example", "demo/example", true, false},
		{"exact match reverse", "demo/example", "demo/example", true, false},
		{"wildcard includes exact", "demo/*", "demo/test", true, false},
		{"wildcard does not include longer", "demo/*", "demo/test/one", false, false},
		{"double wildcard includes exact", "demo/**", "demo/test", true, false},
		{"double wildcard includes longer", "demo/**", "demo/test/one/two", true, false},
		{"different paths no include", "demo/*", "other/*", false, false},
		{"wildcard includes wildcard", "demo/**", "demo/*", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ke1, err := NewKeyExpr(tt.expr1)
			if err != nil {
				t.Fatalf("NewKeyExpr() error = %v", err)
			}
			ke2, err := NewKeyExpr(tt.expr2)
			if err != nil {
				t.Fatalf("NewKeyExpr() error = %v", err)
			}
			got, err := ke1.Includes(ke2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Includes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Includes() = %v, want %v for %s includes %s", got, tt.want, tt.expr1, tt.expr2)
			}
		})
	}
}

func TestKeyExprIntersects(t *testing.T) {
	tests := []struct {
		name    string
		expr1   string
		expr2   string
		want    bool
		wantErr bool
	}{
		{"exact match", "demo/example", "demo/example", true, false},
		{"wildcard intersects exact", "demo/*", "demo/test", true, false},
		{"wildcard intersects wildcard", "demo/*", "demo/*", true, false},
		{"double wildcard intersects", "demo/**", "demo/test", true, false},
		{"different prefixes no intersect", "demo/*", "other/*", false, false},
		{"different exact no intersect", "demo/one", "demo/two", false, false},
		{"double wildcard both", "a/**", "b/**", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ke1, err := NewKeyExpr(tt.expr1)
			if err != nil {
				t.Fatalf("NewKeyExpr() error = %v", err)
			}
			ke2, err := NewKeyExpr(tt.expr2)
			if err != nil {
				t.Fatalf("NewKeyExpr() error = %v", err)
			}
			got, err := ke1.Intersects(ke2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Intersects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Intersects() = %v, want %v for %s intersects %s", got, tt.want, tt.expr1, tt.expr2)
			}
		})
	}
}

func TestKeyExprString(t *testing.T) {
	ke, err := NewKeyExpr("demo/example/*")
	if err != nil {
		t.Fatalf("NewKeyExpr() error = %v", err)
	}
	if got := ke.String(); got != "demo/example/*" {
		t.Errorf("String() = %v, want demo/example/*", got)
	}
}

func TestKeyExprIsValid(t *testing.T) {
	ke, err := NewKeyExpr("demo/example")
	if err != nil {
		t.Fatalf("NewKeyExpr() error = %v", err)
	}
	if !ke.IsValid() {
		t.Error("IsValid() = false, want true for valid KeyExpr")
	}

	var nilKe *KeyExpr
	if nilKe.IsValid() {
		t.Error("IsValid() = true, want false for nil KeyExpr")
	}
}

func TestKeyExprResolve(t *testing.T) {
	ke, err := NewKeyExpr("demo/example")
	if err != nil {
		t.Fatalf("NewKeyExpr() error = %v", err)
	}

	err = ke.Resolve(&Session{ptr: 0})
	if err != ErrInvalidValue {
		t.Errorf("Resolve() with invalid session = %v, want ErrInvalidValue", err)
	}
}

func TestKeyExprCede(t *testing.T) {
	ke, err := NewKeyExpr("demo/example")
	if err != nil {
		t.Fatalf("NewKeyExpr() error = %v", err)
	}
	got := ke.Cede()
	if got != "demo/example" {
		t.Errorf("Cede() = %v, want demo/example", got)
	}

	var nilKe *KeyExpr
	if nilKe.Cede() != "" {
		t.Error("Cede() on nil should return empty string")
	}
}
