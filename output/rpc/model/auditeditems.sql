CREATE TABLE `audited_items` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '已审核条目ID',
  `item_id` bigint(20) UNSIGNED NOT NULL COMMENT '所有条目ID',
  `producer` char(11) NOT NULL DEFAULT '' COMMENT '制作用户的手机号',
  `question_type` tinyint(3) UNSIGNED  NOT NULL DEFAULT 0 COMMENT '条目问题类型',
  `question` varchar(64) NOT NULL DEFAULT '' COMMENT '问题',
  `answer` varchar(12) NOT NULL DEFAULT '' COMMENT '答案',
  `disturb_answer` varchar(20) NOT NULL DEFAULT '' COMMENT '干扰答案',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '制作时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_id` (`item_id`),
  KEY `ix_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='已审核条目表';