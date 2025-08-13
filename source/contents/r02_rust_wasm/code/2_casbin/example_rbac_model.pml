# 请求定义
[request_definition]
r = sub, obj, act
# 策略定义
[policy_definition]
p = sub, obj, act, desc
# 角色定义
[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
