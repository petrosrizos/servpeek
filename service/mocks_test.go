package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/milosgajdos83/servpeek/utils/command"
)

type mockCmd struct {
	out string
	err error
}

func (m *mockCmd) Run() command.Output          { return nil }
func (m *mockCmd) RunCombined() (string, error) { return m.out, m.err }
func (m *mockCmd) AppendArgs(args ...string)    {}

type mockCommander struct {
	StartCmd  *mockCmd
	StopCmd   *mockCmd
	StatusCmd *mockCmd
}

type mockSysInit struct {
	*mockCommander
	sysInitType string
}

func (m *mockSysInit) Type() string { return m.sysInitType }
func (m *mockSysInit) Start(n string) error {
	_, err := m.StartCmd.RunCombined()
	return err
}
func (m *mockSysInit) Stop(n string) error {
	_, err := m.StopCmd.RunCombined()
	return err
}

// this is basically copy-paste from sysinit.go
func (m *mockSysInit) Status(n string) (Status, error) {
	status, err := m.StatusCmd.RunCombined()
	if err != nil {
		return Unknown, err
	}
	switch {
	case strings.Contains(status, svcStatusOut["running"][m.sysInitType]):
		return Running, nil
	case strings.Contains(status, svcStatusOut["stopped"][m.sysInitType]):
		return Stopped, nil
	}
	return Unknown, fmt.Errorf("Unable to determine %s status", n)
}

type mockSvc struct {
	name    string
	sysInit *mockSysInit
}

func newMockSvc(svcStatus, sysInitType, svcName string) (*mockSvc, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fileName := fmt.Sprintf("%s-%s.out", sysInitType, svcStatus)
	fixturesPath := path.Join(currentDir, "test-fixtures", fileName)
	cmdOut, err := ioutil.ReadFile(fixturesPath)
	if err != nil {
		return nil, err
	}
	return &mockSvc{
		name: svcName,
		sysInit: &mockSysInit{
			sysInitType: sysInitType,
			mockCommander: &mockCommander{
				StatusCmd: &mockCmd{
					out: string(cmdOut),
				},
			},
		},
	}, nil
}

func (m *mockSvc) Name() string     { return m.name }
func (m *mockSvc) SysInit() SysInit { return m.sysInit }
