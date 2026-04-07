package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type customerRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *customerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateOrder(ctx context.Context, userID uuid.UUID, order model.Order) (model.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (r.db.BeginTx): %w", err)
	}
	defer tx.Rollback()

	// map requested items to fetch current prices and calculate totalfee
	var itemIDs []uuid.UUID
	qtyMap := make(map[uuid.UUID]int)
	for _, resto := range order.Restaurants {
		for _, item := range resto.Items {
			itemIDs = append(itemIDs, item.ID)
			qtyMap[item.ID] = item.Quantity
		}
	}

	// fetch prices from db
	rows, err := tx.QueryContext(ctx, `SELECT id, price FROM items WHERE id = ANY($1)`, pq.Array(itemIDs))
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (tx.QueryContext): %w", err)
	}
	defer rows.Close()

	totalFee := 0
	for rows.Next() {
		var id uuid.UUID
		var price int
		if err := rows.Scan(&id, &price); err != nil {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (rows.Scan): %w", err)
		}
		qty, ok := qtyMap[id]
		if !ok {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (qtyMap[id]) id not found in qtyMap: %w", model.ErrNotFound)
		}
		totalFee += price * qty
	}
	totalFee += order.DeliveryFee

	// insert order
	var orderID uuid.UUID
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (customer_id, order_status, delivery_fee, total_fee)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		userID, order.OrderStatus, order.DeliveryFee, totalFee,
	).Scan(&orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (tx.QueryRowContext.Scan): %w", err)
	}

	// batch insert order items
	stmt, err := tx.PrepareContext(ctx, "COPY order_items (order_id, item_id, quantity) FROM STDIN")
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (tx.PrepareContext): %w", err)
	}
	defer stmt.Close()

	for itemID, qty := range qtyMap {
		if _, err = stmt.ExecContext(ctx, orderID, itemID, qty); err != nil {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (stmt.ExecContext): %w", err)
		}
	}

	// fetch full populated data for return
	rows, err = tx.QueryContext(
		ctx,
		`SELECT
			cp.address,
			cp.phone_number,
			r.id,
			r.name,
			r.address,
			i.id,
			i.name,
			i.price,
			oi.quantity
		FROM
			orders o
		JOIN
			customer_profiles cp
			ON o.customer_id = cp.user_id
		JOIN
			order_items oi
			ON o.id = oi.order_id
		JOIN
			items i
			ON oi.item_id = i.id
		JOIN
			restaurants r
			ON i.restaurant_id = r.id
		WHERE
			o.id = $1`,
		orderID,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (tx.QueryContext): %w", err)
	}
	defer rows.Close()

	resultOrder := model.Order{
		ID:          orderID,
		OrderStatus: order.OrderStatus,
		DeliveryFee: order.DeliveryFee,
		TotalFee:    totalFee,
	}
	restoMap := make(map[uuid.UUID]*model.Restaurant)

	for rows.Next() {
		var restoName, restoAddr, itemName, deliveryAddr, customerPhoneNumber string
		var restoID, itemID uuid.UUID
		var itemPrice, itemQty int

		err := rows.Scan(
			&deliveryAddr,
			&customerPhoneNumber,
			&restoID,
			&restoName,
			&restoAddr,
			&itemID,
			&itemName,
			&itemPrice,
			&itemQty,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (fetch full rows.Scan): %w", err)
		}

		resultOrder.DeliveryAddress = deliveryAddr
		resultOrder.CustomerPhoneNumber = customerPhoneNumber

		restaurant, ok := restoMap[restoID]
		if !ok {
			restoMap[restoID] = &model.Restaurant{
				ID:      restoID,
				Name:    restoName,
				Address: restoAddr,
				Items: []model.Item{
					{
						ID:       itemID,
						Name:     itemName,
						Price:    itemPrice,
						Quantity: itemQty,
					},
				},
			}
		} else {
			restaurant.Items = append(restaurant.Items, model.Item{
				ID:       itemID,
				Name:     itemName,
				Price:    itemPrice,
				Quantity: itemQty,
			})
		}
	}

	for _, resto := range restoMap {
		resultOrder.Restaurants = append(resultOrder.Restaurants, *resto)
	}

	return resultOrder, tx.Commit()
}
