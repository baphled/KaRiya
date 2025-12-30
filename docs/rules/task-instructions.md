Use `master-task-prompt` and the memory-loaded language rules to process the atomic tasks for: ./tasks/tasks-03-metadata-clarification.md

  - Stick strictly to `master-task-prompt` guidelines
    - This means make sure we use it as a check-list only
    - Check that we follow all the rules below
  - Reference the codebase before task generation
  - Make sure to *strictly* follow the notes and guidelines
  - Embed Go guidelines and rules inline where applicable
  - Use one expectation per `it`, Red→Green→Refactor model
  - Follow existing patterns and avoid side effects
  - Return checklist only under `## Tasks`
  - Do not explain or regenerate later
