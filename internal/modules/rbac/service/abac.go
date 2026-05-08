package service

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：EffectAllow、EffectDeny 声明当前包使用的常量。
// English: EffectAllow、EffectDeny declares constants used by this package.
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// 中文：AccessRequest 定义当前包使用的数据结构或接口。
// English: AccessRequest defines a data structure or interface used by this package.
type AccessRequest struct {
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string
	// 中文：Resource 保存当前结构中的配置或数据值。
	// English: Resource stores a configuration or data value for this struct.
	Resource string
	// 中文：Action 保存当前结构中的配置或数据值。
	// English: Action stores a configuration or data value for this struct.
	Action string
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject map[string]any
	// 中文：ResourceAttributes 保存当前结构中的配置或数据值。
	// English: ResourceAttributes stores a configuration or data value for this struct.
	ResourceAttributes map[string]any
	// 中文：Environment 保存当前结构中的配置或数据值。
	// English: Environment stores a configuration or data value for this struct.
	Environment map[string]any
}

// 中文：ABACPolicy 定义当前包使用的数据结构或接口。
// English: ABACPolicy defines a data structure or interface used by this package.
type ABACPolicy struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string
	// 中文：Effect 保存当前结构中的配置或数据值。
	// English: Effect stores a configuration or data value for this struct.
	Effect string
	// 中文：Resources 保存当前结构中的配置或数据值。
	// English: Resources stores a configuration or data value for this struct.
	Resources []string
	// 中文：Actions 保存当前结构中的配置或数据值。
	// English: Actions stores a configuration or data value for this struct.
	Actions []string
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject map[string]Condition
	// 中文：Resource 保存当前结构中的配置或数据值。
	// English: Resource stores a configuration or data value for this struct.
	Resource map[string]Condition
	// 中文：Environment 保存当前结构中的配置或数据值。
	// English: Environment stores a configuration or data value for this struct.
	Environment map[string]Condition
}

// 中文：Condition 定义当前包使用的数据结构或接口。
// English: Condition defines a data structure or interface used by this package.
type Condition struct {
	// 中文：Operator 保存当前结构中的配置或数据值。
	// English: Operator stores a configuration or data value for this struct.
	Operator string
	// 中文：Value 保存当前结构中的配置或数据值。
	// English: Value stores a configuration or data value for this struct.
	Value any
}

// 中文：ABACDecision 定义当前包使用的数据结构或接口。
// English: ABACDecision defines a data structure or interface used by this package.
type ABACDecision struct {
	// 中文：Allowed 保存当前结构中的配置或数据值。
	// English: Allowed stores a configuration or data value for this struct.
	Allowed bool
	// 中文：MatchedPolicy 保存当前结构中的配置或数据值。
	// English: MatchedPolicy stores a configuration or data value for this struct.
	MatchedPolicy string
	// 中文：Reason 保存当前结构中的配置或数据值。
	// English: Reason stores a configuration or data value for this struct.
	Reason string
}

// 中文：ABACEngine 定义当前包使用的数据结构或接口。
// English: ABACEngine defines a data structure or interface used by this package.
type ABACEngine struct {
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：policies 保存当前结构中的配置或数据值。
	// English: policies stores a configuration or data value for this struct.
	policies []ABACPolicy
}

// 中文：NewABACEngine 创建并返回对应组件实例。
// English: NewABACEngine creates and returns the corresponding component instance.
func NewABACEngine(policies []ABACPolicy) *ABACEngine {
	engine := &ABACEngine{}
	engine.ReplacePolicies(policies)
	return engine
}

// 中文：ReplacePolicies 执行当前包中的对应流程。
// English: ReplacePolicies executes the corresponding workflow in this package.
func (e *ABACEngine) ReplacePolicies(policies []ABACPolicy) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.policies = append([]ABACPolicy(nil), policies...)
}

