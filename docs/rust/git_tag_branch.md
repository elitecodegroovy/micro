创建tag
//创建轻量级tag：没有附带其他信息
git tag v1.0
//带信息的tag：-m注释信息.
git tag -a v1.0 -m 'first version'
共享tag
//我们在执行 git push 的时候，tag是不会上传到服务器的，在github网页上是看不到tag 的，为了共享这些tag:
//针对单个tag v1.0
git push origin v1.0
//将所有tag 一次全部push到github上。
git push origin --tags
删除tag(删除必须两个都要删除)
//删除本地tag
git tag -d v1.0
//删除github远端的指定tag
git push origin :refs/tags/v1.0
创建一个基于指定tag的分支
git checkout -b bratchname v1.0

Branch（分支）推送
1、查看branch分支

命令：git branch

2、创建branch分支

命令：git checkout -b branch分支名

3、将本地branch分支推送到远程branch分支上

命令：git push origin branch名

4、branch分支切换

分支切换命令：git checkout branch名

5、删除远程branch分支

命令：git push origin --delete branch名

6、删除本地branch分支

命令：git branch -d 分支名

7、在对应branch分支上提交相关内容

命令：git push --set-upstream origin branch名

8、将test1分支合并到master（主）分支上

首先要切换到master（主）分支上，然后在进行合并

合并命令：git merge 指定分支名

9、修改本地代码，然后在提交到远程仓库的步骤

（1）先添加到暂存区 git add *

（2）提交暂存区内容 git commit -m "自己注释"

（3）更新远程仓库内容 git push origin master（这一步如果有冲突，先解决冲突，然后在从（1）步开始操作）

（4）提交到远程仓库 git push