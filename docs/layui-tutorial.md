
## LayUi组件编译

  使用了阿里巴巴矢量图标库（iconfont），构建工具采用 Gulp 。
  
  Layui部分模块依赖jQuery（比如layer），但是你并不用去额外加载jQuery。Layui已经将jQuery最稳定的一个版本改为Layui的内部模块，当你去使用 layer 之类的模块时，它会首先判断你的页面是否已经引入了jQuery，如果没有，则加载内部的jQuery模块，如果有，则不会加载。