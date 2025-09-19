# generate-manga-feed

## Overview

以下のマンガサイトの作品一覧のRSSフィード（RSS2.0）を生成します。

- [くらげファーム](https://kuragebunch.com/farm)
- [コミックDAYS 新人賞](https://comic-days.com/newcomer)
- [コミックDAYS 読み切り](https://comic-days.com/oneshot)
- [&Sofa (アンドソファ)](https://andsofa.com)
- [トーチ](https://to-ti.in/product)
- [MATOGROSSO (マトグロッソ)](https://matogrosso.jp)
- [くらげバンチ 漫画賞](https://kuragebunch.com/info/award)
- [コミックエッセイ劇場](https://www.comic-essay.com/comics)
- [コミプレ 読切作品](https://viewer.heros-web.com/series/oneshot)
- [comicブースト 読み切り](https://comic-boost.com/genre/3)
- [ヤングアニマル 読み切り](https://younganimal.com/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A)
- [コミックエッセイ プチ大賞](https://www.comic-essay.com/contest/winner/)
- [コミックバンチKai 漫画賞](https://comicbunch-kai.com/article/award)
- [コミックバンチKai 読切作品](https://comicbunch-kai.com/series#oneshot)
- [アフタヌーン 四季賞](https://afternoon.kodansha.co.jp/award/)
- [少年マガジン 新人漫画大賞](https://debut.shonenmagazine.com/archive/#awards)
- [少年マガジン マガジンライズ](https://debut.shonenmagazine.com/archive/#magazinerise)
- [チャンピオンクロス 読み切り](https://championcross.jp/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A)
- [くらげバンチ 読切](https://kuragebunch.com/series/oneshot)
- [webアクション 読切作品](https://comic-action.com/series/oneshot)

### Note

定期実行などする場合は、`robots.txt`や`meta name="robots"`を確認して使用してください。

生成したRSSフィードの使用については、各サイトの利用規約や法律を遵守してください。

実行された時点のページ情報をRSSフィードとして出力します。フィードアイテムの更新には対応していません。

アクセスが多くなるサイトにはsleep処理を入れているため、RSS生成に時間のかかるサイトがあります。

## Install

実行バイナリをパスの通ったディレクトリに配置するなどして使用してください。

## Usage

サブコマンドを指定することで、標準出力に対象サイトのRSSフィードを出力します。

```bash
$ generate-manga-feed kuragefarm
<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>～～～
```

ファイルに出力する場合は、リダイレクトしてください。

```bash
$ generate-manga-feed kuragefarm > kuragefarm.rss

# エラーログを分ける場合
$ generate-manga-feed kuragefarm > kuragefarm.rss 2> kuragefarm_err.log
```

サブコマンド一覧はヘルプを参照してください。

```bash
$ generate-manga-feed help
Usage: generate-manga-feed [subcommand]

Available subcommands:
  version               Show application version
  help                  Show this help message
  --------
  afternoonaward        Print afternoon award RSS feed
  andsofa               Print andsofa RSS feed
  championcrossoneshot  Print champion cross oneshot RSS feed
  comicactiononeshot    Print comic-acticon oneshot RSS feed
  comicboostoneshot     Print comic-boost oneshot RSS feed
  comicbunckaiaward     Print comicbunchi-kai award RSS feed
  comicbunckaioneshot   Print comicbunchi-kai oneshot RSS feed
  comicdaysnewcomer     Print comic-days newcomer RSS feed
  comicdaysoneshot      Print comic-days oneshot RSS feed
  comicessaycontest     Print comic-essay contest RSS feed
  comicessaygekijo      Print comic-essay gekijo RSS feed
  comiplexoneshot       Print comiplex oneshot RSS feed
  kuragebunchaward      Print kuragebunch award RSS feed
  kuragebunchoneshot    Print kuragebunch oneshot RSS feed
  kuragefarm            Print kurage farm RSS feed
  matogrosso            Print matogrosso RSS feed
  shonenmagazineaward   Print shonen magazine award RSS feed
  shonenmagazinerise    Print shonen magazine rise RSS feed
  toti                  Print toti RSS feed
  younganimaloneshot    Print younganimal oneshot RSS feed
```

対象サイトごとのサブコマンドは以下の通りです。

| 対象サイト                                                                                                       | サブコマンド         |
| ---------------------------------------------------------------------------------------------------------------- | -------------------- |
| [くらげファーム](https://kuragebunch.com/farm)                                                                   | kuragefarm           |
| [コミックDAYS 新人賞](https://comic-days.com/newcomer)                                                           | comicdaysnewcomer    |
| [コミックDAYS 読み切り](https://comic-days.com/oneshot)                                                          | comicdaysoneshot     |
| [&Sofa（アンドソファ）](https://andsofa.com)                                                                     | andsofa              |
| [トーチ](https://to-ti.in/product)                                                                               | toti                 |
| [MATOGROSSO（マトグロッソ）](https://matogrosso.jp)                                                              | matogrosso           |
| [くらげバンチ 漫画賞](https://kuragebunch.com/info/award)                                                        | kuragebunchaward     |
| [コミックエッセイ劇場](https://www.comic-essay.com/comics)                                                       | comicessaygekijo     |
| [コミプレ 読切作品](https://viewer.heros-web.com/series/oneshot)                                                 | comiplexoneshot      |
| [comicブースト 読み切り](https://comic-boost.com/genre/3)                                                        | comicboostoneshot    |
| [ヤングアニマル 読み切り](https://younganimal.com/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A)      | younganimaloneshot   |
| [コミックエッセイ プチ大賞](https://www.comic-essay.com/contest/winner/)                                         | comicessaycontest    |
| [コミックバンチKai 漫画賞](https://comicbunch-kai.com/article/award)                                             | comicbunckaiaward    |
| [コミックバンチKai 読切作品](https://comicbunch-kai.com/series#oneshot)                                          | comicbunckaioneshot  |
| [アフタヌーン 四季賞](https://afternoon.kodansha.co.jp/award/)                                                   | afternoonaward       |
| [少年マガジン 新人漫画大賞](https://debut.shonenmagazine.com/archive/#awards)                                    | shonenmagazineaward  |
| [少年マガジン マガジンライズ](https://debut.shonenmagazine.com/archive/#magazinerise)                            | shonenmagazinerise   |
| [チャンピオンクロス 読み切り](https://championcross.jp/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A) | championcrossoneshot |
| [くらげバンチ 読切](https://kuragebunch.com/series/oneshot)                                                      | kuragebunchoneshot   |
| [webアクション 読切作品](https://comic-action.com/series/oneshot)                                                | comicactiononeshot   |
