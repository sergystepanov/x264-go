package external

/*
#cgo android LDFLAGS: -lx264 -lm
#cgo windows LDFLAGS: -lx264
#cgo linux LDFLAGS: -lx264 -lm
#cgo darwin LDFLAGS: -lx264 -lm
*/
import "C"
