## 代码管理

一般而言，svn/git代码仓库中的每个项目都分为三个目录，分别是trunk、branches、tags。也可以只有trunk和tags或者trunk和branches，视代码管理情况而定。

下面也trunk和branches为例，说明一下如何提交和管理代码。

- trunk：这个目录下面就是项目的最新源代码，所有最新的源代码必须提交到这个目录下面。当测试没有任何问题后，新建一个分支branches（trunk的一个最新稳定版本）。
 
- branches：这个目录下的项目代码均是经过测试后的release版本，为了保证版本的统一性，必须加上项目的版本号。

此为，需要有以下规范：

- 提交代码注释。每次提交到trunk中的代码，需要描述为什么提交代码，特别是修改重要的核心代码，均需要说明修改了什么。
