# 5-Minute Quick Test

## Start the app
```bash
./kariya
```

## Test Scenario: Create an Event

1. **Main Menu** → Select "Capture Event" (press 'c' or navigate with arrows)

2. **Fill Form**:
   - **Title**: "Built REST API for user management"
   - **Date**: "today" (or "2026-01-14")
   - **Company**: "Tech Corp"
   - **Description**: "Implemented RESTful API using Go and PostgreSQL with JWT authentication, Redis caching, and comprehensive test coverage. Included endpoints for user CRUD operations, role management, and audit logging."
   
   Press Tab to move between fields, Enter when done.

3. **Review Screen** (pre-save):
   - Should show your event details
   - NO bursts/facts yet
   - Press Enter or Ctrl+S to submit

4. **Wait for Save**:
   - Should see "Saving event..." or similar
   - Should see "Event saved successfully!" modal
   - Press Enter to dismiss modal

5. **✅ CRITICAL CHECK - Review Screen (post-save)**:
   - **EXPECTED**: You should return to Review screen (NOT main menu)
   - **EXPECTED**: Review should show:
     - Your event details
     - **"Inferred Bursts"** section with 1+ items (e.g., "REST API", "Authentication System")
     - **"Inferred Facts"** section with 1+ items (e.g., "Go", "PostgreSQL", "Redis", "JWT")
   - **NOT EXPECTED**: Going back to main menu immediately

6. **Complete the Intent**:
   - Press Enter or Ctrl+S
   - **EXPECTED**: NOW you go back to main menu
   - **NOT EXPECTED**: Seeing "Event saved successfully!" again (would be double-save bug)

## Success Criteria

✅ After save, you see Review screen with bursts and facts  
✅ Pressing Enter completes the intent (goes to main menu)  
✅ No double-save (modal doesn't appear twice)  

## If it fails

See `docs/development/POST_SAVE_REVIEW_TEST_PLAN.md` for debugging tips.
