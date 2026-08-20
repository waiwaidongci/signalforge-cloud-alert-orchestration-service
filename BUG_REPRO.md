# Bug

Escalation policies with malformed matcher or route JSON are returned as usable policies instead of preserving a decode error through the repository and HTTP layers.

# Trigger

Run the four targeted policy decode and handler tests from `collection.json` against the bug snapshot.

# Error

The tests report that corrupt stored JSON, corrupt route JSON, and a `null` selector were accepted, and the handler failed to return the expected bad-request response.
