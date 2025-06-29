// Down migration to remove the user and collection
print("Starting down migration for logger service...");

// Switch to the logs database
db = db.getSiblingDB('logs');

// Drop the logs collection
try {
    db.logs.drop();
    print("Collection 'logs' dropped successfully");
} catch (e) {
    print("Error dropping logs collection: " + e);
}

// Remove the user
try {
    db.dropUser("logger_user");
    print("User 'logger_user' removed successfully");
} catch (e) {
    print("Error removing user: " + e);
}

print("Down migration completed for logger service"); 