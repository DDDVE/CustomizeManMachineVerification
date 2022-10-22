package logic

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"pkg/conf"
	"pkg/crypto"
	"rpc/types/output"
	"strconv"
	"strings"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	FIRM_ID                        = "testID"
	QUESTION_TYPE_START            = 0
	QUESTION_TYPE_END              = 7
	DIGITAL_SIGNATURE_LENGTH_START = 156
	DIGITAL_SIGNATURE_LENGTH_END   = 160

	DIGITAL_SIGNATURE_CONNECTOR = "@==@"
	KEY_VERIFY_PREFIX           = "cache:output:verify:hash:"
	REQUEST_TIMEOUT_LENGTH      = 3 * 1000 * 1000 * 1000
)

type OutputLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOutputLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OutputLogic {
	return &OutputLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OutputLogic) CheckReq(req *types.OutputReq) error {

	// 参数校验
	if req.QuestionType < QUESTION_TYPE_START || req.QuestionType > QUESTION_TYPE_END {
		return errors.New(conf.GlobalError[conf.REQUEST_PARAM_ERROR] + ": questionType=" + strconv.Itoa(req.QuestionType))
	}
	if req.FirmID != FIRM_ID {
		return errors.New(conf.GlobalError[conf.REQUEST_PARAM_ERROR] + ": firmID=" + req.FirmID)
	}
	if len(req.DigitalSignature) < DIGITAL_SIGNATURE_LENGTH_START || len(req.DigitalSignature) > DIGITAL_SIGNATURE_LENGTH_END {
		return errors.New(conf.GlobalError[conf.REQUEST_PARAM_ERROR] + ": DigitalSignature=" + req.DigitalSignature)
	}

	// 验证请求是否超时
	stamp, err := strconv.ParseInt(req.CurrentTimestamp, 10, 64)
	if err != nil {
		return errors.New(conf.GlobalError[conf.REQUEST_PARAM_ERROR] + ": CurrentTimestamp=" + req.CurrentTimestamp)
	}
	if time.Since(time.Unix(0, stamp)) > REQUEST_TIMEOUT_LENGTH {
		return errors.New(conf.GlobalError[conf.REQUEST_TIMEOUT])
	}

	//验证数字签名
	msg := fmt.Sprintf("%d%v%v", req.QuestionType, req.FirmID, req.CurrentTimestamp)
	hash := sha256.Sum256([]byte(msg)) //生成摘要
	signs := strings.Split(req.DigitalSignature, DIGITAL_SIGNATURE_CONNECTOR)
	if len(signs) != 2 {
		return errors.New(conf.GlobalError[conf.SIGNATURE_CHECK_ERROR])
	}
	if err := crypto.VerifySignECC(hash[:], []byte(signs[0]), []byte(signs[1])); err != nil {
		return err
	}

	//验证请求是否重复发送 _目的是避免黑客截取请求进行攻击而导致服务崩溃及数据泄露，redis版太浪费资源,最好本地缓存
	pool := l.svcCtx.RedisPool
	conn := pool.Get()
	defer conn.Close()
	if r, err := conn.Do("EXISTS", KEY_VERIFY_PREFIX+string(hash[:])); r == 0 || err != nil {
		if r == 0 {
			return errors.New(conf.GlobalError[conf.REPEATED_REQUEST])
		} else {
			return err
		}
	}
	if _, err := conn.Do("SETEX", KEY_VERIFY_PREFIX+string(hash[:]), REQUEST_TIMEOUT_LENGTH, ""); err != nil {
		return err
	}
	return nil
}

func (l *OutputLogic) Output(req *types.OutputReq) (resp *types.OutputReply, err error) {

	//业务逻辑
	or, err := l.svcCtx.OutputRpc.GetItem(l.ctx, &output.OutputRequest{QuestionType: uint32(req.QuestionType)})
	if err != nil {
		return nil, err
	}

	return &types.OutputReply{
		Question:      or.Question,
		Answer:        or.Answer,
		DisturbAnswer: or.DisturbAnswer,
	}, nil
}
