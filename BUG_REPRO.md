# Bug reproduction

- Bug: dedup batch normalization and storage share caller slice memory.
- Trigger: normalize a batch, mutate the derived labels, then inspect the original input and stored batch.
- Error: earlier or later stages contain the mutation unexpectedly.
