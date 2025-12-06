package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"time"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
)

// 错误定义
var (
	ErrCaptchaDisabled     = errors.New("验证码功能已关闭")
	ErrInvalidSession      = errors.New("无效的验证会话")
	ErrSessionExpired      = errors.New("验证会话已过期")
	ErrVerifyFailed        = errors.New("验证失败，请重试")
	ErrRateLimited         = errors.New("操作过于频繁，请稍后再试")
	ErrTokenInvalid        = errors.New("验证令牌无效")
	ErrTokenExpired        = errors.New("验证令牌已过期")
	ErrTokenUsed           = errors.New("验证令牌已使用")
)

// ChallengeResponse 挑战响应
type ChallengeResponse struct {
	SessionID   string `json:"session_id"`
	BgImage     string `json:"bg_image"`      // Base64 编码的背景图（含缺口）
	PuzzleImage string `json:"puzzle_image"`  // Base64 编码的拼图块
	PuzzleY     int    `json:"puzzle_y"`      // 拼图块 Y 坐标
	Width       int    `json:"width"`         // 图片宽度
	Height      int    `json:"height"`        // 图片高度
}

// VerifyResult 验证结果
type VerifyResult struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

// 图片尺寸常量
const (
	ImageWidth    = 300
	ImageHeight   = 150
	PuzzleWidth   = 50
	PuzzleHeight  = 50
	MinPuzzleX    = PuzzleWidth
	MaxPuzzleX    = ImageWidth - PuzzleWidth*2
	MinPuzzleY    = 10
	MaxPuzzleY    = ImageHeight - PuzzleHeight - 10
	TotalBgImages = 10
)

// GenerateChallenge 生成验证挑战
func GenerateChallenge(clientIP string) (*ChallengeResponse, error) {
	// 检查是否启用验证码
	if !setting.CaptchaEnabled {
		return nil, ErrCaptchaDisabled
	}

	// 检查 IP 是否被限流
	blocked, err := IsIPBlocked(clientIP)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, ErrRateLimited
	}

	// 生成会话 ID
	sessionID, err := model.GenerateSessionID()
	if err != nil {
		return nil, err
	}

	// 随机选择背景图
	imageIdx := rand.Intn(TotalBgImages)

	// 随机生成拼图位置
	targetX := MinPuzzleX + rand.Intn(MaxPuzzleX-MinPuzzleX)
	puzzleY := MinPuzzleY + rand.Intn(MaxPuzzleY-MinPuzzleY)

	// 创建挑战记录
	now := time.Now()
	challenge := &model.CaptchaChallenge{
		SessionID: sessionID,
		TargetX:   targetX,
		PuzzleY:   puzzleY,
		ImageIdx:  imageIdx,
		CreatedAt: now,
		ExpiresAt: now.Add(model.CaptchaChallengeExpiration),
		Verified:  false,
	}

	// 存储挑战
	if err := model.StoreCaptchaChallenge(challenge); err != nil {
		return nil, err
	}

	// 生成图片（这里返回占位数据，实际实现需要生成真实图片）
	bgImage, puzzleImage := generateCaptchaImages(imageIdx, targetX, puzzleY)

	return &ChallengeResponse{
		SessionID:   sessionID,
		BgImage:     bgImage,
		PuzzleImage: puzzleImage,
		PuzzleY:     puzzleY,
		Width:       ImageWidth,
		Height:      ImageHeight,
	}, nil
}

