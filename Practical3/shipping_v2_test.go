// shipping_v2_test.go
package shipping

import (
	"math"
	"strings"
	"testing"
)

func TestCalculateShippingFeeV2(t *testing.T) {
	// Parameter Ranges for Updated Shipping Calculator V2:
	// weight: Valid range (0, 50] kg with tiers: (0,10]=Standard, (10,50]=Heavy
	// zone: Valid values {"Domestic", "International", "Express"} - case sensitive
	// insured: Boolean {true, false} - affects final cost calculation

	testCases := []struct {
		name        string
		weight      float64 // Range: (-∞, +∞), Valid: (0, 50], Tiers: (0,10]=Standard, (10,50]=Heavy
		zone        string  // Range: any string, Valid: {"Domestic", "International", "Express"}
		insured     bool    // Range: {true, false} - adds 1.5% if true
		expectedFee float64 // Expected calculated fee for valid inputs
		expectError bool    // True if we expect an error, false otherwise
		errorText   string  // Expected error message substring
	}{
		// ========== EQUIVALENCE PARTITIONING TESTS ==========

		// P1: Invalid Weight - Too Low (weight <= 0)
		{"EP P1: Negative weight", -5.0, "Domestic", false, 0, true, "invalid weight"},
		{"EP P1: Zero weight", 0.0, "International", true, 0, true, "invalid weight"},

		// P2: Valid Weight - Standard (0 < weight <= 10)
		{"EP P2: Standard weight no insurance", 5.0, "Domestic", false, 5.0, false, ""},
		{"EP P2: Standard weight with insurance", 8.0, "International", true, 20.3, false, ""}, // 20.0 * 1.015 = 20.3

		// P3: Valid Weight - Heavy (10 < weight <= 50)
		{"EP P3: Heavy weight no insurance", 25.0, "Express", false, 37.5, false, ""},      // 30.0 + 7.5 = 37.5
		{"EP P3: Heavy weight with insurance", 35.0, "Domestic", true, 12.6875, false, ""}, // (5.0 + 7.5) * 1.015 = 12.6875

		// P4: Invalid Weight - Too High (weight > 50)
		{"EP P4: Overweight", 75.0, "Express", false, 0, true, "invalid weight"},

		// P5: Valid Zones
		{"EP P5: Domestic zone", 15.0, "Domestic", false, 12.5, false, ""},           // 5.0 + 7.5 = 12.5
		{"EP P5: International zone", 20.0, "International", false, 27.5, false, ""}, // 20.0 + 7.5 = 27.5
		{"EP P5: Express zone", 30.0, "Express", false, 37.5, false, ""},             // 30.0 + 7.5 = 37.5

		// P6: Invalid Zones
		{"EP P6: Invalid zone - empty", 10.0, "", false, 0, true, "invalid zone"},
		{"EP P6: Invalid zone - lowercase", 15.0, "domestic", false, 0, true, "invalid zone"},
		{"EP P6: Invalid zone - unknown", 20.0, "Local", false, 0, true, "invalid zone"},

		// P7: Insurance True
		{"EP P7: Insurance enabled", 5.0, "Domestic", true, 5.075, false, ""}, // 5.0 * 1.015 = 5.075

		// P8: Insurance False (already covered above)

		// ========== BOUNDARY VALUE ANALYSIS TESTS ==========

		// Lower Boundary (around 0)
		{"BVA: Weight boundary 0", 0.0, "Domestic", false, 0, true, "invalid weight"},
		{"BVA: Weight just above 0", 0.1, "International", false, 20.0, false, ""}, // 20.0 base, no surcharge

		// Mid Boundary (around 10) - Standard vs Heavy threshold
		{"BVA: Weight exactly 10 (Standard)", 10.0, "Express", false, 30.0, false, ""},  // 30.0 base, no surcharge
		{"BVA: Weight just above 10 (Heavy)", 10.1, "Domestic", false, 12.5, false, ""}, // 5.0 + 7.5 = 12.5

		// Upper Boundary (around 50)
		{"BVA: Weight exactly 50 (valid)", 50.0, "International", false, 27.5, false, ""}, // 20.0 + 7.5 = 27.5
		{"BVA: Weight just above 50 (invalid)", 50.1, "Express", false, 0, true, "invalid weight"},

		// ========== COMPREHENSIVE SCENARIO TESTS ==========

		// All combinations of weight tiers, zones, and insurance
		{"Scenario: Standard Domestic No Insurance", 5.0, "Domestic", false, 5.0, false, ""},
		{"Scenario: Standard Domestic With Insurance", 5.0, "Domestic", true, 5.075, false, ""}, // 5.0 * 1.015
		{"Scenario: Standard International No Insurance", 8.0, "International", false, 20.0, false, ""},
		{"Scenario: Standard International With Insurance", 8.0, "International", true, 20.3, false, ""}, // 20.0 * 1.015
		{"Scenario: Standard Express No Insurance", 9.0, "Express", false, 30.0, false, ""},
		{"Scenario: Standard Express With Insurance", 9.0, "Express", true, 30.45, false, ""}, // 30.0 * 1.015

		{"Scenario: Heavy Domestic No Insurance", 15.0, "Domestic", false, 12.5, false, ""},               // 5.0 + 7.5
		{"Scenario: Heavy Domestic With Insurance", 15.0, "Domestic", true, 12.6875, false, ""},           // 12.5 * 1.015
		{"Scenario: Heavy International No Insurance", 25.0, "International", false, 27.5, false, ""},     // 20.0 + 7.5
		{"Scenario: Heavy International With Insurance", 25.0, "International", true, 27.9125, false, ""}, // 27.5 * 1.015
		{"Scenario: Heavy Express No Insurance", 40.0, "Express", false, 37.5, false, ""},                 // 30.0 + 7.5
		{"Scenario: Heavy Express With Insurance", 40.0, "Express", true, 38.0625, false, ""},             // 37.5 * 1.015
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFeeV2(tc.weight, tc.zone, tc.insured)

			// Check error expectations
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error containing '%s', but got nil", tc.errorText)
				}
				if tc.errorText != "" && !strings.Contains(err.Error(), tc.errorText) {
					t.Errorf("Expected error containing '%s', but got '%s'", tc.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				// Check fee calculation with small tolerance for floating point precision
				if math.Abs(fee-tc.expectedFee) > 0.0001 {
					t.Errorf("Expected fee %f, but got %f", tc.expectedFee, fee)
				}
			}
		})
	}
}
