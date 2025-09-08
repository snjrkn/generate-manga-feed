# generate-manga-feed

## Overview

以下のマンガサイトの作品一覧のRSSフィード（RSS2.0）を生成します。

- [くらげファーム](https://kuragebunch.com/farm)
- [コミックDAYS 新人賞](https://comic-days.com/newcomer)
- [コミックDAYS 読み切り](https://comic-days.com/oneshot)
- [&Sofa (アンドソファ)](https://andsofa.com)
- [トーチ](https://to-ti.in/product)
- [MATOGROSSO (マトグロッソ)](https://matogrosso.jp)
- [くらげバンチ漫画賞](https://kuragebunch.com/info/award)
- [コミックエッセイ劇場](https://www.comic-essay.com/comics)
- [コミプレ 読切作品](https://viewer.heros-web.com/series/oneshot)
- [comicブースト 読み切り](https://comic-boost.com/genre/3)
- [ヤングアニマル 読み切り](https://younganimal.com/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A)
- [コミックエッセイ プチ大賞](https://www.comic-essay.com/contest/winner/)
- [コミックバンチKai 漫画賞](https://comicbunch-kai.com/article/award)

### Note

定期実行などする場合は、`robots.txt`や`meta name="robots"`を確認して使用してください。

生成したRSSフィードの使用については、各サイトの利用規約や法律を遵守してください。

実行された時点のページ情報をRSSフィードとして出力します。フィードアイテムの更新には対応していません。

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
  version             Show application version
  help                Show this help message
  --------
  andsofa             Print andsofa RSS feed
  comicboostoneshot   Print comicboost oneshot RSS feed
  comicbunckaiaward   Print comicbunchi kai award RSS feed
  comicdaysnewcomer   Print comicdays newcomer RSS feed
  comicdaysoneshot    Print comicdays oneshot RSS feed
  comicessayaward     Print comic essay award RSS feed
  comicessaygekijo    Print comic essay gekijo RSS feed
  comiplexoneshot     Print comiplex oneshot RSS feed
  kurageaward         Print kurage award RSS feed
  kuragefarm          Print kurage farm RSS feed
  matogrosso          Print matogrosso RSS feed
  toti                Print toti RSS feed
  younganimaloneshot  Print younganimal oneshot RSS feed
```

対象サイトごとのサブコマンドは以下の通りです。

| 対象サイト                                                                                                  | サブコマンド       |
| ----------------------------------------------------------------------------------------------------------- | ------------------ |
| [くらげファーム](https://kuragebunch.com/farm)                                                              | kuragefarm         |
| [コミックDAYS 新人賞](https://comic-days.com/newcomer)                                                      | comicdaysnewcomer  |
| [コミックDAYS 読み切り](https://comic-days.com/oneshot)                                                     | comicdaysoneshot   |
| [&Sofa（アンドソファ）](https://andsofa.com)                                                                | andsofa            |
| [トーチ](https://to-ti.in/product)                                                                          | toti               |
| [MATOGROSSO（マトグロッソ）](https://matogrosso.jp)                                                         | matogrosso         |
| [くらげバンチ漫画賞](https://kuragebunch.com/info/award)                                                    | kurageaward        |
| [コミックエッセイ劇場](https://www.comic-essay.com/comics)                                                  | comicessaygekijo   |
| [コミプレ 読切作品](https://viewer.heros-web.com/series/oneshot)                                            | comiplexoneshot    |
| [comicブースト 読み切り](https://comic-boost.com/genre/3)                                                   | comicboostoneshot  |
| [ヤングアニマル 読み切り](https://younganimal.com/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A) | younganimaloneshot |
| [コミックエッセイ プチ大賞](https://www.comic-essay.com/contest/winner/)                                    | comicessayaward    |
| [コミックバンチKai 漫画賞](https://comicbunch-kai.com/article/award)                                        | comicbunckaiaward  |
