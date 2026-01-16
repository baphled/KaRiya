// Package behaviors provides embeddable components for table-based UIs with CRUD operations.
//
// # Overview
//
// The behaviors package offers composable, type-safe components that can be embedded in
// intents to eliminate boilerplate code for common UI patterns. Each behavior is designed
// to work independently or in composition with others.
//
// # Architecture
//
// Behaviors follow the same embeddable pattern as BaseIntent:
//
//	type MyIntent struct {
//	    *BaseIntent
//	    *behaviors.TableBehavior[*domain.Item]
//	    *behaviors.CRUDBehavior[*domain.Item]
//	}
//
// # Available Behaviors
//
//   - TableBehavior[T]: Provides data binding, pagination, navigation, filtering, and sorting
//   - CRUDBehavior[T]: Handles create/edit/delete operations with confirmation dialogs
//   - FilterMenuBehavior[T]: Provides a sectioned filter menu with keyboard navigation
//   - SortMenuBehavior[T]: Provides a sort menu with comparator-based sorting
//
// # Design Principles
//
//   - Type Safety: All behaviors use generics for type-safe item handling
//   - Single Source of Truth: TableBehavior owns the data, other behaviors reference it
//   - One-Way References: Behaviors can reference TableBehavior, but not vice versa
//   - No State Caching: Behaviors always delegate to the table for current data
//   - Embeddable: Behaviors are designed to be embedded in intent structs
//
// # Usage Example
//
//	// Create table behavior
//	table := behaviors.NewTableBehavior(
//	    theme,
//	    []behaviors.ColumnDef{
//	        {Title: "Name", Width: 25},
//	        {Title: "Status", Width: 15},
//	    },
//	    func(item *MyItem, idx int) []string {
//	        return []string{item.Name, item.Status}
//	    },
//	).PageSize(15).EmptyMessage("No items yet")
//
//	// Create CRUD behavior (references table)
//	crud := behaviors.NewCRUDBehavior(table, theme).
//	    ItemNamer(func(item *MyItem) string { return item.Name }).
//	    OnCreate(func() tea.Cmd { ... }).
//	    OnEdit(func(item *MyItem) tea.Cmd { ... }).
//	    OnDelete(func(item *MyItem) tea.Cmd { ... })
//
//	// In Update:
//	if cmd, handled := crud.Update(msg); handled {
//	    return cmd
//	}
//	if table.HandleNavigation(keyStr) {
//	    return nil
//	}
//
//	// In View:
//	return table.Render()
//
// # Testing
//
// Each behavior is independently testable without requiring an intent:
//
//	table := behaviors.NewTableBehavior(...)
//	table.SetItems(testItems)
//	Expect(table.GetSelectedItem()).To(Equal(testItems[0]))
//
// For more details, see the individual behavior documentation.
package behaviors
