CREATE TABLE `admin_authorization`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `app_name`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `name`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口名称',
    `url`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口地址',
    `type`         tinyint(1) NOT NULL DEFAULT '1' COMMENT '默认类型 1',
    `created_uid`  varbinary(50) NOT NULL DEFAULT '',
    `operated_uid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '',
    `deleted_at`   datetime                                                               DEFAULT NULL,
    `updated_at`   datetime                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at`   datetime                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='接口权限';