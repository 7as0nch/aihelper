package base

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type qqStateItem struct {
	Redirect string
	ExpireAt time.Time
}

type qqOAuthConfig struct {
	AppID       string
	AppKey      string
	CallbackURL string
	Scope       string
}

var qqHTTPClient = &http.Client{Timeout: 10 * time.Second}

func (s *AuthService) HandleQQLogin(w http.ResponseWriter, r *http.Request) {
	cfg, err := loadQQOAuthConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := strings.TrimSpace(r.URL.Query().Get("redirect"))
	if redirectURL == "" {
		redirectURL = strings.TrimSpace(os.Getenv("QQ_FRONTEND_REDIRECT"))
	}
	if redirectURL == "" {
		redirectURL = "/"
	}

	state, err := newQQState()
	if err != nil {
		http.Error(w, "failed to create qq state", http.StatusInternalServerError)
		return
	}
	s.putQQState(state, redirectURL)

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", cfg.AppID)
	params.Set("redirect_uri", cfg.CallbackURL)
	params.Set("state", state)
	params.Set("scope", cfg.Scope)

	http.Redirect(w, r, "https://graph.qq.com/oauth2.0/authorize?"+params.Encode(), http.StatusFound)
}

func (s *AuthService) HandleQQCallback(w http.ResponseWriter, r *http.Request) {
	cfg, err := loadQQOAuthConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := s.popQQState(strings.TrimSpace(r.URL.Query().Get("state")))
	if redirectURL == "" {
		redirectURL = strings.TrimSpace(os.Getenv("QQ_FRONTEND_REDIRECT"))
	}
	if redirectURL == "" {
		redirectURL = "/"
	}

	if oauthErr := strings.TrimSpace(r.URL.Query().Get("error")); oauthErr != "" {
		http.Redirect(w, r, appendQuery(redirectURL, "qq_error", oauthErr), http.StatusFound)
		return
	}

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		http.Redirect(w, r, appendQuery(redirectURL, "qq_error", "missing_code"), http.StatusFound)
		return
	}

	accessToken, err := exchangeQQAccessToken(r.Context(), cfg, code)
	if err != nil {
		log.Errorf("qq exchange token failed: %v", err)
		http.Redirect(w, r, appendQuery(redirectURL, "qq_error", "exchange_token_failed"), http.StatusFound)
		return
	}

	openID, err := fetchQQOpenID(r.Context(), accessToken)
	if err != nil {
		log.Errorf("qq get openid failed: %v", err)
		http.Redirect(w, r, appendQuery(redirectURL, "qq_error", "get_openid_failed"), http.StatusFound)
		return
	}

	nickname, avatar, err := fetchQQUserInfo(r.Context(), cfg.AppID, accessToken, openID)
	if err != nil {
		log.Warnf("qq get user info failed, continue with openid only: %v", err)
	}

	token, _, err := s.user.LoginByQQ(r.Context(), openID, nickname, avatar)
	if err != nil {
		log.Errorf("qq login failed: %v", err)
		http.Redirect(w, r, appendQuery(redirectURL, "qq_error", "qq_login_failed"), http.StatusFound)
		return
	}

	http.Redirect(w, r, appendQuery(redirectURL, "qq_token", token), http.StatusFound)
}

func loadQQOAuthConfig() (*qqOAuthConfig, error) {
	cfg := &qqOAuthConfig{
		AppID:       strings.TrimSpace(os.Getenv("QQ_APP_ID")),
		AppKey:      strings.TrimSpace(os.Getenv("QQ_APP_KEY")),
		CallbackURL: strings.TrimSpace(os.Getenv("QQ_CALLBACK_URL")),
		Scope:       strings.TrimSpace(os.Getenv("QQ_SCOPE")),
	}
	if cfg.Scope == "" {
		cfg.Scope = "get_user_info"
	}
	if cfg.AppID == "" || cfg.AppKey == "" || cfg.CallbackURL == "" {
		return nil, fmt.Errorf("qq oauth config missing, require QQ_APP_ID/QQ_APP_KEY/QQ_CALLBACK_URL")
	}
	return cfg, nil
}

