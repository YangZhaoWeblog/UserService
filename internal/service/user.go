package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	v1 "github.com/YangZhaoWeblog/UserService/api/user/v1"
	"github.com/YangZhaoWeblog/UserService/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
)

// UserService 是用户服务
type UserService struct {
	v1.UnimplementedUserServer

	uc        *biz.UserUsecase
	logHelper *takin_log.TakinLogger
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase, log *takin_log.TakinLogger) *UserService {
	return &UserService{uc: uc, logHelper: log}
}

// Register 实现注册接口
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	// 创建命名的 span，这会使得追踪链更加丰富
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.Register")
	defer span.End()

	// 添加一些属性到 span 中
	span.SetAttributes(
		attribute.String("user.phone", req.GetPhone().GetPhoneNumber()),
		attribute.String("user.nickname", req.GetNickname()),
	)

	// 记录日志，包含追踪信息
	s.logHelper.InfoContext(ctx, "开始处理用户注册请求", "phone", req.GetPhone().GetPhoneNumber())

	// 检查手机号是否已注册 (模拟调用)
	if err := s.checkPhoneExists(ctx, req.GetPhone().GetPhoneNumber()); err != nil {
		// 设置 span 状态为错误
		span.SetStatus(codes.Error, "手机号已存在")
		span.RecordError(err)
		s.logHelper.ErrorContext(ctx, "手机号已存在", "error", err.Error())
		return &v1.RegisterReply{
			Success: false,
			Message: "手机号已存在",
		}, err
	}

	// 验证验证码 (模拟调用)
	if err := s.verifyCode(ctx, req.GetPhone().GetPhoneNumber(), req.GetPhone().GetVerificationCode()); err != nil {
		span.SetStatus(codes.Error, "验证码无效")
		span.RecordError(err)
		s.logHelper.ErrorContext(ctx, "验证码验证失败", "error", err.Error())
		return &v1.RegisterReply{
			Success: false,
			Message: "验证码无效",
		}, err
	}

	// 创建用户 (模拟调用)
	user, err := s.createUserProfile(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, "创建用户失败")
		span.RecordError(err)
		s.logHelper.ErrorContext(ctx, "创建用户失败", "error", err.Error())
		return &v1.RegisterReply{
			Success: false,
			Message: "创建用户失败",
		}, err
	}

	// 发送欢迎通知 (模拟调用)
	s.sendWelcomeNotification(ctx, user.ID)

	span.SetStatus(codes.Ok, "注册成功")
	s.logHelper.InfoContext(ctx, "用户注册成功", "user_id", user.ID)

	return &v1.RegisterReply{
		Success: true,
		Message: "注册成功",
		UserInfo: &v1.UserInfo{
			UserId:    fmt.Sprintf("%d", user.ID),
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarURL,
		},
	}, nil
}

// Login 实现登录接口
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	// TODO: 实现登录逻辑

	//s.logHelper.Info("测试:", "Login", "是否", "成功")
	return nil, errors.New(521, "LOGIN_FAILED", "登录失败")

	// 不可达代码
	//return &v1.LoginReply{
	//	Success: true,
	//	Message: "登录成功",
	//}, nil
}

// Info 实现获取用户信息接口
func (s *UserService) Info(ctx context.Context, req *v1.InfoRequest) (*v1.InfoReply, error) {
	// TODO: 实现获取用户信息逻辑
	return &v1.InfoReply{}, nil
}

// ----- 以下是追踪链中使用的模拟方法 -----

