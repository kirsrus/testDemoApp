package web

import (
	"TestDemoApp/service"
	"TestDemoApp/service/web/graph"
	"TestDemoApp/service/web/graph/generated"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	port             = 80
	assetsDir        = "./assets"
	jwtTokenLifetime = time.Hour * 72
	jwtSecretString  = "Iju3uH4UO3F"
)

// Web имплементация интерфейса WebSvc. Инициируется через NewWeb
type Web struct {
	ctx context.Context
	log *logrus.Entry
	e   *echo.Echo

	resolver          *graph.Resolver
	graphqlHandler    *handler.Server
	playgroundHandler http.HandlerFunc

	port      int
	assetsDir string
}

// Config конфигурация конструктора NewWeb
type Config struct {
	Log       *logrus.Logger
	Port      int
	AssetsDir string
}

// NewWeb конструктор структуры Web
func NewWeb(ctx context.Context, config Config) (service.WebSvc, error) {
	var err error

	// region Конфигурация сервиса

	if config.Log == nil {
		config.Log = logrus.New()
		config.Log.Out = ioutil.Discard
	}

	web := Web{
		ctx: ctx,
		log: config.Log.WithFields(map[string]interface{}{
			"module": "web",
			"scope":  "service",
		}),
		port:      port,
		assetsDir: assetsDir,
		e:         echo.New(),
	}

	if config.Port != 0 {
		web.port = config.Port
	}
	if config.AssetsDir != "" {
		web.assetsDir = config.AssetsDir
	}

	// todo: при внутреннем хранении всех файлов
	//if _, err := os.Stat(web.assetsDir); os.IsNotExist(err) {
	//	return nil, errors.Annotate(err, "отсутсвует директория с WEB-контентом")
	//}

	tmpLevel := config.Log.Level
	config.Log.Level = logrus.InfoLevel
	web.log.Infof("config: port = %d", web.port)
	config.Log.Level = tmpLevel

	// endregion
	// region Настройка Web-сервера

	web.e.HideBanner = true
	web.e.HidePort = true
	web.e.Use(middleware.Recover())
	web.e.Use(middleware.Logger())
	web.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	web.e.Pre(web.addHeaders)

	// endregion
	// region Настройка GraphQL

	web.resolver, err = graph.NewResolver(graph.ConfigResolver{
		Log: config.Log,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	web.graphqlHandler = handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: web.resolver}))
	web.graphqlHandler.Use(extension.Introspection{})
	web.graphqlHandler.AddTransport(transport.POST{})
	web.graphqlHandler.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second, // Каждые 10 секунд подавать в канал (ping), иначе клиент его закроет
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			},
		})
	web.playgroundHandler = playground.Handler("GraphQL", "/graphql")

	web.e.GET("/graphql", web.handlerGrphQL)
	web.e.POST("/graphql", web.handlerGrphQL)
	web.e.GET("/graphqlplayground", web.handlerGrphQLPlayground)

	// endregion
	// region Настройка роутов

	// Из-за многоязычной поддержки и невозможности подставить в путь RegExp, приходится
	// дублировать запуск нужных роутов для всех поддерживаемых языков.
	// Некрасиво, но походу придётся переходить на другой сервер
	for _, lang := range []string{"ru", "en", "zh"} {

		// Аутентификация.
		web.e.POST("/"+lang+"/login", web.handlerLogin)

		// todo: затычка, пока не разобрался, почему браузер пытается с сервера взять виртуальный адрес.
		web.e.GET("/"+lang+"/login", func(c echo.Context) error {
			web.log.Warn(c.Request().URL)
			if m := regexp.MustCompile(`^(/\w+/).+`).FindStringSubmatch(c.Request().URL.String()); len(m) != 0 {
				return c.Redirect(http.StatusFound, m[1])
			}
			return c.Redirect(http.StatusFound, "/ru/")
		})

		// Группа из закрытой части
		secureGroup := web.e.Group("/" + lang + "/api")
		secureGroup.Use(middleware.JWT([]byte(jwtSecretString)))
		secureGroup.GET("", web.handlerSecretArea)

		// Языкозависимая статика
		//web.e.Static("/"+lang, web.assetsDir+"/"+lang)  // Сатика из файла
		web.e.GET("/"+lang+"/*", web.handlerStaticBindata)

	}

	// Редиректим всех на языкозавсимые точки входа
	web.e.GET("/", web.handlerRedirectRootToLang)

	// endregion

	return &web, nil
}

