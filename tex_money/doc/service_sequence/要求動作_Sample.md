# FIT-A
## 閲覧情報
VSCodeにて、拡張機能[Markdown Preview Enhanced]をインストールすると
シーケンスが確認できます。
記述方法については、マークダウン マーメイドで検索すると出てきます。

## TexMoney
要求に対する他サービスへの呼び出しシーケンスを記載する


### 要求動作名:【イニシャル動作】
    説明

## 通信シーケンス

<style>.mermaid svg {height:100%}</style>
```mermaid
sequenceDiagram
participant ui as ui
participant 1-2 as 1-2
participant 1-3 as 1-3
participant 2-1 as 2-1
participant 2-4 as 2-4
participant 1-1 as 1-1
participant pr as 印刷制御


1-3 ->> 2-1 : request_status
2-1 -->> 1-3 : result_status
1-3 ->> 2-4 : request_get_terminfo_now
2-4 -->> 1-3 : request_get_terminfo_now
1-3 ->> 1-1 : request_status
1-1 -->> 1-3 : result_status

Note over ui : SAMPLE

Note over ui : right of ●, left of ●, over ● どの線の上にメモを置くか
ui ->> ui : ループ

loop ループ条件
    ui ->> 1-3 : くるくる
end

alt is 条件分岐 result true
    1-3 ->> ui : 通知
else is result false
    1-3 ->> ui : 破棄
end 

opt ついでにオプションフラグ
 1-3 ->> ui : 通知
 end



