package middleware

//拦截请求的鉴权中间件，为了验证是否有合法token
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 挑战码结构体
type ChallengeResp struct {
	CodeVerifier  string `json:"code_verifier"`
	CodeChallenge string `json:"code_challenge"`
}

// 挑战码接口返回的结构体
type CodeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CodeVerifier  string `json:"code_verifier"`
		CodeChallenge string `json:"code_challenge"`
	} `json:"data"`
}

func AuthMiddleware(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// 没有 token，则调用生成挑战码的接口
		resp, err := http.Get("http://127.0.0.1:8080/api/auth/get-generate-code")
		if err != nil || resp.StatusCode != http.StatusOK {
			r.Response.WriteStatusExit(http.StatusInternalServerError, g.Map{"error": "获取挑战码失败"})
			return
		}
		defer resp.Body.Close()

		// 解析 JSON 响应
		var codeResp CodeResponse
		if err := json.NewDecoder(resp.Body).Decode(&codeResp); err != nil {
			r.Response.WriteStatusExit(http.StatusInternalServerError, g.Map{"error": "解析挑战码响应失败"})
			return
		}

		// 打印挑战码
		g.Log().Infof(r.GetCtx(), "code_verifier: %s", codeResp.Data.CodeVerifier)
		g.Log().Infof(r.GetCtx(), "code_challenge: %s", codeResp.Data.CodeChallenge)

		// 将 code_verifier 存入 session
		r.Session.Set("code_verifier", codeResp.Data.CodeVerifier)

		// 2.0构建授权 URL，跳转去登录获得授权码code
		clientId := "gf_api"                                                                                                               //客户端ID，固定的。在统一开放授权认证平台生成的。
		redirectUri := "http://10.170.0.96:8001"                                                                                           //这个跳转界面必须是vue做的前端项目。                                                                                         //此地址在统一开放授权认证平台注册了，谁开发改谁的地址                                                                                         //谁测试，就填谁的开发ip；正式部署后，填部署服务器的ip
		authorizeUrl := "http://10.170.1.30:5001/connect/authorize"                                                                        //统一开放授权认证平台，固定的请求接口
		scope := "ApiResourceScope SRoles BranchUnits MainTainDepts StationNames Authoritys Names RealNames openid profile offline_access" //固定的，不要动。
		//用 fmt.Sprintf 函数拼接一个带参数的 URL 字符串
		authUrl := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&code_challenge=%s&code_challenge_method=S256",
			authorizeUrl,
			url.QueryEscape(clientId),
			url.QueryEscape(redirectUri),
			url.QueryEscape(scope),
			url.QueryEscape(codeResp.Data.CodeChallenge),
		)

		// 重定向到授权平台
		r.Response.RedirectTo(authUrl)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"error": "token 格式错误"})
		return
	}

	// 验证token
	resp, err := g.Client().Post(r.GetCtx(), "http://auth-service/api/validate", g.Map{
		"token": token,
	})
	if err != nil {
		r.Response.WriteStatusExit(http.StatusInternalServerError, g.Map{"error": "验证token失败或无token访问"})
		return
	}
	defer resp.Close()

	body := resp.ReadAll()
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		r.Response.WriteStatusExit(http.StatusInternalServerError, g.Map{"error": "解析鉴权响应失败"})
		return
	}

	if result.Code != 0 {
		r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"error": result.Msg})
		return
	}

	// 设置用户上下文变量，方便后续业务使用
	r.SetCtxVar("userID", result.Data.UserID)

	// 继续执行后续中间件或请求处理函数
	r.Middleware.Next()
}
