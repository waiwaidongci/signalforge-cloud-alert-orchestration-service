# Bug

Concurrent HTTP metric observations read and write the recent-path map without a shared synchronization boundary.

# Trigger

Run the targeted recent-path accounting test with the race detector as specified in `collection.json`.

# Error

The race detector reports concurrent access to the recent-path map, and path accounting can lose updates.
