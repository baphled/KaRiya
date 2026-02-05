# KaRiya Import Workflow Demo
# Duration: ~30 seconds
# Shows: View CSV contents -> Import command -> See results
#
# This demo shows how to import career events from a CSV file.
# It first displays the CSV contents, then runs the import.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/09-import-workflow.gif
Output demos/vhs/output/mp4s/09-import-workflow.mp4

# Setup (hidden)
Hide
Type `rm -rf /tmp/kariya-demo && mkdir -p /tmp/kariya-demo`
Enter
Sleep 300ms
Type `cp demos/vhs/minimal-config.yaml /tmp/kariya-demo/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KoRiya"
Enter
Sleep 300ms
Show

# Show the CSV contents first
Type "echo '=== Career Events CSV ===' && head -12 demos/vhs/sample-events.csv | column -t -s ','"
Enter
Sleep 4s

# Show total count
Type "echo '' && echo 'Total events:' $(tail -n +2 demos/vhs/sample-events.csv | wc -l)"
Enter
Sleep 2s

# Run import
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review"
Enter

# Wait for import to complete
Sleep 5s

# Show success message
Type "echo '' && echo 'Career events imported successfully!'"
Enter
Sleep 2s
