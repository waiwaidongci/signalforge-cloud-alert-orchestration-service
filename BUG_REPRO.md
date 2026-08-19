# Bug 是什么
路由解析会原地排序仓储返回的 rules，并返回 channels 的底层引用，导致规则状态被调用方污染。
# 如何触发
连续两次调用 `Service.Resolve`，随后修改第一次返回的 channels。
# 错误信息
第二次解析顺序变化或返回被修改的 destination。
