package shrub

type Task struct {
	Name             string           `json:"name" yaml:"name"`
	PriorityOverride int              `json:"priority,omitempty" yaml:"priority_override,omitempty"`
	Dependencies     []TaskDependency `json:"depends_on,omitempty" yaml:"dependencies,omitempty"`
	Commands         CommandSequence  `json:"commands" yaml:"commands"`
}

type TaskDependency struct {
	Name    string `json:"name" yaml:"name"`
	Variant string `json:"variant" yaml:"variant"`
}

func (t *Task) Command(cmds ...Command) *Task {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}

		t.Commands = append(t.Commands, c.Resolve())
	}

	return t
}

func (t *Task) AddCommand() *CommandDefinition {
	c := &CommandDefinition{}
	t.Commands = append(t.Commands, c)
	return c
}

func (t *Task) Dependency(dep ...TaskDependency) *Task {
	t.Dependencies = append(t.Dependencies, dep...)
	return t
}

func (t *Task) Function(fns ...string) *Task {
	for _, fn := range fns {
		t.Commands = append(t.Commands, &CommandDefinition{
			FunctionName: fn,
		})
	}

	return t
}

func (t *Task) FunctionWithVars(id string, vars map[string]string) *Task {
	t.Commands = append(t.Commands, &CommandDefinition{
		FunctionName: id,
		Vars:         vars,
	})

	return t
}

func (t *Task) Priority(pri int) *Task { t.PriorityOverride = pri; return t }

type TaskGroup struct {
	GroupName             string          `json:"name" yaml:"name"`
	MaxHosts              int             `json:"max_hosts,omitempty" yaml:"max_hosts,omitempty"`
	ShareProcesses        bool            `json:"share_processes,omitempty" yaml:"share_processes,omitemtpy"`
	SetupGroupCmds        CommandSequence `json:"setup_group,omitempty" yaml:"setup_group,omitempty"`
	SetupGroupCanFailTask bool            `json:"setup_group_can_fail_task,omitempty" yaml:"setup_group_can_fail_task,omitempty"`
	SetupGroupTimeoutSecs int             `json:"setup_group_timeout_secs,omitempty" yaml:"setup_group_timeout_secs,omitempty"`
	SetupTaskCmds         CommandSequence `json:"setup_task,omitempty" yaml:"setup_task,omitempty"`
	Tasks                 []string        `json:"tasks" yaml:"tasks"`
	TeardownTaskCmds      CommandSequence `json:"teardown_task,omitempty" yaml:"teardown_task,omitempty"`
	TeardownGroupCmds     CommandSequence `json:"teardown_group,omitempty" yaml:"teardown_group,omitempty"`
	TimeoutCmds           CommandSequence `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Tags                  []string        `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func (g *TaskGroup) Name(id string) *TaskGroup {
	g.GroupName = id
	return g
}

func (g *TaskGroup) SetMaxHosts(num int) *TaskGroup {
	g.MaxHosts = num
	return g
}

func (g *TaskGroup) SetShareProcesses(val bool) *TaskGroup {
	g.ShareProcesses = val
	return g
}

func (g *TaskGroup) SetupGroup(cmds ...Command) *TaskGroup {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		g.SetupGroupCmds = append(g.SetupGroupCmds, c.Resolve())
	}
	return g
}

func (g *TaskGroup) SetSetupGroupCanFailTask(val bool) *TaskGroup {
	g.SetupGroupCanFailTask = val
	return g
}

func (g *TaskGroup) SetSetupGroupTimeoutSecs(timeoutSecs int) *TaskGroup {
	g.SetupGroupTimeoutSecs = timeoutSecs
	return g
}

func (g *TaskGroup) SetupTask(cmds ...Command) *TaskGroup {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		g.SetupTaskCmds = append(g.SetupTaskCmds, c.Resolve())
	}
	return g
}

func (g *TaskGroup) Task(id ...string) *TaskGroup {
	g.Tasks = append(g.Tasks, id...)
	return g
}

func (g *TaskGroup) TeardownTask(cmds ...Command) *TaskGroup {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		g.TeardownTaskCmds = append(g.TeardownTaskCmds, c.Resolve())
	}
	return g
}

func (g *TaskGroup) TeardownGroup(cmds ...Command) *TaskGroup {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		g.TeardownGroupCmds = append(g.TeardownGroupCmds, c.Resolve())
	}
	return g
}

func (g *TaskGroup) Timeout(cmds ...Command) *TaskGroup {
	for _, c := range cmds {
		if err := c.Validate(); err != nil {
			panic(err)
		}
		g.TimeoutCmds = append(g.TimeoutCmds, c.Resolve())
	}
	return g
}

func (g *TaskGroup) Tag(tags ...string) *TaskGroup {
	g.Tags = append(g.Tags, tags...)
	return g
}