// checkPhoneExists 模拟检查手机号是否已存在
func (s *UserService) checkPhoneExists(ctx context.Context, phone string) error {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.checkPhoneExists")
	defer span.End()

	span.SetAttributes(attribute.String("phone", phone))

	// 模拟调用数据库
	s.simulateDBCall(ctx, "查询手机号", 50, 150)

	// 随机模拟手机号已存在的情况 (10%的概率)
	if rand.Intn(10) == 0 {
		err := fmt.Errorf("手机号 %s 已被注册", phone)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// verifyCode 模拟验证码验证过程
func (s *UserService) verifyCode(ctx context.Context, phone, code string) error {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.verifyCode")
	defer span.End()

	span.SetAttributes(
		attribute.String("phone", phone),
		attribute.String("code_length", fmt.Sprintf("%d", len(code))),
	)

	// 模拟调用短信服务
	if err := s.simulateExternalServiceCall(ctx, "短信验证服务"); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "短信服务调用失败")
		return fmt.Errorf("验证码验证服务不可用: %w", err)
	}

	// 随机模拟验证码无效的情况 (5%的概率)
	if rand.Intn(20) == 0 {
		err := fmt.Errorf("验证码 %s 无效或已过期", code)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	s.logHelper.InfoContext(ctx, "验证码验证通过")
	return nil
}

// createUserProfile 模拟创建用户资料流程
func (s *UserService) createUserProfile(ctx context.Context, req *v1.RegisterRequest) (*biz.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.createUserProfile")
	defer span.End()

	// 模拟创建用户基本信息
	if err := s.simulateDBCall(ctx, "插入用户记录", 100, 200); err != nil {
		return nil, err
	}

	// 模拟生成头像缩略图
	ctx, thumbSpan := otel.Tracer("user-service").Start(ctx, "生成头像缩略图")
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
	thumbSpan.End()

	// 模拟用户偏好初始化
	ctx, prefSpan := otel.Tracer("user-service").Start(ctx, "初始化用户偏好")
	time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)
	prefSpan.End()

	// 创建虚拟用户对象
	user := &biz.User{
		ID:        int64(100000 + rand.Intn(900000)),
		Username:  fmt.Sprintf("user_%s", req.GetPhone().GetPhoneNumber()),
		Nickname:  req.GetNickname(),
		AvatarURL: req.GetAvatarUrl(),
		Phone:     req.GetPhone().GetPhoneNumber(),
	}

	span.SetAttributes(attribute.Int64("user.id", user.ID))
	return user, nil
}

// sendWelcomeNotification 模拟发送欢迎通知
func (s *UserService) sendWelcomeNotification(ctx context.Context, userID int64) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.sendWelcomeNotification")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", userID))

	// 模拟发送电子邮件
	_, emailSpan := otel.Tracer("notification-service").Start(ctx, "发送欢迎邮件")
	time.Sleep(time.Duration(100+rand.Intn(150)) * time.Millisecond)
	emailSpan.End()

	// 模拟推送通知
	pushCtx, pushSpan := otel.Tracer("notification-service").Start(ctx, "发送推送通知")
	if err := s.simulateExternalServiceCall(pushCtx, "推送服务"); err != nil {
		pushSpan.RecordError(err)
		pushSpan.SetStatus(codes.Error, "推送服务调用失败")
		s.logHelper.WarnContext(ctx, "推送通知发送失败", "error", err.Error())
	}
	pushSpan.End()

	s.logHelper.InfoContext(ctx, "欢迎通知发送完成", "user_id", userID)
}

// simulateDBCall 模拟数据库调用
func (s *UserService) simulateDBCall(ctx context.Context, operation string, minMs, maxMs int) error {
	ctx, span := otel.Tracer("database").Start(ctx, operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.operation", operation),
	)

	// 模拟数据库响应时间
	time.Sleep(time.Duration(minMs+rand.Intn(maxMs-minMs)) * time.Millisecond)

	// 随机模拟数据库错误 (1%的概率)
	if rand.Intn(100) == 0 {
		err := fmt.Errorf("数据库操作失败: %s", operation)
		span.RecordError(err)
		span.SetStatus(codes.Error, "数据库错误")
		return err
	}

	return nil
}

