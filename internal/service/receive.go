package service

import (
	"context"
	"errors"
	"qqbot/internal/model"
	"qqbot/internal/tool"

	"time"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	newlog "github.com/Lqingzhe/go-common-tool/log"
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
	s.scheduleTask.Do(s.sendCh)
}
func (s *Service) DistributeMessage(logger newlog.Logger, traceID string, oneBotResp *model.OneBotResp) {
	if oneBotResp.Echo != "" || oneBotResp.Status != "" {
		s.logCh <- logger.AddOperation("API响应").AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Service.TimeoutSec)*time.Second)
	ctx = context.WithValue(ctx, "trace_id", traceID)
	event := &oneBotResp.Event
	defer cancel()
	switch oneBotResp.PostType {
	case model.MessagePos:
		//消息分类放入对应的ch
		switch oneBotResp.MessageType {
		case model.PrivateMsg:
			//处理私聊消息
			s.operatorHub.Do(ctx, until.GetRoutKet(model.MessagePos.String(), model.PrivateMsg.String()), event, s.sendCh, logger)
		case model.GroupMsg:
			//处理群聊消息
			s.operatorHub.Do(ctx, until.GetRoutKet(model.MessagePos.String(), model.GroupMsg.String()), event, s.sendCh, logger)
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
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Poke.String()), event, s.sendCh, logger)
			case model.NotifyNt_Honor:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Honor.String()), event, s.sendCh, logger)
			case model.NotifyNt_LuckyKing:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_LuckyKing.String()), event, s.sendCh, logger)
			case model.NotifyNt_InputStatus:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_InputStatus.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.EssenceNt:
			switch oneBotResp.SubType {
			case model.EssenceNt_Add:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.EssenceNt.String(), model.EssenceNt_Add.String()), event, s.sendCh, logger)
			case model.EssenceNt_Delete:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.EssenceNt.String(), model.EssenceNt_Delete.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.ClientStatusNt:
			//处理napcat客户端问题，一般为离线等重要情况
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.ClientStatusNt.String()), event, s.sendCh, logger)
		case model.OfflineFileNt:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.OfflineFileNt.String()), event, s.sendCh, logger)
		case model.GroupCardNt:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupCardNt.String()), event, s.sendCh, logger)
		case model.FriendRecallNt:
			//好友消息撤回，可以删除数据库上下文
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.FriendRecallNt.String()), event, s.sendCh, logger)
		case model.GroupRecallNt:
			//群消息撤回，可以删除数据库上下文
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupRecallNt.String()), event, s.sendCh, logger)
		case model.FriendAddNt:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.FriendAddNt.String()), event, s.sendCh, logger)
		case model.GroupBanNt:
			switch oneBotResp.SubType {
			case model.GroupBan_Ban:
				//处理禁言逻辑
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_Ban.String()), event, s.sendCh, logger)
			case model.GroupBan_ListBan:
				//处理解除禁言逻辑
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_ListBan.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
				return
			}
		case model.GroupIncreaseNt:
			//欢迎新成员，只需判断不是自己，不用判断sub_type的GroupIncrease_Approve和GroupIncrease_Invite
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupIncreaseNt.String()), event, s.sendCh, logger)
		case model.GroupDecreaseNt:
			switch oneBotResp.SubType {
			case model.GroupDecrease_KickMe:
				//bot自己被提出群，清除该群数据
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_KickMe.String()), event, s.sendCh, logger)
			case model.GroupDecrease_Kick:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_Kick.String()), event, s.sendCh, logger)
			case model.GroupDecrease_Leave:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_Leave.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.GroupAdminNt:
			switch oneBotResp.SubType {
			case model.GroupAdmin_Set:
				//更新标记自己为管理员
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Set.String()), event, s.sendCh, logger)
			case model.GroupAdmin_Unset:
				//更新标记自己为成员
				s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Unset.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("NoticeType: "+oneBotResp.NoticeType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.GroupUploadNt:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.NoticePos.String(), model.GroupUploadNt.String()), event, s.sendCh, logger)
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知NoticeType: "+oneBotResp.NoticeType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	case model.RequestPos:
		switch oneBotResp.RequestType {
		case model.FriendReq:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.RequestPos.String(), model.FriendReq.String()), event, s.sendCh, logger)
		case model.GroupReq:
			switch oneBotResp.SubType {
			case model.GroupReq_Add:
				//管理员身份可尝试通过申请
				s.operatorHub.Do(ctx, until.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Add.String()), event, s.sendCh, logger)
			case model.GroupReq_Invite:
				//机器人收到邀请，可根据实际情况选择
				s.operatorHub.Do(ctx, until.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Invite.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("RequestType: "+oneBotResp.RequestType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
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
				s.operatorHub.Do(ctx, until.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Connect.String()), event, s.sendCh, logger)
			case model.LifeCycleME_Disable:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Disable.String()), event, s.sendCh, logger)
			case model.LifeCycleME_Enable:
				s.operatorHub.Do(ctx, until.GetRoutKet(model.MetaEventPos.String(), model.LifeCycleME.String(), model.LifeCycleME_Enable.String()), event, s.sendCh, logger)
			default:
				s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("MetaEventType: "+oneBotResp.MetaEventType.String()+" 下的未知SubType: "+oneBotResp.SubType.String()), newerror.CodeParamValueInvalid).AddTraceID(oneBotResp.Echo).AddLevel(newerror.LevelInfo)
				return
			}
		case model.HeartBeatME:
			s.operatorHub.Do(ctx, until.GetRoutKet(model.MetaEventPos.String(), model.HeartBeatME.String()), event, s.sendCh, logger)
		default:
			s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知MetaEventType: "+oneBotResp.MetaEventType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
			return
		}
	default:
		s.logCh <- logger.AddOperation("处理消息").AddError(errors.New("未知PostType: "+oneBotResp.PostType.String()), newerror.CodeParamValueInvalid).AddLevel(newerror.LevelError)
		return
	}
}
func (s *Service) GetLogCh() chan newlog.Logger {
	return s.logCh
}

//func (s *Service) GetErrCh() chan error {
//	return s.errCh
//}
