namespace go interact

//公共部分
struct Base{
    1:i32 code,
    2:string msg,
}

struct Item{
    1:string video_id,
    2:string user_id;
    3:string video_url,
    4:string cover_url,
    5:string title,
    6:string description,
    7:i32 visit_count,
    8:i32 like_count,
    9:i32 comment_count,
    10:string created_at,
    11:string updated_at,
    12:string deleted_at,
}

struct Data{
    1:Item item,
}

//互动部分
struct LikeActionReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3: string video_id(api.form="video_id", api.vd="len($)>0 || len(@.comment_id)>0"),
    4: string comment_id(api.form="comment_id"),
    5:string action_type(api.form="action_type", api.vd="$=='1' || $=='0'"),//1-点赞，0-取消点赞
}

struct LikeActionResp{
    1:Base base,
}

struct LikeListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string user_id(api.query="user_id"),
    4:i32 page_num(api.query="page_num", api.vd="$>=0"),
    5:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct LikeListResp{
    1:Base base,
    2:Data data,
}

struct CommentPublishReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string video_id(api.form="video_id", api.vd="len($)>0 || len(@.comment_id)>0"),
    4:string comment_id(api.form="comment_id"),
    5:string content(api.form="content", api.vd="len($)>0 && len($)<=200"),    
}

struct CommentPublishResp{
    1:Base base,
}

struct CommentListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string video_id(api.query="video_id", api.vd="len($)>0 || len(@.comment_id)>0"),
    4:string comment_id(api.query="comment_id"),
    5:i32 page_num(api.query="page_num", api.vd="$>=0"),
    6:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct CommentListResp{
    1:Base base,
    2:Data data,
}

struct CommentDeleteReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string video_id(api.query="video_id", api.vd="len($)>0 || len(@.comment_id)>0"),
    4:string comment_id(api.query="comment_id"),
}

struct CommentDeleteResp{
    1:Base base,
}

service InteractService{
    LikeActionResp LikeAction(LikeActionReq req)(api.post="/like/action")
    LikeListResp LikeList(LikeListReq req)(api.get="/like/list")
    CommentPublishResp CommentPublish(CommentPublishReq req)(api.post="/comment/publish")
    CommentListResp CommentList(CommentListReq req)(api.get="/comment/list")
    CommentDeleteResp CommentDelete(CommentDeleteReq req)(api.delete="/comment/delete")
}