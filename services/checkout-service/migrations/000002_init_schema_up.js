// MongoDB Schema Definition for Checkout Service
db = db.getSiblingDB('checkout');

print("Starting MongoDB schema initialization for checkout service...");

// Create collections with validation
db.createCollection("orders", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["id", "user_id", "order_number", "status", "items", "billing_info", "shipping_info", "totals", "created_at", "updated_at"],
      properties: {
        id: {
          bsonType: "string",
          description: "The unique identifier for the order"
        },
        user_id: {
          bsonType: "string",
          description: "The user who placed the order"
        },
        order_number: {
          bsonType: "string",
          description: "Human-readable order identifier"
        },
        status: {
          bsonType: "string",
          enum: ["pending", "processing", "paid", "shipped", "delivered", "cancelled", "refunded"],
          description: "The current status of the order"
        },
        items: {
          bsonType: "array",
          description: "The items in the order",
          items: {
            bsonType: "object",
            required: ["id", "product_id", "name", "price", "quantity", "subtotal"],
            properties: {
              id: { bsonType: "string" },
              product_id: { bsonType: "string" },
              name: { bsonType: "string" },
              price: { bsonType: "double" },
              quantity: { bsonType: "int" },
              subtotal: { bsonType: "double" },
              image_url: { bsonType: "string" }
            }
          }
        },
        billing_info: {
          bsonType: "object",
          required: ["first_name", "last_name", "address", "country", "email", "phone"],
          properties: {
            first_name: { bsonType: "string" },
            last_name: { bsonType: "string" },
            company_name: { bsonType: "string" },
            address: { bsonType: "string" },
            country: { bsonType: "string" },
            region: { bsonType: "string" },
            city: { bsonType: "string" },
            zip_code: { bsonType: "string" },
            email: { bsonType: "string" },
            phone: { bsonType: "string" }
          }
        },
        shipping_info: {
          bsonType: "object",
          required: ["shipping_method", "shipping_cost"],
          properties: {
            ship_to_different_address: { bsonType: "bool" },
            first_name: { bsonType: "string" },
            last_name: { bsonType: "string" },
            company_name: { bsonType: "string" },
            address: { bsonType: "string" },
            country: { bsonType: "string" },
            region: { bsonType: "string" },
            city: { bsonType: "string" },
            zip_code: { bsonType: "string" },
            shipping_method: { bsonType: "string" },
            shipping_cost: { bsonType: "double" }
          }
        },
        payment_info: {
          bsonType: "object",
          properties: {
            payment_method: { bsonType: "string" },
            transaction_id: { bsonType: "string" },
            status: { bsonType: "string" },
            amount: { bsonType: "double" },
            currency: { bsonType: "string" },
            payment_date: { bsonType: "string" }
          }
        },
        totals: {
          bsonType: "object",
          required: ["subtotal", "shipping", "tax", "total"],
          properties: {
            subtotal: { bsonType: "double" },
            shipping: { bsonType: "double" },
            discount: { bsonType: "double" },
            tax: { bsonType: "double" },
            total: { bsonType: "double" }
          }
        },
        notes: { bsonType: "string" },
        created_at: { bsonType: "date" },
        updated_at: { bsonType: "date" }
      }
    }
  }
});

// Create indexes
db.orders.createIndex({ "user_id": 1 });
db.orders.createIndex({ "order_number": 1 }, { unique: true });
db.orders.createIndex({ "created_at": -1 });
db.orders.createIndex({ "status": 1 });

// Create payments collection
db.createCollection("payments", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["id", "order_id", "status", "amount", "currency", "created_at"],
      properties: {
        id: { bsonType: "string" },
        order_id: { bsonType: "string" },
        transaction_id: { bsonType: "string" },
        payment_method: { bsonType: "string" },
        status: { 
          bsonType: "string",
          enum: ["pending", "processing", "completed", "failed", "refunded"]
        },
        amount: { bsonType: "double" },
        currency: { bsonType: "string" },
        created_at: { bsonType: "date" },
        updated_at: { bsonType: "date" }
      }
    }
  }
});

// Create indexes for payments
db.payments.createIndex({ "order_id": 1 });
db.payments.createIndex({ "transaction_id": 1 }, { unique: true, sparse: true });

print("MongoDB schema initialization completed for checkout service"); 