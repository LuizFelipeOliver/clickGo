package cmd

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"clickGo/internal/clickup"
)

var (
	tasksIncludeClosed bool
	tasksStatusFilter  string
	taskViewOpenBrowser bool
)

func renderStatus(s clickup.TaskStatus) string {
	label := strings.ToUpper(s.Status)
	if label == "" {
		label = "SEM STATUS"
	}
	style := lipgloss.NewStyle().Bold(true)
	if s.Color != "" {
		style = style.Foreground(lipgloss.Color(s.Color))
	} else {
		style = style.Foreground(primaryColor)
	}
	return style.Render("[" + label + "]")
}

func formatAssignees(users []clickup.User, currentUserID int) string {
	if len(users) == 0 {
		return mutedStyle.Render("sem responsável")
	}
	var names []string
	for _, u := range users {
		if u.ID == currentUserID {
			names = append(names, successStyle.Render(u.Username+" (você)"))
		} else {
			names = append(names, u.Username)
		}
	}
	return strings.Join(names, ", ")
}

func runListTasks(cmd *cobra.Command, args []string) error {
	client, cfg, err := loadSession()
	if err != nil {
		return err
	}

	tasks, err := client.GetAssignedTasks(cfg.TeamID, cfg.UserID, tasksIncludeClosed)
	if err != nil {
		return err
	}

	if tasksStatusFilter != "" {
		filter := strings.ToLower(strings.TrimSpace(tasksStatusFilter))
		var filtered []clickup.Task
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Status.Status), filter) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	if len(tasks) == 0 {
		lipgloss.Println(mutedStyle.Render("Nenhuma task encontrada para você."))
		return nil
	}

	// Agrupa subtasks pelo ID da task pai
	taskByID := make(map[string]clickup.Task, len(tasks))
	subtasksByParent := make(map[string][]clickup.Task)
	var rootTasks []clickup.Task

	for _, t := range tasks {
		taskByID[t.ID] = t
	}

	for _, t := range tasks {
		if t.Parent != "" {
			subtasksByParent[t.Parent] = append(subtasksByParent[t.Parent], t)
			// Se a task pai não estiver atribuída a você, buscamos ou exibimos a subtask no topo indicando o pai
			if _, parentInList := taskByID[t.Parent]; !parentInList {
				rootTasks = append(rootTasks, t)
			}
		} else {
			rootTasks = append(rootTasks, t)
		}
	}

	header := fmt.Sprintf("Suas Tasks em %s (%d encontradas)", cfg.TeamName, len(tasks))
	lipgloss.Println(titleStyle.Render(header))
	lipgloss.Println("")

	idStyle := lipgloss.NewStyle().Bold(true).Foreground(primaryColor)

	for _, t := range rootTasks {
		isOrphanSubtask := t.Parent != ""
		prefix := "•"
		if isOrphanSubtask {
			prefix = "↳ [Subtask]"
		}

		relatedCount := len(t.RelatedTaskIDs())
		var badges []string
		if t.List.Name != "" {
			badges = append(badges, mutedStyle.Render("📋 "+t.List.Name))
		}
		if isOrphanSubtask {
			badges = append(badges, mutedStyle.Render("pai: #"+t.Parent))
		}
		if relatedCount > 0 {
			badges = append(badges, warnBadge(fmt.Sprintf("🔗 %d relacionada(s)", relatedCount)))
		}

		meta := ""
		if len(badges) > 0 {
			meta = "  " + strings.Join(badges, "  ")
		}

		line := fmt.Sprintf("%s %s %s %s%s",
			prefix,
			idStyle.Render("#"+t.DisplayID()),
			renderStatus(t.Status),
			valStyle.Render(t.Name),
			meta,
		)
		lipgloss.Println(line)

		// Exibe as subtasks cujo pai é esta task
		subs := subtasksByParent[t.ID]
		for i, sub := range subs {
			branch := "├──"
			if i == len(subs)-1 {
				branch = "└──"
			}
			subRelated := len(sub.RelatedTaskIDs())
			subExtra := ""
			if subRelated > 0 {
				subExtra = "  " + warnBadge(fmt.Sprintf("🔗 %d relacionada(s)", subRelated))
			}
			subLine := fmt.Sprintf("    %s %s %s %s%s",
				mutedStyle.Render(branch),
				idStyle.Render("#"+sub.DisplayID()),
				renderStatus(sub.Status),
				valStyle.Render(sub.Name),
				subExtra,
			)
			lipgloss.Println(subLine)
		}
	}

	lipgloss.Println("")
	lipgloss.Println(mutedStyle.Render("Dica: use 'clickGo tasks view <id>' para ver detalhes, todas as subtasks e tasks relacionadas."))
	return nil
}

func warnBadge(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFC15E")).Render(text)
}

