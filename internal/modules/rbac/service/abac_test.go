package service

import (
	"context"
	"testing"

	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestABACEngineAllowsMatchingPolicy 验证相关行为符合预期。
// English: TestABACEngineAllowsMatchingPolicy verifies the related behavior.
func TestABACEngineAllowsMatchingPolicy(t *testing.T) {
	engine := NewABACEngine([]ABACPolicy{
		{
			ID:        "owner-read",
			Effect:    EffectAllow,
			Resources: []string{"orders.*"},
			Actions:   []string{"read"},
			Subject: map[string]Condition{
				"department": {Operator: "eq", Value: "sales"},
			},
			Resource: map[string]Condition{
				"classification": {Operator: "in", Value: []string{"public", "internal"}},
			},
			Environment: map[string]Condition{
				"hour": {Operator: "gte", Value: 9},
			},
		},
	})

	ok, err := engine.Check(context.Background(), AccessRequest{
		Resource: "orders.detail",
		Action:   "read",
		Subject: map[string]any{
			"department": "sales",
		},
		ResourceAttributes: map[string]any{
			"classification": "internal",
		},
		Environment: map[string]any{
			"hour": 10,
		},
	})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if !ok {
		t.Fatal("policy should allow request")
	}
}

// 中文：TestABACEngineDenyOverridesAllow 验证相关行为符合预期。
// English: TestABACEngineDenyOverridesAllow verifies the related behavior.
func TestABACEngineDenyOverridesAllow(t *testing.T) {
	engine := NewABACEngine([]ABACPolicy{
		{ID: "allow-all", Effect: EffectAllow, Resources: []string{"*"}, Actions: []string{"*"}},
		{
			ID:        "deny-confidential",
			Effect:    EffectDeny,
			Resources: []string{"documents"},
			Actions:   []string{"read"},
			Resource: map[string]Condition{
				"classification": {Operator: "eq", Value: "confidential"},
			},
		},
	})

	decision, err := engine.Decide(context.Background(), AccessRequest{
		Resource: "documents",
		Action:   "read",
		ResourceAttributes: map[string]any{
			"classification": "confidential",
		},
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}
	if decision.Allowed || decision.MatchedPolicy != "deny-confidential" {
		t.Fatalf("decision = %+v, want deny-confidential denial", decision)
	}
}

// 中文：TestRBACServiceCheckABACEnrichesContextAttributes 验证相关行为符合预期。
// English: TestRBACServiceCheckABACEnrichesContextAttributes verifies the related behavior.
func TestRBACServiceCheckABACEnrichesContextAttributes(t *testing.T) {
	svc := NewRBACService(nil)
	svc.RegisterABACPolicy(ABACPolicy{
		ID:        "tenant-match",
		Effect:    EffectAllow,
		Resources: []string{"tenant.settings"},
		Actions:   []string{"update"},
		Subject: map[string]Condition{
			"user_id":   {Operator: "eq", Value: "user-1"},
			"tenant_id": {Operator: "eq", Value: "tenant-1"},
		},
		Resource: map[string]Condition{
			"tenant_id": {Operator: "eq", Value: "tenant-1"},
		},
	})

	ctx := types.WithTenantID(types.WithUserID(context.Background(), "user-1"), "tenant-1")
	ok, err := svc.CheckAccess(ctx, AccessRequest{
		Resource: "tenant.settings",
		Action:   "update",
	})
	if err != nil {
		t.Fatalf("CheckAccess returned error: %v", err)
	}
	if !ok {
		t.Fatal("context-enriched ABAC policy should allow request")
	}
}

// 中文：TestABACEngineRejectsUnsupportedOperator 验证相关行为符合预期。
// English: TestABACEngineRejectsUnsupportedOperator verifies the related behavior.
func TestABACEngineRejectsUnsupportedOperator(t *testing.T) {
	engine := NewABACEngine([]ABACPolicy{
		{
			ID:        "bad-policy",
			Effect:    EffectAllow,
			Resources: []string{"orders"},
			Actions:   []string{"read"},
			Subject: map[string]Condition{
				"department": {Operator: "regex", Value: "sales"},
			},
		},
	})

	if _, err := engine.Check(context.Background(), AccessRequest{
		Resource: "orders",
		Action:   "read",
		Subject:  map[string]any{"department": "sales"},
	}); err == nil {
		t.Fatal("Check returned nil error, want unsupported operator error")
	}
}
