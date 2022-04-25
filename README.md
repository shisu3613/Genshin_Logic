# Genshin_Logic_System
原本只想写家园模块的部分，但是这个部分涉及到的实在太多了，干脆就东拼西凑先从其他简单的逻辑写起
思路是先完成部分后端逻辑在本地测试通过后，然后搭载在修改后的zinx框架上（主要参考golang的[Zinx](https://github.com/aceld/zinx)框架,B站golang大海葵的系列[视频](https://space.bilibili.com/30214402/video),以及B站一棵平衡树的抽卡分析[文章](https://www.bilibili.com/read/cv14841352)等 感谢这些大佬的分享！)
目前完成大概60%,由于不是很想学客户端的东西所以测试用的client也是用go写的，目前大概整体长这样：

<img width="1250" alt="image" src="https://user-images.githubusercontent.com/48946918/165088716-948aab82-cd8a-4ea1-9ab4-5e410daef265.png">

项目中目前用到的数据库是mysql和redis,前者用于保存用户信息，后者用于缓存聊天信息

目前进度。更改聊天室模块中，打算改的和原神系统中的更像一点，顺便学习点新的技术。

目前思路思维导图版本：


![golang服务器框架](https://user-images.githubusercontent.com/48946918/165117360-ce90d2f2-7e02-4bb3-9b8f-f917e83346b2.svg)
