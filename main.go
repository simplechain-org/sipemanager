package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"net/http"
	"sipemanager/controllers"
	"sipemanager/dao"
	"strings"
)

var configPath = flag.String("conf", "./etc/config.yaml", "-conf=./etc/config.yaml")

type AppConfig struct {
	Db   *dao.DBConfig `yaml:"db"`
	Port uint          `yaml:"port"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var config AppConfig
	err := configor.Load(&config, path)
	if err != nil {
		return nil, err
	}
	return &config, err
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	flag.Parse()
	appConfig, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := dao.GetDBConnection(appConfig.Db)
	if err != nil {
		fmt.Println(err)
		return
	}
	dao.AutoMigrate(db)

	dao := dao.NewDataBaseAccessObject(db)

	router := gin.Default()
	router.Use(Cors())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	controllers.Register(router, dao)
	controllers.SwaggerDoc(router)

	router.StaticFile("/", "./webroot/dist/index.html")

	//加载静态资源，例如网页的css、js
	router.Static("/static", "./webroot/dist/static")

	router.Run(fmt.Sprintf(":%d", appConfig.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

////// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
