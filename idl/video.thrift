namespace go video

//公共部分
struct Base{
    1:i32 code,
    2:string msg,
}

struct Data{
    1:list<Item> item,
    2:i32 total,
}

struct Item{
    1:string username
    2:string video_id,
    3:string user_id,
    4:string video_url,
    5:string cover_url,
    6:string title,
    7:string description,
    8:i32 visit_count,
    9:i32 like_count,
    10:i32 comment_count,
    11:string created_at,
    12:string updated_at,
    13:string deleted_at,
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
    3:string keyword(api.query="keyword"),
    4:i32 page_num(api.query="page_num", api.vd="$>=0"),
    5:i32 page_size(api.query="page_size", api.vd="$>=1 && $<=100"),
    6:string username(api.query="username"),
}

struct SearchResp{
    1:Base base,
    2:Data data,
}

service VideoService{
    PublishResp Publish(PublishReq req)(api.post="/video/publish")
    ListResp List(ListReq req)(api.get="/video/list")
    PopularResp Popular(PopularReq req)(api.get="/video/popular")
    SearchResp Search(SearchReq req)(api.get="/video/search")
}