// 中文：AddPolicy 执行当前包中的对应流程。
// English: AddPolicy executes the corresponding workflow in this package.
func (e *ABACEngine) AddPolicy(policy ABACPolicy) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.policies = append(e.policies, policy)
}

// 中文：Decide 执行当前包中的对应流程。
// English: Decide executes the corresponding workflow in this package.
func (e *ABACEngine) Decide(_ context.Context, req AccessRequest) (ABACDecision, error) {
	if e == nil {
		return ABACDecision{Reason: "abac engine is not configured"}, nil
	}
	if req.Resource == "" || req.Action == "" {
		return ABACDecision{}, fmt.Errorf("resource and action are required")
	}

	e.mu.RLock()
	policies := append([]ABACPolicy(nil), e.policies...)
	e.mu.RUnlock()

	var allow ABACDecision
	for _, policy := range policies {
		matched, err := policy.matches(req)
		if err != nil {
			return ABACDecision{}, err
		}
		if !matched {
			continue
		}
		id := policy.ID
		if id == "" {
			id = "<anonymous>"
		}
		switch strings.ToLower(policy.Effect) {
		case EffectDeny:
			return ABACDecision{Allowed: false, MatchedPolicy: id, Reason: "deny policy matched"}, nil
		case "", EffectAllow:
			allow = ABACDecision{Allowed: true, MatchedPolicy: id, Reason: "allow policy matched"}
		default:
			return ABACDecision{}, fmt.Errorf("unsupported abac effect: %s", policy.Effect)
		}
	}
	if allow.Allowed {
		return allow, nil
	}
	return ABACDecision{Reason: "no abac policy matched"}, nil
}

// 中文：Check 执行当前包中的对应流程。
// English: Check executes the corresponding workflow in this package.
func (e *ABACEngine) Check(ctx context.Context, req AccessRequest) (bool, error) {
	decision, err := e.Decide(ctx, req)
	if err != nil {
		return false, err
	}
	return decision.Allowed, nil
}

// 中文：RegisterABACPolicy 执行当前包中的对应流程。
// English: RegisterABACPolicy executes the corresponding workflow in this package.
func (s *RBACService) RegisterABACPolicy(policy ABACPolicy) {
	if s.abac == nil {
		s.abac = NewABACEngine(nil)
	}
	s.abac.AddPolicy(policy)
}

// 中文：CheckABAC 执行当前包中的对应流程。
// English: CheckABAC executes the corresponding workflow in this package.
func (s *RBACService) CheckABAC(ctx context.Context, req AccessRequest) (bool, error) {
	if s.abac == nil {
		return false, nil
	}
	return s.abac.Check(ctx, enrichAccessRequest(ctx, req))
}

