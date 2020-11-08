package api

import (
	"encoding/base64"
	"errors"

	"github.com/soulik/lego/v4/acme"
)

type OrderService service

// New Creates a new order.
func (o *OrderService) New(domains []string) (acme.ExtendedOrder, error) {
	var identifiers []acme.Identifier
	for _, domain := range domains {
		identifiers = append(identifiers, acme.Identifier{Type: "dns", Value: domain})
	}

	orderReq := acme.Order{Identifiers: identifiers}

	var order acme.Order
	resp, err := o.core.post(o.core.GetDirectory().NewOrderURL, orderReq, &order)
	if err != nil {
		return acme.ExtendedOrder{}, err
	}

	return acme.ExtendedOrder{
		Order:               order,
		Location:            resp.Header.Get("Location"),
		AlternateChainLinks: getLinks(resp.Header, "alternate"),
	}, nil
}

// Get Gets an order.
func (o *OrderService) Get(orderURL string) (acme.ExtendedOrder, error) {
	if len(orderURL) == 0 {
		return acme.ExtendedOrder{}, errors.New("order[get]: empty URL")
	}

	var order acme.Order
	resp, err := o.core.postAsGet(orderURL, &order)
	if err != nil {
		return acme.ExtendedOrder{}, err
	}

	return acme.ExtendedOrder{
		Order:               order,
		AlternateChainLinks: getLinks(resp.Header, "alternate"),
	}, nil
}

// UpdateForCSR Updates an order for a CSR.
func (o *OrderService) UpdateForCSR(orderURL string, csr []byte) (acme.ExtendedOrder, error) {
	csrMsg := acme.CSRMessage{
		Csr: base64.RawURLEncoding.EncodeToString(csr),
	}

	var order acme.Order
	resp, err := o.core.post(orderURL, csrMsg, &order)
	if err != nil {
		return acme.ExtendedOrder{}, err
	}

	if order.Status == acme.StatusInvalid {
		return acme.ExtendedOrder{}, order.Error
	}

	return acme.ExtendedOrder{
		Order:               order,
		AlternateChainLinks: getLinks(resp.Header, "alternate"),
	}, nil
}
