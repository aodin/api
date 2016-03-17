package api

import "strings"

type Order struct {
	Name string
	Desc bool
}

func (order Order) String() string {
	if order.Desc {
		return "-" + order.Name
	}
	return order.Name
}

type Orders []Order

func (orders Orders) String() string {
	strs := make([]string, len(orders))
	for i, order := range orders {
		strs[i] = order.String()
	}
	// TODO escaping?
	return strings.Join(strs, ",")
}

// Names returns the order names as a slice of strings
func (orders Orders) Names() []string {
	names := make([]string, len(orders))
	for i, order := range orders {
		names[i] = order.Name
	}
	return names
}

func (orders Orders) Has(name string) bool {
	for _, order := range orders {
		if order.Name == name {
			return true
		}
	}
	return false
}

func (orders Orders) Get(name string) (order Order) {
	for _, order := range orders {
		if order.Name == name {
			return order
		}
	}
	return
}

// ParseOrders returns the valid order by columns provided by user
func ParseOrders(get string, whitelisted ...Order) (orders Orders) {
	whitelist := Orders(whitelisted)
	// Save valid order properties to rebuild the URI
	// TODO escaping?
	parts := strings.Split(get, ",")
	for _, part := range parts {
		// TODO error on bad columns / format
		var order Order
		if part != "" && part[0] == '-' {
			order.Desc = true
			part = part[1:]
		}

		// TODO strip whitespace?
		if part == "" {
			continue
		}

		// If a whitelist has been provided and the Order is not present,
		// then silently skip
		if len(whitelist) > 0 && !whitelist.Has(part) {
			continue
		}

		order.Name = part
		orders = append(orders, order)
	}
	return
}
