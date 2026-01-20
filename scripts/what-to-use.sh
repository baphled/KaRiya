#!/bin/bash

# Quick decision helper for AI agents
# Usage: make what-to-use NEED="table"

NEED="${1:-}"

if [ -z "$NEED" ]; then
    echo "Usage: make what-to-use NEED=\"keyword\""
    echo ""
    echo "Keywords: table, form, color, view, footer, modal, list, crud, filter, sort, text, button, input, badge, box, spacing"
    exit 0
fi

NEED_LOWER=$(echo "$NEED" | tr '[:upper:]' '[:lower:]')

echo "================================================"
echo "COMPONENT LOOKUP: $NEED"
echo "================================================"
echo ""

case "$NEED_LOWER" in
    table|list)
        echo "USE: behaviors.TableBehavior[T]"
        echo "NOT: table.New()"
        echo ""
        echo "Location: internal/cli/behaviors/table.go"
        echo ""
        echo "Example:"
        echo "  columns := []behaviors.ColumnDef{{Title: \"Name\", Width: 20}}"
        echo "  formatter := func(item T) []string { return []string{item.Name} }"
        echo "  table := behaviors.NewTableBehavior(theme, columns, formatter)"
        echo "  table.SetItems(items)"
        echo "  // In Update: table.HandleNavigation(key)"
        echo "  // In View: table.Render()"
        ;;
    form)
        echo "USE: models.*Form wrapper (e.g., CaptureForm, SkillForm)"
        echo "NOT: *huh.Form directly in intent"
        echo ""
        echo "Location: internal/cli/models/"
        echo ""
        echo "Why: Direct huh.Form causes alignment issues on resize."
        echo "Wrapper handles WindowSizeMsg correctly."
        echo ""
        echo "Example:"
        echo "  type MyIntent struct {"
        echo "      form *models.CaptureForm  // CORRECT"
        echo "      // form *huh.Form         // WRONG"
        echo "  }"
        ;;
    color|theme|style)
        echo "USE: theme.Primary(), theme.Error(), etc."
        echo "NOT: lipgloss.Color(\"#xxx\")"
        echo ""
        echo "Location: internal/cli/uikit/theme/theme.go"
        echo ""
        echo "Available colors:"
        echo "  theme.Primary()    theme.Secondary()"
        echo "  theme.Success()    theme.Error()     theme.Warning()"
        echo "  theme.Muted()      theme.Background() theme.Foreground()"
        echo "  theme.Border()"
        echo ""
        echo "Example:"
        echo "  style := lipgloss.NewStyle().Foreground(theme.Error())"
        ;;
    view|layout)
        echo "USE: CreateStandardView() or CreateStandardViewWithBreadcrumbs()"
        echo "NOT: Manual view composition"
        echo ""
        echo "Location: internal/cli/intents/view_helpers.go"
        echo ""
        echo "Example:"
        echo "  view := CreateStandardView(i.BaseIntent)"
        echo "  view.WithContent(content)"
        echo "  view.WithHelp(ThemedNavigationFooter(i.Theme()))"
        echo "  return view.Render()"
        ;;
    footer|help)
        echo "USE: ThemedNavigationFooter(), ThemedListFooter(), etc."
        echo "NOT: Hardcoded footer strings"
        echo ""
        echo "Location: internal/cli/intents/view_helpers.go"
        echo ""
        echo "Available:"
        echo "  ThemedNavigationFooter(theme)"
        echo "  ThemedListFooter(theme)"
        echo "  ThemedFormFooter(theme)"
        echo "  CombineThemedFooters(footers...)"
        ;;
    modal|error|loading|success)
        echo "USE: BaseIntent state methods"
        echo ""
        echo "Methods:"
        echo "  i.SetError(err)           // Shows error modal"
        echo "  i.ClearError()"
        echo "  i.SetLoading(\"message\")   // Shows loading modal"
        echo "  i.ClearLoading()"
        echo "  i.SetSuccess(\"message\")   // Shows success (auto-dismiss)"
        echo "  i.SetProgress(title, msg, value)"
        echo ""
        echo "Modals render automatically via CreateStandardView()"
        ;;
    crud|create|edit|delete)
        echo "USE: behaviors.CRUDBehavior[T]"
        echo ""
        echo "Location: internal/cli/behaviors/crud.go"
        echo ""
        echo "Provides: Create/edit/delete with confirmation modals"
        ;;
    filter)
        echo "USE: behaviors.FilterMenuBehavior[T]"
        echo ""
        echo "Location: internal/cli/behaviors/filter_menu.go"
        echo ""
        echo "Provides: Multi-select filter menu with sections"
        ;;
    sort)
        echo "USE: behaviors.SortMenuBehavior[T]"
        echo ""
        echo "Location: internal/cli/behaviors/sort_menu.go"
        echo ""
        echo "Provides: Sort options menu with comparators"
        ;;
    text|title|body)
        echo "USE: uikit/primitives.Text"
        echo ""
        echo "Location: internal/cli/uikit/primitives/text.go"
        echo ""
        echo "Methods:"
        echo "  primitives.Title(text, theme)"
        echo "  primitives.Body(text, theme)"
        echo "  primitives.Muted(text, theme)"
        echo "  primitives.ErrorText(text, theme)"
        echo "  primitives.SuccessText(text, theme)"
        ;;
    button)
        echo "USE: uikit/primitives.Button, ButtonGroup"
        echo ""
        echo "Location: internal/cli/uikit/primitives/button.go"
        ;;
    input)
        echo "USE: uikit/primitives.Input"
        echo ""
        echo "Location: internal/cli/uikit/primitives/input.go"
        ;;
    intent|base)
        echo "USE: Embed *BaseIntent"
        echo ""
        echo "Location: internal/cli/intents/contract.go"
        echo ""
        echo "Example:"
        echo "  type MyIntent struct {"
        echo "      *BaseIntent  // REQUIRED"
        echo "      state   MyState"
        echo "      context *MyContext"
        echo "  }"
        echo ""
        echo "  func NewMyIntent() *MyIntent {"
        echo "      return &MyIntent{"
        echo "          BaseIntent: NewBaseIntent(),"
        echo "      }"
        echo "  }"
        ;;
    badge)
        echo "USE: uikit/primitives.Badge or components.KeyBadge"
        echo ""
        echo "Location: internal/cli/uikit/primitives/badge.go"
        echo ""
        echo "Variants:"
        echo "  primitives.NewBadge(text, theme).Variant(BadgeStatus)"
        echo "  primitives.NewBadge(text, theme).Variant(BadgeKey)"
        echo "  primitives.NewBadge(text, theme).Variant(BadgeTag)"
        echo ""
        echo "Presets:"
        echo "  primitives.NavigateBadge(theme)  // [↑↓ → Navigate]"
        echo "  primitives.SelectBadge(theme)    // [Enter → Select]"
        echo "  primitives.CancelBadge(theme)    // [Esc → Cancel]"
        echo "  primitives.HelpBadge(theme)      // [? → Help]"
        ;;
    box|container)
        echo "USE: uikit/containers.Box"
        echo ""
        echo "Location: internal/cli/uikit/containers/box.go"
        echo ""
        echo "Variants:"
        echo "  containers.NewBox(content, theme).Variant(BoxDefault)"
        echo "  containers.NewBox(content, theme).Variant(BoxEmphasized)"
        echo "  containers.NewBox(content, theme).Variant(BoxDestructive)"
        echo "  containers.NewBox(content, theme).Variant(BoxSubtle)"
        ;;
    spacing|padding|margin)
        echo "Standard spacing values:"
        echo ""
        echo "Padding:"
        echo "  Buttons:  Padding(0, 3)"
        echo "  Inputs:   Padding(0, 1)"
        echo "  Cards:    Padding(1, 2)"
        echo "  Modals:   Padding(1, 2)"
        echo "  Badges:   Padding(0, 1)"
        echo ""
        echo "Margins:"
        echo "  Between buttons: MarginRight(2)"
        echo "  After labels:    MarginBottom(1)"
        echo "  Section spacing: MarginTop(1)"
        echo ""
        echo "See: docs/rules/UI_STYLING.md for full reference"
        ;;
    border)
        echo "Border styles:"
        echo ""
        echo "States:"
        echo "  Default:  lipgloss.RoundedBorder(), theme.BorderColor()"
        echo "  Focused:  lipgloss.ThickBorder(), theme.PrimaryColor()"
        echo "  Error:    lipgloss.RoundedBorder(), theme.ErrorColor()"
        echo ""
        echo "See: docs/rules/UI_STYLING.md for full reference"
        ;;
    *)
        echo "Unknown keyword: $NEED"
        echo ""
        echo "Try: table, form, color, view, footer, modal, list, crud, filter, sort, text, button, input, intent, badge, box, spacing, border"
        ;;
esac

echo ""
echo "================================================"
