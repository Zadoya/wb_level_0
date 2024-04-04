package repository

import (
	"database/sql"
	"wb_level_0/internal/order"
)

type SqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(db *sql.DB) *SqlRepository {
	return &SqlRepository{
		db: db,
	}
}

func (repo *SqlRepository) GetById(orderUID string) (*order.Order, error) {
	var ord order.Order

	if err := repo.db.QueryRow(getOrderByOrderUIDQuery, orderUID).Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerID,
		&ord.DeliveryService,
		&ord.Shardkey,
		&ord.SmID,
		&ord.DateCreated,
		&ord.OofShard,
	); err != nil {
		return nil, err
	}

	if err := repo.db.QueryRow(getDeliveryByOrderUIDQuery, orderUID).Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	); err != nil {
		return nil, err
	}

	if err := repo.db.QueryRow(getPaymentByOrderUIDQuery, orderUID).Scan(
		&ord.Payment.Transaction,
		&ord.Payment.RequestID,
		&ord.Payment.Currency,
		&ord.Payment.Provider,
		&ord.Payment.Amount,
		&ord.Payment.PaymentDt,
		&ord.Payment.Bank,
		&ord.Payment.DeliveryCost,
		&ord.Payment.GoodsTotal,
		&ord.Payment.CustomFee,
	); err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(getItemsByOrderUIDQuery, orderUID)
	if err != nil {
		return nil, err
	}

	item := order.Item{}

	for rows.Next() {
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		ord.Items = append(ord.Items, item)
	}

	return &ord, nil
}

func (repo *SqlRepository) SaveOrder(order *order.Order) error {

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(createOrderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(createDeliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(createPaymentsQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	); err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range order.Items {
		if _, err = tx.Exec(createItemQuery,
			order.OrderUID,
			v.ChrtID,
			v.TrackNumber,
			v.Price,
			v.Rid,
			v.Name,
			v.Sale,
			v.Size,
			v.TotalPrice,
			v.NmID,
			v.Brand,
			v.Status,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *SqlRepository) GetAllOrders() ([]*order.Order, error) {
	rows, err := repo.db.Query(cacheRecoveryQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order

	for rows.Next() {
		var odr order.Order
		var delivery order.Delivery
		var payment order.Payment
		err := rows.Scan(
			&odr.OrderUID,
			&odr.TrackNumber,
			&odr.Entry,
			&odr.Locale,
			&odr.InternalSignature,
			&odr.CustomerID,
			&odr.DeliveryService,
			&odr.Shardkey,
			&odr.SmID,
			&odr.DateCreated,
			&odr.OofShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDt,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}
		odr.Delivery = delivery
		odr.Payment = payment

		orderItems, err := repo.getOrderItemsByOrderUID(odr.OrderUID)
		if err != nil {
			return nil, err
		}
		odr.Items = *orderItems

		orders = append(orders, &odr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (repo *SqlRepository) getOrderItemsByOrderUID(orderUID string) (*[]order.Item, error) {
	rows, err := repo.db.Query(getItemsByOrderUIDQuery, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []order.Item

	for rows.Next() {
		var orderItem order.Item

		err := rows.Scan(
			&orderItem.ChrtID,
			&orderItem.TrackNumber,
			&orderItem.Price,
			&orderItem.Rid,
			&orderItem.Name,
			&orderItem.Sale,
			&orderItem.Size,
			&orderItem.TotalPrice,
			&orderItem.NmID,
			&orderItem.Brand,
			&orderItem.Status,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orderItems, nil
}