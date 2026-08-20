# Bug

Timeline events and metadata escape through storage, pagination, and HTTP response boundaries as shared mutable values. Invalid paging windows can also panic.

# Trigger

Run the four targeted timeline isolation tests from `collection.json` against the bug snapshot.

# Error

The tests show that mutating a returned event changes later reads, metadata is shared, and an invalid window does not return the required error.
