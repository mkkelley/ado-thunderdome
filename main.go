package main

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	err = RunHttpServer(config)
	if err != nil {
		panic(err)
	}
}
