package zenoh

import (
	"errors"
	"log"
	"strings"
)

var ErrInvalidKeyExpr = errors.New("invalid key expression")

type keyExprImpl struct {
	expr     string
	segments []string
	hasWild  bool
}

func newKeyExprImpl(expr string) (*keyExprImpl, error) {
	if expr == "" {
		return nil, ErrInvalidKeyExpr
	}
	if err := validateKeyExpr(expr); err != nil {
		return nil, err
	}
	return &keyExprImpl{
		expr:     expr,
		segments: parseKeySegments(expr),
		hasWild:  strings.Contains(expr, "*"),
	}, nil
}

func validateKeyExpr(expr string) error {
	for _, c := range expr {
		if c == 0 || c == '\n' || c == '\r' {
			return ErrInvalidKeyExpr
		}
	}
	if strings.Contains(expr, "***") {
		return ErrInvalidKeyExpr
	}
	return nil
}

func parseKeySegments(expr string) []string {
	expr = strings.TrimSuffix(expr, "/")
	if expr == "" {
		return []string{}
	}
	return strings.Split(expr, "/")
}

func (k *keyExprImpl) String() string {
	if k == nil {
		return ""
	}
	return k.expr
}

func (k *keyExprImpl) HasWildcard() bool {
	if k == nil {
		return false
	}
	return k.hasWild
}

func NewKeyExpr(expr string) (*KeyExpr, error) {
	impl, err := newKeyExprImpl(expr)
	if err != nil {
		return nil, err
	}
	return &KeyExpr{expr: impl.expr}, nil
}

func KeyExprFromStr(expr string) (*KeyExpr, error) {
	return NewKeyExpr(expr)
}

func Join(a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrInvalidKeyExpr
	}
	a = strings.TrimSuffix(a, "/")
	b = strings.TrimPrefix(b, "/")
	if a == "" {
		return b, nil
	}
	if b == "" {
		return a, nil
	}
	return a + "/" + b, nil
}

func (k *KeyExpr) Canonize() (string, error) {
	if k == nil {
		return "", ErrInvalidKeyExpr
	}
	impl, err := newKeyExprImpl(k.expr)
	if err != nil {
		return "", err
	}
	segments := make([]string, 0)
	for _, seg := range impl.segments {
		if seg != "" {
			segments = append(segments, seg)
		}
	}
	return strings.Join(segments, "/"), nil
}

func CanonizeString(expr string) (string, error) {
	ke, err := NewKeyExpr(expr)
	if err != nil {
		return "", err
	}
	return ke.Canonize()
}

func (k *KeyExpr) Includes(other *KeyExpr) (bool, error) {
	if k == nil || other == nil {
		return false, ErrInvalidKeyExpr
	}
	if k.expr == other.expr {
		return true, nil
	}
	implA, _ := newKeyExprImpl(k.expr)
	implB, _ := newKeyExprImpl(other.expr)
	if !implA.hasWild {
		return k.expr == other.expr, nil
	}
	if !implB.hasWild {
		return keyMatchesPattern(implA.expr, implA.segments, implB.expr), nil
	}
	return keyIncludesPattern(implA.segments, implB.segments), nil
}

func (k *KeyExpr) Intersects(other *KeyExpr) (bool, error) {
	if k == nil || other == nil {
		return false, ErrInvalidKeyExpr
	}
	if k.expr == other.expr {
		return true, nil
	}
	implA, _ := newKeyExprImpl(k.expr)
	implB, _ := newKeyExprImpl(other.expr)
	if !implA.hasWild && !implB.hasWild {
		return k.expr == other.expr, nil
	}
	return keyIntersectsPattern(implA.segments, implB.segments), nil
}

func keyMatchesPattern(pattern string, patternSegs []string, key string) bool {
	keySegs := parseKeySegments(key)
	return matchKeySegments(patternSegs, keySegs)
}

