package clickup

import (
	"fmt"
	"net/http"
	"strconv"
)

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type teamsResponse struct {
	Teams []Team `json:"teams"`
}

// GetTeams retorna os Workspaces (Teams) aos quais o usuário tem acesso.
func (c *Client) GetTeams() ([]Team, error) {
	var res teamsResponse
	if err := c.do(http.MethodGet, "/team", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Teams, nil
}

type TaskStatus struct {
	Status string `json:"status"`
	Color  string `json:"color"`
	Type   string `json:"type"`
}

type TaskPriority struct {
	ID       string `json:"id"`
	Priority string `json:"priority"`
	Color    string `json:"color"`
}

type TaskRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LinkedTask struct {
	TaskID      string `json:"task_id"`
	LinkID      string `json:"link_id"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"userid"`
}

type Dependency struct {
	TaskID    string `json:"task_id"`
	DependsOn string `json:"depends_on"`
	Type      int    `json:"type"`
}

type Task struct {
	ID           string        `json:"id"`
	CustomID     string        `json:"custom_id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	TextContent  string        `json:"text_content"`
	Status       TaskStatus    `json:"status"`
	Priority     *TaskPriority `json:"priority"`
	Parent       string        `json:"parent"`
	URL          string        `json:"url"`
	Assignees    []User        `json:"assignees"`
	List         TaskRef       `json:"list"`
	Folder       TaskRef       `json:"folder"`
	Space        TaskRef       `json:"space"`
	Subtasks     []Task        `json:"subtasks"`
	LinkedTasks  []LinkedTask  `json:"linked_tasks"`
	Dependencies []Dependency  `json:"dependencies"`
}

// DisplayID retorna o CustomID (ex: DEV-123) se existir, junto ao ID original, ou apenas o ID.
func (t Task) DisplayID() string {
	if t.CustomID != "" {
		return fmt.Sprintf("%s (%s)", t.CustomID, t.ID)
	}
	return t.ID
}

// RelatedTaskIDs retorna todos os IDs únicos de tasks vinculadas (linked_tasks e dependencies)
// diferentes do ID da própria task.
func (t Task) RelatedTaskIDs() []string {
	seen := map[string]bool{t.ID: true}
	var ids []string

	add := func(id string) {
		if id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	for _, lt := range t.LinkedTasks {
		add(lt.TaskID)
		add(lt.LinkID)
	}
	for _, dep := range t.Dependencies {
		add(dep.TaskID)
		add(dep.DependsOn)
	}
	return ids
}

type tasksResponse struct {
	Tasks    []Task `json:"tasks"`
	LastPage bool   `json:"last_page"`
}

// GetAssignedTasks busca todas as tasks e subtasks atribuídas ao userID dentro do Workspace (teamID).
func (c *Client) GetAssignedTasks(teamID string, userID int, includeClosed bool) ([]Task, error) {
	var all []Task
	page := 0

	for {
		query := map[string]string{
			"assignees[]":    strconv.Itoa(userID),
			"subtasks":       "true",
			"include_closed": strconv.FormatBool(includeClosed),
			"page":           strconv.Itoa(page),
		}

		var res tasksResponse
		path := fmt.Sprintf("/team/%s/task", teamID)
		if err := c.do(http.MethodGet, path, query, nil, &res); err != nil {
			return nil, err
		}

		all = append(all, res.Tasks...)
		if res.LastPage || len(res.Tasks) == 0 {
			break
		}
		page++
	}

	return all, nil
}

// GetTask busca os detalhes completos de uma task (incluindo suas subtasks diretas).
func (c *Client) GetTask(taskID string) (*Task, error) {
	query := map[string]string{
		"include_subtasks": "true",
	}
	var task Task
	path := fmt.Sprintf("/task/%s", taskID)
	if err := c.do(http.MethodGet, path, query, nil, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
