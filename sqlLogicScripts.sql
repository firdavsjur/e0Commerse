-- 2-task --



SELECT
    st.first_name || ' ' || st.last_name as Sotuvchi,
    ca.category_name, 
    ARRAY_AGG(distinct(pr.product_name)) as product_name,
    count(oi.quantity) soni,
    pr.list_price*count(oi.quantity) as ObshiySumma,
    ARRAY_AGG(ord.order_date) as Sana

    
   
    

FROM staffs AS st
JOIN orders as ord on ord.staff_id = st.staff_id
JOIN order_items as oi on oi.order_id = ord.order_id
JOIN products as pr on pr.product_id = oi.product_id
JOIN categories as ca on ca.category_id = pr.category_id
where ord.order_status=4 and st.staff_id = 1
GROUP BY Sotuvchi,ca.category_name,pr.product_name,pr.list_price

;



-- 1-task --

UPDATE stocks 
set quantity = quantity - $4
where 
store_id = $1 AND product_id=$3 ;


UPDATE stocks
set quantity = quantity+$4

WHERE
store_id = $2 and product_id = $3



-- 4-task --

SELECT
*


FROM orders as ord
JOIN order_items as oi on oi.order_id = ord.order_id



 where order_id = $1





 INSERT INTO order_items(
			order_id, 
			item_id, 
			product_id,
			quantity,
			list_price,
			discount,
			sell_price
		)
		VALUES (
			1, 
			(
				SELECT COALESCE(MAX(item_id), 0) + 1 FROM order_items WHERE order_id = 1
			)
			, 3, 2, 
			(
				SELECT list_price*2 FROM products where product_id=3
			)
			, 0.1,
			(
				SELECT (list_price*2)*(1-0.1) FROM products where product_id=3
			)
		)

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
			1, 
			(
				SELECT COALESCE(MAX(item_id), 0) + 1 FROM order_items WHERE order_id = 1
			)
			, 3, 2, 
			(
				SELECT list_price*2 FROM products where product_id=3
			)
			, 0.2,
			(
				SELECT (list_price*2)*(1-0.2) FROM products where product_id=3
			),1
		)



SELECT 

st.first_name || ' ' || st.last_name as Xodim,
ca.category_name,
pr.product_name,
oi.quantity,
oi.sell_price as Summa,
ord.order_date


FROM orders as ord
JOIN staffs as st on st.staff_id = ord.staff_id
JOIN order_items as oi on oi.order_id = ord.order_id
JOIN products as pr on pr.product_id = oi.product_id
JOIN categories as ca on ca.category_id = pr.category_id


WHERE ord.staff_id =2 and ord.order_status = 4 and ord.order_id = 1

;



SELECT 
st.first_name || ' ' || st.last_name as Xodim,
ca.category_name,
pr.product_name,
SUM(oi.quantity) as JamiSoni,
SUM(oi.sell_price) as Summa



FROM orders as ord
JOIN staffs as st on st.staff_id = ord.staff_id
JOIN order_items as oi on oi.order_id = ord.order_id
JOIN products as pr on pr.product_id = oi.product_id
JOIN categories as ca on ca.category_id = pr.category_id


WHERE ord.order_status = 4
GROUP BY Xodim,ca.category_name,pr.product_name
ORDER BY product_name ASC

;