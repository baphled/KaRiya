package main

func main() {
	// Intentionally empty: test fixture for analyzer validation
}

func init() {
	// Intentionally empty: test fixture for analyzer validation
}

func ExportedInMain() {} // want `exported function ExportedInMain missing doc comment`
