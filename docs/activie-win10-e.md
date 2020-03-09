

## 激活密钥

相关名词解释：

N版：欧盟版，无MediaPlayer

S版：长期服务支持分支LTSB，无Edge浏览器

G版：政府版

GVLK：批量授权许可密钥，用于KMS激活

OEM:NONSLP 非系统锁定预装密钥，跟零售版激活机制一样，联网激活

MAK：大客户版本批量激活密钥，永久激活，有次数限制


一、Win10企业版安装密钥：

Enterprise OEM:NONSLP：XGVPP-NMH47-7TTHJ-W3FW7-8HV2C

企业版（欧盟版，无mediaplayer）Volume:GVLK

EnterpriseN：WGGHN-J84D6-QYCPR-T7PJ7-X766F

企业版（长期服务支持分支2015，无edge浏览器）

EnterpriseS OEM:NONSLP：FWN7H-PF93Q-4GGP8-M8RF3-MDWWW

企业版（欧盟版，无mediaplayer，无edge浏览器）Volume:GVLK

EnterpriseSN：X4R4B-NV6WD-PKTVK-F98BH-4C2J8


二、Windows 10企业版永久激活密钥批量授权版Volume:MAK

Windows 10 Enterprise MAK激活密钥

[Key]：W77WN-TP36G-VBJFX-C77HK-2YT44

[Key]：8QPWX-N9JB8-74W2J-4W2R4-FC2JR

[Key]：9DKNB-8HTHW-V7TBP-77K4V-G6PJR

[Key]：9XN4J-QPX7Q-BYHHV-WWH96-9KXR4

[Key]：BVNMV-QQXW7-K2J8P-TP8HK-66PJR

[Key]：CN6B3-77CYC-X8M77-XCP4G-3V644


Windows 10 EnterpriseS 2016 MAK激活密钥(LTSB )：

[Key]：PGDVN-GKX6K-MGJDR-CRVDV-HXPJQ

[Key]：G36GR-R6N89-WCVTR-T8YKH-GCRCD


Windows 10 EnterpriseSN 2016 MAK激活密钥(LTSB N)：

[Key]：96R6X-N2B92-M4WRQ-HP4KG-3V72F

[Key]：QCDT7-JNX2M-8FBPB-J4VYW-7FT44

三、Wins10企业版VOL版KMS客户端安装密钥GVLK：

Windows 10 企业版：NPPR9-FWDCX-D2C8J-H872K-2YT43

Windows 10 企业版 N：DPH2V-TTNVB-4X9Q3-TJR4H-KHJW4

Windows 10 企业版 G：YYVX9-NTFWV-6MDM3-9PT4T-4M68B

Windows 10 企业版 G N：44RPN-FTY23-9VTTB-MP9BX-T84FV


Windows10 企业版长期服务支持分支版本：

Windows 10 LTSC 2019

Windows 10 企业版 LTSC 2019：M7XTQ-FN8P6-TTKYV-9D4CC-J462D

Windows 10 企业版 N LTSC 2019：2NFX-8DJQP-P6BBQ-THF9C-7CG2H


Windows 10 LTSB 2016

Windows 10 企业版 LTSB 2016：DCPHK-NFMTC-H88MJ-PFHPY-QJ4BJ

Windows 10 企业版 N LTSB 2016：FFDN-GRT3P-VKWWX-X7T3R-8B639


Windows 10 LTSB 2015

Windows 10 企业版 2015 长期服务：WNMTR-4C88C-JK8YV-HQ7T2-76DF9

Windows 10 企业版 2015 长期服务 N：2F77B-TNFGY-69QQF-B8YKP-D69TJ



## 激活步骤

待打开MSDOS界面后，依次输入以下命令：
```
slmgr.vbs /upk
```
此时弹出窗口显未“已成功卸载了产品密钥”。

接着输入以下命令：
```
slmgr /ipk NPPR9-FWDCX-D2C8J-H872K-2YT43
```
弹出窗口提示：“成功的安装了产品密钥”。

输入以下命令：
```
slmgr /ato
```
按回车键后将弹出窗口提示：“成功的激活了产品”。


也可以设置计算机名称，继续输入以下命令：
```
slmgr /skms zh.us.to
```
弹出窗口提示：“密钥管理服务计算机名成功的设置为zh.us.to”。