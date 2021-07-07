# golang-mysql-crud

## 规划

以Restful格式的HTTP对外提供服务，主要以POST请求实现：

- 查：查询数据库(SELECT)并返回JSON格式的数据
- 增：插入数据(INSERT)到数据库相关表中并返回执行结果
- 删：删除(DELETE)相应表中符合条件的数据
- 改：修改(UPDATE)相应表中符合条件的数据