package catomic

// nocmp is an uncomparable struct. Embed this inside another struct to make
// it uncomparable. see:https://github.com/uber-go/atomic/blob/master/nocmp.go
//
//	type Foo struct {
//	  nocmp
//	  // ...
//	}
//
// This DOES NOT:
//
//   - Disallow shallow copies of structs
//   - Disallow comparison of pointers to uncomparable structs
type nocmp [0]func()
