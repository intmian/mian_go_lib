package multi

func SafeGo(f func(), errHandler func(err error)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errHandler(err.(error))
			}
		}()
		f()
	}()
}
