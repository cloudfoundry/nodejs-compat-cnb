package fakes

import "sync"

type ScriptRewriter struct {
	RewriteInstallScriptsCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
}

func (f *ScriptRewriter) RewriteInstallScripts(param1 string) error {
	f.RewriteInstallScriptsCall.Lock()
	defer f.RewriteInstallScriptsCall.Unlock()
	f.RewriteInstallScriptsCall.CallCount++
	f.RewriteInstallScriptsCall.Receives.Path = param1
	if f.RewriteInstallScriptsCall.Stub != nil {
		return f.RewriteInstallScriptsCall.Stub(param1)
	}
	return f.RewriteInstallScriptsCall.Returns.Error
}
