package main

import (
    "fmt"
    "net/http"
    "time"
    "log"
    "os"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection

func addEmail(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    count, err := coll.CountDocuments(context.Background(), bson.M{"email":email})
    if err != nil {
        log.Fatal(err)
    }
    if count == 0 {
        document := bson.M{
            "email": email,
        }
        _, err := coll.InsertOne(context.TODO(), document)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        fmt.Println("Email already added")
    }
    
}

func main() {
    m := http.NewServeMux()

    const addr = ":8080"

    srv := http.Server{
        Handler: m,
        Addr: addr,
        WriteTimeout: 30 * time.Second,
        ReadTimeout: 30 * time.Second,
    }

    fmt.Println("Server started on port ", addr)
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        log.Fatal("You must set your 'MONGODB_URI' env variable")
    }

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }

    defer func() {
        if err:= client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    coll = client.Database("EverestAllegiance").Collection("Emails")
    m.HandleFunc("/newEmail", addEmail)
    err = srv.ListenAndServe()
    log.Fatal(err)
}

