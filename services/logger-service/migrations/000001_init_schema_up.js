// Create user for logs database
print("Starting user creation for logger service...");

// Switch to the logs database
db = db.getSiblingDB('logs');

// Check if the user already exists
let userExists = false;
try {
    // Try to find the user in the system.users collection
    userExists = db.getUser("logger_user") !== null;
} catch (e) {
    print("Error checking user existence: " + e);
}

// If user doesn't exist, create it
if (!userExists) {
    try {
        db.createUser({
            user: "logger_user",
            pwd: "password",
            roles: [
                { role: "readWrite", db: "logs" }
            ]
        });
        print("User 'logger_user' created successfully for logs database");
    } catch (e) {
        print("Error creating user: " + e);
    }
} else {
    print("User 'logger_user' already exists for logs database");
}

// Create logs collection if it doesn't exist
if (db.getCollectionNames().indexOf("logs") == -1) {
    db.createCollection("logs");
    print("Collection 'logs' created successfully");
}

print("User and collection setup completed for logger service"); 