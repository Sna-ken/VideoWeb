namespace go video

//公共部分
struct Base{
    1:i32 code,
    2:string msg,
}

struct Data{
    1:Item item,
    2:i32 total,
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

//视频部分
struct PublishReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:binary video(api.form="video"),
    4:string title(api.form="title", api.vd="len($)>=1 && len($)<=40"),
    5:string description(api.form="description", api.vd="len($)<=200"),
}

struct PublishResp{
    1:Base base,
}

struct ListReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:i32 page_num(api.query="page_num", api.vd="$>=0"),
    4:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
    5:string user_id(api.query="user_id"),
}

struct ListResp{
    1:Base base,
    2:Data data,
}

struct PopularReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:i32 page_num(api.query="page_num", api.vd="$>=0"),
    4:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
}

struct PopularResp{
    1:Base base,
    2:Data data,
}

struct SearchReq{
    1:string access_token(api.header="access_token"),
    2:string refresh_token(api.header="refresh_token"),
    3:string keyword(api.form="keyword"),
    4:i32 page_num(api.form="page_num", api.vd="$>=0"),
    5:i32 page_size(api.form="page_size", api.vd="$>=1 && $<=100"),
    6:string username(api.form="username"),
}

struct SearchResp{
    1:Base base,
    2:Data data,
}

service VideoService{
    PublishResp Publish(PublishReq req)(api.post="/video/publish")
    ListResp List(ListReq req)(api.get="/video/list")
    PopularResp Popular(PopularReq req)(api.get="/video/popular")
    SearchResp Search(SearchReq req)(api.post="/video/search")
}