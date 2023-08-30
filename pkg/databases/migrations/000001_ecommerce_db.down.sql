BEGEN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "users";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON "oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_products_table ON "products";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_orders_table ON "orders";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_images_table ON "images";

DROP FUNCTION IF EXISTS set_updated_at_coloum();

DROP TABLE IF EXISTS "users" CASCADE;

DROP SEQUENCE IF EXISTS  users_id_seq ;
DROP SEQUENCE IF EXISTS  products_id_seq ;
DROP SEQUENCE IF EXISTS  orders_id_seq ;

DROP TYPE IF EXISTS "order_status" ;

COMMIT