package main

func main() {
	a := App{}
	a.Initialize()
	a.Start()

	// clientOptions := options.Client().ApplyURI("mongodb+srv://rovilay:qwertyUp1.@ireporter-cluster-y4nzl.mongodb.net/course_syndicate?authSource=admin&connect=direct")

	// // Connect to MongoDB
	// client, err := mongo.Connect(context.Background(), clientOptions)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Check the connection
	// err = client.Ping(context.Background(), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Connected to MongoDB!")
}
