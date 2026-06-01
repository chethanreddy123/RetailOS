-- name: SearchProducts :many
SELECT p.product_id, p.name, p.company_name, p.sku, p.hsn_code, p.created_at,
       EXISTS (
         SELECT 1 FROM batches b
         WHERE b.product_id = p.product_id
           AND b.expiry_date > CURRENT_DATE
           AND (b.purchase_qty - b.sold_qty) > 0
       ) AS has_active_stock
FROM products p
WHERE $1::text = ''
   OR p.name ILIKE '%' || $1 || '%'
   OR p.company_name ILIKE '%' || $1 || '%'
ORDER BY p.name
LIMIT $2 OFFSET $3;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE $1::text = ''
   OR name ILIKE '%' || $1 || '%'
   OR company_name ILIKE '%' || $1 || '%';

-- name: GetProduct :one
SELECT * FROM products WHERE product_id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, company_name, sku, hsn_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, company_name = $3, sku = $4, hsn_code = $5
WHERE product_id = $1
RETURNING *;
