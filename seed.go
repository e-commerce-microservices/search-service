package main

// func seed(es *elasticsearch.Client) {

// 	var wg sync.WaitGroup
// 	// 2. Index documents
// 	for i, productName := range listProductName {
// 		wg.Add(1)
// 		go func(i int, productName string) {
// 			defer wg.Done()
// 			// Build the request body
// 			data, err := json.Marshal(struct {
// 				ProductName string `json:"product_name"`
// 			}{ProductName: fmt.Sprintf("%s at %d", productName, time.Now().Unix())})
// 			if err != nil {
// 				log.Fatalf("Error marshalling document: %s", err)
// 			}

// 			// Set up the request object
// 			req := esapi.IndexRequest{
// 				Index:      productIndex,
// 				DocumentID: strconv.Itoa(i + 1),
// 				Body:       bytes.NewReader(data),
// 				Refresh:    "true",
// 			}

// 			// Perform the request with the client
// 			res, err := req.Do(context.Background(), es)
// 			if err != nil {
// 				log.Printf("Error getting response: %s", err)
// 			}
// 			defer res.Body.Close()

// 			if res.IsError() {
// 				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
// 			} else {
// 				// Deserialize the response into a map.
// 				var r map[string]interface{}
// 				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 					log.Printf("Error parsing the response body: %s", err)
// 				} else {
// 					// Print the response status and indexed document version.
// 					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
// 				}
// 			}
// 		}(i, productName)
// 	}
// 	wg.Wait()

// 	log.Println(strings.Repeat("-", 37))
// }
