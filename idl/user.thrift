namespace go user

//公共部分
struct Base{
    1:i32 code,
    2:string msg,
}

struct Data{
    1:string user_id,
    2:string username,
    3:string avatar_url,
    4:string created_at,
    5:string updated_at,
    6:string deleted_at,
}

//用户部分
struct RegisterReq{
    1:string username(api.form="username", api.vd="len($)>=1 && len($)<=20"),
    2:string password(api.form="password", api.vd="len($)>=6"),//密码至少6位
}

struct RegisterResp{
   1:Base base,
}

struct LoginReq{
    1:string username(api.form="username"),
    2:string password(api.form="password"),
}

struct LoginResp{
    1:Base base,
    2:Data data,
    3:string access_token(api.form="access_token"),
    4:string refresh_token(api.form="refresh_token"),
}

struct UserInfoReq{
    1:string user_id(api.query="user_id"),
    2:string access_token(api.header="access_token"),
    3:string refresh_token(api.header="refresh_token"),
}

struct UserInfoResp{
    1:Base base,
    2:Data data
}

struct UploadAvatarReq{
   1:string access_token(api.header="access_token"),
   2:string refresh_token(api.header="refresh_token"),
   3:binary avatar(api.form="avatar"),
}

struct UploadAvatarResp{
   1:Base base,
   2:Data data,
}

service UserService{
    RegisterResp Register(1:RegisterReq req)(api.post="/user/register")
    LoginResp Login(1:LoginReq req)(api.post="/user/login")
    UserInfoResp UserInfo(1:UserInfoReq req)(api.get="/user/info")
    UploadAvatarResp UploadAvatar(1:UploadAvatarReq req)(api.put="/user/avatar/upload" )
}