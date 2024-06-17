package dll

import "github.com/CarsonSlovoka/go-pkg/v2/w32"

var User = w32.NewUser32DLL()
var Gdi = w32.NewGdi32DLL()
var Kernel = w32.NewKernel32DLL()
