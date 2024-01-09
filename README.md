实现Restful协议

使用方式见test包下的engine_test

# 对外暴露结构:Engine
将Engine中路由部分抽象成router，同时为了满足用户注册/user/login,/user/logout这些分组路由需求，实际上路由的添加都是通过groupRouter来进行的

路由树的设计思路，首先是看分组。然后树的样子大概如下图

![image](https://github.com/STUDYKING-123/web-/assets/94734516/ae073158-e967-4de7-bac5-5fc105177a55)
