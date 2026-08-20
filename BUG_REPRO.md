# Bug reproduction

- Bug: a disabled source policy provider leaves the default policy unusable.
- Trigger: resolve a policy with the optional provider disabled, then add and evaluate a required label.
- Error: the baseline returns no usable policy or panics while adding a rule.
