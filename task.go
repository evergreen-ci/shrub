package shrub

type Task struct {
	Name             string           `json:"name"`
	PriorityOverride int              `json:"priority,omitempty"`
	Dependencies     []TaskDependency `json:"depends_on,omitempty"`
	Commands         CommandSequence  `json:"commands"`
}

type TaskDependency struct {
	Name    string `json:"name"`
	Variant string `json:"variant"`
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
	GroupName     string          `json:"name"`
	MaxHosts      int             `json:"max_hosts"`
	SetupGroup    CommandSequence `json:"setup_group"`
	SetupTask     CommandSequence `json:"setup_task"`
	Tasks         []string        `json:"tasks"`
	TeardownTask  CommandSequence `json:"teardown_task"`
	TeardownGroup CommandSequence `json:"teardown_group"`
	Timeout       CommandSequence `json:"timeout"`
}

func (g *TaskGroup) Name(id string) *TaskGroup      { g.GroupName = id; return g }
func (g *TaskGroup) SetMaxHosts(num int) *TaskGroup { g.MaxHosts = num; return g }
func (g *TaskGroup) Task(id ...string) *TaskGroup   { g.Tasks = append(g.Tasks, id...); return g }
