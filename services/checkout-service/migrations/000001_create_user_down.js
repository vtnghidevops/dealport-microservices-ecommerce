// Rollback script for checkout service user
print("Starting rollback for checkout service user...");

// Switch to the checkout database
db = db.getSiblingDB('checkout');

// Check if the user exists before attempting to remove
let userExists = false;
try {
    userExists = db.getUser("checkout_user") !== null;
} catch (e) {
    print("Error checking user existence: " + e);
}

// Remove the user if it exists
if (userExists) {
    try {
        db.dropUser("checkout_user");
        print("User 'checkout_user' dropped successfully");
    } catch (e) {
        print("Error dropping user: " + e);
    }
} else {
    print("User 'checkout_user' does not exist, skipping user removal");
}

print("Rollback of user creation completed for checkout service"); 