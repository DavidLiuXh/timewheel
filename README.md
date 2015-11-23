# timewheel
一个通用的timewheel工具类


通用粗精度的timewheel,只启动一个timer,可监听任意多的到期时间，放入的用户数据可以是任何类型。

接口说明：

1. Start:
   开始一个timewheel ;
   
2. SetCallback：
   设置时间到期时的回调函数，回调函数中不可作过于耗时的操作，以免卡住timewheel的正常运行;
   
3. Add:
   添加需要监控的对象，务必在Start调用后再调用此接口;
   
4. Stop:
   停止此timewheel
