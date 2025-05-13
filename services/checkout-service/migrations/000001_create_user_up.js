// Create user for checkout database
print("Starting user creation for checkout service...");

// Switch to the checkout database
db = db.getSiblingDB('checkout');

// Check if the user already exists
let userExists = false;
try {
    // Try to find the user in the system.users collection
    userExists = db.getUser("checkout_user") !== null;
} catch (e) {
    print("Error checking user existence: " + e);
}

// If user doesn't exist, create it
if (!userExists) {
    try {
        db.createUser({
            user: "checkout_user",
            pwd: "password",
            roles: [
                { role: "readWrite", db: "checkout" }
            ]
        });
        print("User 'checkout_user' created successfully for checkout database");
    } catch (e) {
        print("Error creating user: " + e);
    }
} else {
    print("User 'checkout_user' already exists for checkout database");
}

print("User setup completed for checkout service"); 