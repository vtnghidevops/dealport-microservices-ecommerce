# Database Diagram for E-commerce Microservices

## Authentication Service Database (PostgreSQL)

```
┌───────────────────────┐
│      auth_users       │
├───────────────────────┤
│ id (PK)               │
│ email                 │
│ password_hash         │
│ active                │
│ created_at            │
│ updated_at            │
└─────────┬─────────────┘
          │
          │ 1:1
          ▼
┌───────────────────────┐
│     auth_tokens       │
├───────────────────────┤
│ id (PK)               │
│ user_id (FK)          │
│ token                 │
│ token_hash            │
│ expiry                │
│ created_at            │
└───────────────────────┘
```

## User Service Database (PostgreSQL)

```
┌───────────────────────┐
│        users          │
├───────────────────────┤
│ id (PK)               │
│ email                 │
│ first_name            │
│ last_name             │
│ display_name          │
│ phone                 │
│ profile_image         │
│ role                  │
│ status                │
│ active                │
│ created_at            │
│ updated_at            │
│ last_login            │
│ deleted_at            │
│ username              │
│ date_of_birth         │
│ gender                │
│ cart_id               │
│ order_count           │
│ total_spend           │
└─────────┬─────────────┘
          │
          │ 1:N
          ▼
┌───────────────────────┐
│       addresses       │
├───────────────────────┤
│ id (PK)               │
│ user_id (FK)          │
│ name                  │
│ phone                 │
│ line1                 │
│ city                  │
│ state                 │
│ postal_code           │
│ country               │
│ is_default            │
│ address_type          │
└───────────────────────┘
          │
          │ 1:N
┌─────────┴─────────────┐
│    payment_methods    │
├───────────────────────┤
│ id (PK)               │
│ user_id (FK)          │
│ type                  │
│ provider              │
│ account_number        │
│ expiry_date           │
│ is_default            │
└───────────────────────┘
          │
          │ 1:N
┌─────────┴─────────────┐
│     wishlist_items    │
├───────────────────────┤
│ id (PK)               │
│ user_id (FK)          │
│ product_id            │
│ added_at              │
│ notes                 │
└───────────────────────┘
```

## Product Service Database (PostgreSQL)

```
┌───────────────────────┐
│       categories      │
├───────────────────────┤
│ id (PK)               │
│ name                  │
│ slug                  │
│ description           │
│ parent_id (FK)        │
│ image_url             │
│ created_at            │
│ updated_at            │
└─────────┬─────────────┘
          │
          │ 1:N
          ▼
┌───────────────────────┐
│       products        │
├───────────────────────┤
│ id (PK)               │
│ type                  │
│ name                  │
│ description           │
│ slug                  │
│ price                 │
│ original_price        │
│ discount              │
│ image_url             │
│ category_id (FK)      │
│ category_slug         │
│ stock_quantity        │
│ brand                 │
│ created_at            │
│ updated_at            │
└─────────┬─────────────┘
          │
          │ 1:N
          ▼
┌───────────────────────┐
│    product_images     │
├───────────────────────┤
│ id (PK)               │
│ product_id (FK)       │
│ url                   │
│ is_primary            │
│ display_order         │
│ created_at            │
└───────────────────────┘
          │
          │ 1:N
┌─────────┴─────────────┐
│    product_reviews    │
├───────────────────────┤
│ id (PK)               │
│ product_id (FK)       │
│ user_id               │
│ user_name             │
│ rating                │
│ comment               │
│ created_at            │
│ user_avatar           │
└───────────────────────┘
          │
          │ 1:N
┌─────────┴─────────────┐
│     product_tags      │
├───────────────────────┤
│ id (PK)               │
│ product_id (FK)       │
│ tag                   │
└───────────────────────┘
```

## Cart Service Database (Redis)

```
[Key: user:{user_id}:cart]
┌───────────────────────┐
│         cart          │
├───────────────────────┤
│ id                    │
│ user_id               │
│ items (JSON array)    │
│ totals (JSON object)  │
│ coupon_code           │
│ discount_amount       │
│ created_at            │
│ updated_at            │
└───────────────────────┘

[Key: coupons:{coupon_code}]
┌───────────────────────┐
│        coupon         │
├───────────────────────┤
│ id                    │
│ code                  │
│ discount              │
│ discount_type         │
│ min_order_amount      │
│ max_usage             │
│ usage_count           │
│ valid_from            │
│ valid_to              │
│ is_active             │
│ description           │
│ created_at            │
│ updated_at            │
└───────────────────────┘
```

## Checkout Service Database (MongoDB)

```
┌───────────────────────┐
│        orders         │
├───────────────────────┤
│ _id                   │
│ user_id               │
│ order_number          │
│ status                │
│ items (array)         │
│ billing_info (object) │
│ shipping_info (object)│
│ payment_info (object) │
│ totals (object)       │
│ coupon_code           │
│ notes                 │
│ created_at            │
│ updated_at            │
└───────────────────────┘
```

## Logger Service Database (MongoDB)

```
┌───────────────────────┐
│         logs          │
├───────────────────────┤
│ _id                   │
│ name                  │
│ data                  │
│ created_at            │
└───────────────────────┘
```

## Service Relationships

```
┌───────────────────┐            ┌───────────────────┐
│  Authentication   │───────────▶│       User        │
│     Service       │            │     Service       │
└───────────────────┘            └─────────┬─────────┘
                                           │
                                           │
                                           ▼
┌───────────────────┐            ┌───────────────────┐
│       Cart        │◀───────────│      Product      │
│     Service       │            │     Service       │
└─────────┬─────────┘            └───────────────────┘
          │
          │
          ▼
┌───────────────────┐
│     Checkout      │
│     Service       │
└───────────────────┘
``` 