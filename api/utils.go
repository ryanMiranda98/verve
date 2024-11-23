package api

func getUniqueRequestsCount() int {
	mu.Lock()
	defer mu.Unlock()
	return len(requestsMap)
}

func resetUniqueRequests() {
	mu.Lock()
	defer mu.Unlock()
	clear(requestsMap)
}

func updateRequestTracking(id string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := requestsMap[id]; !ok {
		requestsMap[id] = struct{}{}
	}
}
