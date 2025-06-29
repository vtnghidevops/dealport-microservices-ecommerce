// MongoDB Schema Drop for Checkout Service
db = db.getSiblingDB('checkout');

print("Starting MongoDB schema cleanup for checkout service...");

// Drop collections
db.orders.drop();
db.payments.drop();

print("MongoDB schema cleanup completed for checkout service."); 