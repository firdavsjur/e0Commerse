package postgresql

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *orderRepo {
	return &orderRepo{
		db: db,
	}
}

func (r *orderRepo) Create(ctx context.Context, req *models.CreateOrder) (int, error) {
	var (
		query string
		id    int
	)

	query = `
		INSERT INTO orders(
			order_id, 
			customer_id, 
			order_status,
			order_date,
			required_date,
			shipped_date,
			store_id,
			staff_id
		)
		VALUES (
			(
				SELECT MAX(order_id) + 1 FROM orders
			)
			, $1, $2, now()::date, $3, $4, $5, $6) RETURNING order_id
	`
	fmt.Println(query)

	err := r.db.QueryRow(ctx, query,
		helper.NewNullInt32(req.CustomerId),
		req.OrderStatus,
		req.RequiredDate,
		helper.NewNullString(req.ShippedDate),
		req.StoreId,
		req.StaffId,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *orderRepo) GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error) {

	var (
		query string
		order models.Order
	)

	query = `
	WITH order_item_data AS (
		SELECT
			oi.order_id AS order_id,
			JSONB_AGG (
				JSONB_BUILD_OBJECT (
					'order_id', oi.order_id,
					'item_id', oi.item_id,
					'product_id', oi.product_id,
					'quantity', oi.quantity,
					'list_price', oi.list_price,
					'discount', oi.discount,
					'sell_price',COALESCE(oi.sell_price,oi.list_price-(oi.list_price*oi.discount)),
					'store_id',oi.store_id
				)
			) AS order_items
	
		FROM order_items AS oi
		WHERE oi.order_id = $1
		GROUP BY oi.order_id
	)
	SELECT
		o.order_id, 
		o.customer_id,
	
		c.customer_id,
		c.first_name,
		c.last_name,
		COALESCE(c.phone, ''),
		c.email,
		COALESCE(c.street, ''),
		COALESCE(c.city, ''),
		COALESCE(c.state, ''),
		COALESCE(c.zip_code, 0),
		
		o.order_status,
		CAST(o.order_date::timestamp AS VARCHAR),
		CAST(o.required_date::timestamp AS VARCHAR),
		COALESCE(CAST(o.shipped_date::timestamp AS VARCHAR), ''),
		o.store_id,
	
		s.store_id,
		s.store_name,
		COALESCE(s.phone, ''),
		COALESCE(s.email, ''),
		COALESCE(s.street, ''),
		COALESCE(s.city, ''),
		COALESCE(s.state, ''),
		COALESCE(s.zip_code, ''),
	
		o.staff_id,
		st.staff_id,
		st.first_name,
		st.last_name,
		st.email,
		COALESCE(st.phone, ''),
		st.active,
		st.store_id,
		COALESCE(st.manager_id, 0),
	
		oi.order_items
	
	FROM orders AS o
	JOIN customers AS c ON c.customer_id = o.customer_id
	JOIN stores AS s ON s.store_id = o.store_id
	JOIN staffs AS st ON st.staff_id = o.staff_id
	JOIN order_item_data AS oi ON oi.order_id = o.order_id
	WHERE o.order_id = $1;
	`
	order.CustomerData = &models.Customer{}
	order.StoreData = &models.Store{}
	order.StaffData = &models.Staff{}
	orderItemObject := pgtype.JSON{}

	err := r.db.QueryRow(ctx, query, req.OrderId).Scan(
		&order.OrderId,
		&order.CustomerId,
		&order.CustomerData.CustomerId,
		&order.CustomerData.FirstName,
		&order.CustomerData.LastName,
		&order.CustomerData.Phone,
		&order.CustomerData.Email,
		&order.CustomerData.Street,
		&order.CustomerData.City,
		&order.CustomerData.State,
		&order.CustomerData.ZipCode,

		&order.OrderStatus,
		&order.OrderDate,
		&order.RequiredDate,
		&order.ShippedDate,

		&order.StoreId,

		&order.StoreData.StoreId,
		&order.StoreData.StoreName,
		&order.StoreData.Phone,
		&order.StoreData.Email,
		&order.StoreData.Street,
		&order.StoreData.City,
		&order.StoreData.State,
		&order.StoreData.ZipCode,
		&order.StaffId,
		&order.StaffData.StaffId,
		&order.StaffData.FirstName,
		&order.StaffData.LastName,
		&order.StaffData.Email,
		&order.StaffData.Phone,
		&order.StaffData.Active,
		&order.StaffData.StoreId,
		&order.StaffData.ManagerId,

		&orderItemObject,
	)
	if err != nil {
		return nil, err
	}

	orderItemObject.AssignTo(&order.OrderItems)

	return &order, nil
}

