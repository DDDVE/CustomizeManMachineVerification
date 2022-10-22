CREATE TABLE `employee` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '员工ID',
  `mobile_num` char(11) NOT NULL COMMENT '手机号',
  `employee_level` tinyint(1) UNSIGNED  NOT NULL DEFAULT 0 COMMENT '管理员等级',
  `contribution_score` int(10) UNSIGNED  NOT NULL DEFAULT 0 COMMENT '贡献积分',
  `audit_score` int(10) UNSIGNED  NOT NULL DEFAULT 0 COMMENT '审核积分',
  `registration_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_mobile_num` (`mobile_num`),
  KEY `ix_registration_time` (`registration_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='员工表';