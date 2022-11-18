db.createUser(
    {
        user: "root",
        pwd: "rootPassXXX",
        roles: [
            {
                role: "readWrite",
                db: "golang_users"
            }
        ]
    }
);
db.createCollection("test"); //MongoDB creates the database when you first store data in that database