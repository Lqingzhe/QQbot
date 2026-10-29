package until

import "qqbot/internal/model"

func MakeAtMessageSegment(atQQID int64) *model.MessageSegment {
	return &model.MessageSegment{
		Type: "at",
		Data: map[string]any{
			"qq": atQQID,
		},
	}
}
func MakeAllAtMessageSegment() *model.MessageSegment {
	return &model.MessageSegment{
		Type: "at",
		Data: map[string]any{
			"qq": "all",
		},
	}
}
func MakeReplyMessageSegment(messageID int64) *model.MessageSegment {
	return &model.MessageSegment{
		Type: "reply",
		Data: map[string]any{
			"id": messageID,
		},
	}
}
func MakeTextMessageSegment(message string) *model.MessageSegment {
	return &model.MessageSegment{
		Type: "text",
		Data: map[string]any{
			"text": message,
		},
	}
}
func MakeImageMessageSegment(url string) *model.MessageSegment {
	return &model.MessageSegment{
		Type: "image",
		Data: map[string]any{
			"file": url,
		},
	}
}

func MakePrivateMessage(goalUserID int64, message []*model.MessageSegment) model.Message {
	return model.Message{
		Action: "send_private_msg",
		Params: map[string]any{
			"user_id": goalUserID,
			"message": message,
		},
	}
}
func MakeGroupMessage(groupID int64, message *model.MessageSegment) model.Message {
	return model.Message{
		Action: "send_group_msg",
		Params: map[string]any{
			"group_id": groupID,
			"message":  message,
		},
	}
}

func MakeSetGroupBanMessage(goalUserID int64, groupID int64, banTimeSecond int64) model.Message {
	if banTimeSecond <= 0 {
		banTimeSecond = 1
	}
	return model.Message{
		Action: "set_group_ban",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"duration": banTimeSecond,
		},
	}
}
func MakeReleaseGroupBanMessage(goalUserID int64, groupID int64) model.Message {
	return model.Message{
		Action: "set_group_ban",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"duration": 0,
		},
	}
}
func MakeSetGroupWholeBan(groupID int64) model.Message {
	return model.Message{
		Action: "set_group_whole_ban",
		Params: map[string]any{
			"group_id": groupID,
			"enable":   true,
		},
	}
}
func MakeReleaseGroupWholeBan(groupID int64) model.Message {
	return model.Message{
		Action: "set_group_whole_ban",
		Params: map[string]any{
			"group_id": groupID,
			"enable":   false,
		},
	}
}
func MakeGroupKick(groupID int64, goalUserID int64, rejectAddRequest bool) model.Message {
	return model.Message{
		Action: "set_group_kick",
		Params: map[string]any{
			"group_id":           groupID,
			"user_id":            goalUserID,
			"reject_add_request": rejectAddRequest,
		},
	}
}
func MakeSetGroupAdmin(groupID int64, goalUserID int64) model.Message {
	return model.Message{
		Action: "set_group_admin",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"enable":   true,
		},
	}
}
func MakeRemoveGroupAdmin(groupID int64, goalUserID int64) model.Message {
	return model.Message{
		Action: "set_group_admin",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"enable":   false,
		},
	}
}
func MakeSetGroupCard(groupID int64, goalUserID int64, cardContent string) model.Message {
	return model.Message{
		Action: "set_group_card",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"card":     cardContent,
		},
	}
}

func MakeReleaseGroupCard(groupID int64, goalUserID int64) model.Message {
	return model.Message{
		Action: "set_group_card",
		Params: map[string]any{
			"group_id": groupID,
			"user_id":  goalUserID,
			"card":     "",
		},
	}
}
func MakeSetGroupName(groupID int64, groupName string) model.Message {
	return model.Message{
		Action: "set_group_name",
		Params: map[string]any{
			"group_id":   groupID,
			"group_name": groupName,
		},
	}
}
func MakeLeaveGroup(groupID int64) model.Message {
	return model.Message{
		Action: "set_group_leave",
		Params: map[string]any{
			"group_id":   groupID,
			"is_dismiss": false,
		},
	}
}
func MakeDisbandGroup(groupID int64) model.Message {
	return model.Message{
		Action: "set_group_leave",
		Params: map[string]any{
			"group_id":   groupID,
			"is_dismiss": true,
		},
	}
}
func MakeSetGroupSpecialTitle(groupID int64, goalUserID int64, title string, durationSec int64) model.Message {
	if durationSec <= 0 {
		durationSec = -1
	}
	return model.Message{
		Action: "set_group_special_title",
		Params: map[string]any{
			"group_id":      groupID,
			"user_id":       goalUserID,
			"special_title": title,
			"duration":      durationSec,
		},
	}
}
func MakeLeaveGroupSpecialTitle(groupID int64, goalUserID int64) model.Message {
	return model.Message{
		Action: "set_group_special_title",
		Params: map[string]any{
			"group_id":      groupID,
			"user_id":       goalUserID,
			"special_title": "",
		},
	}
}
func MakeSetGroupAddRequest(flag string, subType model.SubType, approve bool, reason string) model.Message {
	params := map[string]any{
		"flag":     flag,
		"sub_type": subType.String(),
		"approve":  approve,
	}
	if reason != "" {
		params["reason"] = reason
	}
	return model.Message{
		Action: "set_group_add_request",
		Params: params,
	}
}
func MakeFriendAddRequest(flag string, approve bool) model.Message {
	return model.Message{
		Action: "set_friend_add_request",
		Params: map[string]any{
			"flag":    flag,
			"approve": approve,
		},
	}
}
func MakeSetFriendRemark(goalUserID int64, remark string) model.Message {
	return model.Message{
		Action: "set_friend_remark",
		Params: map[string]any{
			"user_id": goalUserID,
			"remark":  remark,
		},
	}
}
func MakeDeleteMessage(messageID int64) model.Message {
	return model.Message{
		Action: "delete_message",
		Params: map[string]any{
			"message_id": messageID,
		},
	}
}
func MakeGetMessage(messageID int64) model.Message {
	return model.Message{
		Action: "get_message",
		Params: map[string]any{
			"message_id": messageID,
		},
	}
}
