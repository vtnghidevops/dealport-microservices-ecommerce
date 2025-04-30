// MongoDB Init Data for Checkout Service

print("Starting data initialization for checkout service...");

// Create shipping methods
const shippingMethods = [
  {
    code: "standard",
    name: "Standard Shipping",
    description: "Delivery in 5-7 business days",
    base_price: 5.99,
    free_threshold: 50.00
  },
  {
    code: "express",
    name: "Express Shipping",
    description: "Delivery in 2-3 business days",
    base_price: 12.99,
    free_threshold: 100.00
  },
  {
    code: "overnight",
    name: "Overnight Shipping",
    description: "Next day delivery",
    base_price: 24.99,
    free_threshold: null
  }
];

// Create shipping_methods collection if it doesn't exist
if (db.getCollectionNames().indexOf("shipping_methods") == -1) {
  db.createCollection("shipping_methods");
}

// Insert shipping methods
db.shipping_methods.insertMany(shippingMethods);

// Create tax rates
const taxRates = [
  {
    country: "United States",
    state: "CA",
    rate: 0.0725
  },
  {
    country: "United States",
    state: "NY",
    rate: 0.045
  },
  {
    country: "United States",
    state: "TX",
    rate: 0.0625
  },
  {
    country: "United States",
    state: "Default",
    rate: 0.05
  },
  {
    country: "Canada",
    state: "Default",
    rate: 0.05
  },
  {
    country: "Default",
    state: "Default",
    rate: 0.00
  }
];

// Create tax_rates collection if it doesn't exist
if (db.getCollectionNames().indexOf("tax_rates") == -1) {
  db.createCollection("tax_rates");
}

// Insert tax rates
db.tax_rates.insertMany(taxRates);

// Create payment methods
const paymentMethods = [
  {
    code: "credit_card",
    name: "Credit Card",
    description: "Pay with Visa, Mastercard, Amex, or Discover",
    enabled: true,
    sort_order: 1
  },
  {
    code: "paypal",
    name: "PayPal",
    description: "Pay with your PayPal account",
    enabled: true,
    sort_order: 2
  },
  {
    code: "bank_transfer",
    name: "Bank Transfer",
    description: "Pay via bank transfer",
    enabled: true,
    sort_order: 3
  }
];

// Create payment_methods collection if it doesn't exist
if (db.getCollectionNames().indexOf("payment_methods") == -1) {
  db.createCollection("payment_methods");
}

// Insert payment methods
db.payment_methods.insertMany(paymentMethods);

// Print completion message
print("Sample data initialization completed for checkout service"); 