# Bug

Timeline row iteration, deferred close, error merging, and repository classification discard close failures and break the original iteration error chain.

# Trigger

Run the two targeted timeline row lifecycle tests from `collection.json` against the bug snapshot.

# Error

The close-only case returns `err=<nil>`. The combined case returns only `读取时间线失败`, without the iteration or close errors.
