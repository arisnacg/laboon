package bubletea

import (
	"fmt"
	"os"

	"github.com/arisnacg/laboon/internal/color"
	"github.com/arisnacg/laboon/internal/docker"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dc         = docker.DockerWrapper{}
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().Padding(0, 1).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#2F99EE"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2F99EE"))
)

type Item struct {
	ID       string
	Name     string
	Image    string
	State    string
	Selected bool
}

// FilterValue implements list.Item.
func (i Item) FilterValue() string {
	return i.Name
}

func (i Item) Title() string {
	selectedIcon := ""
	if i.Selected {
		selectedIcon = selectedItemStyle.Render("")
	}

	stateIcon := "⏸︎"
	if i.State == "running" {
		stateIcon = "▶︎"
	} else if i.State == "exited" {
		stateIcon = "⏹"
	}

	// coloring
	sTitle := fmt.Sprintf("   %s (%s) %s", i.Name, i.ID[:12], stateIcon)
	if i.State == "running" {
		sTitle = color.FgGreen(sTitle)
	} else if i.State == "exited" {
		sTitle = color.FgGray(sTitle)
	} else {
		sTitle = color.FgYellow(sTitle)
	}
	return fmt.Sprintf("%s %s", selectedIcon, sTitle)
}

func (i Item) Description() string {
	return fmt.Sprintf("   %s - %s", i.Image, i.State)
}

// key binding
type listKeyMap struct {
	pauseContainer   key.Binding
	unpauseContainer key.Binding
	startContainer   key.Binding
	stopContainer    key.Binding
	toggleSelect     key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		pauseContainer: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause container"),
		),
		unpauseContainer: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "unpause container"),
		),
		startContainer: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start container"),
		),
		stopContainer: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "stop container"),
		),
		toggleSelect: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle select"),
		),
	}
}

type Model struct {
	selectedContainers map[string]int
	list               list.Model
	keys               *listKeyMap
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) getSelectedContainers() ([]string, []int) {
	var containerIDs []string
	var indexes []int
	for id, index := range m.selectedContainers {
		containerIDs = append(containerIDs, id)
		indexes = append(indexes, index)
	}
	return containerIDs, indexes
}

func (m Model) getSelectedContainerSize() int {
	return len(m.selectedContainers)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSelect):
			// toggle selection
			selectedItem, ok := m.list.SelectedItem().(Item)
			if ok {
				selectedItem.Selected = !selectedItem.Selected
				_, keyExist := m.selectedContainers[selectedItem.ID]
				if keyExist {
					delete(m.selectedContainers, selectedItem.ID)
				} else {
					m.selectedContainers[selectedItem.ID] = m.list.Index()
				}
				m.list.SetItem(m.list.Index(), selectedItem)
			}

		case key.Matches(msg, m.keys.pauseContainer):
			selectedCount := m.getSelectedContainerSize()

			if selectedCount > 0 {
				// pause seleted container
				selectedIDs, selectedIndexes := m.getSelectedContainers()
				dc.PauseContainers(selectedIDs)
				items := m.list.Items()
				for _, index := range selectedIndexes {
					curItem := items[index].(Item)
					curItem.State = "paused"
					m.list.SetItem(index, curItem)
				}
			} else {
				// pause current container
				selectedItem, ok := m.list.SelectedItem().(Item)
				if ok {
					dc.PauseContainer(selectedItem.ID)
					selectedItem.State = "paused"
					m.list.SetItem(m.list.Index(), selectedItem)
				}
			}

		case key.Matches(msg, m.keys.unpauseContainer):

			selectedCount := m.getSelectedContainerSize()

			if selectedCount > 0 {
				// unpause seleted container
				selectedIDs, selectedIndexes := m.getSelectedContainers()
				dc.UnpauseContainers(selectedIDs)
				items := m.list.Items()
				for _, index := range selectedIndexes {
					curItem := items[index].(Item)
					curItem.State = "running"
					m.list.SetItem(index, curItem)
				}
			} else {
				// unpause current container
				selectedItem, ok := m.list.SelectedItem().(Item)
				if ok {
					dc.UnpauseContainer(selectedItem.ID)
					selectedItem.State = "running"
					m.list.SetItem(m.list.Index(), selectedItem)
				}

			}

		case key.Matches(msg, m.keys.startContainer):
			selectedCount := m.getSelectedContainerSize()

			if selectedCount > 0 {
				// start seleted container
				selectedIDs, selectedIndexes := m.getSelectedContainers()
				dc.StartContainers(selectedIDs)
				items := m.list.Items()
				for _, index := range selectedIndexes {
					curItem := items[index].(Item)
					curItem.State = "running"
					m.list.SetItem(index, curItem)
				}
			} else {
				// start current container
				selectedItem, ok := m.list.SelectedItem().(Item)
				if ok {
					dc.StartContainer(selectedItem.ID)
					selectedItem.State = "running"
					m.list.SetItem(m.list.Index(), selectedItem)
				}
			}

		case key.Matches(msg, m.keys.stopContainer):
			selectedCount := m.getSelectedContainerSize()

			if selectedCount > 0 {
				// stop seleted container
				selectedIDs, selectedIndexes := m.getSelectedContainers()
				dc.StopContainers(selectedIDs)
				items := m.list.Items()
				for _, index := range selectedIndexes {
					curItem := items[index].(Item)
					curItem.State = "exited"
					m.list.SetItem(index, curItem)
				}
			} else {
				// stop current container
				selectedItem, ok := m.list.SelectedItem().(Item)
				if ok {
					dc.StopContainer(selectedItem.ID)
					selectedItem.State = "exited"
					m.list.SetItem(m.list.Index(), selectedItem)
				}
			}

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func getContainerItems() []list.Item {
	containers := dc.GetContainers()
	var items []list.Item
	for _, container := range containers {
		items = append(
			items,
			Item{
				ID:       container.ID,
				Name:     container.Name,
				Image:    container.Image,
				State:    container.State,
				Selected: false,
			},
		)
	}
	return items
}

func Start() {
	dc.NewClient()
	defer dc.CloseClient()
	items := getContainerItems()
	listKeys := newListKeyMap()
	d := list.NewDefaultDelegate()

	dockerColor := lipgloss.Color("#2F99EE")
	white := lipgloss.Color("#FFFDF5")
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(white).
		BorderLeftForeground(dockerColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy() // reuse the title style here

	list := list.New(items, d, 0, 0)
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.pauseContainer,
			listKeys.unpauseContainer,
			listKeys.startContainer,
			listKeys.stopContainer,
			listKeys.toggleSelect,
		}
	}
	m := Model{
		selectedContainers: make(map[string]int),
		list:               list,
		keys:               listKeys,
	}
	m.list.Title = "Docker Containers"
	m.list.Styles.Title = titleStyle
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
