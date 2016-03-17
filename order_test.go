package api

import (
	"testing"
)

func TestOrder(t *testing.T) {
	nothing := ParseOrders("")
	if len(nothing) != 0 {
		t.Errorf("unexpected length of ParseOrders with an empty string")
	}

	orders := ParseOrders("-year,n,-month") // without white-listing
	if len(orders) != 3 {
		t.Fatalf("unexpected length of ParseOrders: %d != 3", len(orders))
	}
	first := Order{Name: "year", Desc: true}
	second := Order{Name: "n"}
	third := Order{Name: "month", Desc: true}
	if orders[0] != first {
		t.Errorf("unexpected first Order: %v != %v", orders[0], first)
	}
	if orders[1] != second {
		t.Errorf("unexpected second Order: %v != %v", orders[1], second)
	}
	if orders[2] != third {
		t.Errorf("unexpected third Order: %v != %v", orders[2], third)
	}
}