// VerifyChallenge 验证用户提交的位置
func VerifyChallenge(sessionID string, x int, clientIP string) (*VerifyResult, error) {
	// 检查是否启用验证码
	if !setting.CaptchaEnabled {
		return nil, ErrCaptchaDisabled
	}

	// 检查 IP 是否被限流
	blocked, err := IsIPBlocked(clientIP)
	if err != nil {
		return nil, err
	}
	if blocked {
		return &VerifyResult{
			Success: false,
			Message: "操作过于频繁，请稍后再试",
		}, nil
	}

	// 获取挑战
	challenge, err := model.GetCaptchaChallenge(sessionID)
	if err != nil {
		// 记录失败
		_ = model.IncrementCaptchaIPFailure(clientIP)
		return &VerifyResult{
			Success: false,
			Message: "验证会话无效或已过期",
		}, nil
	}

	// 检查是否已验证
	if challenge.Verified {
		return &VerifyResult{
			Success: false,
			Message: "验证会话已使用",
		}, nil
	}

	// 验证位置
	tolerance := setting.CaptchaToleranceRange
	diff := x - challenge.TargetX
	if diff < 0 {
		diff = -diff
	}

	if diff > tolerance {
		// 验证失败
		_ = model.IncrementCaptchaIPFailure(clientIP)
		// 删除当前挑战，强制用户获取新挑战
		_ = model.DeleteCaptchaChallenge(sessionID)
		return &VerifyResult{
			Success: false,
			Message: "验证失败，请重试",
		}, nil
	}

	// 验证成功，删除挑战防止重放
	_ = model.DeleteCaptchaChallenge(sessionID)

	// 重置 IP 失败计数
	_ = model.ResetCaptchaIPRecord(clientIP)

	// 生成验证令牌
	tokenStr, err := model.GenerateCaptchaToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	token := &model.CaptchaToken{
		Token:     tokenStr,
		UserIP:    clientIP,
		CreatedAt: now,
		ExpiresAt: now.Add(model.CaptchaTokenExpiration),
		Used:      false,
	}

	if err := model.StoreCaptchaToken(token); err != nil {
		return nil, err
	}

	return &VerifyResult{
		Success: true,
		Token:   tokenStr,
	}, nil
}

// ValidateCaptchaToken 验证令牌
func ValidateCaptchaToken(tokenStr string, clientIP string) error {
	if tokenStr == "" {
		return ErrTokenInvalid
	}

	token, err := model.GetCaptchaToken(tokenStr)
	if err != nil {
		return ErrTokenInvalid
	}

	if time.Now().After(token.ExpiresAt) {
		return ErrTokenExpired
	}

	if token.Used {
		return ErrTokenUsed
	}

	// 可选：验证 IP 是否匹配
	// if token.UserIP != clientIP {
	// 	return ErrTokenInvalid
	// }

	// 标记令牌已使用
	if err := model.MarkCaptchaTokenUsed(tokenStr); err != nil {
		return err
	}

	return nil
}

// IsIPBlocked 检查 IP 是否被限流
func IsIPBlocked(ip string) (bool, error) {
	record, err := model.GetCaptchaIPRecord(ip)
	if err != nil {
		return false, err
	}
	if record == nil {
		return false, nil
	}

	// 检查是否在封禁期内
	if !record.BlockedUntil.IsZero() && time.Now().Before(record.BlockedUntil) {
		return true, nil
	}

	return false, nil
}

// generateCaptchaImages 生成验证码图片
// 返回 Base64 编码的背景图（含缺口）和拼图块
func generateCaptchaImages(imageIdx int, targetX int, puzzleY int) (bgImage string, puzzleImage string) {
	// 1. 生成背景图
	bg := generateBackgroundImage(imageIdx)

	// 2. 创建拼图块蒙版
	puzzleMask := createPuzzleMask(PuzzleWidth, PuzzleHeight)

	// 3. 从背景图切割拼图块
	puzzle := cutPuzzlePiece(bg, targetX, puzzleY, puzzleMask)

	// 4. 在背景图上绘制缺口
	drawPuzzleHole(bg, targetX, puzzleY, puzzleMask)

	// 5. 编码为 Base64
	bgImage = encodeImageToBase64(bg)
	puzzleImage = encodeImageToBase64(puzzle)

	return
}