// simulateExternalServiceCall 模拟外部服务调用
func (s *UserService) simulateExternalServiceCall(ctx context.Context, serviceName string) error {
	// OpenTelemetry SpanKind.CLIENT 是服务图生成的关键
	ctx, span := otel.Tracer("user-service").Start(ctx, "call:"+serviceName,
		trace.WithSpanKind(trace.SpanKindClient)) // 使用客户端类型的Span

	defer span.End()

	// 添加必要属性来确保服务图生成
	span.SetAttributes(
		attribute.String("service.name", "user-service"),                       // 来源服务名称
		attribute.String("peer.service", serviceName),                          // 目标服务名称-关键属性
		attribute.String("net.peer.name", serviceName+".example.com"),          // 目标主机
		attribute.String("http.url", "http://"+serviceName+".example.com/api"), // 请求URL
		attribute.String("http.method", "POST"),                                // 使用的方法
		attribute.String("rpc.system", "http"),                                 // RPC系统
	)

	// 模拟服务端响应
	time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

	// 为了服务图，创建一个服务端的span
	serverCtx, serverSpan := otel.Tracer(serviceName).Start(ctx, serviceName+".handleRequest",
		trace.WithSpanKind(trace.SpanKindServer)) // 使用服务端类型的Span - 服务图的另一半

	serverSpan.SetAttributes(
		attribute.String("service.name", serviceName), // 服务名
		attribute.String("http.method", "POST"),
		attribute.Int("http.status_code", 200),
	)

	// 模拟服务端处理
	time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

	// 模拟服务端数据库调用
	_ = trace.SpanFromContext(serverCtx) // 仅用于展示获取当前span的方法
	_, childSpan := otel.Tracer(serviceName).Start(serverCtx, "database.query")
	childSpan.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.name", serviceName+"_db"),
		attribute.String("db.statement", "SELECT * FROM users"),
	)
	time.Sleep(time.Duration(rand.Intn(30)+10) * time.Millisecond)
	childSpan.End()

	// 随机模拟错误
	if rand.Intn(10) == 0 {
		err := fmt.Errorf("服务调用失败: %s", serviceName)
		serverSpan.RecordError(err)
		serverSpan.SetStatus(codes.Error, "服务错误")
		span.RecordError(err) // 客户端也记录错误
		span.SetStatus(codes.Error, "远程服务错误")
		serverSpan.End()
		return err
	}

	serverSpan.End()
	return nil
}

// TestServiceGraph 测试服务图功能
func (s *UserService) TestServiceGraph(ctx context.Context, req *v1.TestServiceGraphRequest) (*v1.TestServiceGraphReply, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.TestServiceGraph",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		attribute.String("service.name", "user-service"),
		attribute.String("http.method", "POST"),
		attribute.Int("http.status_code", 200),
	)

	// 模拟调用外部服务A
	if err := s.callExternalService(ctx, "service-a"); err != nil {
		span.RecordError(err)
		return &v1.TestServiceGraphReply{Success: false, Message: err.Error()}, err
	}

	// 模拟调用外部服务B
	if err := s.callExternalService(ctx, "service-b"); err != nil {
		span.RecordError(err)
		return &v1.TestServiceGraphReply{Success: false, Message: err.Error()}, err
	}

	// 数据库操作
	s.simulateDBCall(ctx, "查询用户数据", 50, 100)

	return &v1.TestServiceGraphReply{
		Success: true,
		Message: "Service graph test completed",
	}, nil
}

// callExternalService 模拟调用外部服务
func (s *UserService) callExternalService(ctx context.Context, serviceName string) error {
	ctx, span := otel.Tracer("user-service").Start(ctx, "callExternalService."+serviceName,
		trace.WithSpanKind(trace.SpanKindClient))

	span.SetAttributes(
		attribute.String("service.name", "user-service"),
		attribute.String("peer.service", serviceName),
		attribute.String("http.url", fmt.Sprintf("http://%s.example.com/api", serviceName)),
		attribute.String("http.method", "POST"),
	)

	defer span.End()

	// 模拟服务端响应
	time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

	// 创建服务端span (在真实场景中，这是由被调用服务创建的)
	serverCtx, serverSpan := otel.Tracer(serviceName).Start(ctx, serviceName+".handleRequest",
		trace.WithSpanKind(trace.SpanKindServer))

	serverSpan.SetAttributes(
		attribute.String("service.name", serviceName),
		attribute.String("http.method", "POST"),
		attribute.Int("http.status_code", 200),
	)

	// 模拟服务端处理
	time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

	// 模拟服务端数据库调用
	_ = trace.SpanFromContext(serverCtx) // 仅用于展示获取当前span的方法
	_, childSpan := otel.Tracer(serviceName).Start(serverCtx, "database.query")
	childSpan.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.name", serviceName+"_db"),
		attribute.String("db.statement", "SELECT * FROM users"),
	)
	time.Sleep(time.Duration(rand.Intn(30)+10) * time.Millisecond)
	childSpan.End()

	// 随机模拟错误
	if rand.Intn(10) == 0 {
		err := fmt.Errorf("服务调用失败: %s", serviceName)
		serverSpan.RecordError(err)
		serverSpan.SetStatus(codes.Error, "服务错误")
		return err
	}

	serverSpan.End()

	return nil
}