// 中文：CheckAccess 执行当前包中的对应流程。
// English: CheckAccess executes the corresponding workflow in this package.
func (s *RBACService) CheckAccess(ctx context.Context, req AccessRequest) (bool, error) {
	req = enrichAccessRequest(ctx, req)
	if s.repo != nil && req.UserID != "" && req.Resource != "" && req.Action != "" {
		ok, err := s.CheckPermission(ctx, req.UserID, req.Resource, req.Action)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return s.CheckABAC(ctx, req)
}

// 中文：enrichAccessRequest 执行当前包中的对应流程。
// English: enrichAccessRequest executes the corresponding workflow in this package.
func enrichAccessRequest(ctx context.Context, req AccessRequest) AccessRequest {
	if req.UserID == "" {
		req.UserID = types.GetUserID(ctx)
	}
	if req.TenantID == "" {
		req.TenantID = types.GetTenantID(ctx)
	}
	req.Subject = cloneAnyMap(req.Subject)
	req.ResourceAttributes = cloneAnyMap(req.ResourceAttributes)
	req.Environment = cloneAnyMap(req.Environment)
	if req.UserID != "" {
		req.Subject["user_id"] = req.UserID
	}
	if req.TenantID != "" {
		req.Subject["tenant_id"] = req.TenantID
		req.ResourceAttributes["tenant_id"] = req.TenantID
	}
	req.ResourceAttributes["resource"] = req.Resource
	req.Environment["action"] = req.Action
	return req
}

// 中文：matches 执行当前包中的对应流程。
// English: matches executes the corresponding workflow in this package.
func (p ABACPolicy) matches(req AccessRequest) (bool, error) {
	if !patternListMatches(p.Resources, req.Resource) || !patternListMatches(p.Actions, req.Action) {
		return false, nil
	}
	if ok, err := matchConditions(req.Subject, p.Subject); !ok || err != nil {
		return false, err
	}
	if ok, err := matchConditions(req.ResourceAttributes, p.Resource); !ok || err != nil {
		return false, err
	}
	if ok, err := matchConditions(req.Environment, p.Environment); !ok || err != nil {
		return false, err
	}
	return true, nil
}

// 中文：patternListMatches 执行当前包中的对应流程。
// English: patternListMatches executes the corresponding workflow in this package.
func patternListMatches(patterns []string, value string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, pattern := range patterns {
		if patternMatches(pattern, value) {
			return true
		}
	}
	return false
}

// 中文：patternMatches 执行当前包中的对应流程。
// English: patternMatches executes the corresponding workflow in this package.
func patternMatches(pattern, value string) bool {
	if pattern == "*" || pattern == value {
		return true
	}
	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(value, strings.TrimSuffix(pattern, "*")) {
		return true
	}
	return false
}

// 中文：matchConditions 执行当前包中的对应流程。
// English: matchConditions executes the corresponding workflow in this package.
func matchConditions(attrs map[string]any, conditions map[string]Condition) (bool, error) {
	for key, cond := range conditions {
		actual, exists := attrs[key]
		ok, err := matchCondition(actual, exists, cond)
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

// 中文：matchCondition 执行当前包中的对应流程。
// English: matchCondition executes the corresponding workflow in this package.
func matchCondition(actual any, exists bool, cond Condition) (bool, error) {
	op := strings.ToLower(cond.Operator)
	if op == "" {
		op = "eq"
	}
	switch op {
	case "exists":
		want, ok := cond.Value.(bool)
		if !ok {
			return false, fmt.Errorf("exists condition requires bool value")
		}
		return exists == want, nil
	case "eq", "equals":
		return exists && reflect.DeepEqual(actual, cond.Value), nil
	case "ne", "not_equals":
		return !exists || !reflect.DeepEqual(actual, cond.Value), nil
	case "in":
		return exists && contains(cond.Value, actual), nil
	case "contains":
		return exists && contains(actual, cond.Value), nil
	case "gt", "gte", "lt", "lte":
		left, ok := asFloat(actual)
		if !exists || !ok {
			return false, nil
		}
		right, ok := asFloat(cond.Value)
		if !ok {
			return false, fmt.Errorf("%s condition requires numeric value", op)
		}
		switch op {
		case "gt":
			return left > right, nil
		case "gte":
			return left >= right, nil
		case "lt":
			return left < right, nil
		default:
			return left <= right, nil
		}
	default:
		return false, fmt.Errorf("unsupported abac condition operator: %s", cond.Operator)
	}
}

// 中文：contains 执行当前包中的对应流程。
// English: contains executes the corresponding workflow in this package.
func contains(container, item any) bool {
	if container == nil {
		return false
	}
	if s, ok := container.(string); ok {
		needle, ok := item.(string)
		return ok && strings.Contains(s, needle)
	}
	value := reflect.ValueOf(container)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return reflect.DeepEqual(container, item)
	}
	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(value.Index(i).Interface(), item) {
			return true
		}
	}
	return false
}

// 中文：asFloat 执行当前包中的对应流程。
// English: asFloat executes the corresponding workflow in this package.
func asFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

// 中文：cloneAnyMap 执行当前包中的对应流程。
// English: cloneAnyMap executes the corresponding workflow in this package.
func cloneAnyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in)+4)
	for k, v := range in {
		out[k] = v
	}
	return out
}
