package group

import (
	"net/http"

	"GoChat/apps/social/api/internal/logic/group"
	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 群成员列表
func GroupUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupUserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupUserListLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
