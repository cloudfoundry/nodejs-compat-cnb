package fakes

import "sync"

type PrePostParser struct {
	ParseCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			ScriptsExist bool
			Err          error
		}
		Stub func(string) (bool, error)
	}
}

func (f *PrePostParser) Parse(param1 string) (bool, error) {
	f.ParseCall.Lock()
	defer f.ParseCall.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Path = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.ScriptsExist, f.ParseCall.Returns.Err
}