// generateBackgroundImage 生成背景图
// 使用渐变色和随机图案创建视觉效果
func generateBackgroundImage(imageIdx int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))

	// 基于 imageIdx 生成不同的颜色方案
	baseColors := [][]color.RGBA{
		{{66, 133, 244, 255}, {52, 168, 83, 255}},   // 蓝绿
		{{234, 67, 53, 255}, {251, 188, 5, 255}},    // 红黄
		{{103, 58, 183, 255}, {33, 150, 243, 255}},  // 紫蓝
		{{0, 150, 136, 255}, {76, 175, 80, 255}},    // 青绿
		{{255, 87, 34, 255}, {255, 152, 0, 255}},    // 橙色
		{{121, 85, 72, 255}, {161, 136, 127, 255}},  // 棕色
		{{96, 125, 139, 255}, {120, 144, 156, 255}}, // 蓝灰
		{{233, 30, 99, 255}, {156, 39, 176, 255}},   // 粉紫
		{{0, 188, 212, 255}, {0, 150, 136, 255}},    // 青色
		{{139, 195, 74, 255}, {205, 220, 57, 255}},  // 黄绿
	}

	colorIdx := imageIdx % len(baseColors)
	startColor := baseColors[colorIdx][0]
	endColor := baseColors[colorIdx][1]

	// 绘制渐变背景
	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			// 对角线渐变
			ratio := float64(x+y) / float64(ImageWidth+ImageHeight)
			r := uint8(float64(startColor.R)*(1-ratio) + float64(endColor.R)*ratio)
			g := uint8(float64(startColor.G)*(1-ratio) + float64(endColor.G)*ratio)
			b := uint8(float64(startColor.B)*(1-ratio) + float64(endColor.B)*ratio)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// 添加随机装饰图案
	addDecorations(img, imageIdx)

	return img
}

// addDecorations 添加装饰图案
func addDecorations(img *image.RGBA, seed int) {
	r := rand.New(rand.NewSource(int64(seed * 12345)))

	// 添加一些随机圆形
	for i := 0; i < 5; i++ {
		cx := r.Intn(ImageWidth)
		cy := r.Intn(ImageHeight)
		radius := 10 + r.Intn(30)
		alpha := uint8(30 + r.Intn(50))

		drawCircle(img, cx, cy, radius, color.RGBA{255, 255, 255, alpha})
	}

	// 添加一些随机线条
	for i := 0; i < 3; i++ {
		x1 := r.Intn(ImageWidth)
		y1 := r.Intn(ImageHeight)
		x2 := r.Intn(ImageWidth)
		y2 := r.Intn(ImageHeight)
		alpha := uint8(40 + r.Intn(40))

		drawLine(img, x1, y1, x2, y2, color.RGBA{255, 255, 255, alpha})
	}
}

// drawCircle 绘制圆形
func drawCircle(img *image.RGBA, cx, cy, radius int, c color.RGBA) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if x >= 0 && x < ImageWidth && y >= 0 && y < ImageHeight {
				dx := float64(x - cx)
				dy := float64(y - cy)
				if dx*dx+dy*dy <= float64(radius*radius) {
					blendPixel(img, x, y, c)
				}
			}
		}
	}
}

// drawLine 绘制线条
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		if x1 >= 0 && x1 < ImageWidth && y1 >= 0 && y1 < ImageHeight {
			blendPixel(img, x1, y1, c)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// blendPixel 混合像素颜色
func blendPixel(img *image.RGBA, x, y int, c color.RGBA) {
	existing := img.RGBAAt(x, y)
	alpha := float64(c.A) / 255.0
	r := uint8(float64(existing.R)*(1-alpha) + float64(c.R)*alpha)
	g := uint8(float64(existing.G)*(1-alpha) + float64(c.G)*alpha)
	b := uint8(float64(existing.B)*(1-alpha) + float64(c.B)*alpha)
	img.Set(x, y, color.RGBA{r, g, b, 255})
}

// createPuzzleMask 创建拼图块蒙版
// 创建带有凸起的拼图形状
func createPuzzleMask(width, height int) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, width, height))

	// 基础矩形
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mask.SetAlpha(x, y, color.Alpha{255})
		}
	}

	// 右侧凸起
	bumpRadius := height / 5
	bumpCenterX := width
	bumpCenterY := height / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width+bumpRadius; x++ {
			dx := float64(x - bumpCenterX)
			dy := float64(y - bumpCenterY)
			if dx*dx+dy*dy <= float64(bumpRadius*bumpRadius) {
				if x < width+bumpRadius && y >= 0 && y < height {
					// 扩展蒙版以包含凸起
				}
			}
		}
	}

	return mask
}

