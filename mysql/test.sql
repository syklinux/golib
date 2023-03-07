create database test;
create table `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `age` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '年龄',
  `name` varchar(2048) COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '姓名',
  `created_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='用户表';
insert into users(name, age, created_at) values ('cai', 27, '2019-08-07 10:50:01');