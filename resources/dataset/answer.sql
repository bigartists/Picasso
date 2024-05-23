

CREATE TABLE `answer` (
                           `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
                           `question_id` BIGINT NOT NULL DEFAULT 0 COMMENT '问题id',
                           `context` TEXT DEFAULT NULL COMMENT '内容',
                           `author_id` BIGINT NOT NULL DEFAULT 0 COMMENT '作者id',
                           `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           INDEX `idx_question_id` (`question_id`),
                           INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='回答表';