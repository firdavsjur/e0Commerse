CREATE TABLE promoCodes(
    promo_code_id int not null,
    promo_code_name varchar not null,
    discount numeric not null,
    discount_type varchar not null,
    order_limit_price numeric not null,
);