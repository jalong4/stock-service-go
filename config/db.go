package config
// Path: config/db.go

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

	"github.com/jalong4/stock-service-go/models"
    "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func LoadEnv() {
    err := godotenv.Load("./config/config.env") // Loads values from .env into the system
    if err != nil {
        log.Printf("Error loading config.env file")
    }
}

func ConnectDB() {
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        log.Fatal("MONGO_URI not set")
    }

    clientOptions := options.Client().ApplyURI(uri)
    var err error
    MongoDB, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection with a temporary context
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = MongoDB.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
}

func GetMongoDB() *mongo.Client {
    return MongoDB
}


// This function creates a new context for each call instead of using a global context
func NewMongoDBContext(duration time.Duration) context.Context {
    ctx, _ := context.WithTimeout(context.Background(), duration)
    return ctx
}

// Users

func GetUsersCollection() *mongo.Collection {
	collection := MongoDB.Database(os.Getenv("MONGO_DB")).Collection("users")
    return collection
}

func GetAllUsers() ([]models.User, error) {
    var users []models.User
    collection := GetUsersCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Failed to retrieve users: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        if err = cursor.Decode(&user); err != nil {
            log.Printf("Failed to decode user: %v", err)
            continue
        }
        users = append(users, user)
    }

    if err = cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil, err
    }

    return users, nil
}

func GetUserByEmail(email string) (*models.User, error) {
    collection := GetUsersCollection()
    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func GetUserByID(id string) (*models.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("User ID: %s not found", id)
    }

	log.Println("ID:", id)
	log.Println("UserID:", userID)

	collection := GetUsersCollection()
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

    err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
		dbName := os.Getenv("MONGO_DB")
        return nil, fmt.Errorf("User ID: %s not found in Database: %q", id, dbName)
    }
	return &user, nil
}

func InsertUser(user *models.User) (*mongo.InsertOneResult, error) {
    collection := GetUsersCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, &user)
    if err != nil {
        return nil, err
    }
    return result, nil
}


func DeleteUserByID(id string) (*models.User, error) {

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    user, err := GetUserByID(id)
    if err != nil {
        return nil, err
    }

    collection := GetUsersCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return nil, err
    }
    return user, nil
}

func UpdateUserByID(id string, updatedUser models.User) (*mongo.UpdateResult, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    collection := GetUsersCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.ReplaceOne(ctx, bson.M{"_id": oid}, updatedUser)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// Holdings

// GetHoldingsCollection returns the holdings collection from the MongoDB
func GetHoldingsCollection() *mongo.Collection {
    dbName := os.Getenv("MONGO_DB")
    collection := MongoDB.Database(dbName).Collection("holdings")
    return collection
}

// AddHolding inserts a new holding into the database
func AddHolding(holding models.Holding) (string, error) {

    // Ensure ID is not set for new holdings
    if holding.ID != primitive.NilObjectID {
        return "", fmt.Errorf("ID should not be provided for a new holding")
    }

    collection := GetHoldingsCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, holding)
    if err != nil {
        return "", err
    }

    // Retrieve the new document ID and convert it to a string
    newID := result.InsertedID.(primitive.ObjectID).Hex()
    return newID, nil
}

// GetAllHoldings retrieves all holdings from the database
func GetAllHoldings() ([]models.Holding, error) {
    var holdings []models.Holding
    collection := GetHoldingsCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Failed to retrieve holdings: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var holding models.Holding
        if err = cursor.Decode(&holding); err != nil {
            log.Printf("Failed to decode holding: %v", err)
            continue
        }
        holdings = append(holdings, holding)
    }

    if err = cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil, err
    }

    return holdings, nil
}

// GetHoldingsByTicker retrieves all holdings that match a specific ticker
func GetHoldingsByTicker(ticker string) ([]models.Holding, error) {
    var holdings []models.Holding
    collection := GetHoldingsCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    filter := bson.M{"ticker": ticker}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        log.Printf("Failed to retrieve holdings for ticker %s: %v", ticker, err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var holding models.Holding
        if err = cursor.Decode(&holding); err != nil {
            log.Printf("Failed to decode holding: %v", err)
            continue
        }
        holdings = append(holdings, holding)
    }

    if err = cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil, err
    }

    return holdings, nil
}

// GetHoldingsByAccount retrieves holdings that match a regex pattern on the account field
func GetHoldingsByAccount(accountPattern string) ([]models.Holding, error) {
    var holdings []models.Holding
    collection := GetHoldingsCollection()

    // Create a regex pattern to filter the account
    regexPattern := fmt.Sprintf(".*%s.*", accountPattern)
    filter := bson.M{"account": bson.M{"$regex": regexPattern, "$options": "i"}}  // Case insensitive matching

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        log.Printf("Failed to retrieve holdings by account pattern %s: %v", accountPattern, err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var holding models.Holding
        if err = cursor.Decode(&holding); err != nil {
            log.Printf("Failed to decode holding: %v", err)
            continue
        }
        holdings = append(holdings, holding)
    }

    if err = cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil, err
    }

    return holdings, nil
}

func GetHoldingByID(id string) (*models.Holding, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    collection := GetHoldingsCollection()
    var holding models.Holding
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&holding)
    if err != nil {
        return nil, err
    }
    return &holding, nil
}

func DeleteHoldingByID(id string) (*mongo.DeleteResult, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    collection := GetHoldingsCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return nil, err
    }
    return result, nil
}

func UpdateHoldingByID(id string, updatedHolding models.Holding) (*mongo.UpdateResult, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    collection := GetHoldingsCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.ReplaceOne(ctx, bson.M{"_id": oid}, updatedHolding)
    if err != nil {
        return nil, err
    }
    return result, nil
}
