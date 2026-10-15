package service

import (
	"context"
	"demo/internal/model"
	newerror "demo/pkg/error"
	newlog "demo/pkg/log"
	"errors"
	"time"
)



func (s *Service) Run(config *model.Config) {
	s.logCh = make(chan newlog.Logger, config.Service.ChannelLength)
	l := config.ResponserList
	p := make([]int64, 0, len(l))
	g := make([]int64, 0, len(l))
	for _, info := range l {
		if info.IsPrivate {
			p = append(p, info.QQid)
		} else {
			g = append(g, info.QQid)
		}
	}
}
func (s *Service) DistributeMessage(logger newlog.Logger, traceID string, oneBotResp *model.OneBotResp) {
	if oneBotResp.Echo != "" || oneBotResp.Status != "" {
		s.logCh <- logger.AddOperation("API响应").AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
		return
	}
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Service.TimeoutSec)*time.Second)
	context.WithValue(ctx, "trace_id", traceID)
	defer cancel()
	switch oneBotResp.PostType {
	case model.MessagePos:
		//消息分类放入对应的ch
		switch oneBotResp.MessageType {
		case model.PrivateMsg:处理私聊消息
		case model.GroupMsg:处理群聊消息
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知MessageType:"+oneBotResp.MessageType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.NoticePos:
		switch oneBotResp.NoticeType {
		case model.NotifyNt:
			switch oneBotResp.SubType {
			case model.NotifyNt_Poke:处理戳一戳逻辑
			case model.NotifyNt_Honor:
			case model.NotifyNt_LuckyKing:
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.EssenceNt:
			switch oneBotResp.SubType {
			case model.EssenceNt_Add:
			case model.EssenceNt_Delete:
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.ClientStatusNt:处理napcat客户端问题，一般为离线等重要情况
		case model.OfflineFileNt:
		case model.GroupCardNt:
		case model.FriendRecallNt:好友消息撤回，可以删除数据库上下文
		case model.GroupRecallNt:群消息撤回，可以删除数据库上下文
		case model.FriendAddNt:
		case model.GroupBanNt:
			switch oneBotResp.SubType {
			case model.GroupBan_Ban:处理禁言逻辑
			case model.GroupBan_ListBan:处理解除禁言逻辑
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.GroupIncreaseNt:欢迎新成员，只需判断不是自己//不用判断sub_type的GroupIncrease_Approve和GroupIncrease_Invite
		case model.GroupDecreaseNt:
			switch oneBotResp.SubType {
			case model.GroupDecrease_KickMe:bot自己被提出群，清除该群数据
			case model.GroupDecrease_Kick:
			case model.GroupDecrease_Leave:
			default:
					s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
					return
			}
		case model.GroupAdminNt:
			switch oneBotResp.SubType {
			case model.GroupAdmin_Set:更新标记自己为管理员
			case model.GroupAdmin_Unset:更新标记自己为成员
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.GroupUploadNt:
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知NoticeType: "+oneBotResp.NoticeType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.RequestPos:
		switch oneBotResp.RequestType {
		case model.FriendReq:
		case model.GroupReq:
			switch oneBotResp.SubType {
			case model.GroupReq_Add:管理员身份可尝试通过申请
			case model.GroupReq_Invite:机器人收到邀请，可根据实际情况选择
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("RequestType: "+oneBotResp.RequestType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知RequestType: "+oneBotResp.RequestType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.MetaEventPos:
		switch oneBotResp.MetaEventType {
		case model.LifeCycleME://后续考虑是否需要处理
			switch oneBotResp.SubType {
			case model.LifeCycleME_Connect:
			case model.LifeCycleME_Disable:
			case model.LifeCycleME_Enable:
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("MetaEventType: "+oneBotResp.MetaEventType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.HeartBeatME:
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知MetaEventType: "+oneBotResp.MetaEventType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	default:
		s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知PostType: "+oneBotResp.PostType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
		return
	}
	if err!=nil{
		err2:=newerror.TranslateError(err)
		s.logCh<-logger.AddOperation("处理消息").
			AddError(err2,err2.StatusCode).
			AddLevel(err2.LogLevel)
		return
	}
}
func (s *Service) GetLogCh() chan newlog.Logger {
	return s.logCh
}

//func (s *Service) GetErrCh() chan error {
//	return s.errCh
//}
