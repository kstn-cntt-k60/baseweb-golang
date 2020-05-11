CREATE OR REPLACE FUNCTION vn_unaccent(TEXT)
RETURNS TEXT AS $$
    SELECT LOWER(TRANSLATE($1,
    '¹²³ÀÁẢẠÂẤẦẨẬẪÃÄÅÆàáảạâấầẩẫậãäåæĀāĂẮẰẲẴẶăắằẳẵặĄąÇçĆćĈĉĊċČčĎďĐđÈÉẸÊẾỀỄỆËèéẹêềếễệëĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨÌÍỈỊÎÏìíỉịîïĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłÑñŃńŅņŇňŉŊŋÒÓỎỌÔỐỒỔỖỘỐỒỔỖỘƠỚỜỞỠỢÕÖòóỏọôốồổỗộơớờỡợởõöŌōŎŏŐőŒœØøŔŕŖŗŘřßŚśŜŝŞşŠšŢţŤťŦŧÙÚỦỤƯỪỨỬỮỰÛÜùúủụûưứừửữựüŨũŪūŬŭŮůŰűŲųŴŵÝýÿŶŷŸŹźŻżŽžёЁ',
    '123AAAAAAAAAAAAAAaaaaaaaaaaaaaaAaAAAAAAaaaaaaAaCcCcCcCcCcDdDdEEEEEEEEEeeeeeeeeeEeEeEeEeEeGgGgGgGgHhHhIIIIIIIiiiiiiiIiIiIiIiIiJjKkkLlLlLlLlLlNnNnNnNnnNnOOOOOOOOOOOOOOOOOOOOOOOooooooooooooooooooOoOoOoEeOoRrRrRrSSsSsSsSsTtTtTtUUUUUUUUUUUUuuuuuuuuuuuuUuUuUuUuUuUuWwYyyYyYZzZzZzеЕ'));
$$ LANGUAGE sql IMMUTABLE;

CREATE OR REPLACE FUNCTION updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION person_full_name_tsvector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.full_name_tsvector =
        setweight(to_tsvector(vn_unaccent(NEW.first_name)), 'A') ||
        setweight(to_tsvector(vn_unaccent(NEW.last_name)), 'B') ||
        setweight(to_tsvector(vn_unaccent(NEW.middle_name)), 'C');
    RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION customer_name_tsvector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.name_tsvector =
        to_tsvector(vn_unaccent(NEW.name));
    RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

DROP TABLE IF EXISTS salesman_checkin_history;
DROP TABLE IF EXISTS sales_route_detail;
DROP TABLE IF EXISTS sales_route_planning_period;
DROP TABLE IF EXISTS sales_route_config_day;
DROP TABLE IF EXISTS sales_route_config;
DROP TABLE IF EXISTS day_of_week;
DROP TABLE IF EXISTS salesman;

DROP TABLE IF EXISTS inventory_item_detail;

DROP TABLE IF EXISTS sale_order_item;
DROP TABLE IF EXISTS sale_order;
DROP TABLE IF EXISTS sale_order_status;

DROP TABLE IF EXISTS warehouse_product_statistics;
DROP TABLE IF EXISTS inventory_item;
DROP TABLE IF EXISTS facility_customer;
DROP TABLE IF EXISTS facility_warehouse;
DROP TABLE IF EXISTS facility;
DROP TABLE IF EXISTS facility_type;

DROP TABLE IF EXISTS product_price;
DROP TABLE IF EXISTS product;

DROP TABLE IF EXISTS security_group_permission;
DROP TABLE IF EXISTS user_login_security_group;
DROP TABLE IF EXISTS security_permission;
DROP TABLE IF EXISTS security_group;
DROP TABLE IF EXISTS user_login;

DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS gender;
DROP TABLE IF EXISTS party;
DROP TABLE IF EXISTS party_type;

DROP TABLE IF EXISTS unit_uom;
DROP TABLE IF EXISTS weight_uom;
DROP TABLE IF EXISTS currency_uom;

CREATE TABLE currency_uom(
    id VARCHAR PRIMARY KEY
);

CREATE TABLE weight_uom(
    id VARCHAR PRIMARY KEY
);

CREATE TABLE unit_uom(
    id VARCHAR PRIMARY KEY
);

