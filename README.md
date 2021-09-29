<h1 align="center">Starcharts Action</h1>

灵感来自 [caarlos0/starcharts](https://github.com/caarlos0/starcharts)，用 Actions 避免了 [GitHub API 的速率限制](https://github.com/caarlos0/starcharts/issues/125)

## 入参

|       参数       |            描述            | 是否必传 |                默认值                |
| :--------------: | :------------------------: | :------: | :----------------------------------: |
|  `github_token`  | 用于提交时身份验证的 token |    是    |                                      |
|    `svg_path`    |       星图的保存路径       |    否    |           `STARCHARTS.svg`           |
| `commit_message` |          提交信息          |    否    | `chore: update starcharts [skip ci]` |

## 示例

新建 **.github/workflows/starcharts.yml**，内容如下：

```yml
name: Starcharts

on:
  schedule:
    - cron: "0 0 * * 0"
  workflow_dispatch:

jobs:
  update-readme:
    name: Generate starcharts
    runs-on: ubuntu-latest
    steps:
      - uses: MaoLongLong/actions-starcharts@main
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          svg_path: STARCHARTS.svg
```
