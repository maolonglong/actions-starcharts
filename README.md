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
  starcharts:
    name: Generate starcharts
    runs-on: ubuntu-latest
    steps:
      - uses: MaoLongLong/actions-starcharts@main
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          svg_path: STARCHARTS.svg
```

## 效果

[doocs/leetcode](https://github.com/doocs/leetcode) Stars 趋势（2021-09-30 生成）

![](./images/doocs_leetcode_2021_09_30.svg)

## TODO

- [x] 修复由于 GitHub V3 API 分页限制，无法获取 40K stars 以上数据的问题
- [ ] 部分操作仍然依赖 GitHub API V3，打算全部替换为 V4
- [ ] 由于 Actions 中调用 V4 API 有 1000 的次数限制，所以它暂时只支持到 100K stars 的仓库
- [ ] 为了文明使用 GitHub API，暂时没有使用多 goroutine，所以生成速度较慢
