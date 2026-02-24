namespace go social

//公共部分
struct Base{
    1:i32 code,
    2:string msg,
}

struct Item{
    1:string user_id,
    2:string username,
    3:string avatar_url,
}

struct Data{
    1:Item item,
    2:i32 total,
}

//社交部分
struct FollowActionReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string to_user_id(api.form="to_user_id"),
    4:string action_type(api.form="action_type", api.vd="$=='1' || $=='0'"),//1-关注，0-取消关注
}

struct FollowActionResp{
    1:Base base,
}

struct FollowListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string user_id(api.query="user_id"),
    4:i32 page_num(api.query="page_num", api.vd="$>=0"),
    5:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct FollowListResp{
    1:Base base,
    2:Data data,
}

struct FollowerListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string user_id(api.query="user_id"),
    4:i32 page_num(api.query="page_num", api.vd="$>=0"),
    5:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct FollowerListResp{
    1:Base base,
    2:Data data,
}

struct FriendListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:i32 page_num(api.query="page_num", api.vd="$>=0"),
    4:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct FriendListResp{
    1:Base base,
    2:Data data,
}

service SocialService{
    FollowActionResp FollowAction(1:FollowActionReq req)(api.post="/follow/action")
    FollowListResp FollowList(1:FollowListReq req)(api.get="/follow/list")
    FollowerListResp FollowerList(1:FollowerListReq req)(api.get="/follower/list")
    FriendListResp FriendList(1:FriendListReq req)(api.get="/friend/list")
}