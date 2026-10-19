package service

import (
	"context"
	"errors"
	"qqbot/internal/model"
	"qqbot/internal/tool"
	newerror "qqbot/pkg/error"
	newlog "qqbot/pkg/log"
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
	event := &oneBotResp.Event
	defer cancel()
	switch oneBotResp.PostType {
	case model.MessagePos:
		//消息分类放入对应的ch
		switch oneBotResp.MessageType {
		case model.PrivateMsg:
			//处理私聊消息
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MessagePos.String(), model.PrivateMsg.String()), event)
		case model.GroupMsg:
			//处理群聊消息
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MessagePos.String(), model.GroupMsg.String()), event)
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知MessageType:"+oneBotResp.MessageType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.NoticePos:
		switch oneBotResp.NoticeType {
		case model.NotifyNt:
			switch oneBotResp.SubType {
			case model.NotifyNt_Poke:
				// 处理戳一戳逻辑
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Poke.String()), event)
			case model.NotifyNt_Honor:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Honor.String()), event)
			case model.NotifyNt_LuckyKing:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_LuckyKing.String()), event)
			case model.NotifyNt_InputStatus:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_InputStatus.String()), event)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.EssenceNt:
			switch oneBotResp.SubType {
			case model.EssenceNt_Add:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.EssenceNt.String(), model.EssenceNt_Add.String()), event)
			case model.EssenceNt_Delete:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.EssenceNt.String(), model.EssenceNt_Delete.String()), event)
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.ClientStatusNt:
			//处理napcat客户端问题，一般为离线等重要情况
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.ClientStatusNt.String()), event)
		case model.OfflineFileNt:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.OfflineFileNt.String()), event)
		case model.GroupCardNt:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupCardNt.String()), event)
		case model.FriendRecallNt:
			//好友消息撤回，可以删除数据库上下文
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.FriendRecallNt.String()), event)
		case model.GroupRecallNt:
			//群消息撤回，可以删除数据库上下文
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupRecallNt.String()), event)
		case model.FriendAddNt:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.FriendAddNt.String()), event)
		case model.GroupBanNt:
			switch oneBotResp.SubType {
			case model.GroupBan_Ban:
				//处理禁言逻辑
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_Ban.String()), event)
			case model.GroupBan_ListBan:
				//处理解除禁言逻辑
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_ListBan.String()), event)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.GroupIncreaseNt:
			//欢迎新成员，只需判断不是自己，不用判断sub_type的GroupIncrease_Approve和GroupIncrease_Invite
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupIncreaseNt.String()), event)
		case model.GroupDecreaseNt:
			switch oneBotResp.SubType {
			case model.GroupDecrease_KickMe:
				//bot自己被提出群，清除该群数据
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_KickMe.String()), event)
			case model.GroupDecrease_Kick:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_Kick.String()), event)
			case model.GroupDecrease_Leave:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_Leave.String()), event)
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.GroupAdminNt:
			switch oneBotResp.SubType {
			case model.GroupAdmin_Set:
				//更新标记自己为管理员
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Set.String()), event)
			case model.GroupAdmin_Unset:
				//更新标记自己为成员
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupAdmin_Unset.String()), event)
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.GroupUploadNt:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.NoticePos.String(), model.GroupUploadNt.String()), event)
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知NoticeType: "+oneBotResp.NoticeType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.RequestPos:
		switch oneBotResp.RequestType {
		case model.FriendReq:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.RequestPos.String(), model.FriendReq.String()), event)
		case model.GroupReq:
			switch oneBotResp.SubType {
			case model.GroupReq_Add:
				//管理员身份可尝试通过申请
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Add.String()), event)
			case model.GroupReq_Invite:
				//机器人收到邀请，可根据实际情况选择
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Invite.String()), event)
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
		case model.LifeCycleME: //后续考虑是否需要处理
			switch oneBotResp.SubType {
			case model.LifeCycleME_Connect:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Connect.String()), event)
			case model.LifeCycleME_Disable:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Disable.String()), event)
			case model.LifeCycleME_Enable:
				err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Enable.String()), event)
			default:
				s.logCh <- logger.AddOperation("API响应").AddError(errors.New("MetaEventType: "+oneBotResp.MetaEventType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.HeartBeatME:
			err = s.operatorHub.do(ctx, tool.GetRoutKet(model.MetaEventPos.String(), model.HeartBeatME.String()), event)
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知MetaEventType: "+oneBotResp.MetaEventType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	default:
		s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知PostType: "+oneBotResp.PostType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
		return
	}
	if err != nil {
		err2 := newerror.TranslateError(err)
		s.logCh <- logger.AddOperation("处理消息").
			AddError(err2, err2.StatusCode).
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