func (m Web) Run() error {
	done := make(chan error)
	go func() {
		err := m.e.Start(fmt.Sprintf(":%d", m.port))
		if err != nil {
			err = errors.Trace(err)
		}
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-m.ctx.Done():
		// Ошибку от становки сервера только логируем, т.к. она нам уже не важна
		err := m.Stop()
		if err != nil {
			m.log.Error(err)
		}

		return m.ctx.Err()
	}
}

func (m Web) Stop() error {
	return m.e.Close()
}

// Устанавливаем заголовки для выводимого контента
func (m Web) addHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Response().Header()
		headers.Set("Cache-Control", "no-store, no-cache, must-revalidate")
		headers.Set("Expires", strconv.Itoa(int(time.Now().Unix())))
		return next(c)
	}
}

// Точка подключения GraphQL
func (m Web) handlerGrphQL(c echo.Context) error {
	req := c.Request()
	res := c.Response()
	m.graphqlHandler.ServeHTTP(res, req)
	return nil
}

// Точка подключения GraphQL Playground
func (m Web) handlerGrphQLPlayground(c echo.Context) error {
	req := c.Request()
	res := c.Response()
	m.playgroundHandler.ServeHTTP(res, req)
	return nil
}

// Получаем файлы из внутренней памяти Bindata
func (m Web) handlerStaticBindata(c echo.Context) error {
	var file []byte
	var err error

	file, err = Asset("assets" + c.Request().URL.Path)
	if err != nil {
		// Возможно это обращение к корневой директории. Нужно проверить на наличие "index.html"
		file, _ = Asset("assets" + c.Request().URL.Path + "index.html")
	}

	if file == nil {
		return c.NoContent(http.StatusNotFound)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(c.Request().URL.Path))
	return c.Blob(http.StatusOK, mimeType, file)
}

// Хэндлер переадресации входящих соединений в корневой раздел на
// раздел с нужным языком
func (m Web) handlerRedirectRootToLang(c echo.Context) error {
	// Определяем предпочитаемый язык вызываемого браузера. Если он
	// неизвестен, пересылаем на русскую страницу
	prefix := "ru"
	if headLand := c.Request().Header.Get("Accept-Language"); headLand != "" {
		hList := strings.Split(headLand, ",")
		switch hList[0] {
		case "ru", "en", "zh":
			prefix = hList[0]
		}
	}

	// Переадресовываем на нужный префикс
	m.log.Debugf("переадресация на /%s по заголовку '%s'", prefix, c.Request().Header.Get("Accept-Language"))
	return c.Redirect(http.StatusFound, "/"+prefix+"/")
}

func (m Web) handlerLogin(c echo.Context) error {
	type loginData struct {
		Username string
		Password string
	}

	u := new(loginData)
	if err := c.Bind(u); err != nil {
		return err
	}

	username := u.Username
	password := u.Password

	// todo: заглушка места проверки авторизации. Далее сделать и БД
	if username != "admin" || password != "admin" {
		m.log.Infof("пользователь %s не прошёл авторизацию", username)
		return echo.ErrUnauthorized
	}

	fmt.Println("username=", username, "password=", password)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Admin"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(jwtTokenLifetime).Unix()

	t, err := token.SignedString([]byte(jwtSecretString))
	if err != nil {
		m.log.Warnf("ошибка создание JWT токена: %s", err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (m Web) handlerSecretArea(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
