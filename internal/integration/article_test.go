package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike_gin/internal/repository"
	"jike_gin/internal/repository/dao"
	"jike_gin/internal/service"
	"jike_gin/internal/web"
	ijwt "jike_gin/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// ArticleSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// SetupSuite 执行所有的测试之前执行初始化逻辑
func (s *ArticleTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(c *gin.Context) {
		c.Set("claims", &ijwt.UserClaims{Uid: 123})
	})
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	s.db = db
	d := dao.NewGORMArticleDAO(db)
	repo := repository.NewArticleRepository(d)
	svc := service.NewArticleService(repo)
	articleHandler := web.NewArticleHandler(svc)
	// 注册路由
	articleHandler.RegisterRoutes(s.server)

}

func (s *ArticleTestSuite) TearDownTest() {
	//s.db.Exec("TRUNCATE TABLE articles;")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		wantCode int
		wantRes  Result[int64]
		art      Article
	}{
		{
			name: "插入数据",
			before: func(t *testing.T) {
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			after: func(t *testing.T) {
				var a dao.Article
				err := s.db.Table("article").Where("id = ?", 1).Find(&a).Error
				require.NoError(t, err)

				assert.Equal(t, dao.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
				}, a)
			},
			wantCode: http.StatusOK,
			wantRes:  Result[int64]{Data: 1},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)

			reqBody, err := json.Marshal(tt.art)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req

			resp := httptest.NewRecorder()

			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRes, webRes)

			tt.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello，这是测试套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}