var tasksCmd = &cobra.Command{
	Use:     "tasks",
	Aliases: []string{"task", "t"},
	Short:   "Lista as tasks e subtasks onde você é responsável",
	RunE:    runListTasks,
}

var tasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista as tasks e subtasks onde você é responsável",
	RunE:  runListTasks,
}

var tasksViewCmd = &cobra.Command{
	Use:     "view <task_id>",
	Aliases: []string{"show", "get"},
	Short:   "Mostra detalhes de uma task, suas subtasks e tasks relacionadas",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, cfg, err := loadSession()
		if err != nil {
			return err
		}

		taskID := strings.TrimPrefix(strings.TrimSpace(args[0]), "#")
		task, err := client.GetTask(taskID)
		if err != nil {
			return err
		}

		if taskViewOpenBrowser && task.URL != "" {
			_ = clickup.OpenBrowser(task.URL)
		}

		idStyle := lipgloss.NewStyle().Bold(true).Foreground(primaryColor)

		var rows []string
		rows = append(rows,
			titleStyle.Render(fmt.Sprintf("#%s — %s", task.DisplayID(), task.Name)),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Status:"), renderStatus(task.Status)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Lista:"), valStyle.Render(task.List.Name)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Responsável:"), formatAssignees(task.Assignees, cfg.UserID)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Link:"), mutedStyle.Render(task.URL)),
		)

		if task.Priority != nil && task.Priority.Priority != "" {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top,
				keyStyle.Render("Prioridade:"),
				valStyle.Render(strings.ToUpper(task.Priority.Priority)),
			))
		}

		// Se for uma subtask, busca a task pai para mostrar contexto completo
		if task.Parent != "" {
			parentTask, err := client.GetTask(task.Parent)
			if err == nil {
				parentInfo := fmt.Sprintf("#%s %s %s",
					parentTask.DisplayID(),
					renderStatus(parentTask.Status),
					parentTask.Name,
				)
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top,
					keyStyle.Render("Task Pai:"),
					valStyle.Render(parentInfo),
				))
			} else {
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top,
					keyStyle.Render("Task Pai:"),
					valStyle.Render("#"+task.Parent),
				))
			}
		}

		if strings.TrimSpace(task.TextContent) != "" {
			rows = append(rows,
				"",
				titleStyle.Render("Descrição"),
				valStyle.Render(strings.TrimSpace(task.TextContent)),
			)
		}

		// Subtasks
		rows = append(rows, "", titleStyle.Render(fmt.Sprintf("Subtasks (%d)", len(task.Subtasks))))
		if len(task.Subtasks) == 0 {
			rows = append(rows, mutedStyle.Render("  Nenhuma subtask."))
		} else {
			for _, sub := range task.Subtasks {
				rows = append(rows, fmt.Sprintf("  • %s %s %s  (%s)",
					idStyle.Render("#"+sub.DisplayID()),
					renderStatus(sub.Status),
					valStyle.Render(sub.Name),
					formatAssignees(sub.Assignees, cfg.UserID),
				))
			}
		}

		// Tasks Relacionadas (linked_tasks e dependencies)
		relatedIDs := task.RelatedTaskIDs()
		rows = append(rows, "", titleStyle.Render(fmt.Sprintf("Tasks Relacionadas (%d)", len(relatedIDs))))
		if len(relatedIDs) == 0 {
			rows = append(rows, mutedStyle.Render("  Nenhuma task relacionada."))
		} else {
			for _, relID := range relatedIDs {
				relTask, err := client.GetTask(relID)
				if err != nil {
					rows = append(rows, fmt.Sprintf("  • %s %s",
						idStyle.Render("#"+relID),
						mutedStyle.Render("(não foi possível carregar detalhes)"),
					))
					continue
				}
				rows = append(rows, fmt.Sprintf("  • %s %s %s  (%s)",
					idStyle.Render("#"+relTask.DisplayID()),
					renderStatus(relTask.Status),
					valStyle.Render(relTask.Name),
					formatAssignees(relTask.Assignees, cfg.UserID),
				))
			}
		}

		lipgloss.Println(boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...)))
		return nil
	},
}

func init() {
	tasksCmd.PersistentFlags().BoolVarP(&tasksIncludeClosed, "all", "a", false, "Incluir tasks concluídas/fechadas")
	tasksCmd.PersistentFlags().StringVarP(&tasksStatusFilter, "status", "s", "", "Filtrar por nome do status (ex: 'in progress')")

	tasksViewCmd.Flags().BoolVarP(&taskViewOpenBrowser, "open", "o", false, "Abrir a task no navegador")

	tasksCmd.AddCommand(tasksListCmd, tasksViewCmd)
	rootCmd.AddCommand(tasksCmd)
}
