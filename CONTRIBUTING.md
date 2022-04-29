# Contribution to the AtomCI

## 前置条件
* 暂无

## Submitting a pull request


### 1. 代码贡献
代码贡献请遵循如下工作流:
1. 首次，你需要 [fork the repository on GitHub](https://github.com/go-atomci/atomci), 检出你 __fork__ 的代码至本地，然后 [将本地fork同步到上游](https://docs.github.com/cn/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork), 并且确保在每次开始变更前已经同步了最新的上游. 

2. 为你打算做的改动创建一个分支。如果你的贡献是基于已知的Issue，建议用以下方式命名你的分支："[feat|fix]-<issue-id>-<简短描述>"。
   在你本地创建分支：
   ```sh
   git checkout -b feat-17-support-statefulset
   ```

   当你的修改就绪时，就可以提交，如果与issues相关，可以在提交信息中提及Issue ID。

   ```sh
   git add .
   git commit -m "feat: #17 support statefulset ......"
   git push -u origin feat-17-support-statefulset
   ```

   __注意__: 如果在自己分支开发阶段，如有需要合并最新的`master`, 请在本地环境配置 `git rebase` ，从而避免出现无效的提交信息，具体配置指令如下: 

   ```sh
   # branch.autosetuprebase only changes the default pull “mode” for new branches that have an upstream to track. 
   git config branch.autosetuprebase always
   # use the pull.rebase config option to change the behavior for every git pull (instead of only newly-created branches)
   git config pull.rebase true
   ```
3. 如果您提交的代码涉及到数据迁移，请在`internal/migrations`目录下添加迁移脚本，并在`migration.go`中注册该脚本
   ```go
   type MigrationXXXXX struct {
   }

   func (m MigrationXXXXX) GetCreateAt() time.Time{
      //时间请设定为提交代码的当前时间，以便于系统自动记录迁移历史
      return time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
   }
   func (m Migration20220101) Upgrade(ormer orm.Ormer) error{ 
      //样例代码如下，err不为nil时，会自动回滚所有操作
      /*
      err := ormer.Raw(xxxx).Exec()
      return err
      */
      
   }

   ```

   ```go
   //在migrations.go中
   func InitMigration() {
      migrationTypes := MigrationTypes{
         //注册
         new(MigrationXXXXX),
         //...
      }

   }
   ```
4. Create a pull request to the main repository on GitHub.
5. When the reviewer makes some comments, address any feedback that comes and update the pull request.
6. When your contribution is accepted, your pull request will be approved and merged to the main branch.

### 2. 文档贡献

The workflow for documentation is similar. Please take into account a few things:

1. All documentation is written using the Markdown.
2. We store the documentation as *.md files in the [atomci-docs](https://github.com/go-atomci/atomci-press). The documentation is licensed under the [MIT License](https://github.com/go-atomci/atomci-press/blob/master/LICENSE).

After you make changes, you can use `yarn run dev` commands have a preview locally. 


## 代码审核

### 1. 自动代码审核

### 2. 通过AtomCI开发者进行代码审核


## 构建和自动化测试