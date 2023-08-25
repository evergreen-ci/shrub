// Package shrub provides a simple, low-overhead interface for
// generating Evergreen project configurations.
//
// For most use cases, you can start with a Configuration created via
// BuildConfiguration and add data such as new tasks, task groups, build
// variants, and functions to that Configuration using the provided setters. For
// example:
//
//	conf := &Configuration{}
//	// The following function definitions is equivalent to the YAML
//	// configuration:
//	// functions:
//	//   - name: my-new-func
//	//     commands:
//	//       - name: git.get_project
//	//         params:
//	//           directory: my-working-directory
//	//       - name: shell.exec
//	//         params:
//	//           script: echo hello world!
//	newFunc := conf.Function("my-new-func")
//	newFunc.Command().Command("git.get_project").Param("directory", "my-working-directory")
//	newFunc.Command().Command("shell.exec").Param("script", "echo hello world!")
//
//	// The following task definition is equivalent to the YAML configuration:
//	// tasks:
//	//   - name: my-new-task
//	//     commands:
//	//        - func: my-new-func
//	newTask := conf.Task("my-new-task")
//	newTask.Function("my-new-func")
//
//	// The following build variant definitions is equivalent to the YAML
//	// configuration:
//	// buildvariants:
//	//   - name: my-new-build-variant
//	//     run_on:
//	//       - some-distro
//	//     tasks:
//	//       - name: my-new-task
//	newBV := conf.Variant("my-new-build-variant")
//	newBV.RunOn("some-distro")
//	newBV.AddTasks("my-new-task")
//
// Be aware that some command methods will panic if you attempt to
// construct an invalid command. You can wrap your configuration logic with
// BuildConfiguration to convert any panic into an error.
package shrub

// Configuration is the top-level representation of the components of
// an evergreen project configuration.
type Configuration struct {
	Functions map[string]*CommandSequence `json:"functions,omitempty" yaml:"functions,omitempty"`
	Tasks     []*Task                     `json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Groups    []*TaskGroup                `json:"task_groups,omitempty" yaml:"task_groups,omitempty"`
	Variants  []*Variant                  `json:"buildvariants,omitempty" yaml:"buildvariants,omitempty"`
}

// Task returns a task of the specified name. If the task already
// exists, then it returns the existing task of that name, and
// otherwise returns a new task of the specified name.
func (c *Configuration) Task(name string) *Task {
	for _, t := range c.Tasks {
		if t.Name == name {
			return t
		}
	}

	t := new(Task)
	t.Name = name
	c.Tasks = append(c.Tasks, t)
	return t
}

// TaskGroup returns a task group configuration of the specified
// name. If the taskgroup already exists, then it returns the existing
// task group of that name, and otherwise returns a new task group of
// the specified name.
func (c *Configuration) TaskGroup(name string) *TaskGroup {
	for _, g := range c.Groups {
		if g.GroupName == name {
			return g
		}
	}

	g := new(TaskGroup)
	c.Groups = append(c.Groups, g)
	return g.Name(name)
}

// Function creates a new function of the specific name and returns a
// CommandSequence builder for use in adding commands to the function.
func (c *Configuration) Function(name string) *CommandSequence {
	if c.Functions == nil {
		c.Functions = make(map[string]*CommandSequence)
	}

	seq, ok := c.Functions[name]
	if ok {
		return seq
	}

	seq = new(CommandSequence)
	c.Functions[name] = seq
	return seq
}

// Variant returns a build variant of the specified name. If the
// variant already exists, then it returns the existing variant of
// that name, and otherwise returns a new variant of the specified
// name.
func (c *Configuration) Variant(id string) *Variant {
	for _, v := range c.Variants {
		if v.BuildName == id {
			return v
		}
	}

	v := new(Variant)
	c.Variants = append(c.Variants, v)
	return v.Name(id)
}