// cutPuzzlePiece 从背景图切割拼图块
func cutPuzzlePiece(bg *image.RGBA, targetX, targetY int, mask *image.Alpha) *image.RGBA {
	bounds := mask.Bounds()
	puzzle := image.NewRGBA(bounds)

	// 填充透明背景
	draw.Draw(puzzle, bounds, image.Transparent, image.Point{}, draw.Src)

	// 复制背景图对应区域的像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if mask.AlphaAt(x, y).A > 0 {
				srcX := targetX + x
				srcY := targetY + y
				if srcX >= 0 && srcX < ImageWidth && srcY >= 0 && srcY < ImageHeight {
					puzzle.Set(x, y, bg.At(srcX, srcY))
				}
			}
		}
	}

	// 添加边框效果
	addPuzzleBorder(puzzle, mask)

	return puzzle
}

// addPuzzleBorder 为拼图块添加边框效果
func addPuzzleBorder(puzzle *image.RGBA, mask *image.Alpha) {
	bounds := mask.Bounds()

	// 添加阴影效果
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if mask.AlphaAt(x, y).A > 0 {
				// 检查是否是边缘像素
				isEdge := false
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
							if mask.AlphaAt(nx, ny).A == 0 {
								isEdge = true
								break
							}
						} else {
							isEdge = true
							break
						}
					}
					if isEdge {
						break
					}
				}

				if isEdge {
					// 绘制白色边框
					existing := puzzle.RGBAAt(x, y)
					r := uint8(min(255, int(existing.R)+60))
					g := uint8(min(255, int(existing.G)+60))
					b := uint8(min(255, int(existing.B)+60))
					puzzle.Set(x, y, color.RGBA{r, g, b, 255})
				}
			}
		}
	}
}

// drawPuzzleHole 在背景图上绘制缺口
func drawPuzzleHole(bg *image.RGBA, targetX, targetY int, mask *image.Alpha) {
	bounds := mask.Bounds()

	// 用半透明深色填充缺口区域
	holeColor := color.RGBA{0, 0, 0, 100}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if mask.AlphaAt(x, y).A > 0 {
				bgX := targetX + x
				bgY := targetY + y
				if bgX >= 0 && bgX < ImageWidth && bgY >= 0 && bgY < ImageHeight {
					existing := bg.RGBAAt(bgX, bgY)
					// 混合颜色
					alpha := float64(holeColor.A) / 255.0
					r := uint8(float64(existing.R)*(1-alpha) + float64(holeColor.R)*alpha)
					g := uint8(float64(existing.G)*(1-alpha) + float64(holeColor.G)*alpha)
					b := uint8(float64(existing.B)*(1-alpha) + float64(holeColor.B)*alpha)
					bg.Set(bgX, bgY, color.RGBA{r, g, b, 255})
				}
			}
		}
	}

	// 添加缺口边框
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if mask.AlphaAt(x, y).A > 0 {
				// 检查是否是边缘像素
				isEdge := false
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
							if mask.AlphaAt(nx, ny).A == 0 {
								isEdge = true
								break
							}
						} else {
							isEdge = true
							break
						}
					}
					if isEdge {
						break
					}
				}

				if isEdge {
					bgX := targetX + x
					bgY := targetY + y
					if bgX >= 0 && bgX < ImageWidth && bgY >= 0 && bgY < ImageHeight {
						// 绘制深色边框
						bg.Set(bgX, bgY, color.RGBA{0, 0, 0, 180})
					}
				}
			}
		}
	}
}

// encodeImageToBase64 将图片编码为 Base64
func encodeImageToBase64(img image.Image) string {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// abs 返回绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min 返回最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 确保 math 包被使用
var _ = math.Pi

// IsCaptchaRequiredForOperation 检查指定操作是否需要验证码
func IsCaptchaRequiredForOperation(operation string) bool {
	if !setting.CaptchaEnabled {
		return false
	}

	switch operation {
	case "login":
		return setting.CaptchaRequireOnLogin
	case "register":
		return setting.CaptchaRequireOnRegister
	case "checkin":
		return setting.CaptchaRequireOnCheckin
	default:
		return false
	}
}
