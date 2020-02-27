package fakes

import "sync"

type PrePostParser struct {
	ContainsScriptsCall struct {
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

func (f *PrePostParser) ContainsScripts(param1 string) (bool, error) {
	f.ContainsScriptsCall.Lock()
	defer f.ContainsScriptsCall.Unlock()
	f.ContainsScriptsCall.CallCount++
	f.ContainsScriptsCall.Receives.Path = param1
	if f.ContainsScriptsCall.Stub != nil {
		return f.ContainsScriptsCall.Stub(param1)
	}
	return f.ContainsScriptsCall.Returns.ScriptsExist, f.ContainsScriptsCall.Returns.Err
}
