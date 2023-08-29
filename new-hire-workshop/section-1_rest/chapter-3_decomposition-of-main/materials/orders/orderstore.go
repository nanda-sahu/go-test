package orders

type OrderStore struct {
	orders []Order
}

func NewOrderStore() *OrderStore {
	// Logic for initialising an OrderStore with the 'starting' orders goes here
}

func (o *OrderStore) GetOrders(customerID string) []Order {
	// Logic to fetch an order goes here
}
