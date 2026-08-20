# Bug reproduction

- Bug: a zero-value ID allocator dereferences a nil generator or writes to a nil map.
- Trigger: allocate an identifier from an uninitialized allocator.
- Error: the baseline panics instead of returning and recording a prefixed identifier.