func (r *orderRepo) GetList(ctx context.Context, req *models.GetListOrderRequest) (resp *models.GetListOrderResponse, err error) {

	resp = &models.GetListOrderResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			o.order_id, 
			o.customer_id,
			c.customer_id,

			c.first_name,
			c.last_name,
			COALESCE(c.phone, ''),
			c.email,
			COALESCE(c.street, ''),
			COALESCE(c.city, ''),
			COALESCE(c.state, ''),
			COALESCE(c.zip_code, 0),

			o.order_status,
			CAST(o.order_date::timestamp AS VARCHAR),
			CAST(o.required_date::timestamp AS VARCHAR),
			COALESCE(CAST(o.shipped_date::timestamp AS VARCHAR), ''),
			o.store_id,

			s.store_id,
			s.store_name,
			COALESCE(s.phone, ''),
			COALESCE(s.email, ''),
			COALESCE(s.street, ''),
			COALESCE(s.city, ''),
			COALESCE(s.state, ''),
			COALESCE(s.zip_code, ''),

			o.staff_id,
			st.staff_id,
			st.first_name,
			st.last_name,
			st.email,
			COALESCE(st.phone, ''),
			st.active,
			st.store_id,
			COALESCE(st.manager_id, 0)

		FROM orders AS o
		JOIN customers AS c ON c.customer_id = o.customer_id
		JOIN stores AS s ON s.store_id = o.store_id
		JOIN staffs AS st ON st.staff_id = o.staff_id
	`

	if len(req.Search) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		order.CustomerData = &models.Customer{}
		order.StoreData = &models.Store{}
		order.StaffData = &models.Staff{}

		err = rows.Scan(
			&resp.Count,
			&order.OrderId,
			&order.CustomerId,
			&order.CustomerData.CustomerId,
			&order.CustomerData.FirstName,
			&order.CustomerData.LastName,
			&order.CustomerData.Phone,
			&order.CustomerData.Email,
			&order.CustomerData.Street,
			&order.CustomerData.City,
			&order.CustomerData.State,
			&order.CustomerData.ZipCode,

			&order.OrderStatus,
			&order.OrderDate,
			&order.RequiredDate,
			&order.ShippedDate,

			&order.StoreId,

			&order.StoreData.StoreId,
			&order.StoreData.StoreName,
			&order.StoreData.Phone,
			&order.StoreData.Email,
			&order.StoreData.Street,
			&order.StoreData.City,
			&order.StoreData.State,
			&order.StoreData.ZipCode,
			&order.StaffId,
			&order.StaffData.StaffId,
			&order.StaffData.FirstName,
			&order.StaffData.LastName,
			&order.StaffData.Email,
			&order.StaffData.Phone,
			&order.StaffData.Active,
			&order.StaffData.StoreId,
			&order.StaffData.ManagerId,
		)
		if err != nil {
			return nil, err
		}

		resp.Orders = append(resp.Orders, &order)
	}

	return resp, nil
}

func (r *orderRepo) Update(ctx context.Context, req *models.UpdateOrder) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
		orders
		SET
			order_id = :order_id, 
			customer_id = :customer_id, 
			order_status = :order_status, 
			order_date = :order_date,
			required_date = :required_date,
			shipped_date = :shipped_date,
			store_id = :store_id,
			staff_id = :staff_id
		WHERE order_id = :order_id
	`

	params = map[string]interface{}{
		"order_id":      req.OrderId,
		"customer_id":   req.CustomerId,
		"order_status":  req.OrderStatus,
		"order_date":    req.OrderDate,
		"required_date": req.RequiredDate,
		"shipped_date":  helper.NewNullString(req.ShippedDate),
		"store_id":      req.StoreId,
		"staff_id":      req.StaffId,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *orderRepo) UpdatePatch(ctx context.Context, req *models.PatchRequest) (int64, error) {
	var (
		query string
		set   string
	)

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	i := 0
	for key := range req.Fields {
		if i == len(req.Fields)-1 {
			set += fmt.Sprintf(" %s = :%s ", key, key)
		} else {
			set += fmt.Sprintf(" %s = :%s, ", key, key)
		}
		i++
	}

	query = `
		UPDATE
		orders
		SET
		` + set + `
		WHERE order_id = :order_id
	`

	req.Fields["order_id"] = req.ID

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	fmt.Println(query)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *orderRepo) Delete(ctx context.Context, req *models.OrderPrimaryKey) (int64, error) {
	query := `
		DELETE 
		FROM orders
		WHERE order_id = $1
	`

	result, err := r.db.Exec(ctx, query, req.OrderId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// ------------------------------------------------------------------------------------------------------------
func (r *orderRepo) AddOrderItem(ctx context.Context, req *models.CreateOrderItem) error {
	haveProducQuantity := 0
	newOrderItemId := 0
	var sellPrice float64
	var listPrice float64
	query := `
		SELECT 
		quantity
		FROM stocks WHERE store_id = $1 and product_id = $2
	`
	err := r.db.QueryRow(ctx, query, req.StoreId, req.ProductId).Scan(
		&haveProducQuantity,
	)
	if err != nil {
		return errors.New("skladdan malumot olishda error")
	}
	query = `
	(
		SELECT COALESCE(MAX(item_id), 0) + 1 FROM order_items WHERE order_id = $1
	)
	`
	err = r.db.QueryRow(ctx, query, req.OrderId).Scan(
		&newOrderItemId,
	)
	if err != nil {
		return errors.New("skladdan malumot olishda error")
	}

	query = `
		(
			SELECT list_price*$1 FROM products where product_id=$2
		)
	`
	err = r.db.QueryRow(ctx, query, req.Quantity, req.ProductId).Scan(
		&listPrice,
	)
	if err != nil {
		return errors.New("skladdan malumot olishda error")
	}

	query = `
		(
			
				SELECT (list_price*$1)*(1-$2) FROM products where product_id=$3
			
		)
	`
	err = r.db.QueryRow(ctx, query, req.Quantity, req.Discount, req.ProductId).Scan(
		&sellPrice,
	)
	if err != nil {
		return errors.New("skladdan malumot olishda error")
	}

	if haveProducQuantity >= req.Quantity {
		query = `
		INSERT INTO order_items(
			order_id, 
			item_id, 
			product_id,
			quantity,
			list_price,
			discount,
			sell_price,
			store_id
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8
		)
	`
		_, err = r.db.Exec(ctx, query,
			req.OrderId,
			newOrderItemId,
			req.ProductId,
			req.Quantity,
			listPrice,
			req.Discount,
			sellPrice,
			req.StoreId,
		)
		fmt.Println(query)
		if err != nil {
			return err
		}
		query = `
			UPDATE stocks
			set quantity = quantity-$1
			WHERE store_id = $2 and product_id = $3
		`
		_, err = r.db.Exec(ctx, query,
			req.Quantity,
			req.StoreId,
			req.ProductId,
		)

		if err != nil {
			return err
		}

	} else {
		return errors.New("skladda product yetarlimas")
	}

	return nil
}

func (r *orderRepo) RemoveOrderItem(ctx context.Context, req *models.OrderItemPrimaryKey) error {

	query := `
		DELETE FROM order_items WHERE order_id = $1 AND item_id = $2
	`
	_, err := r.db.Exec(ctx, query, req.OrderId, req.ItemId)

	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepo) TotalOrderSum(ctx context.Context, req *models.TotalOrderSumRequest) (*models.TotalOrderSumResponse, string, error) {
	var (
		query             string
		resp              models.TotalOrderSumResponse
		discount_pr       float64
		discount_type_pr  string
		order_limit_price float64
	)

	query = `
		SELECT 
		discount,
		discount_type,
		order_limit_price
		FROM promoCodes WHERE promo_code_name=$1
	`
	err := r.db.QueryRow(ctx, query, req.PromoCodeName).Scan(
		&discount_pr,
		&discount_type_pr,
		&order_limit_price,
	)
	if err != nil {
		log.Println("Promocode didnt work")

	}
	query = `
		SELECT
			SUM(oi.sell_price)
			

		FROM order_items as oi
		WHERE oi.order_id = $1
	`

	err = r.db.QueryRow(ctx, query, req.OrderId).Scan(
		&resp.Total,
	)
	if err != nil {
		log.Println("OrderItem summ error")
		return nil, "orderItem summ error", err
	}

	if resp.Total >= order_limit_price && discount_type_pr == "fixed" {
		resp.Total -= discount_pr
		return &resp, "fixed promocode used", nil
	}
	if resp.Total >= order_limit_price && discount_type_pr == "percent" {
		resp.Total *= 1 - discount_pr
		return &resp, "persent promocode used", nil
	}

	return &resp, "PromoCode ishlatilmadi", nil
}
