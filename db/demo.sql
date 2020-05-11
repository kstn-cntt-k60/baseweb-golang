INSERT INTO party(id, party_type_id, created_by_user_login_id, updated_by_user_login_id)
VALUES
    ('164f0e68-5a01-11ea-b26d-14dda9bea6d7', 2, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7',
        'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    ('47aacbac-5c39-11ea-98a0-14dda9bea6d7', 2, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7',
        'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    ('c2e65572-5d18-11ea-a7e6-14dda9bea6d7', 1, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7',
        'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    ('514377ae-90dc-11ea-9357-40167e8ca7b6', 1, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7',
        'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7');

INSERT INTO person(id, first_name, middle_name, last_name, gender_id, birth_date)
VALUES
    ('c2e65572-5d18-11ea-a7e6-14dda9bea6d7', 'Tùng', 'Thanh', 'Cao', 1, '1997-1-1'),
    ('514377ae-90dc-11ea-9357-40167e8ca7b6', 'Dung', 'Trung', 'Hoang', 1, '1997-12-30');

INSERT INTO customer(id, name)
VALUES
    ('164f0e68-5a01-11ea-b26d-14dda9bea6d7', 'test customer 1'),
    ('47aacbac-5c39-11ea-98a0-14dda9bea6d7', 'test customer 2');

INSERT INTO user_login(id, username, password, person_id)
VALUES
    ('51074f18-5851-11ea-98c8-14dda9bea6d7', 'tungcao',
    '$2a$10$eInQwsiNJW9ZRQPb7aXYIOYIEYJ4TLsYFuTvHcaAd.XDqJ.b/dkR.',
    'c2e65572-5d18-11ea-a7e6-14dda9bea6d7'),
    ('c618b8fc-5872-11ea-adab-14dda9bea6d7', 'test',
    '$2a$10$eInQwsiNJW9ZRQPb7aXYIOYIEYJ4TLsYFuTvHcaAd.XDqJ.b/dkR.',
    '8ed51e8e-59fe-11ea-b26c-14dda9bea6d7');

INSERT INTO user_login_security_group(user_login_id, security_group_id)
VALUES
    ('51074f18-5851-11ea-98c8-14dda9bea6d7', 2),
    ('51074f18-5851-11ea-98c8-14dda9bea6d7', 3);


INSERT INTO facility(id, name, facility_type_id, address)
VALUES
    ('28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 'warehouse', 1, 'Ha Noi'),
    ('3e5b7814-5a02-11ea-b26f-14dda9bea6d7', 'test customer facility 1', 2, 'Ha Noi'),
    ('c0721b28-5a02-11ea-b272-14dda9bea6d7', 'test customer facility 2', 2, 'Ha Noi');

INSERT INTO facility_warehouse(id)
VALUES
    ('28fb8f4a-5a02-11ea-b26e-14dda9bea6d7');

INSERT INTO facility_customer(id, customer_id)
VALUES
    ('3e5b7814-5a02-11ea-b26f-14dda9bea6d7', '164f0e68-5a01-11ea-b26d-14dda9bea6d7'),
    ('c0721b28-5a02-11ea-b272-14dda9bea6d7', '47aacbac-5c39-11ea-98a0-14dda9bea6d7');

INSERT INTO product(id, name, created_by_user_login_id, unit_uom_id)
VALUES
    (1, 'test product 1', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7', 'package'),
    (2, 'test product 2', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7', 'box');

SELECT setval('product_id_seq', 2, true);

INSERT INTO inventory_item(
    product_id, warehouse_id, quantity, quantity_on_hand, unit_cost, currency_uom_id)
VALUES
    (1, '28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 10, 10, 100000, 'vnd'),
    (2, '28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 23, 23, 55000, 'vnd'),
    (1, '28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 13, 13, 55000, 'vnd'),
    (1, '28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 17, 17, 58000, 'vnd');

INSERT INTO warehouse_product_statistics(
    warehouse_id, product_id, inventory_item_count,
    quantity_total, quantity_on_hand, quantity_available)
VALUES
    ('28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 1, 3, 40, 40, 40),
    ('28fb8f4a-5a02-11ea-b26e-14dda9bea6d7', 2, 1, 23, 23, 23);

INSERT INTO salesman (id, created_by_user_login_id) 
VALUES 
    ('c618b8fc-5872-11ea-adab-14dda9bea6d7', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    ('51074f18-5851-11ea-98c8-14dda9bea6d7', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7');


INSERT INTO sales_route_config (id, repeat_week, created_by_user_login_id)
VALUES 
    (1, 1, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    (2, 1, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    (3, 2, 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7');

SELECT setval('sales_route_config_id_seq', 3, true);

INSERT INTO sales_route_config_day (config_id, day)
VALUES
    (1, 1),
    (1, 2),
    (2, 2),
    (3, 1),
    (3, 5);

INSERT INTO sales_route_planning_period (id, from_date, thru_date, created_by_user_login_id)
VALUES
    (1, '2020-12-17', '2020-12-30', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    (2, '2020-12-17', '2021-5-6', 'e66c1e0c-59fb-11ea-b26b-14dda9bea6d7'),
    (3, '2020-12-17', '2020-12-30', '51074f18-5851-11ea-98c8-14dda9bea6d7');

SELECT setval('sales_route_planning_period_id_seq', 3, true);

INSERT INTO sales_route_detail (id, planning_period_id, config_id, customer_id, salesman_id)
VALUES
    (1, 1, 1, '164f0e68-5a01-11ea-b26d-14dda9bea6d7', '51074f18-5851-11ea-98c8-14dda9bea6d7'),
    (2, 2, 2, '164f0e68-5a01-11ea-b26d-14dda9bea6d7', 'c618b8fc-5872-11ea-adab-14dda9bea6d7'),
    (3, 2, 2, '47aacbac-5c39-11ea-98a0-14dda9bea6d7', 'c618b8fc-5872-11ea-adab-14dda9bea6d7');

SELECT setval('sales_route_detail_id_seq', 3, true);

INSERT INTO salesman_checkin_history(sales_route_detail_id)
VALUES
    (1), (2), (3);


    