func matchKeySegments(pattern, key []string) bool {
	if len(pattern) == 0 {
		return len(key) == 0
	}
	dWildIdx := -1
	for i, seg := range pattern {
		if seg == "**" {
			dWildIdx = i
			break
		}
	}
	if dWildIdx == -1 {
		if len(pattern) != len(key) {
			return false
		}
		for i := range pattern {
			if !matchWildSegment(pattern[i], key[i]) {
				return false
			}
		}
		return true
	}
	for i := 0; i < dWildIdx; i++ {
		if i >= len(key) || !matchWildSegment(pattern[i], key[i]) {
			return false
		}
	}
	for i := dWildIdx + 1; i < len(pattern); i++ {
		segIdx := len(key) - (len(pattern) - i)
		if segIdx < dWildIdx || !matchWildSegment(pattern[i], key[segIdx]) {
			return false
		}
	}
	return true
}

func matchWildSegment(pattern, key string) bool {
	if pattern == "*" || pattern == "**" {
		return true
	}
	return pattern == key
}

func keyIncludesPattern(pattern, other []string) bool {
	dWildIdx := -1
	for i, seg := range pattern {
		if seg == "**" {
			dWildIdx = i
			break
		}
	}
	if dWildIdx == -1 {
		if len(pattern) != len(other) {
			return false
		}
		for i := range pattern {
			if pattern[i] != "*" && pattern[i] != other[i] {
				return false
			}
		}
		return true
	}
	for i := 0; i < dWildIdx; i++ {
		if i >= len(other) {
			return false
		}
		if pattern[i] != "*" && pattern[i] != other[i] {
			return false
		}
	}
	for i := dWildIdx + 1; i < len(pattern); i++ {
		segIdx := len(other) - (len(pattern) - i)
		if segIdx < dWildIdx {
			continue
		}
		if pattern[i] != "*" && pattern[i] != other[segIdx] {
			return false
		}
	}
	return true
}

func keyIntersectsPattern(a, b []string) bool {
	aDWild := findDWild(a)
	bDWild := findDWild(b)
	if aDWild != -1 && bDWild != -1 {
		return true
	}
	if aDWild != -1 {
		return intersectWithDWild(a, b, aDWild)
	}
	if bDWild != -1 {
		return intersectWithDWild(b, a, bDWild)
	}
	return intersectNoWildcards(a, b)
}

func findDWild(seg []string) int {
	for i, s := range seg {
		if s == "**" {
			return i
		}
	}
	return -1
}

func intersectWithDWild(wildP, pat []string, wIdx int) bool {
	for i := 0; i < wIdx; i++ {
		if i >= len(pat) {
			return false
		}
		if wildP[i] != "*" && wildP[i] != pat[i] {
			return false
		}
	}
	for i := wIdx + 1; i < len(wildP); i++ {
		segIdx := len(pat) - (len(wildP) - i)
		if segIdx < wIdx {
			continue
		}
		if wildP[i] != "*" && wildP[i] != pat[segIdx] {
			return false
		}
	}
	return true
}

func intersectNoWildcards(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != "*" && b[i] != "*" && a[i] != b[i] {
			return false
		}
	}
	return true
}

func (k *KeyExpr) Resolve(session *Session) error {
	if k == nil {
		return ErrInvalidKeyExpr
	}
	if session == nil || !session.IsValid() {
		return ErrInvalidValue
	}
	log.Printf("[zenoh] keyexpr resolved: %s", k.expr)
	return nil
}

func ResolveString(session *Session, expr string) error {
	ke, err := NewKeyExpr(expr)
	if err != nil {
		return err
	}
	return ke.Resolve(session)
}

func (k *KeyExpr) Cede() string {
	if k == nil {
		return ""
	}
	return k.expr
}

func (k *OwnedKeyExpr) Cede() string {
	if k == nil || k.ptr == 0 {
		return ""
	}
	return ""
}
