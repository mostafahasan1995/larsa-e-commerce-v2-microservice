package call_center_models

import (
	"encoding/json"
	"testing"
)

func TestDeliveryFeasibilityResponse_HasDeliveryPricingModAndEveryFeeKm(t *testing.T) {
	everyFeeKm := 2.5
	resp := DeliveryFeasibilityResponse{
		Feasible:           true,
		DeliveryFee:        10.0,
		EstimatedMinutes:   30,
		DeliveryPricingMod: "distance",
		EveryFeeKm:         everyFeeKm,
	}
	if resp.DeliveryPricingMod != "branch_distance" {
		t.Errorf("expected deliveryPricingMod branch_distance, got %s", resp.DeliveryPricingMod)
	}
	if resp.EveryFeeKm == 0 || resp.EveryFeeKm != 2.5 {
		t.Errorf("expected everyFeeKm 2.5, got %v", resp.EveryFeeKm)
	}
}

func TestDeliveryFeasibilityResponse_JSONIncludesNewFields(t *testing.T) {
	everyFeeKm := 1.0
	resp := DeliveryFeasibilityResponse{
		Feasible:           true,
		DeliveryPricingMod: "zone",
		EveryFeeKm:         everyFeeKm,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := decoded["deliveryPricingMod"]; !ok {
		t.Error("JSON missing deliveryPricingMod")
	}
	if _, ok := decoded["everyFeeKm"]; !ok {
		t.Error("JSON missing everyFeeKm")
	}
}

func TestDeliveryFeasibilityResponse_ZoneModeEveryFeeKmNil(t *testing.T) {
	resp := DeliveryFeasibilityResponse{
		Feasible:           true,
		DeliveryPricingMod: "zone",
		EveryFeeKm:         0,
	}
	if resp.EveryFeeKm != 0 {
		t.Errorf("zone mode should have everyFeeKm nil, got %v", resp.EveryFeeKm)
	}
}