func exchangeQQAccessToken(ctx context.Context, cfg *qqOAuthConfig, code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", cfg.AppID)
	params.Set("client_secret", cfg.AppKey)
	params.Set("code", code)
	params.Set("redirect_uri", cfg.CallbackURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.qq.com/oauth2.0/token?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	resp, err := qqHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("qq token http status=%d body=%s", resp.StatusCode, truncateForLog(string(body)))
	}

	bodyText := strings.TrimSpace(string(body))
	parsed, err := url.ParseQuery(bodyText)
	if err == nil {
		if accessToken := strings.TrimSpace(parsed.Get("access_token")); accessToken != "" {
			return accessToken, nil
		}
	}

	return "", fmt.Errorf("qq token response invalid: %s", truncateForLog(bodyText))
}

func fetchQQOpenID(ctx context.Context, accessToken string) (string, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("fmt", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.qq.com/oauth2.0/me?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	resp, err := qqHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("qq me http status=%d body=%s", resp.StatusCode, truncateForLog(string(body)))
	}

	var payload struct {
		OpenID string `json:"openid"`
		Error  int    `json:"error"`
		Msg    string `json:"error_description"`
	}
	if err = json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.Error != 0 {
		return "", fmt.Errorf("qq me error=%d msg=%s", payload.Error, payload.Msg)
	}
	if strings.TrimSpace(payload.OpenID) == "" {
		return "", fmt.Errorf("qq openid missing")
	}
	return payload.OpenID, nil
}

func fetchQQUserInfo(ctx context.Context, appID, accessToken, openID string) (nickname, avatar string, err error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("oauth_consumer_key", appID)
	params.Set("openid", openID)
	params.Set("format", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.qq.com/user/get_user_info?"+params.Encode(), nil)
	if err != nil {
		return "", "", err
	}

	resp, err := qqHTTPClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("qq user info http status=%d body=%s", resp.StatusCode, truncateForLog(string(body)))
	}

	var payload struct {
		Ret         int    `json:"ret"`
		Msg         string `json:"msg"`
		Nickname    string `json:"nickname"`
		FigureurlQQ string `json:"figureurl_qq_1"`
		FigureurlQQ2 string `json:"figureurl_qq_2"`
	}
	if err = json.Unmarshal(body, &payload); err != nil {
		return "", "", err
	}
	if payload.Ret != 0 {
		return "", "", fmt.Errorf("qq user info ret=%d msg=%s", payload.Ret, payload.Msg)
	}

	avatar = strings.TrimSpace(payload.FigureurlQQ2)
	if avatar == "" {
		avatar = strings.TrimSpace(payload.FigureurlQQ)
	}
	return strings.TrimSpace(payload.Nickname), avatar, nil
}

func (s *AuthService) putQQState(state, redirectURL string) {
	now := time.Now()
	s.qqStateMu.Lock()
	defer s.qqStateMu.Unlock()

	for k, v := range s.qqState {
		if now.After(v.ExpireAt) {
			delete(s.qqState, k)
		}
	}

	s.qqState[state] = qqStateItem{
		Redirect: redirectURL,
		ExpireAt: now.Add(10 * time.Minute),
	}
}

func (s *AuthService) popQQState(state string) string {
	if state == "" {
		return ""
	}
	s.qqStateMu.Lock()
	defer s.qqStateMu.Unlock()

	item, ok := s.qqState[state]
	if !ok {
		return ""
	}
	delete(s.qqState, state)
	if time.Now().After(item.ExpireAt) {
		return ""
	}
	return item.Redirect
}

func appendQuery(rawURL, key, value string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u == nil {
		u = &url.URL{Path: "/"}
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String()
}

func newQQState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func truncateForLog(s string) string {
	const maxLen = 240
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