CREATE TABLE party_type(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE party(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    party_type_id SMALLINT NOT NULL REFERENCES party_type(id),
    description VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by_user_login_id UUID NOT NULL,
    updated_by_user_login_id UUID NOT NULL
);

CREATE TRIGGER party_updated_at BEFORE UPDATE ON
    party FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE gender(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE person(
    id UUID PRIMARY KEY REFERENCES party(id),
    first_name VARCHAR NOT NULL,
    middle_name VARCHAR NOT NULL DEFAULT '',
    last_name VARCHAR NOT NULL,
    full_name_tsvector TSVECTOR NOT NULL DEFAULT to_tsvector(''),
    gender_id SMALLINT NOT NULL REFERENCES gender(id),
    birth_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER person_updated_at BEFORE UPDATE ON
    person FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TRIGGER person_full_name_tsvector
    BEFORE INSERT OR UPDATE ON person FOR EACH ROW
    EXECUTE PROCEDURE person_full_name_tsvector();

CREATE INDEX idx_person_full_name_tsvector ON person 
    USING gin(full_name_tsvector);

CREATE TABLE customer(
    id UUID PRIMARY KEY REFERENCES party(id),
    name VARCHAR NOT NULL,
    name_tsvector TSVECTOR NOT NULL DEFAULT to_tsvector(''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER customer_updated_at BEFORE UPDATE ON
    customer FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TRIGGER customer_name_tsvector
    BEFORE INSERT OR UPDATE ON customer FOR EACH ROW
    EXECUTE PROCEDURE customer_name_tsvector();

CREATE INDEX idx_customer_name_tsvector ON customer 
    USING gin(name_tsvector);

CREATE TABLE user_login(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    person_id UUID NOT NULL REFERENCES person(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER user_login_updated_at BEFORE UPDATE ON
    user_login FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE security_group(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE security_permission(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_login_security_group(
    user_login_id UUID REFERENCES user_login(id),
    security_group_id SMALLINT REFERENCES security_group(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_login_id, security_group_id)
);

CREATE TABLE security_group_permission(
    security_group_id SMALLINT REFERENCES security_group(id),
    security_permission_id SMALLINT REFERENCES security_permission(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (security_group_id, security_permission_id)
);

CREATE TABLE product(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),
    description VARCHAR NOT NULL DEFAULT '',

    weight DECIMAL,
    weight_uom_id VARCHAR NOT NULL REFERENCES weight_uom(id) DEFAULT 'kg',

    unit_uom_id VARCHAR NOT NULL REFERENCES unit_uom(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER product_updated_at BEFORE UPDATE ON
    product FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE product_price(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    product_id INTEGER NOT NULL REFERENCES product(id),

    currency_uom_id VARCHAR NOT NULL REFERENCES currency_uom(id),
    price DECIMAL NOT NULL,

    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),

    effective_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER product_price_updated_at BEFORE UPDATE ON
    product_price FOR EACH ROW EXECUTE PROCEDURE updated_at_column();


CREATE TABLE facility_type(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE facility(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    name VARCHAR NOT NULL UNIQUE,
    facility_type_id SMALLINT NOT NULL REFERENCES facility_type(id),
    address VARCHAR NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER facility_updated_at BEFORE UPDATE ON
    facility FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE facility_warehouse(
    id UUID PRIMARY KEY REFERENCES facility(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER facility_warehouse_updated_at BEFORE UPDATE ON
    facility_warehouse FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE facility_customer(
    id UUID PRIMARY KEY REFERENCES facility(id),
    customer_id UUID NOT NULL REFERENCES customer(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER facility_customer_updated_at BEFORE UPDATE ON
    facility_customer FOR EACH ROW EXECUTE PROCEDURE updated_at_column();


CREATE TABLE inventory_item(
    id BIGSERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES product(id),
    warehouse_id UUID NOT NULL REFERENCES facility_warehouse(id),

    quantity DECIMAL NOT NULL,
    quantity_on_hand DECIMAL NOT NULL,
    unit_cost DECIMAL NOT NULL,
    currency_uom_id VARCHAR NOT NULL REFERENCES currency_uom(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_inventory_item_export ON
    inventory_item(product_id, warehouse_id, quantity_on_hand);

CREATE TRIGGER inventory_item_updated_at BEFORE UPDATE ON
    inventory_item FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE warehouse_product_statistics(
    warehouse_id UUID REFERENCES facility_warehouse(id),
    product_id INTEGER REFERENCES product(id),

    inventory_item_count BIGINT NOT NULL,
    quantity_total DECIMAL NOT NULL,
    quantity_on_hand DECIMAL NOT NULL,
    quantity_available DECIMAL NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT pk_warehouse_product PRIMARY KEY (warehouse_id, product_id)
);

CREATE TRIGGER warehouse_product_statistics_updated_at BEFORE UPDATE ON
    warehouse_product_statistics FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

ALTER TABLE warehouse_product_statistics
    ADD CONSTRAINT quantity_total_non_neg CHECK (quantity_total >= 0);

ALTER TABLE warehouse_product_statistics
    ADD CONSTRAINT quantity_on_hand_non_neg CHECK (quantity_on_hand >= 0);

ALTER TABLE warehouse_product_statistics
    ADD CONSTRAINT quantity_available_non_neg CHECK (quantity_available >= 0);

CREATE TABLE sale_order_status(
    id SMALLINT PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sale_order(
    id BIGSERIAL PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customer(id),
    original_warehouse_id UUID NOT NULL REFERENCES facility_warehouse(id),
    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),

    ship_to_address VARCHAR NOT NULL,
    ship_to_facility_customer_id UUID REFERENCES facility_customer(id),

    sale_order_status_id SMALLINT NOT NULL REFERENCES sale_order_status(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER sale_order_updated_at BEFORE UPDATE ON
    sale_order FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE sale_order_item(
    sale_order_id BIGINT REFERENCES sale_order(id),
    sale_order_seq SMALLINT,

    product_price_id UUID NOT NULL REFERENCES product_price(id),
    quantity DECIMAL NOT NULL,
    exported BOOL NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT pk_sale_order_item PRIMARY KEY (sale_order_id, sale_order_seq)
);

CREATE TRIGGER sale_order_item_updated_at BEFORE UPDATE ON
    sale_order_item FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE inventory_item_detail(
    id BIGSERIAL PRIMARY KEY,
    inventory_item_id BIGINT NOT NULL REFERENCES inventory_item(id),
    exported_quantity DECIMAL NOT NULL,
    effective_from TIMESTAMPTZ NOT NULL DEFAULT now(),

    sale_order_id BIGINT NOT NULL,
    sale_order_seq SMALLINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_inventory_item_detail_sale_order_item
        FOREIGN KEY (sale_order_id, sale_order_seq)
        REFERENCES sale_order_item(sale_order_id, sale_order_seq)
);

CREATE TRIGGER inventory_item_detail_updated_at BEFORE UPDATE ON
    inventory_item_detail FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE salesman(
    id UUID PRIMARY KEY REFERENCES user_login(id),
    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER salesman_updated_at BEFORE UPDATE ON
    salesman FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE day_of_week(
    day SMALLINT PRIMARY KEY
);

CREATE TABLE sales_route_config(
    id SERIAL PRIMARY KEY,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    repeat_week SMALLINT NOT NULL,
    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER sales_route_config_updated_at BEFORE UPDATE ON
    sales_route_config FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE sales_route_config_day(
    config_id INT REFERENCES sales_route_config(id),
    day SMALLINT REFERENCES day_of_week(day),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_sales_route_config_day PRIMARY KEY (config_id, day)
);

CREATE TABLE sales_route_planning_period(
    id SERIAL PRIMARY KEY,
    from_date DATE NOT NULL,
    thru_date DATE NOT NULL,
    created_by_user_login_id UUID NOT NULL REFERENCES user_login(id),
    description VARCHAR NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER sales_route_planning_period_updated_at BEFORE UPDATE ON
    sales_route_planning_period FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE sales_route_detail(
    id SERIAL PRIMARY KEY,
    config_id INT NOT NULL REFERENCES sales_route_config(id),
    planning_period_id INT NOT NULL REFERENCES sales_route_planning_period(id),
    customer_id UUID NOT NULL REFERENCES customer(id),
    salesman_id UUID NOT NULL REFERENCES salesman(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT u_sales_route_detail UNIQUE (planning_period_id, customer_id, salesman_id)
);

CREATE TRIGGER sales_route_detail_updated_at BEFORE UPDATE ON
    sales_route_detail FOR EACH ROW EXECUTE PROCEDURE updated_at_column();

CREATE TABLE salesman_checkin_history(
    sales_route_detail_id INT PRIMARY KEY REFERENCES sales_route_detail(id),
    checkin_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

