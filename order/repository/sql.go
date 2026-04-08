package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *sqlRepository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) CreateOrder(ctx context.Context, userID uuid.UUID, order model.OrderIn) (model.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (r.db.BeginTx): %w", err)
	}
	defer tx.Rollback()

	// map requested items to fetch current prices and calculate totalfee
	var itemIDs []uuid.UUID
	var itemQtys []int
	qtyMap := make(map[uuid.UUID]int)

	for _, item := range order.ItemsIn {
		itemIDs = append(itemIDs, item.ID)
		itemQtys = append(itemQtys, item.Quantity)
		qtyMap[item.ID] = item.Quantity
	}

	// fetch prices from db
	rows, err := tx.QueryContext(
		ctx,
		`UPDATE
			items AS i
		SET
			stock = i.stock - v.qty
		FROM
			(
				SELECT
					unnest($1::uuid[]) AS id,
					unnest($2::int[]) AS qty
			) AS v
		WHERE
			i.id = v.id
		RETURNING
			i.id,
			i.price,
			i.stock`,
		pq.Array(itemIDs),
		pq.Array(itemQtys),
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (tx.QueryContext): %w", err)
	}
	defer rows.Close()

	totalFee := 0
	updatedCount := 0
	for rows.Next() {
		var id uuid.UUID
		var price, updatedStock int

		if err := rows.Scan(&id, &price, &updatedStock); err != nil {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (rows.Scan): %w", err)
		}

		if updatedStock < 0 {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder: insufficient stock for %v", id)
		}

		qty, ok := qtyMap[id]
		if !ok {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder: id %v not found in qtyMap: %w", id, model.ErrNotFound)
		}

		totalFee += price * qty
		updatedCount++
	}

	if updatedCount != len(itemIDs) {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder: one or more item not found: %w", model.ErrNotFound)
	}

	if err := rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (rows.Err): %w", err)
	}

	totalFee += order.DeliveryFee

	// insert order
	var orderID uuid.UUID
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (customer_id, order_status, delivery_fee, total_fee)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		userID, model.OrderStatusSearchingForDriver, order.DeliveryFee, totalFee,
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

	// flush
	if _, err := stmt.ExecContext(ctx); err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (stmt.ExecContext flush): %w", err)
	}
	
	resultOrder, err := r.GetOrderByOrderID(ctx, tx, orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (GetOrderByOrderID): %w", err)
	}

	return resultOrder, tx.Commit()
}

func (r *sqlRepository) GetOrderByOrderID(ctx context.Context, tx *sql.Tx, orderID uuid.UUID) (model.Order, error) {
	rows, err := tx.QueryContext(
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
			oi.quantity,
			o.order_status,
			o.delivery_fee,
			o.total_fee
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

	resultOrder := model.Order{ID: orderID}
	restoMap := make(map[uuid.UUID]*model.Restaurant)

	for rows.Next() {
		var restoName, restoAddr, itemName, deliveryAddr, customerPhoneNumber, orderStatus string
		var restoID, itemID uuid.UUID
		var itemPrice, itemQty, deliveryFee, totalFee int

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
			&orderStatus,
			&deliveryFee,
			&totalFee,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("order.customer_repo.CreateOrder (fetch full rows.Scan): %w", err)
		}

		resultOrder.DeliveryAddress = deliveryAddr
		resultOrder.CustomerPhoneNumber = customerPhoneNumber
		resultOrder.OrderStatus = model.OrderStatus(orderStatus)
		resultOrder.DeliveryFee = deliveryFee
		resultOrder.TotalFee = totalFee

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

	return resultOrder, nil
}

func (r *sqlRepository) GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			dp.user_id,
			COALESCE((
			      SELECT
					ROUND(AVG(r.rating)::numeric, 1)
			      FROM
					orders o
			      JOIN
					ratings r
					ON r.order_id = o.id
			      WHERE
					o.driver_id = dp.user_id
			), 0.0) AS avg_rating
			dp.first_name,
			dp.last_name,
			dp.bike,
			dp.license_plate,
			dp.phone_number,
		FROM
		    	order_applicants op
		JOIN
		    	driver_profiles dp
			ON dp.user_id = op.driver_id
		WHERE
		    	op.order_id = $1;`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("order.customer_repo.GetDrivers (r.db.QueryContext): %w", err)
	}
	defer rows.Close()

	drivers := make([]model.Driver, 0, 8)

	for rows.Next() {
		var drv model.Driver
		var firstName, lastName string

		if err := rows.Scan(
			&drv.ID,
			&drv.AverageRating,
			&firstName,
			&lastName,
			&drv.Bike,
			&drv.LicensePlate,
			&drv.PhoneNumber,
		); err != nil {
			return nil, fmt.Errorf("order.customer_repo.GetDrivers (rows.Scan): %w", err)
		}

		drv.Name = firstName + " " + lastName

		drivers = append(drivers, drv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("order.customer_repo.GetDrivers (rows.Err): %w", err)
	}
	if len(drivers) == 0 {
		return nil, fmt.Errorf("order.customer_repo.GetDrivers: no drivers found: %w", model.ErrNotFound)
	}

	return drivers, nil
}

func (r *sqlRepository) UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error {
	finalAmount := amount

	if reason == model.LedgerReasonCustomerOrder {
		finalAmount = -amount
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO
			ledgers (user_id, amount, reason)
		VALUES
			($1, $2, $3)`,
		userID, finalAmount, string(reason),
	)
	if err != nil {
		return fmt.Errorf("user.repository.UpdateLedger: %w", err)
	}

	return nil
}
