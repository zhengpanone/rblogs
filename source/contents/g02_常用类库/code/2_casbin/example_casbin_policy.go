// 数据库定义
type CasbinPolicy struct {
	PType  string `json:"p_type" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
	Desc   string `json:"desc" binding:"required"`
}

// 初始化执行器
func InitCasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		TPLogger.Error("初始化casbin策略管理器失败：", err)
		panic(err)
	}
	e.EnableAutoSave(true)
	CasbinEnforcer = e
}

// 模型加载
import(
		"github.com/casbin/gorm-adapter/v3"
		"github.com/casbin/casbin/v2"
)
func mysqlCasbin() (*casbin.Enforcer, error) {
	// 1. 初始化适配器Adapter
	// 方式一：默认方式
	adapter, err := gormadapter.NewAdapterByDB(GORM)
	if err != nil {
		TPLogger.Error("casbin adapter gorm failed: ", err)
		return nil, err
	}

	// 方式二：使用自定义表明的方式
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db,"","my_casbin_rules")
	if err !=nil{
	panic("failed to initialize adapter")
	}
	// 2. 初始化Enforcer。创建Enforcer时自动从数据库中加载策略。
	e, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic("failed to create enforcer")
	}
	// 3. 加载策略,将策略加载到内存中
	if err = e.LoadPolicy(); err != nil {
		panic("failed to load policy failed")
	}
	return e, nil
}

// 中间件定义
func CasbinMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取username
		ctxUser := c.GetString("username")
		if ctxUser == "admin" {
			c.Next()
		} else {
			// TODO 从缓存中获取用户相关的信息，例如：role、dept、menu
			// 从数据库中获取用户角色信息sub
			usersInfo, err := system.NewUserInterface().GetUserFromUserName(ctxUser)
			if err != nil {
				global.TPLogger.Error("从数据库中获取用户角色信息sub失败:", err)
				global.ReturnContext(c).Failed("failed", "权限访问失败,请联系管理员")
				c.Abort()
				return
			}
			sub := usersInfo.RoleId
			//获取请求路径
			obj := strings.Split(c.Request.RequestURI, "?")[0]
			// 获取请求方法
			act := c.Request.Method
			success, err := global.CasbinEnforcer.Enforce(strconv.Itoa(int(sub)), obj, act)
			if err != nil || !success {
				global.TPLogger.Error("权限验证失败：", err, success)
				global.ReturnContext(c).Failed("failed", "权限验证失败")
				c.Abort()
				return
			} else {
				c.Next()
			}
		}

	}
}